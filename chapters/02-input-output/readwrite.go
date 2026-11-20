package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// numberLines 從 r 一行一行讀，在每一行前面加上行號，寫到 w。
// 參數用 io.Reader 和 io.Writer，就不用管資料從哪裡來、要寫到哪裡去。
func numberLines(r io.Reader, w io.Writer) error {
	sc := bufio.NewScanner(r)
	n := 0
	for sc.Scan() {
		n++
		if _, err := fmt.Fprintf(w, "%02d %s\n", n, sc.Text()); err != nil {
			return err
		}
	}
	return sc.Err()
}

// numberLinesDemo 讓同一個函式一次寫到 stdout，一次寫到 strings.Builder。
func numberLinesDemo() {
	if err := numberLines(strings.NewReader("排骨飯\n雞腿飯\n滷蛋\n"), os.Stdout); err != nil {
		fmt.Println("錯誤：", err)
	}

	var sb strings.Builder
	if err := numberLines(strings.NewReader("紅茶\n綠茶\n"), &sb); err != nil {
		fmt.Println("錯誤：", err)
	}
	fmt.Printf("%q\n", sb.String())
}

// bufferedWriter 示範 bufio.Writer 先把資料存在記憶體，Flush 時才真的寫出去。
func bufferedWriter() {
	var sb strings.Builder
	bw := bufio.NewWriter(&sb)
	fmt.Fprint(bw, "叫號：A01")
	fmt.Printf("Flush 前：%q\n", sb.String())
	if err := bw.Flush(); err != nil {
		fmt.Println("錯誤：", err)
	}
	fmt.Printf("Flush 後：%q\n", sb.String())
}

// copyDemo 示範 io.Copy 把 Reader 的內容全部搬到 Writer，回傳搬了多少 byte。
func copyDemo() {
	n, err := io.Copy(os.Stdout, strings.NewReader("整份菜單原封不動地複製過去\n"))
	fmt.Println(n, err)
}
