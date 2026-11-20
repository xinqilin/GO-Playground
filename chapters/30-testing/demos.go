package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/xinqilin/GO-Playground/chapters/30-testing/parking"
)

var taipei = time.FixedZone("UTC+8", 8*60*60)

// feeDemo 印出幾種停法的停車費。
func feeDemo() {
	at := func(day, hour, min int) time.Time {
		return time.Date(2026, 11, day, hour, min, 0, 0, taipei)
	}
	stays := []struct {
		name    string
		in, out time.Time
	}{
		{"停 8 分鐘", at(20, 9, 0), at(20, 9, 8)},
		{"停 95 分鐘", at(20, 8, 5), at(20, 9, 40)},
		{"23:50 停到 00:20", at(20, 23, 50), at(21, 0, 20)},
		{"停一整天", at(20, 0, 0), at(21, 0, 0)},
		{"出場時間打錯", at(20, 9, 0), at(20, 8, 0)},
	}
	for _, s := range stays {
		fee, err := parking.Fee(s.in, s.out)
		if err != nil {
			fmt.Printf("%s：%v\n", s.name, err)
			continue
		}
		fmt.Printf("%s：%d 元\n", s.name, fee)
	}
}

// printPayments 把收費印出來，代替真正的金流。
type printPayments struct{}

func (printPayments) Charge(ctx context.Context, plate string, amount int) error {
	fmt.Printf("刷卡 %s %d 元\n", plate, amount)
	return nil
}

// gateDemo：時間由 now 這個變數決定，想讓時間往前走就直接改它。
func gateDemo() {
	now := time.Date(2026, 11, 20, 8, 5, 0, 0, taipei)
	g := parking.NewGate(func() time.Time { return now }, printPayments{})
	exit := func(plate string) {
		rec, err := g.Exit(context.Background(), plate)
		if err != nil {
			fmt.Printf("出場失敗：%v\n", err)
			return
		}
		fmt.Printf("%s 出場，%d 元\n", rec.Plate, rec.Fee)
	}

	fmt.Println(g.Enter("abc 1234"), g.Enter("1234ab"))
	fmt.Println(g.Enter("ABC-1234"))
	now = now.Add(5 * time.Minute)
	exit("1234-AB")
	now = now.Add(90 * time.Minute)
	exit("abc1234")
	exit("abc1234")

	parking.WriteReport(os.Stdout, g.Records())
}
