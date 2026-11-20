package main

import (
	"reflect"
	"testing"
)

func Example_parseDemo() {
	parseDemo()
	// Output:
	// 排骨飯 2 份 "少飯"
	// 滷肉飯 1 份 ""
	// 珍珠奶茶 1 份 "微糖, 去冰"
	// 第 2 行：數量 "二"：strconv.Atoi: parsing "二": invalid syntax
}

func Example_numErrorDemo() {
	numErrorDemo()
	// Output:
	// Atoi 二 true false
	// Atoi 99999999999999999999 false true
	// strconv.ParseInt: parsing "300": value out of range
	// true true true
	// 1 true true
	// T true true
	// yes false false
}

func Example_matchDemo() {
	matchDemo()
	// Output:
	// "AB-12345678"：完整比對 true，部分比對 true
	// "AB-123456789"：完整比對 false，部分比對 true
	// "發票 AB-12345678"：完整比對 false，部分比對 true
	// "ab-12345678"：完整比對 false，部分比對 false
}

func Example_submatchDemo() {
	submatchDemo()
	// Output:
	// 08:15 進站 台北車站 車資=""
	// 08:41 出站 市政府 車資="25"
	// 看不懂：2026/11/19 08:41 轉乘 市政府
	// ["NT$85" "NT$120" "NT$1,280"]
}

func Example_replaceDemo() {
	replaceDemo()
	// Output:
	// 訂位日期 2026-11-19，取消期限 2026-11-17
	// "訂位日期 ，取消期限 "
	// 訂位日期 2026年11月19日，取消期限 2026年11月17日
}

func Example_replaceFuncDemo() {
	replaceFuncDemo()
	// Output:
	// 本店1979年開幕，2026年搬到現址
	// NT\$1\.5 true
}

func Example_re2Demo() {
	re2Demo()
	// Output:
	// error parsing regexp: invalid or unsupported Perl syntax: `(?=`
	// error parsing regexp: invalid escape sequence: `\1`
	// error parsing regexp: missing closing ]: `[A-Z`
}

func Example_receiptDemo() {
	receiptDemo()
	// Output:
	// 阿嬤便當 第 7 桌
	// 排骨飯 x2 180 元
	// 滷蛋 x3 45 元
	// 合計 225 元（已付款）
	// 備註：統編 12345678
	// 阿嬤便當 第 3 桌
	// （沒有品項）
	// 合計 0 元
}

func Example_execErrorDemo() {
	execErrorDemo()
	// Output:
	// template: typo:1:2: executing "typo" at <.Shp>: can't evaluate field Shp in type main.Receipt
	// 阿嬤便當 <no value>
	// 寫出了 "阿嬤便當 "
	// template: strict:1:12: executing "strict" at <.phone>: map has no entry for key "phone"
}

func Example_noticeDemo() {
	noticeDemo()
	// Output:
	// 【阿嬤便當】今天正常營業。
	// 營業時間 10:30-14:00、17:00-20:30，外送滿 1,500 元免運。
	// 【阿嬤便當】11/24（一） 公休一天，造成不便請見諒。
	// 營業時間 10:30-14:00、17:00-20:30，外送滿 1,500 元免運。
}

func Example_escapeDemo() {
	escapeDemo()
	// Output:
	// text/template：
	// <p><a href="javascript:alert(1)"><b>阿明</b></a>：好吃<script>steal()</script></p>
	// <button onclick="reply(<b>阿明</b>)">回覆</button>
	// html/template：
	// <p><a href="#ZgotmplZ">&lt;b&gt;阿明&lt;/b&gt;</a>：好吃&lt;script&gt;steal()&lt;/script&gt;</p>
	// <button onclick="reply(&#34;\u003cb\u003e阿明\u003c/b\u003e&#34;)">回覆</button>
}

func Example_trustedDemo() {
	trustedDemo()
	// Output:
	// <p>&lt;em&gt;招牌&lt;/em&gt;</p>
	// <p><em>招牌</em></p>
}

// TestParseOrder 合法的點餐單和各種寫錯的情況。
func TestParseOrder(t *testing.T) {
	got, err := parseOrder("# 外帶\n\n雞腿飯 x3 不要辣 多菜\n紅茶\n")
	want := []OrderLine{{"雞腿飯", 3, "不要辣 多菜"}, {"紅茶", 1, ""}}
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Errorf("parseOrder = %+v（%v），want %+v", got, err, want)
	}
	for _, in := range []string{"排骨飯 x0", "排骨飯 x", "排骨飯 x-1", "排骨飯 x二"} {
		if _, err := parseOrder(in); err == nil {
			t.Errorf("parseOrder(%q) 應該回傳錯誤", in)
		}
	}
	if got, err := parseOrder("# 只有註解\n\n"); err != nil || len(got) != 0 {
		t.Errorf("只有註解：%+v（%v），want 沒有品項", got, err)
	}
}

// TestMoney 每三位一個逗號。
func TestMoney(t *testing.T) {
	for n, want := range map[int]string{0: "0", 999: "999", 1000: "1,000", 1500: "1,500", 1234567: "1,234,567"} {
		if got := money(n); got != want {
			t.Errorf("money(%d) = %q，want %q", n, got, want)
		}
	}
}
