package main

import (
	"fmt"
	"math"
)

// floatError 示範浮點數無法精確表示大部分的小數。
func floatError() {
	a, b := 0.1, 0.2
	fmt.Println(a+b, a+b == 0.3)

	sum := 0.0
	for range 10 {
		sum += 0.1
	}
	fmt.Println(sum)
}

// floatCompare 示範比較浮點數時要允許一點誤差。
func floatCompare() {
	a, b := 0.1, 0.2
	const epsilon = 1e-9
	fmt.Println(math.Abs(a+b-0.3) < epsilon)
}

// moneyInTenths 示範金額用整數存：電費每度 2.1 元，先換成 21 角再計算。
func moneyInTenths() {
	rate := 2.1
	fmt.Println(3 * rate)

	rateTenths := 21
	total := 3 * rateTenths
	fmt.Printf("%d.%d 元\n", total/10, total%10)
}

// floatSpecial 示範浮點數除以 0 不會 panic，而是得到無限大或 NaN。
func floatSpecial() {
	var zero float64
	fmt.Println(1/zero, -1/zero)
	nan := zero / zero
	fmt.Println(nan, nan == nan, math.IsNaN(nan))
}
