package main

import "fmt"

// Convert 把一種數字的 slice 換成另一種。From 推論得出來，To 推論不出來。
func Convert[To, From Number](s []From) []To {
	out := make([]To, len(s))
	for i, v := range s {
		out[i] = To(v)
	}
	return out
}

func inferDemo() {
	prices := []NTD{85, 95}

	fs := Convert[float64](prices)
	fmt.Printf("%T %v\n", fs, fs)

	all := Convert[float64, NTD](prices)
	fmt.Println(all)

	var sumPrices func([]NTD) NTD = Total
	fmt.Println(sumPrices(prices))
}
