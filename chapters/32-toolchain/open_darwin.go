package main

// 檔名以 _darwin.go 結尾：只有 GOOS=darwin（macOS）時才會被編譯，不需要寫 //go:build。

// openCommand 回傳在這個作業系統上用預設瀏覽器打開網址的指令。
func openCommand() []string {
	return []string{"open"}
}
