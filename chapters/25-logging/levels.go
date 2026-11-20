package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// levelsDemo 層級是整數，數字越大越嚴重；不在四個名稱上的值，印成「最接近的名稱 + 差距」。
func levelsDemo() {
	fmt.Println(slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError)
	fmt.Println(int(slog.LevelDebug), int(slog.LevelInfo), int(slog.LevelWarn), int(slog.LevelError))
	fmt.Println(slog.LevelInfo+2, slog.Level(-8))

	for _, s := range []string{"warn", "INFO+2", "verbose"} {
		var lv slog.Level
		if err := lv.UnmarshalText([]byte(s)); err != nil {
			fmt.Println(s, "→", err)
			continue
		}
		fmt.Println(s, "→", lv)
	}
}

// levelVarDemo HandlerOptions.Level 決定最低要記哪個層級。
// 傳 *slog.LevelVar 的話，程式執行中可以隨時調整，不用重新建立 Logger。
func levelVarDemo() {
	var level slog.LevelVar // 零值是 Info
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: &level, ReplaceAttr: dropTime}))

	logger.Debug("看不到", "step", 1)
	logger.Info("接單", "table", 7)

	level.Set(slog.LevelDebug) // 例如收到管理指令，暫時打開 debug
	logger.Debug("看得到了", "step", 2)

	level.Set(slog.LevelWarn)
	logger.Info("看不到", "table", 8)
	logger.Warn("排骨快賣完了", "left", 3)
	fmt.Println("Info 要不要記：", logger.Enabled(context.Background(), slog.LevelInfo))
}
