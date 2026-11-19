package main

import (
	"fmt"
	"time"
	_ "time/tzdata" // 把時區資料庫編進執行檔，沒有安裝時區資料的系統也能 LoadLocation
)

var (
	taipei = mustLoad("Asia/Taipei")
	tokyo  = mustLoad("Asia/Tokyo")
)

// mustLoad 載入時區。名稱寫錯是程式的錯，啟動時就該發現，所以直接 panic。
func mustLoad(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}

// zonesDemo 同一個瞬間，在不同時區有不同的「牆上時間」。
func zonesDemo() {
	paid := time.Date(2026, time.November, 19, 10, 30, 0, 0, time.UTC) // 資料庫裡存 UTC
	const layout = time.DateTime + " MST"

	fmt.Println(paid.Format(layout))
	fmt.Println(paid.In(taipei).Format(layout))
	fmt.Println(paid.In(tokyo).Format(layout))
	fmt.Println("是同一個瞬間？", paid.Equal(paid.In(tokyo)))
}

// parseZoneDemo 字串裡沒有時區時，Parse 當成 UTC，ParseInLocation 當成你指定的時區。
func parseZoneDemo() {
	const layout = "2006-01-02 15:04"
	a, err := time.Parse(layout, "2026-11-19 18:30")
	if err != nil {
		fmt.Println(err)
		return
	}
	b, err := time.ParseInLocation(layout, "2026-11-19 18:30", taipei)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(a.Location(), a.Format(time.RFC3339))
	fmt.Println(b.Location(), b.Format(time.RFC3339))
	fmt.Println("兩個差了", a.Sub(b))
}
