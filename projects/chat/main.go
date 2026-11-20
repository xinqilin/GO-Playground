// 專題 6：TCP 聊天室。
//
// 每個連線一個 goroutine 負責讀、一個負責寫；聊天室用一把鎖保護線上名單，
// 廣播只把訊息放進每個人的待送區，跟不上的人會被請出去。Ctrl-C 時先向所有人道別再關閉。
//
// 執行：go run ./projects/chat
//
// 再開幾個終端機連進來：nc localhost 9000
package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"os"
	"os/signal"
)

func main() {
	addr := flag.String("addr", "localhost:9000", "監聽的位址")
	flag.Parse()
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	context.AfterFunc(ctx, stop) // 第一次 Ctrl-C 之後恢復預設行為：道別卡住時，再按一次就直接結束
	err := run(ctx, *addr, logger)
	stop()
	if err != nil {
		logger.Error("聊天室停止", "err", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, addr string, logger *slog.Logger) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	logger.Info("聊天室開張", "addr", ln.Addr().String())
	return NewServer(logger).Serve(ctx, ln)
}
