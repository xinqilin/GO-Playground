package main

import (
	"errors"
	"fmt"

	"github.com/xinqilin/GO-Playground/chapters/33-idiomatic-go/ticket"
)

var errNoCounter = errors.New("沒有指定櫃台")

// serveNested 叫下一號，回傳櫃台螢幕要顯示的文字。
// 每一個條件都把正常的流程往右推一層，錯誤處理離它的條件越來越遠。
func serveNested(m *ticket.Machine, counter string) (string, error) {
	if counter != "" {
		n, err := m.Call()
		if err == nil {
			msg := fmt.Sprintf("請 %d 號到%s櫃台", n, counter)
			if w := m.Waiting(); w >= 3 {
				msg += fmt.Sprintf("，後面還有 %d 位", w)
			}
			return msg, nil
		} else {
			return "", fmt.Errorf("%s櫃台：%w", counter, err)
		}
	} else {
		return "", errNoCounter
	}
}

// serve 和 serveNested 做一樣的事：先處理例外，正常的流程沿著左邊一路往下。
func serve(m *ticket.Machine, counter string) (string, error) {
	if counter == "" {
		return "", errNoCounter
	}
	n, err := m.Call()
	if err != nil {
		return "", fmt.Errorf("%s櫃台：%w", counter, err)
	}
	msg := fmt.Sprintf("請 %d 號到%s櫃台", n, counter)
	if w := m.Waiting(); w >= 3 {
		msg += fmt.Sprintf("，後面還有 %d 位", w)
	}
	return msg, nil
}

func flowDemo() {
	var m ticket.Machine
	for range 4 {
		m.Take()
	}
	for _, counter := range []string{"結帳", "結帳", "", "取餐", "取餐", "取餐"} {
		msg, err := serve(&m, counter)
		switch {
		case errors.Is(err, ticket.ErrNoOneWaiting):
			fmt.Printf("沒有人在等：%v\n", err)
		case err != nil:
			fmt.Printf("錯誤：%v\n", err)
		default:
			fmt.Println(msg)
		}
	}
}
