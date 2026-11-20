package main

import "fmt"

// indexOfMax 回傳最大值的索引；有好幾個一樣大時，回傳最前面的那一個。nums 至少要有一個數。
func indexOfMax(nums []int) int {
	best := 0
	for i := range nums {
		if nums[i] > nums[best] {
			best = i
		}
	}
	return best
}

// weeklyReport 把「一週營業報表」拆成幾個小問題，每個小問題交給一個函式。
func weeklyReport() {
	sales := []int{3200, 4100, 2800, 5600, 4100, 6300, 5000}
	days := []string{"一", "二", "三", "四", "五", "六", "日"}

	total := sum(sales)
	fmt.Println("總營業額：", total)
	fmt.Println("平均每天：", total/len(sales))

	d := indexOfMax(sales)
	fmt.Println("最好的一天：星期"+days[d], sales[d])

	fmt.Println("達到 4500 的天數：", countAtLeast(sales, 4500))
}

// buggyTotal 想把清單加起來，結果少了。迴圈裡加一行 Println，看看每一圈發生了什麼事。
func buggyTotal(nums []int) int {
	total := 0
	for i := 1; i < len(nums); i++ {
		total += nums[i]
		fmt.Println("i =", i, "total =", total)
	}
	return total
}

func debugDemo() {
	fmt.Println(buggyTotal([]int{85, 90, 75}))
}
