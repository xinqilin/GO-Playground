package main

import "fmt"

// passFail 依分數決定印哪一句：條件成立走 if 的大括號，不成立走 else 的大括號。
func passFail(score int) {
	if score >= 60 {
		fmt.Println(score, "分：及格")
	} else {
		fmt.Println(score, "分：不及格")
	}
}

func passFailDemo() {
	passFail(72)
	passFail(59)
	passFail(60)
}

// umbrella 只有 if、沒有 else：條件不成立就跳過大括號，直接往下走。
func umbrella(rain bool) {
	fmt.Println("準備出門")
	if rain {
		fmt.Println("帶傘")
	}
	fmt.Println("鎖門出發")
}

func umbrellaDemo() {
	umbrella(true)
	umbrella(false)
}

// grade 從最高的等第開始問，第一個成立的條件決定結果，後面的就不再問。
func grade(score int) {
	level := ""
	if score >= 90 {
		level = "優"
	} else if score >= 80 {
		level = "甲"
	} else if score >= 70 {
		level = "乙"
	} else if score >= 60 {
		level = "丙"
	} else {
		level = "丁"
	}
	fmt.Println(score, level)
}

func gradeDemo() {
	grade(95)
	grade(83)
	grade(61)
	grade(40)
}

// wrongGrade 把條件的順序寫反了：95 分在第一個條件就成立，後面的條件根本沒機會問。
func wrongGrade(score int) {
	level := ""
	if score >= 60 {
		level = "丙"
	} else if score >= 70 {
		level = "乙"
	} else if score >= 80 {
		level = "甲"
	} else if score >= 90 {
		level = "優"
	} else {
		level = "丁"
	}
	fmt.Println(score, level)
}

func wrongGradeDemo() {
	wrongGrade(95)
}

// ticket 依年齡和是不是學生決定票種：未滿 6 歲免費；未滿 12 歲或 65 歲以上半票；
// 其他人是學生就買學生票，不是就買全票。
func ticket(age int, student bool) {
	kind := ""
	price := 0
	if age < 6 {
		kind = "免費"
		price = 0
	} else if age < 12 || age >= 65 {
		kind = "半票"
		price = 60
	} else if student {
		kind = "學生票"
		price = 90
	} else {
		kind = "全票"
		price = 120
	}
	fmt.Println(age, "歲", kind, price, "元")
}

func ticketDemo() {
	ticket(3, false)
	ticket(8, false)
	ticket(20, true)
	ticket(20, false)
	ticket(70, true)
}
