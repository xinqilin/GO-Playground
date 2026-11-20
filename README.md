# Go 圖解手冊

[![ci](https://github.com/xinqilin/GO-Playground/actions/workflows/ci.yml/badge.svg)](https://github.com/xinqilin/GO-Playground/actions/workflows/ci.yml)

繁體中文的 Go 手冊：先修 4 章給從來沒寫過程式的人，正文 34 章從基礎、型別、並行、標準庫一路到工程實務，最後是 7 個專題。每一章先用圖把觀念講懂，再一步一步看可以執行、有測試的程式碼，附上常見錯誤、Java 對照，以及有自動測試的練習（全書 243 題）。

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

第 28 章（資料庫）是獨立的 module，根目錄的 `./...` 不會包含它，要到它的目錄裡跑：

```sh
go -C chapters/28-database test .                      # 範例程式的測試（./... 會連練習一起跑）
go -C chapters/28-database test ./exercises/           # 這一章的練習
```

每一章都是獨立的程式：`chapters/NN-主題/`（先修是 `chapters/pre-N-主題/`）放程式碼和測試，`docs/` 裡同名的 `.html` 是對應的圖解頁。每一段示範都有一個 `Example` 測試比對它印出來的內容。

`exercises/NN-主題/` 是練習題：`exercises.go` 是題目，`solution.go` 是參考解答（加上 `-tags solution` 才會編譯），`exercises_test.go` 是測試。每一題的測試一開始都會失敗，寫對一題就通過一題；「找錯」題的題目檔裡是有 bug 的程式，要找出來修好。因為練習還沒寫完，直接在根目錄跑 `go test ./...` 會看到練習題的紅燈，這是預期的。

頁面上的程式碼和輸出由 `go run ./tools/docsync` 從 repo 裡的程式同步過去，CI 會用 `-check` 確認兩邊一致。

## 目錄

### 先修：從來沒寫過程式的人從這裡開始

| 章 | 主題 | 圖解 | 程式碼 |
|---|---|---|---|
| 先1 | 程式怎麼跑起來 | [看圖](https://bill-lin.dev/GO-Playground/pre-1-running.html) | [chapters/pre-1-running](chapters/pre-1-running) |
| 先2 | 值、變數與運算 | [看圖](https://bill-lin.dev/GO-Playground/pre-2-values.html) | [chapters/pre-2-values](chapters/pre-2-values) |
| 先3 | 判斷與重複 | [看圖](https://bill-lin.dev/GO-Playground/pre-3-decisions.html) | [chapters/pre-3-decisions](chapters/pre-3-decisions) |
| 先4 | 拆解問題：函式與清單 | [看圖](https://bill-lin.dev/GO-Playground/pre-4-problems.html) | [chapters/pre-4-problems](chapters/pre-4-problems) |

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
| P1 | 費氏數列：五種寫法、整數溢位、math/big、benchmark | [看圖](https://bill-lin.dev/GO-Playground/project-fibonacci.html) | [projects/fibonacci](projects/fibonacci) |
| P2 | 字詞統計：多檔輸入、中英文斷詞、排行、自訂旗標型別 | [看圖](https://bill-lin.dev/GO-Playground/project-wordcount.html) | [projects/wordcount](projects/wordcount) |
| P3 | 並行網址檢查器：worker pool、逾時、限速、退避重試、synctest | [看圖](https://bill-lin.dev/GO-Playground/project-urlchecker.html) | [projects/urlchecker](projects/urlchecker) |
| P4 | 待辦事項 API：三層架構、驗證、並行修改、原子存檔、優雅關機 | [看圖](https://bill-lin.dev/GO-Playground/project-todoapi.html) | [projects/todoapi](projects/todoapi) |
| P5 | 日誌分析：regexp、gzip、百分位數、三種輸出、exit code | [看圖](https://bill-lin.dev/GO-Playground/project-logstats.html) | [projects/logstats](projects/logstats) |
| P6 | TCP 聊天室：每連線兩個 goroutine、鎖與待送區、deadline、優雅關機、net.Pipe + synctest | [看圖](https://bill-lin.dev/GO-Playground/project-chat.html) | [projects/chat](projects/chat) |
| P7 | TTL＋LRU 快取：泛型、map 加雙向串列、Mutex 與 RLock 的陷阱、只載入一次、synctest、benchmark | [看圖](https://bill-lin.dev/GO-Playground/project-cache.html) | [projects/cache](projects/cache) |

### 附錄

- [術語表](https://bill-lin.dev/GO-Playground/glossary.html)：98 個術語，各用一兩句話說明，再連回詳細介紹它的那一節。

## 每一章的結構

- **本章內容、本章你會學到**：這一章要解決什麼問題，需要先讀哪幾章。
- **編號的小節**：一節一個主題，先看圖，再看說明、程式碼和實際的輸出。標著「深入」的段落講原理和少見的情況，第一次讀可以先跳過。
- **常見錯誤**：錯誤和正確的寫法並排；編譯不過的例子附上真的編譯錯誤訊息。
- **從 Java 來的你**：同一件事在 Java 和 Go 的寫法對照（先修沒有這一段）。
- **練習**：5 到 7 題，分成入門、進階、挑戰和找錯，每一題都有測試和提示。
- **重點回顧、延伸閱讀**：一份檢查清單，以及官方文件的連結。
- **動手跑**：在 repo 裡執行的指令。

專題頁則是走讀一個完整的程式：從需求、設計到測試，最後是讓你自己動手改的延伸挑戰。

## 參考資料

- [A Tour of Go](https://go.dev/tour/)
- [Effective Go](https://go.dev/doc/effective_go)
- [The Go Programming Language Specification](https://go.dev/ref/spec)
- [Go 版本說明](https://go.dev/doc/devel/release)
- [Go by Example](https://gobyexample.com/)
- [jaywcjlove/golang-tutorial](https://github.com/jaywcjlove/golang-tutorial)：這個 repo 最早的學習參考

最早的練習程式（`Basic/`、`Exercise/`）留在 git 歷史裡，可以從 commit `8bf0c12` 找到。
