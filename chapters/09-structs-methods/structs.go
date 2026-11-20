package main

import "fmt"

// Bento 是菜單上的一個便當。
type Bento struct {
	Name  string
	Price int
	Spicy bool
}

// structLiterals 示範建立 struct 的幾種寫法，以及三種印法。
func structLiterals() {
	var zero Bento
	a := Bento{Name: "排骨飯", Price: 85}
	b := Bento{"雞腿飯", 95, true}
	fmt.Printf("%v\n", zero)
	fmt.Printf("%+v\n", a)
	fmt.Printf("%#v\n", b)

	a.Price = 90
	fmt.Println(a.Name, a.Price, a.Spicy)
}

// structPointer 透過指標存取欄位時，Go 會自動解參考。
func structPointer() {
	p := &Bento{Name: "素食餐", Price: 80}
	p.Price += 5
	fmt.Println(p.Price)
	fmt.Printf("%+v\n", p)
}

// anonymousStruct 只用一次的資料，不必替它取型別名稱。
func anonymousStruct() {
	stats := []struct {
		Day   string
		Sales int
	}{
		{"一", 3200},
		{"二", 4100},
	}
	for _, s := range stats {
		fmt.Println(s.Day, s.Sales)
	}
}

// copyStruct 指定 struct 會把每個欄位複製一份。
func copyStruct() {
	a := Bento{Name: "排骨飯", Price: 85}
	b := a
	b.Price = 100
	fmt.Println(a.Price, b.Price)
}

// Combo 是套餐，Sides 是配菜清單。
type Combo struct {
	Name  string
	Sides []string
}

// shallowCopy 複製 struct 時，slice 欄位複製的只是 slice header，底層陣列還是共用的。
func shallowCopy() {
	a := Combo{Name: "A 餐", Sides: []string{"滷蛋", "青菜"}}
	b := a
	b.Name = "B 餐"
	b.Sides[0] = "豆干"
	fmt.Println(a.Name, a.Sides)
	fmt.Println(b.Name, b.Sides)
}

// compareStructs 每個欄位都能比較時，struct 可以用 == 比較，也可以當 map 的 key。
func compareStructs() {
	a := Bento{Name: "排骨飯", Price: 85}
	b := Bento{Name: "排骨飯", Price: 85}
	c := Bento{Name: "排骨飯", Price: 90}
	fmt.Println(a == b, a == c)

	sold := map[Bento]int{}
	sold[a]++
	sold[b]++
	sold[c]++
	fmt.Println(len(sold), sold[a])
}
