// 專題 1：費氏數列。
//
// 改寫自這個 repo 最早的 Exercise/fibonacci.go。舊版有三個問題：
//  1. 隨機抽到 n = 0 時，fibonacci(0) 會呼叫 fibonacci(-1)、fibonacci(-2)……永遠停不下來，最後 stack overflow。
//  2. fibonacci(2) 回傳 2，整個數列錯開一格（正確是 0, 1, 1, 2, 3, 5…）。
//  3. rand.Seed 自 Go 1.20 起已經 deprecated，現在不用自己設種子。
//
// 新版用五種方法算同一個數列：遞迴、記憶化、迴圈（加上溢位檢查）、疊代器、math/big。
//
// 執行：go run ./projects/fibonacci -n 30
//
//	go run ./projects/fibonacci -list 15
//	go run ./projects/fibonacci -n 100
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// slowN 以上就不跑直接遞迴：呼叫次數會超過一千萬次。
const slowN = 35

// run 解析參數並輸出結果；和 main 分開，測試時就能傳入自己的參數和輸出。
func run(args []string, out io.Writer) error {
	fs := flag.NewFlagSet("fibonacci", flag.ContinueOnError)
	fs.SetOutput(out)
	n := fs.Int("n", 10, "要算第幾項（0 以上）")
	list := fs.Int("list", 0, fmt.Sprintf("改成列出前幾項（1 到 %d）", maxN+1))
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil // -h：用法已經印出來了，不是錯誤
		}
		return err
	}
	if *list != 0 {
		return printList(out, *list)
	}
	if *n < 0 {
		return fmt.Errorf("n 不能是負數，收到 %d", *n)
	}

	v, err := Checked(*n)
	if err != nil {
		fmt.Fprintf(out, "F(%d) = %s\n", *n, Big(*n))
		fmt.Fprintf(out, "超過 int 的範圍（最大到 F(%d)），只用 math/big 計算\n", maxN)
		return nil
	}
	fmt.Fprintf(out, "F(%d) = %d\n", *n, v)
	if *n <= slowN {
		fmt.Fprintf(out, "遞迴：%d（呼叫了 %d 次）\n", Recursive(*n), countCalls(*n))
	} else {
		fmt.Fprintf(out, "遞迴：略過（n 大於 %d 會跑很久）\n", slowN)
	}
	fmt.Fprintf(out, "記憶化：%d\n", Memoized()(*n))
	fmt.Fprintf(out, "迴圈：%d\n", Iterative(*n))
	fmt.Fprintf(out, "math/big：%s\n", Big(*n))
	return nil
}

// printList 用 Seq 列出前 count 項，以空白分隔。
func printList(out io.Writer, count int) error {
	if count < 1 || count > maxN+1 {
		return fmt.Errorf("list 必須介於 1 到 %d，收到 %d", maxN+1, count)
	}
	var terms []string
	for v := range Seq() {
		terms = append(terms, strconv.Itoa(v))
		if len(terms) == count {
			break
		}
	}
	fmt.Fprintln(out, strings.Join(terms, " "))
	return nil
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "錯誤：", err)
		os.Exit(2)
	}
}
