package main

import "fmt"

// stringBasics 示範字串的長度、索引、切片、串接和比較。
func stringBasics() {
	s := "Hello, Go"
	fmt.Println(len(s), s[0], string(s[0])) // 索引拿到的是 byte
	fmt.Println(s[7:], s[:5])               // 切片的結果還是字串
	greeting := s[:5] + "，便當店"
	fmt.Println(greeting)
	fmt.Println("排骨" == "排骨", "apple" < "banana", "Zoo" < "apple")
}

// literals 示範兩種字串字面值：雙引號會處理跳脫字元，反引號是原始字串，內容照原樣保留，也可以跨行。
func literals() {
	fmt.Println("第一行\n第二行")
	fmt.Println("價格\t85")
	fmt.Println(`C:\menu\new.txt`)
	query := `SELECT name, price
FROM menu
WHERE price < 100`
	fmt.Println(query)
	fmt.Println("\u8a9e\u8a00")
}
