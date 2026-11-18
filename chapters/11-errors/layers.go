package main

import (
	"errors"
	"fmt"
	"net/http"
)

// statusFor 在最外層把錯誤翻譯成 HTTP 狀態碼：只在這裡判斷錯誤的種類。
func statusFor(err error) int {
	switch {
	case err == nil:
		return http.StatusOK
	case errors.Is(err, ErrNotOnMenu):
		return http.StatusNotFound
	case errors.Is(err, ErrSoldOut):
		return http.StatusConflict
	}
	if _, ok := errors.AsType[*StockError](err); ok {
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}

func layersDemo() {
	errs := []error{
		nil,
		fmt.Errorf("訂單 19：%w", ErrNotOnMenu),
		&OrderError{No: 20, Err: reserve("排骨飯", 5)},
		errors.New("資料庫連不上"),
	}
	for _, err := range errs {
		fmt.Println(statusFor(err), err)
	}
}
