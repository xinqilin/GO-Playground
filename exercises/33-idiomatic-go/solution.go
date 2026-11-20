//go:build solution

package idiom

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Stutters 回傳匯出名稱 name 從呼叫端看起來是不是重複了 package 名稱 pkg。
func Stutters(pkg, name string) bool {
	if len(name) <= len(pkg) || !strings.EqualFold(name[:len(pkg)], pkg) {
		return false
	}
	r, _ := utf8.DecodeRuneInString(name[len(pkg):])
	return unicode.IsUpper(r)
}

// DocProblems 檢查匯出名稱 name 的 doc comment，依序回傳發現的問題。
func DocProblems(name, doc string) []string {
	doc = strings.TrimSpace(doc)
	if doc == "" {
		return []string{ProblemMissing}
	}
	var problems []string
	if !strings.HasPrefix(doc, name+" ") {
		problems = append(problems, ProblemName)
	}
	if !strings.HasSuffix(doc, "。") && !strings.HasSuffix(doc, ".") {
		problems = append(problems, ProblemPeriod)
	}
	return problems
}

// Add 記錄 dish 又賣了 n 份。
func (t *Tally) Add(dish string, n int) {
	if t.counts == nil {
		t.counts = make(map[string]int)
	}
	t.counts[dish] += n
}

// Count 回傳 dish 賣了幾份。
func (t *Tally) Count(dish string) int {
	return t.counts[dish]
}

// Top 回傳賣最多的 n 道菜。
func (t *Tally) Top(n int) []string {
	dishes := slices.SortedFunc(maps.Keys(t.counts), func(a, b string) int {
		return cmp.Or(cmp.Compare(t.counts[b], t.counts[a]), strings.Compare(a, b))
	})
	return dishes[:min(n, len(dishes))]
}

// WithPrefix 設定號碼的前綴。
func WithPrefix(p string) Option {
	return func(m *Machine) { m.prefix = p }
}

// WithStart 設定第一個號碼。
func WithStart(n int) Option {
	return func(m *Machine) { m.next = n }
}

// WithWidth 設定號碼補零後的位數。
func WithWidth(w int) Option {
	return func(m *Machine) { m.width = w }
}

// New 用預設值建立 Machine，再依序套用 opts。
func New(opts ...Option) *Machine {
	m := &Machine{prefix: "A", next: 1, width: 3}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Next 回傳下一張號碼牌。
func (m *Machine) Next() string {
	s := fmt.Sprintf("%s%0*d", m.prefix, m.width, m.next)
	m.next++
	return s
}

// MixedCaps 把 snake_case 的名稱轉成 Go 的寫法。
func MixedCaps(snake string, exported bool) string {
	var b strings.Builder
	for i, part := range strings.FieldsFunc(snake, func(r rune) bool { return r == '_' }) {
		switch {
		case i == 0 && !exported:
			b.WriteString(part)
		case Initialisms[part]:
			b.WriteString(strings.ToUpper(part))
		default:
			b.WriteString(strings.ToUpper(part[:1]) + part[1:])
		}
	}
	return b.String()
}

// QueueOf 用 nums 的複本建立一條隊伍。
func QueueOf(nums []int) *Queue {
	return &Queue{waiting: slices.Clone(nums)}
}

// Join 讓號碼 n 排到隊伍的最後面。
func (q *Queue) Join(n int) {
	q.waiting = append(q.waiting, n)
}

// Waiting 回傳目前排隊號碼的複本。
func (q *Queue) Waiting() []int {
	return slices.Clone(q.waiting)
}
