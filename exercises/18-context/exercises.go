//go:build !solution

package courier

import (
	"context"
	"time"
)

// Sleep 等 d 這麼久之後回傳 nil；ctx 先被取消的話，立刻回傳 ctx.Err()。
func Sleep(ctx context.Context, d time.Duration) error {
	// TODO: 實作
	return nil
}

// ShouldStart 回報還來不來得及開始一件要花 need 的工作：
// ctx 沒有期限，或離期限還有至少 need，回傳 true；
// ctx 已經被取消，或剩下的時間不到 need，回傳 false。
func ShouldStart(ctx context.Context, need time.Duration) bool {
	// TODO: 實作
	return false
}

// WithStop 回傳一個新的 ctx，記錄外送路線在 ctx 原本的路線後面多經過 stop 這一站。
func WithStop(ctx context.Context, stop string) context.Context {
	// TODO: 實作
	return ctx
}

// Route 依序回傳 ctx 記錄的每一站；沒有任何一站時回傳 nil。
// 從同一個 ctx 分出去的路線互不影響；呼叫的人修改回傳的 slice，也不會改到 ctx 裡的紀錄。
func Route(ctx context.Context) []string {
	// TODO: 實作
	return nil
}

// Pump 把從 in 收到的值依序送到 out，直到 in 被 close，回傳 nil；
// 或是 ctx 被取消，回傳 ctx.Err()。不管正在等收還是等送，ctx 一取消就要回來。
func Pump[T any](ctx context.Context, in <-chan T, out chan<- T) error {
	// TODO: 實作
	return nil
}

// Fastest 同時向每一間 kitchens 下單，回傳最先成功的那一份。
// 有人成功之後，取消其他還在做的，並等它們都結束才回傳。
// 全部都失敗時，回傳 "" 和所有錯誤用 errors.Join 合在一起的結果。kitchens 至少有一個。
func Fastest(ctx context.Context, kitchens ...func(context.Context) (string, error)) (string, error) {
	// TODO: 實作
	return "", nil
}

// Quote 向廚房詢問一道菜要等多久。ask 會聽傳給它的 ctx；Quote 最多只等 1 秒。
//
// 找錯：呼叫的人已經取消了，它還是會等滿 1 秒。
func Quote(ctx context.Context, ask func(context.Context) (time.Duration, error)) (time.Duration, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return ask(ctx)
}
