package main

import (
	"fmt"
	"strings"
	"time"
)

// equalDemo 同一個瞬間換了時區，== 就不相等了；比較時間要用 Equal、Before、After、Compare。
func equalDemo() {
	a := time.Date(2026, time.November, 19, 10, 30, 0, 0, time.UTC)
	b := a.In(taipei)
	fmt.Println(a == b, a.Equal(b), a.Compare(b))

	c := a.Add(time.Minute)
	fmt.Println(a.Before(c), c.After(a), a.Compare(c))
}

// monotonicDemo time.Now() 除了牆上時間，還帶著一個單調時鐘的讀數，專門用來量經過多久。
// 印出來的字串裡 m= 後面那一段就是它。
func monotonicDemo() {
	start := time.Now()
	fmt.Println("Now 帶著單調時鐘：", strings.Contains(start.String(), " m="))
	fmt.Println("Round(0) 之後：", strings.Contains(start.Round(0).String(), " m="))
	fmt.Println("經過的時間不會是負的：", time.Since(start) >= 0)
}

// unixDemo Unix 時間是從 1970-01-01 00:00:00 UTC 起算的秒數，和時區無關。
func unixDemo() {
	t := time.Date(2026, time.November, 19, 10, 30, 0, 0, time.UTC)
	fmt.Println(t.Unix(), t.UnixMilli())

	back := time.Unix(t.Unix(), 0) // 回來的時間用的是本機時區
	fmt.Println(back.UTC().Format(time.RFC3339), back.Equal(t))
}
