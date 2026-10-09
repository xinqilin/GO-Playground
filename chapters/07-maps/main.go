// 第 07 章：map。
//
// map 是一排有名牌的置物櫃：用 key 直接找到櫃子，不用從頭找。
// 執行：go run ./chapters/07-maps
package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

func main() {
	fmt.Println("== 基本操作 ==")
	basics()

	fmt.Println("== 找不到的 key ==")
	missingKey()

	fmt.Println("== 走訪順序不固定 ==")
	sortedIteration()

	fmt.Println("== 數字數 ==")
	fmt.Println(wordCount("go is fun and go is fast"))

	fmt.Println("== nil map ==")
	nilMap()
}

// basics 示範新增、覆寫、刪除與長度。
func basics() {
	ages := map[string]int{"小明": 18, "小華": 20}
	ages["小美"] = 22 // 新增
	ages["小明"] = 19 // 覆寫
	delete(ages, "小華")
	fmt.Println(ages, len(ages)) // fmt 印 map 時會把 key 排序，方便閱讀
}

// missingKey 示範找不到 key 時拿到零值，以及用 comma-ok 分辨「不存在」和「值剛好是零」。
func missingKey() {
	stock := map[string]int{"蘋果": 0}
	fmt.Println(stock["香蕉"])
	_, ok1 := stock["蘋果"]
	_, ok2 := stock["香蕉"]
	fmt.Println(ok1, ok2)
}

// sortedIteration 示範走訪順序每次都可能不同；需要固定順序就先排序 key。
func sortedIteration() {
	m := map[string]int{"c": 3, "a": 1, "b": 2}
	for _, k := range slices.Sorted(maps.Keys(m)) {
		fmt.Println(k, m[k])
	}
}

// wordCount 計算每個字出現幾次。
func wordCount(s string) map[string]int {
	counts := make(map[string]int)
	for _, w := range strings.Fields(s) {
		counts[w]++ // key 不存在時，從零值 0 開始加
	}
	return counts
}

// nilMap 示範 nil map：可以讀，但寫入會 panic。recover 在第 11 章細講。
func nilMap() {
	var m map[string]int
	fmt.Println(m["x"], len(m), m == nil)
	defer func() {
		fmt.Println("寫入 nil map：", recover())
	}()
	m["x"] = 1
}
