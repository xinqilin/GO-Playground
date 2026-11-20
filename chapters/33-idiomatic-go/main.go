// 第 33 章：寫出道地的 Go。
//
// 這一章的主角是 ticket 這個 package：一台排隊叫號機。它示範了命名、doc comment、
// 零值可用、回傳具體型別、在邊界上複製這些慣例；main 裡的程式則是它的使用端：
// 在使用的地方定義小 interface、先處理例外讓正常流程靠左。
//
// 執行：go run ./chapters/33-idiomatic-go
package main

import "fmt"

func main() {
	section("3. 零值就能用")
	zeroDemo()

	section("4. 接受 interface，回傳具體型別")
	interfaceDemo()

	section("5. 正常流程靠左")
	flowDemo()

	section("6. 在邊界上複製")
	calledDemo()
}

func section(title string) {
	fmt.Println("==", title, "==")
}
