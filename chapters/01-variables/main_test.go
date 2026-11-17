package main

import "testing"

func Example_declareVar() {
	declareVar()
	// Output:
	// 85 3 false
}

func Example_declareShort() {
	declareShort()
	// Output:
	// 排骨便當 90
}

func Example_declareMulti() {
	declareMulti()
	// Output:
	// 10 20 雞腿便當 2
	// A02 A01
}

func Example_redeclare() {
	redeclare()
	// Output:
	// 85 <nil>
	// 0 strconv.Atoi: parsing "三": invalid syntax
}

func Example_typeInference() {
	typeInference()
	// Output:
	// int float64 int32 string bool
}

func Example_packageLevel() {
	packageLevel()
	// Output:
	// 店名：阿姨便當
}

func Example_zeroValues() {
	zeroValues()
	// Output:
	// int=0 float64=0 bool=false string="" pointer=<nil>
}

func Example_zeroUsable() {
	zeroUsable()
	// Output:
	// 2 個：排骨、雞腿
}

func Example_intRanges() {
	intRanges()
	// Output:
	// int8   -128 ~ 127
	// uint8  0 ~ 255
	// int16  -32768 ~ 32767
	// int32  -2147483648 ~ 2147483647
	// int64  -9223372036854775808 ~ 9223372036854775807
	// uint64 0 ~ 18446744073709551615
}

func Example_overflowWrap() {
	overflowWrap()
	// Output:
	// -128
	// 255
}

func Example_integerDivision() {
	integerDivision()
	// Output:
	// 3 1
	// -3 -1
	// 3.5
	// 每人 333 元，剩 1 元
}

func Example_divideByZero() {
	divideByZero()
	// Output:
	// panic: runtime error: integer divide by zero
}

func Example_floatError() {
	floatError()
	// Output:
	// 0.30000000000000004 false
	// 0.9999999999999999
}

func Example_floatCompare() {
	floatCompare()
	// Output:
	// true
}

func Example_moneyInTenths() {
	moneyInTenths()
	// Output:
	// 6.300000000000001
	// 6.3 元
}

func Example_floatSpecial() {
	floatSpecial()
	// Output:
	// +Inf -Inf
	// NaN false true
}

func Example_convertNumbers() {
	convertNumbers()
	// Output:
	// 66 66.66666666666667
	// 85 -85
	// 44
}

func Example_convertStrings() {
	convertStrings()
	// Output:
	// 85 元
	// 121 <nil>
	// strconv.Atoi: parsing "12a": invalid syntax
	// A 65
}

func Example_constBasics() {
	constBasics()
	// Output:
	// 阿姨便當 85 4.25
}

func Example_untypedConst() {
	untypedConst()
	// Output:
	// 3 3 3
	// 4
}

func Example_iotaWeekday() {
	iotaWeekday()
	// Output:
	// 0 1 6
	// true
}

func Example_iotaPatterns() {
	iotaPatterns()
	// Output:
	// 0 1 2 3
	// 1 2 4 8
}

func Example_scopeLevels() {
	scopeLevels()
	// Output:
	// 套件層級： 阿姨便當
	// 區塊裡： 255
	// 函式裡： 3
}

func Example_shadowPackage() {
	shadowPackage()
	// Output:
	// 阿姨便當 二號店
	// 阿姨便當
}

func Example_precedence() {
	precedence()
	// Output:
	// 14 20
	// 5
	// true
	// true
}

func Example_shortCircuit() {
	shortCircuit()
	// Output:
	// true false
}

func Example_incDec() {
	incDec()
	// Output:
	// 8 170
}

func Example_bitFlags() {
	bitFlags()
	// Output:
	// 0101
	// true false
	// 1100
	// 1110
}

func TestExpensive(t *testing.T) {
	tests := []struct {
		total, qty int
		want       bool
	}{
		{450, 3, true},
		{300, 3, false},
		{303, 3, true},
		{0, 0, false},
		{500, 0, false},
	}
	for _, tt := range tests {
		if got := expensive(tt.total, tt.qty); got != tt.want {
			t.Errorf("expensive(%d, %d) = %v, want %v", tt.total, tt.qty, got, tt.want)
		}
	}
}
