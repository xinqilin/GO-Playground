// 專題 4：待辦事項 REST API。
//
// 只用標準庫做出一個可以上線的小服務：路由、JSON、驗證、log、panic 防護，以及優雅關機。
// 用到：struct 與 interface（第 09、10 章）、錯誤處理（第 11 章）、Mutex（第 17 章）、
// context（第 18 章）、測試（第 21 章）、net/http（第 23 章）。
//
// 執行：go run ./projects/todoapi，按 Ctrl-C 會等進行中的請求做完才結束。
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("server error", "err", err)
		os.Exit(1)
	}
}

// run 啟動 server，收到 Ctrl-C 或 SIGTERM 時優雅關機。
func run(logger *slog.Logger) error {
	addr := ":8080"
	if v := os.Getenv("ADDR"); v != "" {
		addr = v
	}
	a := &api{store: NewStore(), logger: logger}
	srv := &http.Server{
		Addr:              addr,
		Handler:           a.routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// 收到中斷訊號時，ctx 會被取消。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		logger.Info("listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err // 例如 port 被占用
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	// 不再接新連線，最多等 5 秒讓進行中的請求做完。
	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
