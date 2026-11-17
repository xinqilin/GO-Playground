package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const demoGo = `package main

import "fmt"

// greet 的 doc comment 不會出現在頁面上。
func greet() {
	fmt.Println("哈囉 <世界> & 大家")
}

type pair struct{ a, b int }

func (p *pair) sum() int { return p.a + p.b }

const (
	small = 1
	large = 2
)
`

const demoTestGo = `package main

func Example_greet() {
	greet()
	// Output:
	// 哈囉 <世界> & 大家
}
`

const pageHTML = `<h2 id="top">示範</h2>
<pre class="code" data-src="chapters/demo/demo.go#greet"><code></code></pre>
<pre class="out" data-example="chapters/demo#Example_greet"><code>舊的輸出</code></pre>
<pre class="code" data-src="chapters/demo/demo.go#pair,pair.sum,large"><code></code></pre>
<pre class="code"><code>沒有標記的區塊不會被動到</code></pre>
<pre class="code" data-src="chapters/demo/demo_test.go"><code></code></pre>
<a href="#top">本頁</a> <a href="other.html#sec">別頁</a> <a href="https://go.dev/">外部</a>
`

const wantGreet = `func greet() {
	fmt.Println("哈囉 &lt;世界&gt; &amp; 大家")
}`

const wantPair = `type pair struct{ a, b int }

func (p *pair) sum() int { return p.a + p.b }

const (
	small = 1
	large = 2
)`

// setup 在暫存目錄建一個迷你 repo：一個 Go 檔、一個測試檔、兩個頁面。
func setup(t *testing.T, page string) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"chapters/demo/demo.go":      demoGo,
		"chapters/demo/demo_test.go": demoTestGo,
		"docs/page.html":             page,
		"docs/other.html":            `<h2 id="sec">別頁</h2>`,
	}
	for name, content := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func readPage(t *testing.T, root string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "docs", "page.html"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestSyncFillsMarkedBlocks(t *testing.T) {
	root := setup(t, pageHTML)
	problems, err := run(root, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) > 0 {
		t.Fatalf("problems = %q, want none", problems)
	}
	got := readPage(t, root)
	for _, want := range []string{
		`data-src="chapters/demo/demo.go#greet"><code>` + wantGreet + `</code></pre>`,
		`data-example="chapters/demo#Example_greet"><code>哈囉 &lt;世界&gt; &amp; 大家</code></pre>`,
		`data-src="chapters/demo/demo.go#pair,pair.sum,large"><code>` + wantPair + `</code></pre>`,
		`<pre class="code"><code>沒有標記的區塊不會被動到</code></pre>`,
		`data-src="chapters/demo/demo_test.go"><code>` + escapeHTML(strings.TrimRight(demoTestGo, "\n")) + `</code></pre>`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("page is missing:\n%s\n\npage:\n%s", want, got)
		}
	}
	if strings.Contains(got, "doc comment") {
		t.Error("doc comment should not be copied into the page")
	}
}

func TestCheckReportsDriftWithoutWriting(t *testing.T) {
	root := setup(t, pageHTML)
	if _, err := run(root, false); err != nil {
		t.Fatal(err)
	}
	if problems, err := run(root, true); err != nil || len(problems) > 0 {
		t.Fatalf("after sync: problems = %q, err = %v; want none", problems, err)
	}

	path := filepath.Join(root, "docs", "page.html")
	stale := strings.Replace(readPage(t, root), "哈囉 &lt;世界&gt; &amp; 大家</code>", "過期的輸出</code>", 1)
	if err := os.WriteFile(path, []byte(stale), 0o644); err != nil {
		t.Fatal(err)
	}
	problems, err := run(root, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(problems) != 1 || !strings.Contains(problems[0], `data-example="chapters/demo#Example_greet"`) {
		t.Fatalf("problems = %q, want one drift on Example_greet", problems)
	}
	if readPage(t, root) != stale {
		t.Error("-check must not rewrite the page")
	}
}

func TestReportsMissingSourcesAndBrokenLinks(t *testing.T) {
	page := `<h2 id="here">x</h2>
<pre class="code" data-src="chapters/demo/demo.go#nope"><code></code></pre>
<pre class="out" data-example="chapters/demo#Example_nope"><code></code></pre>
<a href="#here">ok</a> <a href="#gone">壞</a> <a href="missing.html">壞</a> <a href="other.html#nope">壞</a>
`
	root := setup(t, page)
	problems, err := run(root, true)
	if err != nil {
		t.Fatal(err)
	}
	wants := []string{"找不到 nope", "找不到 Example_nope", "#gone", "missing.html", "other.html#nope"}
	if len(problems) != len(wants) {
		t.Fatalf("got %d problems %q, want %d", len(problems), problems, len(wants))
	}
	for _, want := range wants {
		found := false
		for _, p := range problems {
			found = found || strings.Contains(p, want)
		}
		if !found {
			t.Errorf("no problem mentions %q; got %q", want, problems)
		}
	}
}
