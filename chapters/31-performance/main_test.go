package main

import (
	"maps"
	"testing"
)

func Example_revenueDemo() {
	revenueDemo()
	// Output:
	// 排骨飯 36000 元
	// 滷肉飯 19950 元
	// 素食便當 32080 元
	// 雞腿飯 40000 元
	// 魚排飯 37905 元
	// 兩種結果相同： true <nil> <nil>
	// 訂單格式不對："10:05,排骨飯,兩個,90"
}

func Example_allocDemo() {
	allocDemo()
	// Output:
	// fmt.Sprintf：每次配置 1 次
	// 字串相加：每次配置 2 次
	// strings.Builder：每次配置 2 次
	// append 到舊的 buf：每次配置 0 次
}

func Example_escapeDemo() {
	escapeDemo()
	// Output:
	// 請 15 號到櫃台
	// 20 3
}

func Example_preallocDemo() {
	preallocDemo()
	// Output:
	// 沒有預先配置：len 1000 cap 1023
	// 預先配置：len 1000 cap 1000
}

func Example_poolDemo() {
	poolDemo()
	// Output:
	// 12:01 排骨飯 x2
	// 12:03 滷肉飯 x1
	// 合計 230
	// 相同： true
}

// TestParsersAgree：改寫成比較快的版本之前，先確認兩個版本對每一行的結果都一樣。
func TestParsersAgree(t *testing.T) {
	lines := []string{
		"10:05,排骨飯,2,90", "13:59,素食便當,1,80",
		"", "10:05,排骨飯,2", "10:05,,2,90", "10:05,排骨飯,兩個,90", "1005,排骨飯,2,90",
	}
	for _, line := range lines {
		a, errA := parseRegexp(line)
		b, errB := parseCut(line)
		if a != b || (errA == nil) != (errB == nil) {
			t.Errorf("%q：regexp 得到 %+v, %v；Cut 得到 %+v, %v", line, a, errA, b, errB)
		}
	}
	log := sampleLog(5000)
	slow, _ := revenue(log, parseRegexp)
	fast, _ := revenue(log, parseCut)
	if !maps.Equal(slow, fast) {
		t.Errorf("營業額不同：%v vs %v", slow, fast)
	}
}

// BenchmarkRevenue 用子 benchmark 比較兩種解析方法，輸入完全相同。
func BenchmarkRevenue(b *testing.B) {
	log := sampleLog(1000) // 準備資料不算在時間裡：b.Loop 會在第一次呼叫時重設計時
	for _, bc := range []struct {
		name  string
		parse func(string) (Order, error)
	}{
		{"regexp", parseRegexp},
		{"cut", parseCut},
	} {
		b.Run(bc.name, func(b *testing.B) {
			for b.Loop() {
				revenue(log, bc.parse)
			}
		})
	}
}

func BenchmarkCollectNames(b *testing.B) {
	orders := make([]Order, 1000)
	b.Run("append", func(b *testing.B) {
		for b.Loop() {
			collectNames(orders)
		}
	})
	b.Run("prealloc", func(b *testing.B) {
		for b.Loop() {
			collectNamesPrealloc(orders)
		}
	})
}

func BenchmarkReceipt(b *testing.B) {
	orders := make([]Order, 20)
	for i := range orders {
		orders[i] = Order{"12:00", "排骨飯", 2, 90}
	}
	b.Run("new", func(b *testing.B) {
		for b.Loop() {
			receipt(orders)
		}
	})
	b.Run("pool", func(b *testing.B) {
		for b.Loop() {
			receiptPooled(orders)
		}
	})
}
