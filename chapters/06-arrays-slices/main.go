// 第 06 章：陣列與 slice。
//
// 陣列是一排固定數量的格子；slice 是看向陣列的一扇窗，記著從哪裡開始看、看幾格、最多能看到哪。
// 每個示範都是一個函式，main 依照頁面上的順序呼叫它們；
// main_test.go 用 Example 鎖住每個示範的輸出，頁面上的輸出就是從那裡同步過去的。
// append 擴充後的容量是 runtime 的實作細節，所以示範裡只比較容量的大小關係，不印出確切的數字。
//
// 執行：go run ./chapters/06-arrays-slices
package main

import "fmt"

func main() {
	section("1. 陣列")
	arrayBasics()
	arrayIsValue()
	indexOutOfRange()

	section("2. slice 是一扇窗")
	sliceLiteral()
	sliceWindow()
	sliceBounds()

	section("3. append")
	appendBasics()
	appendGrowth()
	appendPrealloc()

	section("4. 共用底層陣列")
	sharedBacking()
	fullSliceExpr()
	passSlice()

	section("5. copy、插入與刪除")
	copyBasics()
	insertDelete()
	filterInPlace()

	section("6. nil slice 與空 slice")
	nilVsEmpty()

	section("7. slices 套件")
	slicesPackage()

	section("8. 二維 slice")
	seatGrid()
	gridCopy()
	retainBig()
}

func section(title string) {
	fmt.Println("==", title, "==")
}
