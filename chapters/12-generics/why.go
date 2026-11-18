package main

import "fmt"

// sumInts 和 sumFloats：沒有泛型時，每一種型別都要寫一次幾乎一樣的程式。
func sumInts(s []int) int {
	total := 0
	for _, v := range s {
		total += v
	}
	return total
}

func sumFloats(s []float64) float64 {
	total := 0.0
	for _, v := range s {
		total += v
	}
	return total
}

// sumAny 只寫一次，但什麼都收：放錯型別要到執行時才知道，這裡只能默默略過。
func sumAny(s []any) float64 {
	total := 0.0
	for _, v := range s {
		switch x := v.(type) {
		case int:
			total += float64(x)
		case float64:
			total += x
		}
	}
	return total
}

// Sum 是泛型函式：N 是型別參數，只能是 int 或 float64。
func Sum[N int | float64](s []N) N {
	var total N
	for _, v := range s {
		total += v
	}
	return total
}

func whyDemo() {
	fmt.Println(sumInts([]int{85, 95}), sumFloats([]float64{1.5, 2.25}))
	fmt.Println(sumAny([]any{85, 2.5, "排骨飯"}))
	fmt.Println(Sum([]int{85, 95}), Sum([]float64{1.5, 2.25}))
}
