package mapping

import (
	"maps"
	"slices"
	"testing"
)

func TestTally(t *testing.T) {
	counts, top := Tally([]string{"排骨", "雞腿", "排骨", "素食", "雞腿"})
	if want := map[string]int{"排骨": 2, "雞腿": 2, "素食": 1}; !maps.Equal(counts, want) {
		t.Errorf("counts = %v, want %v", counts, want)
	}
	if top != "排骨" {
		t.Errorf("top = %q, want %q（次數相同時取排序較前面的）", top, "排骨")
	}
	if _, top := Tally([]string{"素食", "魚排", "魚排"}); top != "魚排" {
		t.Errorf("top = %q, want %q", top, "魚排")
	}
	if counts, top := Tally(nil); len(counts) != 0 || top != "" {
		t.Errorf("Tally(nil) = %v, %q; want 空的, \"\"", counts, top)
	}
}

func TestInvert(t *testing.T) {
	got, ok := Invert(map[string]string{"A1": "排骨", "A2": "雞腿"})
	if want := map[string]string{"排骨": "A1", "雞腿": "A2"}; !ok || !maps.Equal(got, want) {
		t.Errorf("Invert = %v, %v; want %v, true", got, ok, want)
	}
	if got, ok := Invert(map[string]string{"A1": "排骨", "A2": "排骨", "A3": "雞腿"}); ok || got != nil {
		t.Errorf("有重複的品名時 Invert = %v, %v; want nil, false", got, ok)
	}
	if got, ok := Invert(nil); !ok || len(got) != 0 {
		t.Errorf("Invert(nil) = %v, %v; want 空的, true", got, ok)
	}
}

func TestGroupByLine(t *testing.T) {
	got := GroupByLine([]Station{
		{"忠孝復興", "板南"},
		{"大安", "文湖"},
		{"市政府", "板南"},
		{"忠孝復興", "文湖"},
		{"中山", "淡水信義"},
	})
	want := map[string][]string{
		"板南":   {"市政府", "忠孝復興"},
		"文湖":   {"大安", "忠孝復興"},
		"淡水信義": {"中山"},
	}
	if !maps.EqualFunc(got, want, slices.Equal) {
		t.Errorf("GroupByLine = %v\nwant %v", got, want)
	}
	if got := GroupByLine(nil); len(got) != 0 {
		t.Errorf("GroupByLine(nil) = %v, want 空的", got)
	}
}

func TestRestock(t *testing.T) {
	stock := map[string]Item{"排骨": {Price: 85, Qty: 3}}
	if got := Restock(stock, "排骨", 5); got != 8 {
		t.Errorf("Restock(排骨, 5) = %d, want 8", got)
	}
	if got, want := stock["排骨"], (Item{Price: 85, Qty: 8}); got != want {
		t.Errorf("stock[排骨] = %+v, want %+v", got, want)
	}
	if got := Restock(stock, "魚排", 2); got != 2 {
		t.Errorf("Restock(魚排, 2) = %d, want 2", got)
	}
	if got, want := stock["魚排"], (Item{Qty: 2}); got != want {
		t.Errorf("stock[魚排] = %+v, want %+v", got, want)
	}
}

func TestSummary(t *testing.T) {
	menu := map[string]int{"雞腿": 90, "排骨": 85, "魚排": 100, "素食": 70, "滷肉": 60, "咖哩": 80}
	want := "咖哩=80, 排骨=85, 滷肉=60, 素食=70, 雞腿=90, 魚排=100"
	for range 50 {
		if got := Summary(menu); got != want {
			t.Fatalf("Summary = %q\nwant %q", got, want)
		}
	}
	if got := Summary(nil); got != "" {
		t.Errorf("Summary(nil) = %q, want \"\"", got)
	}
}

func TestCommon(t *testing.T) {
	got := Common([]string{"排骨", "雞腿", "素食", "排骨"}, []string{"魚排", "排骨", "素食", "素食"})
	if want := []string{"排骨", "素食"}; !slices.Equal(got, want) {
		t.Errorf("Common = %v, want %v", got, want)
	}
	if got := Common([]string{"排骨"}, []string{"雞腿"}); got != nil {
		t.Errorf("沒有共同品項時 Common = %v, want nil", got)
	}
	if got := Common(nil, []string{"雞腿"}); got != nil {
		t.Errorf("Common(nil, …) = %v, want nil", got)
	}
}

func TestHops(t *testing.T) {
	links := map[string][]string{
		"A1": {"A2"},
		"A2": {"A1", "A3", "B2"},
		"A3": {"A2", "A4"},
		"A4": {"A3", "B4"},
		"B1": {"B2"},
		"B2": {"B1", "A2", "B3"},
		"B3": {"B2", "B4"},
		"B4": {"B3", "A4"},
		"C1": {"C2"},
		"C2": {"C1"},
		"D1": {"D2"}, // 單向
	}
	tests := []struct {
		from, to string
		want     int
	}{
		{"A1", "A4", 3},
		{"B1", "A4", 4},
		{"A1", "B3", 3},
		{"A1", "A2", 1},
		{"A1", "A1", 0},
		{"A1", "C1", -1},
		{"A1", "Z9", -1},
		{"D1", "D2", 1},
		{"D2", "D1", -1},
	}
	for _, tt := range tests {
		if got := Hops(links, tt.from, tt.to); got != tt.want {
			t.Errorf("Hops(%s, %s) = %d, want %d", tt.from, tt.to, got, tt.want)
		}
	}
}
