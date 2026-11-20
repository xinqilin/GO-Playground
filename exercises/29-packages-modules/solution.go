//go:build solution

package modver

import (
	"cmp"
	"fmt"
	"strconv"
	"strings"
)

func Parse(s string) (Version, error) {
	rest, ok := strings.CutPrefix(s, "v")
	if !ok {
		return Version{}, fmt.Errorf("%q 要以 v 開頭", s)
	}
	core, pre, hasPre := strings.Cut(rest, "-")
	if hasPre && !validPre(pre) {
		return Version{}, fmt.Errorf("%q 的預發行標記不合法", s)
	}
	parts := strings.Split(core, ".")
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("%q 要有三個數字", s)
	}
	var nums [3]int
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 || (len(p) > 1 && p[0] == '0') || p[0] == '+' {
			return Version{}, fmt.Errorf("%q 的第 %d 個數字不合法", s, i+1)
		}
		nums[i] = n
	}
	return Version{Major: nums[0], Minor: nums[1], Patch: nums[2], Pre: pre}, nil
}

func validPre(pre string) bool {
	if pre == "" {
		return false
	}
	for _, r := range pre {
		if !('a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9' || r == '.' || r == '-') {
			return false
		}
	}
	return true
}

func Compare(a, b Version) int {
	if c := cmp.Or(cmp.Compare(a.Major, b.Major), cmp.Compare(a.Minor, b.Minor), cmp.Compare(a.Patch, b.Patch)); c != 0 {
		return c
	}
	switch {
	case a.Pre == b.Pre:
		return 0
	case a.Pre == "":
		return 1
	case b.Pre == "":
		return -1
	}
	as, bs := strings.Split(a.Pre, "."), strings.Split(b.Pre, ".")
	for i := range min(len(as), len(bs)) {
		if c := comparePart(as[i], bs[i]); c != 0 {
			return c
		}
	}
	return cmp.Compare(len(as), len(bs))
}

// comparePart 比較預發行標記的一段：數字比數字大小，數字比非數字舊，非數字比字串。
func comparePart(x, y string) int {
	xn, xerr := strconv.Atoi(x)
	yn, yerr := strconv.Atoi(y)
	switch {
	case xerr == nil && yerr == nil:
		return cmp.Compare(xn, yn)
	case xerr == nil:
		return -1
	case yerr == nil:
		return 1
	}
	return strings.Compare(x, y)
}

func Latest(current Version, available []Version) (Version, bool) {
	best, found := current, false
	for _, v := range available {
		if v.Major == current.Major && v.Pre == "" && Compare(v, best) > 0 {
			best, found = v, true
		}
	}
	return best, found
}

func ParseRequires(gomod string) ([]Require, error) {
	var out []Require
	inBlock := false
	for n, line := range strings.Split(gomod, "\n") {
		code, comment, _ := strings.Cut(line, "//")
		fields := strings.Fields(code)
		switch {
		case len(fields) == 0:
			continue
		case inBlock && fields[0] == ")":
			inBlock = false
			continue
		case !inBlock && fields[0] == "require" && len(fields) == 2 && fields[1] == "(":
			inBlock = true
			continue
		case !inBlock && fields[0] == "require":
			fields = fields[1:]
		case !inBlock:
			continue // module、go、replace 這些指令
		}
		if len(fields) != 2 {
			return nil, fmt.Errorf("第 %d 行：require 要是「路徑 版本」：%q", n+1, strings.TrimSpace(line))
		}
		out = append(out, Require{Path: fields[0], Version: fields[1], Indirect: strings.TrimSpace(comment) == "indirect"})
	}
	return out, nil
}

// Select 先把所有走得到的 path@version 走一遍，每個 path 記下看過的最新版本。
func Select(reqs map[string][]string) (map[string]string, error) {
	best := map[string]Version{}
	seen := map[string]bool{}
	queue := []string{"main"}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, req := range reqs[node] {
			if seen[req] {
				continue
			}
			seen[req] = true
			path, ver, ok := strings.Cut(req, "@")
			if !ok {
				return nil, fmt.Errorf("%q 要寫成 path@version", req)
			}
			v, err := Parse(ver)
			if err != nil {
				return nil, err
			}
			if old, ok := best[path]; !ok || Compare(v, old) > 0 {
				best[path] = v
			}
			queue = append(queue, req)
		}
	}
	out := make(map[string]string, len(best))
	for path, v := range best {
		out[path] = v.String()
	}
	return out, nil
}

// InternalAllowed 找最後一段 internal，並且比較的是完整的路徑段，不是字串開頭。
func InternalAllowed(importer, imported string) bool {
	parts := strings.Split(imported, "/")
	last := -1
	for i, p := range parts {
		if p == "internal" {
			last = i
		}
	}
	if last < 0 {
		return true
	}
	parent := strings.Join(parts[:last], "/")
	return importer == parent || strings.HasPrefix(importer, parent+"/")
}
