// Command docsync 讓 docs/*.html 裡的程式碼和輸出，跟 repo 裡測試過的 Go 程式保持一致。
//
// 在 repo 根目錄執行：
//
//	go run ./tools/docsync          依照程式碼改寫 HTML
//	go run ./tools/docsync -check   只檢查；內容不一致或連結壞掉就 exit 1（CI 用）
//
// HTML 裡用兩種標記：
//
//	<pre class="code" data-src="chapters/03-control-flow/if.go#canShip,order"><code></code></pre>
//	<pre class="out" data-example="chapters/03-control-flow#Example_guardClause"><code></code></pre>
//
// data-src 是「檔案#名稱」。名稱可以是函式、型別、常數或變數，方法寫成 Type.Method，
// 多個名稱用逗號分開；取出的是宣告本身，不含上方的 doc comment。只寫「檔案」時取出整個檔案。
// data-example 是「目錄#Example 函式名」，取出該 Example 的 // Output: 內容。
// 另外會檢查每個相對連結指到的檔案和 #anchor 都存在。
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	check := flag.Bool("check", false, "只檢查，不改寫檔案")
	root := flag.String("root", ".", "repo 根目錄")
	flag.Parse()

	problems, err := run(*root, *check)
	if err != nil {
		fmt.Fprintln(os.Stderr, "docsync:", err)
		os.Exit(2)
	}
	for _, p := range problems {
		fmt.Fprintln(os.Stderr, p)
	}
	if len(problems) > 0 {
		os.Exit(1)
	}
}

// run 處理 root/docs 底下的每個 HTML 檔。check 為 true 時不寫檔，把需要更新的區塊當成問題回報。
func run(root string, check bool) ([]string, error) {
	pages, err := filepath.Glob(filepath.Join(root, "docs", "*.html"))
	if err != nil {
		return nil, err
	}
	src := newSources(root)
	var problems []string
	for _, page := range pages {
		data, err := os.ReadFile(page)
		if err != nil {
			return nil, err
		}
		name := filepath.Base(page)
		out, changed, probs := syncPage(string(data), src)
		for _, p := range probs {
			problems = append(problems, name+": "+p)
		}
		if len(changed) == 0 {
			continue
		}
		if check {
			for _, ref := range changed {
				problems = append(problems, fmt.Sprintf("%s: %s 和程式碼不一致，請執行 go run ./tools/docsync", name, ref))
			}
			continue
		}
		if err := os.WriteFile(page, []byte(out), 0o644); err != nil {
			return nil, err
		}
	}
	linkProblems, err := checkLinks(filepath.Join(root, "docs"), pages)
	if err != nil {
		return nil, err
	}
	return append(problems, linkProblems...), nil
}

var (
	blockRE   = regexp.MustCompile(`(?s)<pre\b([^>]*)><code>(.*?)</code></pre>`)
	srcRE     = regexp.MustCompile(`\sdata-src="([^"]*)"`)
	exampleRE = regexp.MustCompile(`\sdata-example="([^"]*)"`)
)

// syncPage 把頁面裡有標記的區塊換成最新內容，回傳新頁面、內容有變動的標記，以及找不到來源之類的問題。
func syncPage(page string, src *sources) (out string, changed, problems []string) {
	var b strings.Builder
	last := 0
	for _, m := range blockRE.FindAllStringSubmatchIndex(page, -1) {
		attrs, inner := page[m[2]:m[3]], page[m[4]:m[5]]
		var ref, want string
		var err error
		if s := srcRE.FindStringSubmatch(attrs); s != nil {
			ref = `data-src="` + s[1] + `"`
			want, err = src.code(s[1])
		} else if s := exampleRE.FindStringSubmatch(attrs); s != nil {
			ref = `data-example="` + s[1] + `"`
			want, err = src.output(s[1])
		} else {
			continue
		}
		if err != nil {
			problems = append(problems, err.Error())
			continue
		}
		want = escapeHTML(want)
		if want == inner {
			continue
		}
		changed = append(changed, ref)
		b.WriteString(page[last:m[4]])
		b.WriteString(want)
		last = m[5]
	}
	b.WriteString(page[last:])
	return b.String(), changed, problems
}

var htmlEscaper = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")

func escapeHTML(s string) string { return htmlEscaper.Replace(s) }

// sources 讀取並快取 Go 原始檔和 Example 的輸出。
type sources struct {
	root     string
	fset     *token.FileSet
	files    map[string]*parsedFile
	examples map[string]map[string]string
}

type parsedFile struct {
	src  []byte
	file *ast.File
}

func newSources(root string) *sources {
	return &sources{
		root:     root,
		fset:     token.NewFileSet(),
		files:    map[string]*parsedFile{},
		examples: map[string]map[string]string{},
	}
}

// code 依 data-src 取出原始碼：「檔案#名稱1,名稱2」取出這些宣告，之間空一行；
// 只寫「檔案」就是整個檔案，這時不解析，所以範本這類不是 Go 的檔案也可以。
func (s *sources) code(ref string) (string, error) {
	path, names, ok := strings.Cut(ref, "#")
	if ok && names == "" {
		return "", fmt.Errorf("data-src %q 的 # 後面要寫名稱", ref)
	}
	if !ok {
		data, err := os.ReadFile(filepath.Join(s.root, filepath.FromSlash(path)))
		if err != nil {
			return "", err
		}
		return strings.TrimRight(string(data), "\n"), nil
	}
	pf, err := s.parse(path)
	if err != nil {
		return "", err
	}
	var parts []string
	for _, name := range strings.Split(names, ",") {
		d := findDecl(pf.file, name)
		if d == nil {
			return "", fmt.Errorf("%s 裡找不到 %s", path, name)
		}
		start := s.fset.Position(d.Pos()).Offset
		end := s.fset.Position(d.End()).Offset
		parts = append(parts, string(pf.src[start:end]))
	}
	return strings.Join(parts, "\n\n"), nil
}

