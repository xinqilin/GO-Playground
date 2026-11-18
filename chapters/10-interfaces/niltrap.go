package main

import "fmt"

// FormError 是表單驗證失敗的錯誤。
type FormError struct {
	Field string
}

func (e *FormError) Error() string {
	return e.Field + "沒有填"
}

// checkNameBuggy 想在沒有錯誤時回傳 nil，回傳的卻是「裝著 nil 指標的 interface」。
func checkNameBuggy(name string) error {
	var err *FormError
	if name == "" {
		err = &FormError{Field: "姓名"}
	}
	return err
}

// checkName 沒有錯誤時直接寫 return nil。
func checkName(name string) error {
	if name == "" {
		return &FormError{Field: "姓名"}
	}
	return nil
}

func nilTrapDemo() {
	err := checkNameBuggy("王小明")
	fmt.Println(err == nil)
	fmt.Printf("%T\n", err)

	fmt.Println(checkName("王小明") == nil)
	fmt.Println(checkName(""))
}
