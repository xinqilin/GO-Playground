package main

import (
	"slices"
	"testing"
)

func Example_arrayBasics() {
	arrayBasics()
	// Output:
	// [3200 0 0 0 0 0 4100] 7
	// 5 五
	// [false true false true false]
}

func Example_arrayIsValue() {
	arrayIsValue()
	// Output:
	// [85 90 15] [100 90 15]
	// [85 90 15]
	// true false
}

func Example_indexOutOfRange() {
	indexOutOfRange()
	// Output:
	// panic: runtime error: index out of range [3] with length 3
}

func Example_sliceLiteral() {
	sliceLiteral()
	// Output:
	// [排骨 雞腿 素食] 3 3
	// [0 0 0] 3 3
	// [] 0 10
}

func Example_sliceWindow() {
	sliceWindow()
	// Output:
	// [20 30 40] 3 4
	// [10 99 30 40 50]
}

func Example_sliceBounds() {
	sliceBounds()
	// Output:
	// [0 1] [4 5] [0 1 2 3 4 5]
	// [1 2] 2 5
	// [1 2 3 4]
}

func Example_appendBasics() {
	appendBasics()
	// Output:
	// [排骨 雞腿 滷蛋 紅茶 味噌湯] 5
}

func Example_appendGrowth() {
	appendGrowth()
	// Output:
	// [0 0 0] [99 0 0 4] true
	// [99 0 0] [99 0 0 4] true
}

func Example_appendPrealloc() {
	appendPrealloc()
	// Output:
	// 沒有保留容量，換過陣列： true
	// 保留 1000 格，換陣列的次數： 0
}

func Example_sharedBacking() {
	sharedBacking()
	// Output:
	// [1 2 99 4] [1 2 99]
}

func Example_fullSliceExpr() {
	fullSliceExpr()
	// Output:
	// [1 2 3 4] [1 2 99]
}

func Example_passSlice() {
	passSlice()
	// Output:
	// [90 95]
	// [排骨] 1
	// [排骨 滷蛋] 2
}

func Example_copyBasics() {
	copyBasics()
	// Output:
	// 2 [85 90]
	// [85 90 15] [0 90 15]
}

func Example_insertDelete() {
	insertDelete()
	// Output:
	// [A01 A03 A04]
	// [A01 V01 A03 A04]
}

func Example_filterInPlace() {
	filterInPlace()
	// Output:
	// [85 90 70]
	// [85 90 70 150 70]
	// [85 90 70]
}

func Example_nilVsEmpty() {
	nilVsEmpty()
	// Output:
	// true false false
	// 0 0 0
	// [] []
	// []int(nil) []int{}
	// [1] 1
}

func Example_slicesPackage() {
	slicesPackage()
	// Output:
	// true 1 -1
	// 70 120
	// [70 85 85 90 120]
	// 3 true
	// [70 85 90 120]
	// [120 90 85 70]
	// true
}

func Example_seatGrid() {
	seatGrid()
	// Output:
	// . . . .
	// . . X .
	// . . . .
	// 0 1
	// 1 3
	// 2 0
}

func Example_gridCopy() {
	gridCopy()
	// Output:
	// X
	// ""
}

func Example_retainBig() {
	retainBig()
	// Output:
	// 16 1048576
	// 16 true
}

func TestCountMoves(t *testing.T) {
	if got := countMoves(make([]int, 0, 50), 50); got != 0 {
		t.Errorf("容量足夠時 countMoves = %d, want 0", got)
	}
	if got := countMoves(nil, 50); got == 0 {
		t.Error("從 nil 開始 append 50 個，應該至少換過一次陣列")
	}
}

func TestWithItem(t *testing.T) {
	cart := make([]string, 1, 4)
	cart[0] = "排骨"
	got := withItem(cart, "滷蛋")
	if !slices.Equal(got, []string{"排骨", "滷蛋"}) {
		t.Errorf("withItem = %v", got)
	}
	if len(cart) != 1 {
		t.Errorf("呼叫端的長度變成 %d, want 1", len(cart))
	}
}

func TestNewSeats(t *testing.T) {
	seats := newSeats(2, 3)
	if len(seats) != 2 || len(seats[0]) != 3 || len(seats[1]) != 3 {
		t.Fatalf("newSeats(2, 3) 的大小錯誤：%v", seats)
	}
	seats[0][0] = "X"
	if seats[1][0] != "." {
		t.Error("每一列應該是各自獨立的 slice")
	}
}
