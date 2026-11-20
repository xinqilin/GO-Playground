//go:build !solution

package idiom

// Stutters 回傳匯出名稱 name 從呼叫端看起來是不是重複了 package 名稱 pkg，
// 也就是 name 以 pkg 開頭（不分大小寫），而且後面接著一個大寫字母，
// 例如 ticket.TicketMachine、http.HTTPServer。
// 和 package 同名（time.Time）或只是字首相同（ticket.Tickets）都不算。
func Stutters(pkg, name string) bool {
	// TODO: 實作
	return false
}

// DocProblems 檢查匯出名稱 name 的 doc comment（已去掉 // 和換行，接成一段文字），
// 依序回傳發現的問題：
//   - doc 去掉前後空白後是空的：只回傳 ProblemMissing。
//   - 不是以「name 加一個空白」開頭：ProblemName。
//   - 最後一個字不是「。」或「.」：ProblemPeriod。
//
// 沒有問題時回傳 nil。
func DocProblems(name, doc string) []string {
	// TODO: 實作
	return nil
}

// Add 記錄 dish 又賣了 n 份。零值的 Tally 也要能用。
func (t *Tally) Add(dish string, n int) {
	// TODO: 實作
}

// Count 回傳 dish 賣了幾份，沒賣過是 0。
func (t *Tally) Count(dish string) int {
	// TODO: 實作
	return 0
}

// Top 回傳賣最多的 n 道菜：份數多的在前，份數相同時依菜名排序。
// n 比菜的種類多時全部回傳。
func (t *Tally) Top(n int) []string {
	// TODO: 實作
	return nil
}

// WithPrefix 設定號碼的前綴，預設是 "A"。
func WithPrefix(p string) Option {
	// TODO: 實作
	return func(*Machine) {}
}

// WithStart 設定第一個號碼，預設是 1。
func WithStart(n int) Option {
	// TODO: 實作
	return func(*Machine) {}
}

// WithWidth 設定號碼補零後的位數，預設是 3。
func WithWidth(w int) Option {
	// TODO: 實作
	return func(*Machine) {}
}

// New 用預設值建立 Machine，再依序套用 opts。
func New(opts ...Option) *Machine {
	// TODO: 實作
	return &Machine{}
}

// Next 回傳下一張號碼牌，例如 "A001"。號碼超過位數時照實寫出，例如 "A1000"。
func (m *Machine) Next() string {
	// TODO: 實作
	return ""
}

// MixedCaps 把 snake_case 的名稱轉成 Go 的寫法。exported 是 true 時轉成匯出名稱：
// order_count 變成 OrderCount、user_id 變成 UserID、api_base_url 變成 APIBaseURL；
// false 時第一段整個小寫：orderCount、userID、apiBaseURL。
// Initialisms 裡的縮寫整個大寫（在開頭而且 exported 是 false 時整個小寫），
// 其他段落只把第一個字母大寫。連續或前後多餘的底線忽略。輸入只有小寫英文字母、數字和底線。
func MixedCaps(snake string, exported bool) string {
	// TODO: 實作
	return ""
}

// QueueOf 用 nums 建立一條隊伍。之後呼叫端修改 nums，不能影響這條隊伍。
//
// 找錯：呼叫端改了 nums，或改了 Waiting 回傳的 slice，隊伍裡的號碼也跟著變了。
func QueueOf(nums []int) *Queue {
	return &Queue{waiting: nums}
}

// Join 讓號碼 n 排到隊伍的最後面。
func (q *Queue) Join(n int) {
	q.waiting = append(q.waiting, n)
}

// Waiting 回傳目前排隊的號碼。呼叫端拿到的是複本，怎麼改都不影響隊伍。
func (q *Queue) Waiting() []int {
	return q.waiting
}
