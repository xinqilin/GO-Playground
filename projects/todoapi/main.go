// 專題 4：待辦事項 REST API。
//
// 只用標準庫做出一個可以上線的小服務，分成三層：api 處理 HTTP，service 放規則，Store 保存資料。
// 資料存在 JSON 檔裡，每次修改都先寫好檔案才回應；按 Ctrl-C 會等進行中的請求做完才結束。
//
// 執行：go run ./projects/todoapi -data todos.json
//
//	curl -i -X POST localhost:8080/todos -H 'Content-Type: application/json' -d '{"title":"訂明天的便當"}'
package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "監聽的位址")
	data := flag.String("data", "todos.json", "存資料的 JSON 檔；空字串代表只放在記憶體")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err := run(ctx, *addr, *data, logger)
	stop()
	if err != nil {
		logger.Error("server 停止", "err", err)
		os.Exit(1)
	}
}

// run 組合三層、開始服務；ctx 被取消時優雅關機。
func run(ctx context.Context, addr, dataPath string, logger *slog.Logger) error {
	store, err := OpenStore(dataPath)
	if err != nil {
		return err
	}
	svc := &service{store: store, now: time.Now}
	a := &api{svc: svc, logger: logger}
	srv := &http.Server{
		Addr:              addr,
		Handler:           a.routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	errc := make(chan error, 1)
	go func() { errc <- srv.ListenAndServe() }()
	logger.Info("開始服務", "addr", addr, "data", dataPath)

	select {
	case err := <-errc:
		return err // 例如 port 已經被佔用
	case <-ctx.Done():
	}
	logger.Info("關機中，等進行中的請求做完")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	if err := <-errc; !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
