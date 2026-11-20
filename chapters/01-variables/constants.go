package main

import "fmt"

// constBasics 示範常數：值在編譯時就決定，之後不能修改。
func constBasics() {
	const shop = "阿姨便當"
	const bentoPrice = 85
	const taxRate = 0.05
	fmt.Println(shop, bentoPrice, bentoPrice*taxRate)
}

// untypedConst 示範沒有型別的常數：用在哪裡，就變成那裡需要的型別。
func untypedConst() {
	const n = 3
	var qty int = n
	var weight float64 = n
	var small int8 = n
	fmt.Println(qty, weight, small)

	const huge = 1 << 100 // 遠超過 int64，但只要最後用到的結果放得進去就好
	fmt.Println(huge >> 98)
}

// Weekday 是自訂型別，底層型別是 int。Go 沒有 enum 關鍵字，慣例是 const 搭配 iota。
type Weekday int

const (
	Sunday  Weekday = iota // 0
	Monday                 // 1：沒寫的話沿用上一行的寫法，iota 自動加 1
	Tuesday                // 2
	Wednesday
	Thursday
	Friday
	Saturday
)

// iotaWeekday 示範 iota 從 0 開始自動編號。
func iotaWeekday() {
	fmt.Println(Sunday, Monday, Saturday)
	today := Friday
	fmt.Println(today+1 == Saturday) // 可以和沒有型別的常數 1 運算
}

// Priority 從 1 開始編號，讓零值 0 代表「還沒設定」。
type Priority int

const (
	Low Priority = iota + 1
	Normal
	High
)

// Note 是便當的備註。每個備註佔一個位元，同一份便當可以有好幾個備註。
type Note uint8

const (
	Spicy     Note = 1 << iota // 0001
	NoOnion                    // 0010
	ExtraRice                  // 0100
	Takeout                    // 1000
)

// iotaPatterns 示範 iota 的兩種常見寫法：從 1 開始，以及位元旗標。
func iotaPatterns() {
	var p Priority
	fmt.Println(p, Low, Normal, High)
	fmt.Println(Spicy, NoOnion, ExtraRice, Takeout)
}
