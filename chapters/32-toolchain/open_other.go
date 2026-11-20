//go:build !darwin && !windows

package main

// 其他系統（Linux、各種 BSD…）沒有對應的檔名，用 build 條件排除 darwin 和 windows。

// openCommand 回傳在這個作業系統上用預設瀏覽器打開網址的指令。
func openCommand() []string {
	return []string{"xdg-open"}
}
