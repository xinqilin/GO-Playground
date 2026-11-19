package main

import (
	"fmt"
	"time"
)

// etaDemo 算外送的預計送達時間：下單的時間點，加上做菜和送餐要花的時間長度。
func etaDemo() {
	ordered := time.Date(2026, time.November, 19, 11, 52, 0, 0, time.UTC)
	cook := 15 * time.Minute
	ride := 25*time.Minute + 30*time.Second

	eta := ordered.Add(cook + ride)
	fmt.Println("總共要花", cook+ride)
	fmt.Println("預計送達", eta.Format(time.TimeOnly))

	promised := time.Date(2026, time.November, 19, 12, 30, 0, 0, time.UTC)
	fmt.Println("比答應的晚了", eta.Sub(promised))
}

// durationDemo 看 Duration 的幾種寫法和換算。
func durationDemo() {
	d := 90 * time.Minute
	fmt.Println(d, d.Hours(), int64(d))

	wait, err := time.ParseDuration("1h15m30s")
	fmt.Println(wait, wait.Minutes(), err)

	people := 7
	perPerson := 3 * time.Minute
	fmt.Println("排隊要等", time.Duration(people)*perPerson)

	fmt.Println((95 * time.Second).Round(time.Minute), (95 * time.Second).Truncate(time.Minute))
}
