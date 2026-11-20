package main

import (
	"fmt"
	"slices"
	"strings"
)

// newSeats 建立 rows 列、cols 欄的座位表，每個位置一開始是 "."。
func newSeats(rows, cols int) [][]string {
	seats := make([][]string, rows)
	for r := range seats {
		seats[r] = make([]string, cols)
		for c := range seats[r] {
			seats[r][c] = "."
		}
	}
	return seats
}

// printSeats 一列印一行。
func printSeats(seats [][]string) {
	for _, row := range seats {
		fmt.Println(strings.Join(row, " "))
	}
}

// seatGrid 示範二維 slice：外層 slice 的每個元素又是一個 slice。
func seatGrid() {
	seats := newSeats(3, 4)
	seats[1][2] = "X"
	printSeats(seats)

	daily := [][]string{{"排骨"}, {"雞腿", "素食", "排骨"}, {}} // 每一列可以不一樣長
	for day, orders := range daily {
		fmt.Println(day, len(orders))
	}
}

// gridCopy 示範 slices.Clone 只複製外層：兩份座位表還是共用同樣的列。
// 二維陣列則是一整塊值，指定時會全部複製。
func gridCopy() {
	seats := newSeats(2, 3)
	shallow := slices.Clone(seats)
	shallow[0][0] = "X"
	fmt.Println(seats[0][0])

	var board [2][3]string
	copied := board
	copied[0][0] = "X"
	fmt.Printf("%q\n", board[0][0])
}

// retainBig 示範從大 slice 切出的小 slice 會讓整個大陣列留在記憶體裡；複製出來，大陣列才能被回收。
func retainBig() {
	big := make([]byte, 1<<20) // 1 MB，例如整個讀進來的檔案
	header := big[:16]
	fmt.Println(len(header), cap(header))
	kept := slices.Clone(header)
	fmt.Println(len(kept), cap(kept) < 1024)
}
