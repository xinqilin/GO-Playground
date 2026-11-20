package main

import (
	"fmt"
	"testing"
)

func Example_ifElse() {
	ifElse()
	// Output:
	// 33 度：開冷氣
	// 24 度：開窗就好
	// 12 度：穿外套
}

func Example_ifInit() {
	ifInit()
	// Output:
	// 訂 3 個便當，共 255 元
	// "三" 不是數字
	// 至少要訂 1 個
}

func Example_guardClause() {
	guardClause()
	// Output:
	// 還沒付款 | 還沒付款
	// 缺少地址 | 缺少地址
	// 購物車是空的 | 購物車是空的
	// 可以出貨 | 可以出貨
}

func Example_shadowing() {
	shadowing()
	// Output:
	// 區塊裡： 50
	// 區塊外： 100
}

func Example_forClassic() {
	forClassic()
	// Output:
	// 請 1 號到櫃台
	// 請 2 號到櫃台
	// 請 3 號到櫃台
}

func Example_forCondition() {
	forCondition()
	// Output:
	// 4 個月後存到 12000 元
}

func Example_forInfinite() {
	forInfinite()
	// Output:
	// 叫號 A01
	// 叫號 A02
	// 叫號 A03
	// 沒有人在排隊了
}

func Example_forRangeInt() {
	forRangeInt()
	// Output:
	// 第 0 圈
	// 第 1 圈
	// 第 2 圈
}

func Example_forContinue() {
	forContinue()
	// Output:
	// 0 是 3 的倍數
	// 3 是 3 的倍數
	// 6 是 3 的倍數
	// 9 是 3 的倍數
}

func Example_rangeCopy() {
	rangeCopy()
	// Output:
	// [{雞腿飯 3} {排骨飯 5}]
	// [{雞腿飯 0} {排骨飯 0}]
}

func Example_rangeArrayCopy() {
	rangeArrayCopy()
	// Output:
	// 陣列： [空 空 空]
	// slice： [空 空 有人]
}

func Example_rangeString() {
	rangeString()
	// Output:
	// 位置 0：G（U+0047）
	// 位置 1：o（U+006F）
	// 位置 2：語（U+8A9E）
	// 位置 5：言（U+8A00）
	// len = 8，字元數 = 4
}

func Example_rangeMap() {
	rangeMap()
	// Output:
	// 排骨 5
	// 雞腿 3
	// 魚排 0
}

func Example_rangeChannel() {
	rangeChannel()
	// Output:
	// 叫號 A01
	// 叫號 A02
	// 叫號 A03
}

func Example_rangeFunc() {
	rangeFunc()
	// Output:
	// 3
	// 2
	// 1
	// 發射
}

func Example_rangeOnce() {
	rangeOnce()
	// Output:
	// 跑了 3 圈，nums = [1 2 3 10 20 30]
}

func Example_loopVarPerIteration() {
	loopVarPerIteration()
	// Output:
	// 0 3
	// 1 3
	// 2 3
}

func Example_labeledBreak() {
	labeledBreak()
	// Output:
	// 第 0 和第 1 個加起來剛好 100 元
}

func Example_labeledContinue() {
	labeledContinue()
	// Output:
	// 第 1 班全部及格
	// 第 2 班有人不及格： 55
	// 第 3 班全部及格
}

func Example_breakInSwitch() {
	breakInSwitch()
	// Output:
	// 沒有標籤： 3
	// 有標籤： 2
}

func Example_switchValue() {
	switchValue()
	// Output:
	// 星期六：假日
	// 星期二：上班
	// 星期日：假日
}

func Example_switchInit() {
	switchInit()
	// Output:
	// 王明：一般長度
	// 陳小美：一般長度
	// 歐陽靖文：可能是複姓
	// 李：1 個字，請再確認
}

func Example_switchTrue() {
	switchTrue()
	// Output:
	// AQI 35：良好
	// AQI 120：對敏感族群不健康
	// AQI 320：危害
}

func Example_switchOrder() {
	switchOrder()
	// Output:
	// 檢查 n > 10
	// 檢查 n > 5
	// 大於 5
}

func Example_switchFallthrough() {
	switchFallthrough()
	// Output:
	// 金卡 [機場接送 免運 生日禮]
	// 銀卡 [免運 生日禮]
	// 一般 [生日禮]
	// 訪客 []
}

func Example_typeSwitch() {
	typeSwitch()
	// Output:
	// 整數 42，兩倍是 84
	// 字串 "便當"，長度 6
	// 其他型別 float64
	// nil
}

func Example_shadowPitfall() {
	shadowPitfall()
	// Output:
	// if 裡讀到： dark
	// 錯誤寫法回傳： light
	// 正確寫法回傳： dark
}

// Example_reverseBuggy 用 recover 接住 panic，把錯誤訊息印出來給頁面使用。
func Example_reverseBuggy() {
	defer func() { fmt.Println("panic:", recover()) }()
	reverseBuggy([]string{"A01", "A02"})
	// Output:
	// panic: runtime error: index out of range [18446744073709551615] with length 2
}

func Example_reverse() {
	fmt.Println(reverse([]string{"A01", "A02", "A03"}))
	fmt.Println(reverse(nil))
	// Output:
	// [A03 A02 A01]
	// []
}

func TestAqiLevel(t *testing.T) {
	tests := []struct {
		aqi  int
		want string
	}{
		{0, "良好"},
		{50, "良好"},
		{51, "普通"},
		{100, "普通"},
		{101, "對敏感族群不健康"},
		{150, "對敏感族群不健康"},
		{151, "對所有族群不健康"},
		{200, "對所有族群不健康"},
		{201, "非常不健康"},
		{300, "非常不健康"},
		{301, "危害"},
	}
	for _, tt := range tests {
		if got := aqiLevel(tt.aqi); got != tt.want {
			t.Errorf("aqiLevel(%d) = %q, want %q", tt.aqi, got, tt.want)
		}
	}
}

func TestCountUntilStop(t *testing.T) {
	tests := []struct {
		name string
		cmds []string
		want int
	}{
		{"沒有 stop", []string{"a", "b"}, 2},
		{"中間有 stop", []string{"a", "stop", "b"}, 1},
		{"一開始就 stop", []string{"stop", "a"}, 0},
		{"空的", nil, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countUntilStop(tt.cmds); got != tt.want {
				t.Errorf("countUntilStop(%q) = %d, want %d", tt.cmds, got, tt.want)
			}
		})
	}
}
