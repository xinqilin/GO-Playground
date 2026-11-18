package main

import "fmt"

// Shop 是所有店家共有的資料。
type Shop struct {
	Name  string
	Phone string
}

// Contact 回傳店名和電話。
func (s Shop) Contact() string {
	return s.Name + " " + s.Phone
}

// BentoShop 嵌入 Shop：Shop 的欄位和方法都被提升到 BentoShop。
type BentoShop struct {
	Shop
	Menu []string
}

func embedDemo() {
	b := BentoShop{
		Shop: Shop{Name: "阿姨便當", Phone: "02-2345-6789"},
		Menu: []string{"排骨飯", "雞腿飯"},
	}
	fmt.Println(b.Name)
	fmt.Println(b.Contact())
	fmt.Println(b.Shop.Name)
}

// promotedKeys：Go 1.27 起，struct 字面值的 key 可以直接寫被提升的欄位。
func promotedKeys() {
	b := BentoShop{Name: "巷口便當", Phone: "02-8765-4321"}
	fmt.Println(b.Contact())
}

// DrinkShop 自己也定義了 Contact：外層的方法會遮住被嵌入的 Shop.Contact。
type DrinkShop struct {
	Shop
	LineID string
}

// Contact 先呼叫被遮住的 Shop.Contact，再補上 LINE 帳號。
func (d DrinkShop) Contact() string {
	return d.Shop.Contact() + " LINE:" + d.LineID
}

func shadowDemo() {
	d := DrinkShop{Shop: Shop{Name: "手搖飲", Phone: "02-1111-2222"}, LineID: "@drink"}
	fmt.Println(d.Contact())
	fmt.Println(d.Shop.Contact())
}

// Greeter 的 Welcome 呼叫自己的 Title。
type Greeter struct{}

// Title 是一般客人的稱呼。
func (Greeter) Title() string { return "客人" }

// Welcome 呼叫的是 Greeter 自己的 Title。
func (g Greeter) Welcome() string { return "歡迎光臨，" + g.Title() }

// VIPGreeter 嵌入 Greeter 並定義自己的 Title。
type VIPGreeter struct {
	Greeter
}

// Title 是貴賓的稱呼。
func (VIPGreeter) Title() string { return "貴賓" }

// notInheritance：提升上來的 Welcome 收到的 receiver 是裡面那個 Greeter，
// 所以它呼叫的是 Greeter.Title，不是 VIPGreeter.Title。
func notInheritance() {
	v := VIPGreeter{}
	fmt.Println(v.Title())
	fmt.Println(v.Welcome())
}
