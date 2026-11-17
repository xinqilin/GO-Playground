package pointers

import (
	"slices"
	"testing"
)

func TestSwap(t *testing.T) {
	a, b := 85, 90
	Swap(&a, &b)
	if a != 90 || b != 85 {
		t.Errorf("Swap 之後 a, b = %d, %d; want 90, 85", a, b)
	}
	c := 7
	Swap(&c, &c)
	if c != 7 {
		t.Errorf("Swap(&c, &c) 之後 c = %d; want 7", c)
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		v, lo, hi, want int
	}{
		{5, 0, 10, 5},
		{-3, 0, 10, 0},
		{15, 0, 10, 10},
		{0, 0, 10, 0},
		{10, 0, 10, 10},
	}
	for _, tt := range tests {
		v := tt.v
		Clamp(&v, tt.lo, tt.hi)
		if v != tt.want {
			t.Errorf("Clamp(%d, %d, %d) 之後是 %d; want %d", tt.v, tt.lo, tt.hi, v, tt.want)
		}
	}
	Clamp(nil, 0, 10) // 不能 panic
}

func TestApplyPatch(t *testing.T) {
	p := Profile{Name: "小明", Phone: "0912", Points: 120}
	newPhone := "0987"
	zero := 0
	ApplyPatch(&p, ProfilePatch{Phone: &newPhone, Points: &zero})
	want := Profile{Name: "小明", Phone: "0987", Points: 0}
	if p != want {
		t.Errorf("ApplyPatch 之後 = %+v; want %+v", p, want)
	}

	ApplyPatch(&p, ProfilePatch{})
	if p != want {
		t.Errorf("空的 patch 不應該改變資料，卻變成 %+v", p)
	}

	name := "阿美"
	ApplyPatch(&p, ProfilePatch{Name: &name})
	if p.Name != "阿美" || p.Phone != "0987" {
		t.Errorf("只改名字之後 = %+v", p)
	}
	name = "改掉"
	if p.Name != "阿美" {
		t.Errorf("ApplyPatch 應該複製值，不應該在之後跟著 patch 變動：%+v", p)
	}
}

func TestCheapest(t *testing.T) {
	items := []Item{
		{"雞腿", 90, 5},
		{"素食", 70, 3},
		{"排骨", 85, 4},
		{"滷肉", 70, 6},
	}
	got := Cheapest(items)
	if got == nil {
		t.Fatal("Cheapest 回傳 nil")
	}
	if got.Name != "素食" {
		t.Errorf("Cheapest = %+v; want 素食", *got)
	}
	got.Stock = 0
	if items[1].Stock != 0 {
		t.Errorf("透過回傳的指標把庫存改成 0，items[1].Stock 卻還是 %d", items[1].Stock)
	}
	if got != &items[1] {
		t.Error("回傳的指標應該指向 items[1] 本身")
	}
	if Cheapest(nil) != nil {
		t.Error("空的 items 應該回傳 nil")
	}
}

// build 依序建立一條隊伍，給測試用。
func build(names ...string) *Node {
	var head, tail *Node
	for _, n := range names {
		node := &Node{Name: n}
		if head == nil {
			head = node
		} else {
			tail.Next = node
		}
		tail = node
	}
	return head
}

// collect 依序收集名字，給測試用。最多走 100 步，避免寫錯成環狀時測試卡住。
func collect(head *Node) []string {
	var names []string
	for n := head; n != nil && len(names) < 100; n = n.Next {
		names = append(names, n.Name)
	}
	return names
}

func TestPushBack(t *testing.T) {
	head := PushBack(nil, "A01")
	if head == nil {
		t.Fatal("PushBack(nil, \"A01\") 回傳 nil")
	}
	first := head
	head = PushBack(head, "A02")
	head = PushBack(head, "A03")
	if head != first {
		t.Error("隊伍不是空的時，PushBack 應該回傳原本的第一個人")
	}
	if got, want := collect(head), []string{"A01", "A02", "A03"}; !slices.Equal(got, want) {
		t.Errorf("隊伍 = %v; want %v", got, want)
	}
}

func TestNames(t *testing.T) {
	if got, want := Names(build("A01", "A02", "A03")), []string{"A01", "A02", "A03"}; !slices.Equal(got, want) {
		t.Errorf("Names = %v; want %v", got, want)
	}
	if got := Names(nil); got != nil {
		t.Errorf("Names(nil) = %v; want nil", got)
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		names  []string
		remove string
		want   []string
	}{
		{[]string{"A01", "A02", "A03"}, "A01", []string{"A02", "A03"}},
		{[]string{"A01", "A02", "A03"}, "A02", []string{"A01", "A03"}},
		{[]string{"A01", "A02", "A03"}, "A03", []string{"A01", "A02"}},
		{[]string{"A01", "A02", "A01"}, "A01", []string{"A02", "A01"}},
		{[]string{"A01"}, "A01", nil},
		{[]string{"A01", "A02"}, "B01", []string{"A01", "A02"}},
		{nil, "A01", nil},
	}
	for _, tt := range tests {
		got := collect(Remove(build(tt.names...), tt.remove))
		if !slices.Equal(got, tt.want) {
			t.Errorf("從 %v 移除 %s = %v; want %v", tt.names, tt.remove, got, tt.want)
		}
	}
}

func TestDuplicate(t *testing.T) {
	note := "不要香菜"
	orig := Order{ID: 1, Note: &note}
	dup := Duplicate(orig, 2)
	if dup.ID != 2 || dup.Note == nil || *dup.Note != "不要香菜" {
		t.Fatalf("Duplicate = {ID:%d Note:%v}; want ID 2、備註「不要香菜」", dup.ID, dup.Note)
	}
	*dup.Note = "加辣"
	if *orig.Note != "不要香菜" {
		t.Errorf("改了複本的備註，原本的備註變成 %q", *orig.Note)
	}

	plain := Duplicate(Order{ID: 3}, 4)
	if plain.ID != 4 || plain.Note != nil {
		t.Errorf("沒有備註的訂單複製後 = {ID:%d Note:%v}; want ID 4、Note nil", plain.ID, plain.Note)
	}
}
