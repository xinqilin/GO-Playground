package main

import (
	"maps"
	"testing"
)

func Example_mapBasics() {
	mapBasics()
	// Output:
	// map[排骨:95 素食:70] 2
	// 95
}

func Example_makeMap() {
	makeMap()
	// Output:
	// map[排骨:10 雞腿:6] 2
}

func Example_nilMap() {
	nilMap()
	// Output:
	// true 0 0
	// panic: assignment to entry in nil map
}

func Example_commaOK() {
	commaOK()
	// Output:
	// 0 0
	// 0 true
	// 0 false
	// 有賣雞腿
}

func Example_countOrders() {
	countOrders()
	// Output:
	// map[排骨:3 素食:1 雞腿:1]
}

func Example_groupOrders() {
	groupOrders()
	// Output:
	// map[1:[排骨 素食] 2:[排骨] 3:[雞腿 滷蛋]]
	// [排骨 素食] 0
}

func Example_rangeMap() {
	rangeMap()
	// Output:
	// 245
}

func Example_sortedRange() {
	sortedRange()
	// Output:
	// 排骨 85
	// 素食 70
	// 雞腿 90
	// 魚排 100
	// [排骨 素食 雞腿 魚排]
}

func Example_sortByValue() {
	sortByValue()
	// Output:
	// 雞腿 30
	// 排骨 12
	// 素食 12
	// 魚排 7
}

func Example_deleteWhileRange() {
	deleteWhileRange()
	// Output:
	// map[雞腿:3 魚排:5]
	// map[魚排:5]
	// map[] 0
}

func Example_structKey() {
	structKey()
	// Output:
	// 小明 2
	// false
	// 2 2
}

func Example_nanKey() {
	nanKey()
	// Output:
	// 2 false
}

func Example_structValue() {
	structValue()
	// Output:
	// 10 9
	// 9
}

func Example_pointerValue() {
	pointerValue()
	// Output:
	// 9
	// 沒有魚排 true
}

func Example_mapIsReference() {
	mapIsReference()
	// Output:
	// 3
	// map[排骨:90 素食:75 雞腿:95]
	// 90 0
	// false
}

func Example_mapAsSet() {
	mapAsSet()
	// Output:
	// true 2
	// [雞腿 魚排]
	// [A01 A02]
}

func Example_indexByName() {
	indexByName()
	// Output:
	// 2 3
}

func TestBuildIndex(t *testing.T) {
	names := []string{"排骨", "雞腿", "素食"}
	idx := buildIndex(names)
	if len(idx) != len(names) {
		t.Fatalf("len = %d, want %d", len(idx), len(names))
	}
	for i, name := range names {
		if got, ok := idx[name]; !ok || got != i {
			t.Errorf("idx[%q] = %d, %v; want %d, true", name, got, ok, i)
		}
	}
	if got := buildIndex(nil); len(got) != 0 {
		t.Errorf("buildIndex(nil) = %v, want 空的", got)
	}
}

func TestRaisePrices(t *testing.T) {
	menu := map[string]int{"排骨": 85, "雞腿": 90}
	raisePrices(menu, 10)
	if want := map[string]int{"排骨": 95, "雞腿": 100}; !maps.Equal(menu, want) {
		t.Errorf("raisePrices 之後 = %v, want %v", menu, want)
	}
}
