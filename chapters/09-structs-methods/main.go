// 第 09 章：struct 與方法。
//
// struct 是一張表格，把相關的欄位放在一起；方法是綁在某個型別上的函式。
// 執行：go run ./chapters/09-structs-methods
package main

import "fmt"

// Point 是平面上的一個點。
type Point struct {
	X, Y int
}

// Dist2 是 value receiver：拿到的是複本，適合不需要修改內容的方法。
func (p Point) Dist2() int { return p.X*p.X + p.Y*p.Y }

// Move 是 pointer receiver：要修改 struct 本身就用指標。
func (p *Point) Move(dx, dy int) {
	p.X += dx
	p.Y += dy
}

// Animal 會被嵌入到 Dog 裡。
type Animal struct{ Name string }

// Hello 是 Animal 的方法，嵌入後 Dog 也能直接呼叫。
func (a Animal) Hello() string { return "我是" + a.Name }

// Dog 嵌入 Animal：Animal 的欄位和方法會「提升」到 Dog 身上。這是組合，不是繼承。
type Dog struct {
	Animal
	Breed string
}

func main() {
	fmt.Println("== 建立 struct ==")
	literals()

	fmt.Println("== value 與 pointer receiver ==")
	receivers()

	fmt.Println("== 嵌入 ==")
	embedding()

	fmt.Println("== 比較與匿名 struct ==")
	compare()
}

// literals 示範建立 struct 的幾種寫法。
func literals() {
	var zero Point         // 零值：每個欄位都是零值
	p := Point{X: 1, Y: 2} // 建議寫出欄位名
	q := Point{3, 4}       // 依順序填，欄位一多就難讀
	r := &Point{X: 5}      // 取地址得到 *Point；沒寫的欄位是零值
	fmt.Println(zero, p, q, *r)
	fmt.Printf("%+v\n", p)
}

// receivers 示範兩種 receiver。
func receivers() {
	p := Point{X: 3, Y: 4}
	fmt.Println(p.Dist2())
	p.Move(1, 1) // p 是可以取地址的變數，Go 自動改寫成 (&p).Move(1, 1)
	fmt.Println(p)
}

// embedding 示範嵌入後，欄位和方法都可以直接用。
func embedding() {
	d := Dog{Animal: Animal{Name: "小黑"}, Breed: "柴犬"}
	fmt.Println(d.Name, d.Breed) // Name 被提升，可以直接用
	fmt.Println(d.Hello())       // 方法也被提升
	fmt.Println(d.Animal.Name)   // 也可以寫完整路徑

	// Go 1.27 起，struct 字面值的 key 也可以直接寫被提升的欄位。
	w := Dog{Name: "小白", Breed: "米克斯"}
	fmt.Println(w.Hello())
}

// compare 示範 struct 的比較，以及只用一次的匿名 struct。
func compare() {
	a := Point{1, 2}
	b := Point{1, 2}
	fmt.Println(a == b) // 欄位都能比較時，struct 可以直接用 ==

	cfg := struct {
		Host string
		Port int
	}{"localhost", 8080}
	fmt.Println(cfg.Host, cfg.Port)
}
