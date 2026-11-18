package main

import "fmt"

// sum 是累加器：把清單裡的數全部加起來。
func sum(nums []int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// countAtLeast 是計數器：數一數有幾個數大於或等於 limit。
func countAtLeast(nums []int, limit int) int {
	count := 0
	for _, n := range nums {
		if n >= limit {
			count++
		}
	}
	return count
}

// maxOf 記住「目前為止最大的」，遇到更大的就換掉。nums 至少要有一個數。
func maxOf(nums []int) int {
	best := nums[0]
	for _, n := range nums {
		if n > best {
			best = n
		}
	}
	return best
}

// wrongMax 從 0 開始比，清單裡全是負數時，答案就錯了。
func wrongMax(nums []int) int {
	best := 0
	for _, n := range nums {
		if n > best {
			best = n
		}
	}
	return best
}

// indexOf 搜尋 target 第一次出現的位置，找到就立刻 return；整份清單都找完還沒有，回傳 -1。
func indexOf(names []string, target string) int {
	for i, name := range names {
		if name == target {
			return i
		}
	}
	return -1
}

func sumDemo() {
	sales := []int{3200, 4100, 2800, 5600, 4100}
	fmt.Println(sum(sales))
}

func countDemo() {
	sales := []int{3200, 4100, 2800, 5600, 4100}
	fmt.Println(countAtLeast(sales, 4000))
}

func maxDemo() {
	sales := []int{3200, 4100, 2800, 5600, 4100}
	fmt.Println(maxOf(sales))
}

func wrongMaxDemo() {
	temps := []int{-3, -5, -1}
	fmt.Println(wrongMax(temps), maxOf(temps))
}

func searchDemo() {
	queue := []string{"小明", "阿美", "小華"}
	fmt.Println(indexOf(queue, "阿美"))
	fmt.Println(indexOf(queue, "小芬"))
}
