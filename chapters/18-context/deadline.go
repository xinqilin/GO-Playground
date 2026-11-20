package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// errDemo：被取消和超過期限是兩種不同的錯誤，用 errors.Is 分辨。
func errDemo() {
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	expired, cancel2 := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel2()

	for _, ctx := range []context.Context{canceled, expired} {
		<-ctx.Done()
		err := ctx.Err()
		fmt.Printf("%v：Canceled=%t DeadlineExceeded=%t\n",
			err, errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded))
	}
}

// deadlineDemo：衍生出來的 ctx 期限只能比上層早，不能比上層晚。
func deadlineDemo() {
	parent, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	long, cancelLong := context.WithTimeout(parent, time.Hour)
	defer cancelLong()
	short, cancelShort := context.WithTimeout(parent, 10*time.Millisecond)
	defer cancelShort()

	pd, _ := parent.Deadline()
	ld, _ := long.Deadline()
	sd, _ := short.Deadline()
	fmt.Printf("要 1 小時，實際期限和上層一樣：%t\n", ld.Equal(pd))
	fmt.Printf("要 10 毫秒，期限比上層早：%t\n", sd.Before(pd))

	_, ok := context.Background().Deadline()
	fmt.Printf("Background 有期限：%t\n", ok)
}
