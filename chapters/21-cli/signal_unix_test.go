//go:build unix

package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
)

// TestInterruptSignal 對自己送一個真的 SIGINT（就是按 Ctrl-C 時作業系統送的訊號）：
// NotifyContext 把它變成 ctx 的取消，run 讀下一行時發現了，回傳 130。
// Windows 沒有這種訊號，所以這個檔案只在 Unix 系統上編譯。
func TestInterruptSignal(t *testing.T) {
	ctx, stop := signal.NotifyContext(t.Context(), os.Interrupt)
	defer stop()

	stdin, w := io.Pipe()
	var stderr strings.Builder
	done := make(chan int)
	go func() {
		done <- run(ctx, []string{"total"}, sys{
			stdin: stdin, stdout: io.Discard, stderr: &stderr,
			getenv: func(string) string { return "" },
		})
	}()

	fmt.Fprintln(w, "排骨飯 90") // 等 run 讀走第一行才會回來
	if err := syscall.Kill(syscall.Getpid(), syscall.SIGINT); err != nil {
		t.Fatal(err)
	}
	<-ctx.Done()
	fmt.Fprintln(w, "雞腿飯 100") // run 讀到這一行之前先看 ctx，就停下來了
	w.Close()

	if code := <-done; code != exitInterrupted || !strings.Contains(stderr.String(), "interrupt signal received") {
		t.Errorf("結束碼 %d，stderr %q；want 130 而且說明收到 interrupt", code, stderr.String())
	}
}
