// 第 10 章：interface。
//
// interface 是一張「會做什麼」的清單。任何型別只要有清單上的方法，就自動符合，不用宣告 implements。
// 執行：go run ./chapters/10-interfaces
package main

import (
	"fmt"
	"math"
)

// Shape 只描述行為：能算面積和周長。
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rect 是長方形。
type Rect struct{ W, H float64 }

// Circle 是圓形。
type Circle struct{ R float64 }

func (r Rect) Area() float64      { return r.W * r.H }
func (r Rect) Perimeter() float64 { return 2 * (r.W + r.H) }

func (c Circle) Area() float64      { return math.Pi * c.R * c.R }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.R }

// String 讓 Rect 符合 fmt.Stringer，fmt 印出 Rect 時會呼叫它。
func (r Rect) String() string { return fmt.Sprintf("Rect(%gx%g)", r.W, r.H) }

// 編譯期檢查：少實作一個方法，這兩行就編譯不過。
var (
	_ Shape = Rect{}
	_ Shape = Circle{}
)

func main() {
	fmt.Println("== 隱式實作 ==")
	shapes := []Shape{Rect{3, 4}, Circle{1}}
	for _, s := range shapes {
		fmt.Printf("%v 面積 %.2f 周長 %.2f\n", s, s.Area(), s.Perimeter())
	}
	fmt.Printf("總面積 %.2f\n", totalArea(shapes))

	fmt.Println("== type switch 與 type assertion ==")
	describe(Rect{1, 2})
	describe(Circle{2})
	describe(42)
	describe("hi")

	fmt.Println("== nil interface 陷阱 ==")
	nilTrap()
}

// totalArea 只依賴 Shape 介面，不在乎拿到的是哪一種形狀。
func totalArea(shapes []Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

// describe 用 type switch 依實際型別做不同的事，再用 type assertion 取出具體型別。
func describe(v any) {
	switch x := v.(type) {
	case Shape:
		fmt.Printf("形狀，面積 %.2f\n", x.Area())
	case int:
		fmt.Println("整數", x+1)
	default:
		fmt.Printf("其他型別 %T\n", x)
	}
	if r, ok := v.(Rect); ok { // 帶 ok 的寫法：型別不符時不會 panic
		fmt.Println("也是 Rect，寬", r.W)
	}
}

// MyErr 是自訂錯誤型別，用 pointer receiver 實作 error。
type MyErr struct{}

func (*MyErr) Error() string { return "出事了" }

// mayFail 回傳一個 nil 的 *MyErr。包進 error 介面後，介面記著「型別 *MyErr、值 nil」，所以介面不是 nil。
func mayFail() error {
	var p *MyErr
	return p
}

// nilTrap 示範 interface 陷阱：裝了 nil 指標的介面不等於 nil。
func nilTrap() {
	err := mayFail()
	fmt.Println(err == nil)
	fmt.Printf("%T %v\n", err, err == (*MyErr)(nil))
}
