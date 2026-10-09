// 第 22 章：標準庫精選。
//
// Go 的標準庫是一個工具箱：時間、JSON、檔案、log、UUID 都有，大部分工作不用裝第三方套件。
// 執行：go run ./chapters/22-stdlib
package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
	"uuid"
)

func main() {
	fmt.Println("== time：用 2006-01-02 當格式範本 ==")
	timeDemo()

	fmt.Println("== encoding/json ==")
	jsonDemo()

	fmt.Println("== os：讀寫檔案 ==")
	fileDemo()

	fmt.Println("== log/slog：結構化 log ==")
	slogDemo()

	fmt.Println("== uuid（Go 1.27）==")
	uuidDemo()
}

// timeDemo 示範時間的格式化與計算。格式不是 yyyy-MM-dd，而是用一個固定的參考時間：
// 2006-01-02 15:04:05（1 月 2 日 下午 3 點 4 分 5 秒 2006 年）。
func timeDemo() {
	t := time.Date(2026, time.October, 9, 14, 30, 0, 0, time.UTC)
	fmt.Println(t.Format("2006-01-02 15:04:05"))
	fmt.Println(t.Format(time.RFC3339))

	end := t.Add(90 * time.Minute)
	fmt.Println(end.Format("15:04"), end.Sub(t))

	d, err := time.ParseDuration("1h15m")
	fmt.Println(d.Minutes(), err)

	xmas, err := time.Parse("2006/01/02", "2026/12/25")
	fmt.Println(xmas.Weekday(), err)
}

// User 的 struct tag 決定 JSON 欄位名稱；omitempty 讓空值不輸出。
type User struct {
	Name  string   `json:"name"`
	Age   *int     `json:"age,omitempty"` // 用指標：nil 代表不知道，和 0 歲不一樣
	Tags  []string `json:"tags,omitempty"`
	email string   // 小寫開頭沒有匯出，JSON 看不到
}

// jsonDemo 示範 struct 和 JSON 互轉。
func jsonDemo() {
	u := User{Name: "小明", Age: new(30), email: "ming@example.com"} // new(30)：Go 1.26 起
	b, err := json.Marshal(u)
	fmt.Println(string(b), err)

	var got User
	err = json.Unmarshal([]byte(`{"name":"小華","tags":["go","java"],"extra":1}`), &got)
	fmt.Println(got.Name, got.Age == nil, got.Tags, err) // 不認得的欄位 extra 會被忽略
}

// fileDemo 在暫存目錄寫一個檔案再讀回來。
func fileDemo() {
	dir, err := os.MkdirTemp("", "ch22-*")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.RemoveAll(dir) // 示範完就清掉暫存目錄

	path := filepath.Join(dir, "note.txt")
	if err := os.WriteFile(path, []byte("第一行\n第二行\n"), 0o644); err != nil {
		fmt.Println(err)
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(string(data))
	fmt.Println(len(strings.Split(strings.TrimSpace(string(data)), "\n")), "行")
}

// slogDemo 示範 log/slog 的結構化 log：訊息之外，再附上一組組 key=value。
func slogDemo() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				return slog.Attr{} // 範例拿掉時間，輸出才固定
			}
			return a
		},
	}))
	logger.Info("使用者登入", "user", "bill", "attempt", 2)
	logger.Warn("額度快用完", slog.Int("remaining", 3))
	logger.Debug("預設等級是 Info，這行不會印出來")
}

// uuidDemo 示範 Go 1.27 加入標準庫的 uuid 套件。
func uuidDemo() {
	id := uuid.NewV7() // v7 開頭是時間戳記，依產生順序遞增，適合當資料庫主鍵
	fmt.Println(len(id.String()))

	a, b := uuid.NewV7(), uuid.NewV7()
	fmt.Println(a.Compare(b) < 0)

	fixed := uuid.MustParse("017F22E2-79B0-7CC3-98C4-DC0C0C07398F")
	fmt.Println(fixed)
}
