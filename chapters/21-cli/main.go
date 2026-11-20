// 第 21 章：命令列工具。
//
// 這一章的程式本身就是一個完整的命令列工具 tally，加總便當店的營業額：
//
//	printf '排骨飯 90\n雞腿飯 100\n' | go run ./chapters/21-cli total -tax 5
//
// main 只負責把作業系統的東西（參數、標準輸入輸出、環境變數、訊號）接上 run；
// run 不碰 os.Args、os.Exit，所以 main_test.go 可以用 Example 呼叫它，
// 鎖住每一種用法的輸出和結束碼，頁面上的輸出就是從那裡同步過去的。
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ctx.Done()
		stop() // 第一次 Ctrl-C 之後恢復預設行為：再按一次就直接結束
	}()

	code := run(ctx, os.Args[1:], sys{
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
		getenv: os.Getenv,
	})
	stop()
	os.Exit(code) // os.Exit 不會執行 defer，所以 main 裡不放任何 defer
}
