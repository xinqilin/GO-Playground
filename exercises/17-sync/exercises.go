//go:build !solution

package stall

import "fmt"

// Jar 是櫃台的小費箱，可以給多個 goroutine 同時使用，零值就能直接用。
type Jar struct {
	// TODO: 加上需要的欄位
}

// Put 投入 amount 元。
func (j *Jar) Put(amount int) {
	// TODO: 實作
}

// Count 回傳投過幾次、總共多少元。
// 兩個數字必須是同一個時間點的：不可以讀到「次數加了、金額還沒加」的狀態。
func (j *Jar) Count() (times, total int) {
	// TODO: 實作
	return 0, 0
}

// Shutter 是店門口的鐵捲門，零值代表開著。
type Shutter struct {
	// TODO: 加上需要的欄位
}

// Close 拉下鐵捲門。好幾個 goroutine 同時呼叫時，只有真正把門拉下來的那一個回傳 true，
// 其他人（包括門已經拉下來之後才呼叫的）都回傳 false。
func (s *Shutter) Close() bool {
	// TODO: 實作
	return false
}

// Closed 回報鐵捲門是不是已經拉下來了。
func (s *Shutter) Closed() bool {
	// TODO: 實作
	return false
}

// Reservations 是餐廳的訂位表，可以給多個 goroutine 同時使用。用 NewReservations 建立。
type Reservations struct {
	// TODO: 加上需要的欄位
}

func NewReservations() *Reservations {
	// TODO: 實作
	return &Reservations{}
}

// Book 用 name 訂 table 號桌。那一桌已經有人訂了，就回傳 ErrTaken，原本的訂位不變。
func (r *Reservations) Book(table int, name string) error {
	// TODO: 實作
	return nil
}

// Who 回傳 table 號桌是誰訂的；沒有人訂，回傳 "" 和 false。
func (r *Reservations) Who(table int) (string, bool) {
	// TODO: 實作
	return "", false
}

// Current 保存目前的 Settings，可以給多個 goroutine 同時使用，零值代表 Settings 的零值。
type Current struct {
	// TODO: 加上需要的欄位
}

// Get 回傳目前的設定。
func (c *Current) Get() Settings {
	// TODO: 實作
	return Settings{}
}

// Update 用 fn 從目前的設定算出新的設定，再換上去。
// 好幾個 goroutine 同時 Update 時，每一個人的修改都不可以遺失。
// fn 可能被呼叫不只一次，所以 fn 只能根據參數算出結果，不可以有其他副作用。
func (c *Current) Update(fn func(Settings) Settings) {
	// TODO: 實作
}

// OnceOK 回傳一個函式：f 成功之前，每次呼叫都會執行 f，並回傳 f 的結果；
// f 成功一次之後，就一直回傳那一次的值和 nil，不再執行 f。
// 回傳的函式可以被好幾個 goroutine 同時呼叫，同一時間最多只有一個 f 在執行。
func OnceOK[T any](f func() (T, error)) func() (T, error) {
	// TODO: 實作
	return func() (T, error) {
		var zero T
		return zero, nil
	}
}

// Summary 回傳像「3 組候位，下一組：王先生」這樣的說明；沒有人在等時回傳「目前不用等」。
//
// 找錯：呼叫之後永遠不會回傳。
func (q *Queue) Summary() string {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.Len() == 0 {
		return "目前不用等"
	}
	return fmt.Sprintf("%d 組候位，下一組：%s", q.Len(), q.names[0])
}
