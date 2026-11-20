package main

// 檔名以 _windows.go 結尾：只有 GOOS=windows 時才會被編譯。

// openCommand 回傳在這個作業系統上用預設瀏覽器打開網址的指令。
func openCommand() []string {
	return []string{"rundll32", "url.dll,FileProtocolHandler"}
}