func (s *sources) parse(path string) (*parsedFile, error) {
	if pf, ok := s.files[path]; ok {
		return pf, nil
	}
	data, err := os.ReadFile(filepath.Join(s.root, filepath.FromSlash(path)))
	if err != nil {
		return nil, err
	}
	f, err := parser.ParseFile(s.fset, path, data, parser.SkipObjectResolution)
	if err != nil {
		return nil, err
	}
	pf := &parsedFile{src: data, file: f}
	s.files[path] = pf
	return pf, nil
}

// findDecl 找出名稱對應的頂層宣告；Type.Method 代表方法。型別、常數、變數回傳整個宣告區塊。
func findDecl(f *ast.File, name string) ast.Decl {
	recv, method, isMethod := strings.Cut(name, ".")
	for _, d := range f.Decls {
		switch d := d.(type) {
		case *ast.FuncDecl:
			if isMethod && d.Recv != nil && d.Name.Name == method && recvTypeName(d.Recv) == recv {
				return d
			}
			if !isMethod && d.Recv == nil && d.Name.Name == name {
				return d
			}
		case *ast.GenDecl:
			if !isMethod && declares(d, name) {
				return d
			}
		}
	}
	return nil
}

func declares(d *ast.GenDecl, name string) bool {
	for _, spec := range d.Specs {
		switch sp := spec.(type) {
		case *ast.TypeSpec:
			if sp.Name.Name == name {
				return true
			}
		case *ast.ValueSpec:
			for _, n := range sp.Names {
				if n.Name == name {
					return true
				}
			}
		}
	}
	return false
}

// recvTypeName 取出 receiver 的型別名稱，去掉指標和型別參數：*Stack[T] → Stack。
func recvTypeName(fl *ast.FieldList) string {
	t := fl.List[0].Type
	for {
		switch x := t.(type) {
		case *ast.StarExpr:
			t = x.X
		case *ast.IndexExpr:
			t = x.X
		case *ast.IndexListExpr:
			t = x.X
		case *ast.Ident:
			return x.Name
		default:
			return ""
		}
	}
}

// output 依 data-example（目錄#Example名稱）取出 Example 的 // Output: 內容。
func (s *sources) output(ref string) (string, error) {
	dir, name, ok := strings.Cut(ref, "#")
	if !ok || name == "" {
		return "", fmt.Errorf("data-example %q 要寫成「目錄#Example名稱」", ref)
	}
	exs, err := s.loadExamples(dir)
	if err != nil {
		return "", err
	}
	out, ok := exs[name]
	if !ok {
		return "", fmt.Errorf("%s 的測試檔裡找不到 %s（或它沒有 // Output:）", dir, name)
	}
	// go test 比對輸出時會去掉頭尾空白，這裡用同樣的規則。
	return strings.TrimSpace(out), nil
}

func (s *sources) loadExamples(dir string) (map[string]string, error) {
	if m, ok := s.examples[dir]; ok {
		return m, nil
	}
	paths, err := filepath.Glob(filepath.Join(s.root, filepath.FromSlash(dir), "*_test.go"))
	if err != nil {
		return nil, err
	}
	var files []*ast.File
	for _, p := range paths {
		f, err := parser.ParseFile(s.fset, p, nil, parser.ParseComments|parser.SkipObjectResolution)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	m := map[string]string{}
	for _, ex := range doc.Examples(files...) {
		if ex.Output != "" || ex.EmptyOutput {
			m["Example"+ex.Name] = ex.Output
		}
	}
	s.examples[dir] = m
	return m, nil
}

var (
	hrefRE = regexp.MustCompile(`\shref="([^"]*)"`)
	idRE   = regexp.MustCompile(`\sid="([^"]*)"`)
)

// checkLinks 檢查每個相對連結：檔案要存在，指到 HTML 的 #anchor 也要存在。
// pre 區塊裡的 href 是程式碼或輸出的一部分，不是連結，不檢查。
func checkLinks(docs string, pages []string) ([]string, error) {
	ids := map[string]map[string]bool{}
	idsOf := func(name string) (map[string]bool, error) {
		if set, ok := ids[name]; ok {
			return set, nil
		}
		data, err := os.ReadFile(filepath.Join(docs, filepath.FromSlash(name)))
		if err != nil {
			return nil, err
		}
		set := map[string]bool{}
		for _, m := range idRE.FindAllStringSubmatch(string(data), -1) {
			set[m[1]] = true
		}
		ids[name] = set
		return set, nil
	}

	var problems []string
	for _, page := range pages {
		data, err := os.ReadFile(page)
		if err != nil {
			return nil, err
		}
		base := filepath.Base(page)
		for _, m := range hrefRE.FindAllStringSubmatch(blockRE.ReplaceAllString(string(data), ""), -1) {
			href := m[1]
			if strings.Contains(href, ":") {
				continue // https:、mailto: 之類的外部連結不檢查
			}
			path, frag, _ := strings.Cut(href, "#")
			target := base
			if path != "" {
				target = path
				if _, err := os.Stat(filepath.Join(docs, filepath.FromSlash(path))); err != nil {
					problems = append(problems, fmt.Sprintf("%s: 連結 %s 指到不存在的檔案", base, href))
					continue
				}
			}
			if frag == "" || !strings.HasSuffix(target, ".html") {
				continue
			}
			set, err := idsOf(target)
			if err != nil {
				return nil, err
			}
			if !set[frag] {
				problems = append(problems, fmt.Sprintf("%s: 連結 %s 的 #%s 不存在", base, href, frag))
			}
		}
	}
	return problems, nil
}
