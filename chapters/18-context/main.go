// 第 18 章：context。
//
// context 是一條往下傳的取消繩：上層一拉（取消或逾時），所有拿著這條繩子的工作都知道該停了。
// 執行：go run ./chapters/18-context
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// slowWork 模擬一個要花 d 的工作；ctx 被取消就提早放棄。
// 慣例：會阻塞或呼叫外部服務的函式，第一個參數放 ctx context.Context。
func slowWork(ctx context.Context, d time.Duration) error {
	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// requestIDKey 是 context 值的 key。用不匯出的型別，就不會和別的套件撞名。
type requestIDKey struct{}

// handle 從 context 取出 request ID。
func handle(ctx context.Context) {
	id, _ := ctx.Value(requestIDKey{}).(string)
	fmt.Println("處理請求", id)
}

func main() {
	fmt.Println("== 逾時 ==")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel() // 一定要呼叫 cancel，釋放計時器等資源
	err := slowWork(ctx, time.Hour)
	fmt.Println(err, errors.Is(err, context.DeadlineExceeded))

	fmt.Println("== 手動取消 ==")
	ctx2, cancel2 := context.WithCancel(context.Background())
	cancel2()
	fmt.Println(slowWork(ctx2, time.Hour))

	fmt.Println("== 取消會往下傳 ==")
	parent, cancelParent := context.WithCancel(context.Background())
	child, cancelChild := context.WithTimeout(parent, time.Hour)
	defer cancelChild()
	cancelParent() // 取消上層，下層跟著被取消
	<-child.Done()
	fmt.Println("child:", child.Err())

	fmt.Println("== 取消原因 ==")
	ctx3, cancel3 := context.WithCancelCause(context.Background())
	cancel3(errors.New("使用者按了取消"))
	fmt.Println(ctx3.Err(), "/", context.Cause(ctx3))

	fmt.Println("== 帶值 ==")
	ctx4 := context.WithValue(context.Background(), requestIDKey{}, "req-42")
	handle(ctx4)
}
