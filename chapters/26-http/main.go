// 第 26 章：HTTP 服務。
//
// 這一章的程式是一個真的可以跑的便當店訂單 API。頁面上的示範是各個檔案裡的 ...Demo 函式，
// 它們用 httptest 直接呼叫 handler，不用開 port；main_test.go 用 Example 鎖住它們的輸出。
//
// 執行：go run ./chapters/26-http，再開另一個終端機用 curl 測試，按 Ctrl-C 優雅關機。
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
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	a := &api{shop: newShop(), logger: logger}
	srv := newServer(":8080", a.handler(), logger)
	logger.Info("開始營業", "addr", srv.Addr)
	if err := run(ctx, srv, logger); err != nil {
		logger.Error("server 停止", "err", err)
		os.Exit(1)
	}
	logger.Info("打烊")
}

// newServer 設好逾時。直接用 http.ListenAndServe 的話，這些值全部是 0，也就是沒有逾時。
func newServer(addr string, h http.Handler, logger *slog.Logger) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,  // 讀完請求標頭的期限
		ReadTimeout:       10 * time.Second, // 讀完整個請求（含 body）的期限
		WriteTimeout:      15 * time.Second, // 從讀完標頭到寫完回應的期限
		IdleTimeout:       60 * time.Second, // keep-alive 連線閒置多久就關掉
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
}

// run 開始服務，ctx 被取消（例如收到 Ctrl-C）時優雅關機：不再接新的連線，等進行中的請求做完。
func run(ctx context.Context, srv *http.Server, logger *slog.Logger) error {
	errc := make(chan error, 1)
	go func() { errc <- srv.ListenAndServe() }()

	select {
	case err := <-errc:
		return err // 還沒收到關機訊號就停了，例如 port 已經被佔用
	case <-ctx.Done():
	}

	logger.Info("關機中，等進行中的請求做完")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err // 10 秒內沒做完
	}
	if err := <-errc; !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
