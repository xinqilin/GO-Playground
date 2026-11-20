//go:build !solution

package toolbox

import (
	"bytes"
	"runtime/debug"
)

// BinaryName 回傳交叉編譯時執行檔的檔名：「名稱-GOOS-GOARCH」，Windows 再加上 .exe。
// 例如 BinaryName("pos", "linux", "arm64") 是 "pos-linux-arm64"，
// BinaryName("pos", "windows", "amd64") 是 "pos-windows-amd64.exe"。
func BinaryName(name, goos, goarch string) string {
	// TODO: 實作
	return ""
}

// Matches 回傳 //go:build 後面的條件 expr 在 tags 這組標籤下成不成立。
// 例如 Matches("linux && !cgo", []string{"linux"}) 是 true。
// expr 的語法不對時回傳錯誤。
func Matches(expr string, tags []string) (bool, error) {
	// TODO: 實作
	return false, nil
}

// FileMatches 回傳 name 這個 .go 檔在 goos/goarch 下會不會被編譯，只看檔名：
//   - 以 _ 或 . 開頭的檔案一律不編譯。
//   - 去掉第一個 . 和後面的部分（.go、.pb.go），再去掉結尾的 _test，
//     剩下的檔名裡第一個 _ 前面的部分不算。
//   - 剩下的部分以 _ 分段，最後兩段是「已知的 GOOS_已知的 GOARCH」時兩個都要符合；
//     否則最後一段是已知的 GOOS 或 GOARCH 時，要符合它。
//   - 其他檔名都會編譯。
//
// 已知的名稱在 KnownOS 和 KnownArch。不用處理 android、ios、illumos 這些
// 「也算是另一個作業系統」的特例。
func FileMatches(name, goos, goarch string) bool {
	// TODO: 實作
	return false
}

// Version 決定程式要顯示的版本，依序：
//  1. injected（-ldflags -X 注入的值）不是空字串也不是 "dev" 時，用它。
//  2. info.Main.Version 不是空字串也不是 "(devel)" 時，用它。
//  3. 建置設定裡有 vcs.revision 時，用它的前 7 個字元；vcs.modified 是 "true" 時再加上 "-dirty"。
//  4. 都沒有時回傳 "dev"。
//
// info 可能是 nil（debug.ReadBuildInfo 失敗時）。
func Version(info *debug.BuildInfo, injected string) string {
	// TODO: 實作
	return ""
}

// GenerateArgs 解析一行 //go:generate 指令，回傳要執行的指令和參數。
// 行首必須剛好是 "//go:generate" 再接一個空白或 tab，否則回傳 ErrNotDirective。
// 參數以空白或 tab 分開；用雙引號括起來的參數依 Go 的字串語法解讀
// （例如 "a b" 是一個參數 a b，"\t" 是 tab），右引號後面必須是空白或行尾。
// 引號沒有成對、字串語法不對、沒有任何參數時回傳錯誤。不用處理 $GOFILE 這類變數。
func GenerateArgs(line string) ([]string, error) {
	// TODO: 實作
	return nil, nil
}

// IsGenerated 回傳 src 是不是程式產生的 Go 原始檔。
// 規則：package 子句之前有一行註解完全符合「// Code generated 任意文字 DO NOT EDIT.」。
//
// 找錯：手寫的檔案只是在程式裡提到這幾個字，也被當成產生的檔案。
func IsGenerated(src []byte) bool {
	return bytes.Contains(src, []byte("DO NOT EDIT"))
}
