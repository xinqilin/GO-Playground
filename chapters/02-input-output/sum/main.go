// sum 從標準輸入一行一行讀整數，印出個數、總和與平均。
//
// 提示和錯誤寫到 stderr，結果寫到 stdout，所以可以把結果單獨導到檔案或下一個程式。
//
// 執行：printf '3\n4\n5\n' | go run ./chapters/02-input-output/sum
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
	fmt.Fprintln(os.Stderr, "請輸入整數，每行一個，Ctrl-D 結束：")
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "錯誤：", err)
		os.Exit(1)
	}
}

// run 從 in 一行一行讀整數，把個數、總和與平均寫到 out。
// 參數用 io.Reader 和 io.Writer，而不是直接用 os.Stdin 和 os.Stdout，測試時就能換成字串。
func run(in io.Reader, out io.Writer) error {
	sc := bufio.NewScanner(in)
	sum, count, lineNo := 0, 0, 0
	for sc.Scan() {
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
	if err := sc.Err(); err != nil {
		return err
	}
	if count == 0 {
		_, err := fmt.Fprintln(out, "沒有輸入任何數字")
		return err
	}
	_, err := fmt.Fprintf(out, "共 %d 個數字，總和 %d，平均 %.2f\n", count, sum, float64(sum)/float64(count))
	return err
}
