# Go 圖解手冊

[![ci](https://github.com/xinqilin/GO-Playground/actions/workflows/ci.yml/badge.svg)](https://github.com/xinqilin/GO-Playground/actions/workflows/ci.yml)

給工程師的 Go 入門手冊（繁體中文）。每一章先用大圖把觀念講懂，再附上可以執行、有測試的程式碼，並對照 Java 的寫法。

**線上閱讀：<https://bill-lin.dev/GO-Playground/>**

## 怎麼使用

需要 Go 1.21 以上。`go.mod` 宣告的是 `go 1.27`，比較舊的 Go 會自動下載需要的版本。

```sh
git clone https://github.com/xinqilin/GO-Playground.git
cd GO-Playground
go run ./chapters/00-hello                  # 執行某一章
go test ./chapters/... ./projects/...       # 跑章節和專題的測試
go test ./exercises/03-control-flow/        # 做某一章的練習（寫完之前會失敗）
go test -race -tags solution ./...          # 全部測試，練習題改用參考解答
```

每一章都是獨立的程式：`chapters/NN-主題/` 放程式碼和測試，`docs/NN-主題.html` 是對應的圖解頁。測試會比對每一章印出來的內容，所以頁面上的輸出和程式碼不會脫節。

`exercises/NN-主題/` 是練習題：`exercises.go` 是題目，`solution.go` 是參考解答（加上 `-tags solution` 才會編譯）。題目還沒寫完時測試本來就會失敗，所以直接在根目錄跑 `go test ./...` 會看到練習題的紅燈。

頁面上的程式碼和輸出由 `go run ./tools/docsync` 從 repo 裡的程式同步過去，CI 會用 `-check` 確認兩邊一致。

## 目錄

### Part 0 開始

| 章 | 主題 | 圖解 | 程式碼 |
|---|---|---|---|
| 00 | 第一支程式 | [看圖](https://bill-lin.dev/GO-Playground/00-hello.html) | [chapters/00-hello](chapters/00-hello) |

### Part 1 基礎

| 章 | 主題 | 圖解 | 程式碼 |
|---|---|---|---|
| 01 | 變數與型別 | [看圖](https://bill-lin.dev/GO-Playground/01-variables.html) | [chapters/01-variables](chapters/01-variables) |
| 02 | 輸入與輸出 | [看圖](https://bill-lin.dev/GO-Playground/02-input-output.html) | [chapters/02-input-output](chapters/02-input-output) |
| 03 | 流程控制 | [看圖](https://bill-lin.dev/GO-Playground/03-control-flow.html) | [chapters/03-control-flow](chapters/03-control-flow) |
| 04 | 函式 | [看圖](https://bill-lin.dev/GO-Playground/04-functions.html) | [chapters/04-functions](chapters/04-functions) |
| 05 | 指標 | [看圖](https://bill-lin.dev/GO-Playground/05-pointers.html) | [chapters/05-pointers](chapters/05-pointers) |
| 06 | 陣列與切片 | [看圖](https://bill-lin.dev/GO-Playground/06-arrays-slices.html) | [chapters/06-arrays-slices](chapters/06-arrays-slices) |
| 07 | map | [看圖](https://bill-lin.dev/GO-Playground/07-maps.html) | [chapters/07-maps](chapters/07-maps) |
| 08 | 字串與 rune | [看圖](https://bill-lin.dev/GO-Playground/08-strings-runes.html) | [chapters/08-strings-runes](chapters/08-strings-runes) |

### Part 2 型別系統

| 章 | 主題 | 圖解 | 程式碼 |
|---|---|---|---|
| 09 | struct 與方法 | [看圖](https://bill-lin.dev/GO-Playground/09-structs-methods.html) | [chapters/09-structs-methods](chapters/09-structs-methods) |
| 10 | interface | [看圖](https://bill-lin.dev/GO-Playground/10-interfaces.html) | [chapters/10-interfaces](chapters/10-interfaces) |
| 11 | 錯誤處理 | [看圖](https://bill-lin.dev/GO-Playground/11-errors.html) | [chapters/11-errors](chapters/11-errors) |
| 12 | 泛型 | [看圖](https://bill-lin.dev/GO-Playground/12-generics.html) | [chapters/12-generics](chapters/12-generics) |
| 13 | 疊代器 | [看圖](https://bill-lin.dev/GO-Playground/13-iterators.html) | [chapters/13-iterators](chapters/13-iterators) |

### Part 3 並行

| 章 | 主題 | 圖解 | 程式碼 |
|---|---|---|---|
| 14 | goroutine | [看圖](https://bill-lin.dev/GO-Playground/14-goroutines.html) | [chapters/14-goroutines](chapters/14-goroutines) |
| 15 | channel | [看圖](https://bill-lin.dev/GO-Playground/15-channels.html) | [chapters/15-channels](chapters/15-channels) |
| 16 | select | [看圖](https://bill-lin.dev/GO-Playground/16-select.html) | [chapters/16-select](chapters/16-select) |
| 17 | sync 與 data race | [看圖](https://bill-lin.dev/GO-Playground/17-sync.html) | [chapters/17-sync](chapters/17-sync) |
| 18 | context | [看圖](https://bill-lin.dev/GO-Playground/18-context.html) | [chapters/18-context](chapters/18-context) |
| 19 | 並行模式 | [看圖](https://bill-lin.dev/GO-Playground/19-concurrency-patterns.html) | [chapters/19-concurrency-patterns](chapters/19-concurrency-patterns) |

### Part 4 標準庫實戰

| 章 | 主題 | 圖解 | 程式碼 |
|---|---|---|---|
| 20 | 檔案與 I/O | [看圖](https://bill-lin.dev/GO-Playground/20-files.html) | [chapters/20-files](chapters/20-files) |
| 21 | 命令列工具 | [看圖](https://bill-lin.dev/GO-Playground/21-cli.html) | [chapters/21-cli](chapters/21-cli) |
| 22 | 時間 | [看圖](https://bill-lin.dev/GO-Playground/22-time.html) | [chapters/22-time](chapters/22-time) |
| 23 | JSON 與資料編碼 | [看圖](https://bill-lin.dev/GO-Playground/23-json.html) | [chapters/23-json](chapters/23-json) |
| 24 | 文字處理 | [看圖](https://bill-lin.dev/GO-Playground/24-text.html) | [chapters/24-text](chapters/24-text) |
| 25 | 日誌 | [看圖](https://bill-lin.dev/GO-Playground/25-logging.html) | [chapters/25-logging](chapters/25-logging) |
| 26 | HTTP 服務 | [看圖](https://bill-lin.dev/GO-Playground/26-http.html) | [chapters/26-http](chapters/26-http) |
| 27 | HTTP 用戶端 | [看圖](https://bill-lin.dev/GO-Playground/27-http-client.html) | [chapters/27-http-client](chapters/27-http-client) |
| 28 | 資料庫 | [看圖](https://bill-lin.dev/GO-Playground/28-database.html) | [chapters/28-database](chapters/28-database) |

### Part 5 工程實務

| 章 | 主題 | 圖解 | 程式碼 |
|---|---|---|---|
| 29 | package 與 module | [看圖](https://bill-lin.dev/GO-Playground/29-packages-modules.html) | [chapters/29-packages-modules](chapters/29-packages-modules) |
| 30 | 測試 | [看圖](https://bill-lin.dev/GO-Playground/30-testing.html) | [chapters/30-testing](chapters/30-testing) |
| 31 | 效能與 profiling | [看圖](https://bill-lin.dev/GO-Playground/31-performance.html) | [chapters/31-performance](chapters/31-performance) |
| 32 | 工具鏈 | [看圖](https://bill-lin.dev/GO-Playground/32-toolchain.html) | [chapters/32-toolchain](chapters/32-toolchain) |
| 33 | 寫出道地的 Go | [看圖](https://bill-lin.dev/GO-Playground/33-idiomatic-go.html) | [chapters/33-idiomatic-go](chapters/33-idiomatic-go) |

### 專題

| 專題 | 內容 | 圖解 | 程式碼 |
|---|---|---|---|
| P1 | 費氏數列：修好舊版的 bug，比較三種寫法 | [看圖](https://bill-lin.dev/GO-Playground/project-fibonacci.html) | [projects/fibonacci](projects/fibonacci) |
| P2 | 字詞統計：Scanner、map、排序 | [看圖](https://bill-lin.dev/GO-Playground/project-wordcount.html) | [projects/wordcount](projects/wordcount) |
| P3 | 並行網址檢查器：semaphore、逾時、interface | [看圖](https://bill-lin.dev/GO-Playground/project-urlchecker.html) | [projects/urlchecker](projects/urlchecker) |
| P4 | 待辦事項 API：REST、middleware、優雅關機 | [看圖](https://bill-lin.dev/GO-Playground/project-todoapi.html) | [projects/todoapi](projects/todoapi) |

## 每一章的結構

- **看圖**：大圖少字，先把觀念看懂。
- **最小範例**：可以直接執行的短程式，輸出是實際跑出來的結果。
- **工程師筆記**：實務上會踩的坑和慣例。
- **從 Java 來的你**：同一件事在 Java 和 Go 的寫法對照。
- **動手跑**：在 repo 裡執行的指令。

## 參考資料

- [A Tour of Go](https://go.dev/tour/)
- [Effective Go](https://go.dev/doc/effective_go)
- [The Go Programming Language Specification](https://go.dev/ref/spec)
- [Go 版本說明](https://go.dev/doc/devel/release)
- [Go by Example](https://gobyexample.com/)
- [jaywcjlove/golang-tutorial](https://github.com/jaywcjlove/golang-tutorial)：這個 repo 最早的學習參考

最早的練習程式（`Basic/`、`Exercise/`）留在 git 歷史裡，可以從 commit `8bf0c12` 找到。
