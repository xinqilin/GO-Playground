package main

import (
	"fmt"
	"strconv"
)

// convertNumbers 示範數字型別之間的轉換，寫法是「型別(值)」。
func convertNumbers() {
	bill, people := 200, 3
	fmt.Println(bill/people, float64(bill)/float64(people))

	price := 85.9
	fmt.Println(int(price), int(-price)) // 直接截掉小數，往 0 的方向

	big := 300
	fmt.Println(int8(big)) // 放不下時只留下最低的 8 個位元
}

// convertStrings 示範數字和字串的轉換要用 strconv。
func convertStrings() {
	fmt.Println(strconv.Itoa(85) + " 元")

	n, err := strconv.Atoi("120")
	fmt.Println(n+1, err)

	_, err = strconv.Atoi("12a")
	fmt.Println(err)

	fmt.Println(string(rune(65)), strconv.Itoa(65))
}
