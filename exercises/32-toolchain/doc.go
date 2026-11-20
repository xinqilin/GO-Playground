// Package toolbox 是第 32 章「工具鏈」的練習題：自己實作 go 指令裡的幾個小規則，
// 看懂 build tag、檔名、go generate 和建置資訊是怎麼運作的。
//
// 題目在 exercises.go：把每個 TODO 換成你的實作，或修好標著「找錯」的函式，然後執行：
//
//	go test ./exercises/32-toolchain/
//
// 只想跑某一題時加上 -run：
//
//	go test ./exercises/32-toolchain/ -run TestFileMatches
//
// 參考解答在 solution.go，用 -tags solution 才會編譯它：
//
//	go test -tags solution ./exercises/32-toolchain/
package toolbox

import "errors"

// ErrNotDirective 表示這一行不是 //go:generate 指令。
var ErrNotDirective = errors.New("不是 //go:generate 指令")

// KnownOS 是 go 指令認得的作業系統名稱（GOOS），用在檔名規則。
var KnownOS = map[string]bool{
	"aix": true, "android": true, "darwin": true, "dragonfly": true, "freebsd": true,
	"illumos": true, "ios": true, "js": true, "linux": true, "netbsd": true,
	"openbsd": true, "plan9": true, "solaris": true, "wasip1": true, "windows": true,
}

// KnownArch 是 go 指令認得的處理器架構名稱（GOARCH），用在檔名規則。
var KnownArch = map[string]bool{
	"386": true, "amd64": true, "arm": true, "arm64": true, "loong64": true,
	"mips": true, "mipsle": true, "mips64": true, "mips64le": true,
	"ppc64": true, "ppc64le": true, "riscv64": true, "s390x": true, "wasm": true,
}
