package main

import "fmt"

type item struct {
	Name  string
	Price int
}

// verbsGeneral 示範萬用的 %v 和它的兩種變化，以及印出型別的 %T。
func verbsGeneral() {
	it := item{"排骨飯", 85}
	fmt.Printf("%v | %+v\n", it, it)
	fmt.Printf("%#v\n", it)
	fmt.Printf("%T\n", it)

	menu := map[string]int{"雞腿": 90, "排骨": 85, "滷蛋": 15}
	fmt.Printf("%v\n", menu)
	queue := []string{"A01", "A02"}
	fmt.Printf("%v %q\n", queue, queue)
}

// verbsInt 示範整數和字元的動詞。
func verbsInt() {
	n := 255
	fmt.Printf("%d %+d %b %o %x %X %#x\n", n, n, n, n, n, n, n)
	r := '飯'
	fmt.Printf("%c %q %U %d\n", r, r, r, r)
}

// verbsFloat 示範浮點數的動詞。
func verbsFloat() {
	f := 1234.5678
	fmt.Printf("%f | %.2f | %.0f\n", f, f, f)
	fmt.Printf("%e | %g\n", f, f)
	fmt.Printf("%g | %g\n", 0.000012345, 1e21)
}

// verbsString 示範字串、布林和百分比符號。
func verbsString() {
	s := "Go語"
	fmt.Printf("%s | %q | %x | % x\n", s, s, s, s)
	fmt.Printf("%t | %v\n", true, false)
	fmt.Printf("%d%%\n", 85)
}
