// 第 02 章：輸入與輸出。
//
// 程式像一台果汁機：從 io.Reader 倒進原料，從 io.Writer 倒出結果。
// 執行：printf '3\n4\n5\n' | go run ./chapters/02-input-output
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	showVerbs(os.Stdout)

	// 提示訊息寫到 stderr，stdout 只留結果，方便接到其他程式。
	fmt.Fprintln(os.Stderr, "請輸入整數，每行一個，Ctrl-D 結束：")
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "錯誤：", err)
		os.Exit(1)
	}
}

type point struct{ X, Y int }

// showVerbs 示範 Printf 常用的格式化動詞。
func showVerbs(w io.Writer) {
	p := point{1, 2}
	fmt.Fprintf(w, "%v | %+v | %#v | %T\n", p, p, p, p)
	fmt.Fprintf(w, "%d | %5d | %-5d | %05d\n", 42, 42, 42, 42)
	fmt.Fprintf(w, "%f | %.2f | %e\n", 3.14159, 3.14159, 314.159)
	fmt.Fprintf(w, "%s | %q | %x\n", "Go", "Go", "Go")
	fmt.Fprintf(w, "%t | %c | %U\n", true, '世', '世')
}

// run 從 in 一行一行讀整數，把個數、總和與平均寫到 out。
// 參數用 io.Reader / io.Writer 而不是直接用 os.Stdin / os.Stdout，測試時就能換成字串。
func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sum, count, lineNo := 0, 0, 0
	for sc.Scan() { // 每次讀一行，讀完或出錯時回傳 false
		lineNo++
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		n, err := strconv.Atoi(line)
		if err != nil {
			return fmt.Errorf("第 %d 行 %q 不是整數", lineNo, line)
		}
		sum += n
		count++
	}
	if err := sc.Err(); err != nil { // Scan 回傳 false 之後，要確認是讀完還是出錯
		return err
	}
	if count == 0 {
		fmt.Fprintln(out, "沒有輸入任何數字")
		return nil
	}
	fmt.Fprintf(out, "共 %d 個數字，總和 %d，平均 %.2f\n", count, sum, float64(sum)/float64(count))
	return nil
}
