package rush

import (
	"bytes"
	"errors"
	"math"
	"slices"
	"strings"
	"sync"
	"testing"
)

func TestSummarize(t *testing.T) {
	tests := []struct {
		name           string
		samples        []float64
		median, spread float64
	}{
		{"三個", []float64{104, 100, 98}, 100, 4},
		{"偶數個取平均", []float64{210, 190, 200, 202}, 201, 100.0 * 11 / 201},
		{"都一樣", []float64{50, 50, 50}, 50, 0},
		{"一個", []float64{80}, 80, 0},
		{"沒有", nil, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := slices.Clone(tt.samples)
			median, spread := Summarize(tt.samples)
			if math.Abs(median-tt.median) > 1e-9 || math.Abs(spread-tt.spread) > 1e-9 {
				t.Errorf("Summarize(%v) = %v, %v；want %v, %v", tt.samples, median, spread, tt.median, tt.spread)
			}
			if !slices.Equal(tt.samples, before) {
				t.Errorf("Summarize 改動了 samples：%v → %v", before, tt.samples)
			}
		})
	}
}

func TestAppendLine(t *testing.T) {
	got := AppendLine([]byte("12:00\n"), "排骨飯", 2, 90)
	got = AppendLine(got, "滷肉飯", 1, 50)
	if want := "12:00\n排骨飯 x2 180\n滷肉飯 x1 50\n"; string(got) != want {
		t.Fatalf("AppendLine = %q，want %q", got, want)
	}
	buf := make([]byte, 0, 64)
	if n := testing.AllocsPerRun(100, func() { buf = AppendLine(buf[:0], "雞腿飯", 12, 100) }); n != 0 {
		t.Errorf("容量夠時配置了 %v 次，want 0", n)
	}
	if string(buf) != "雞腿飯 x12 1200\n" {
		t.Errorf("AppendLine = %q", buf)
	}
}

func TestParseOrder(t *testing.T) {
	good := map[string]Order{
		"12:05,排骨飯,2,90":   {"12:05", "排骨飯", 2, 90},
		"09:00,素食便當,99,80": {"09:00", "素食便當", 99, 80},
	}
	for line, want := range good {
		if got, err := ParseOrder(line); got != want || err != nil {
			t.Errorf("ParseOrder(%q) = %+v, %v；want %+v", line, got, err, want)
		}
	}
	bad := []string{
		"", "12:05,排骨飯,2", "12:05,排骨飯,2,90,備註", "12:05,,2,90", "1205,排骨飯,2,90",
		"12-05,排骨飯,2,90", "ab:cd,排骨飯,2,90", "12:05,排骨飯,0,90", "12:05,排骨飯,100,90",
		"12:05,排骨飯,2,0", "12:05,排骨飯,兩個,90",
	}
	for _, line := range bad {
		if got, err := ParseOrder(line); !errors.Is(err, ErrBadLine) {
			t.Errorf("ParseOrder(%q) = %+v, %v；want ErrBadLine", line, got, err)
		}
	}
	if n := testing.AllocsPerRun(100, func() { ParseOrder("12:05,排骨飯,2,90") }); n != 0 {
		t.Errorf("解析一行正確的紀錄配置了 %v 次，want 0", n)
	}
}

func TestBufferPool(t *testing.T) {
	p := &BufferPool{MaxCap: 1024}
	b := p.Get()
	if b == nil {
		t.Fatal("Get() 回傳 nil")
	}
	b.WriteString("上一張收據")
	p.Put(b)
	for range 10 {
		b := p.Get()
		if b == nil || b.Len() != 0 {
			t.Fatalf("Get() 拿到的緩衝區不是空的：%v", b)
		}
		p.Put(b)
	}

	big := bytes.NewBuffer(make([]byte, 0, 1<<20))
	p.Put(big)
	for range 10 {
		if b := p.Get(); b == big || b.Cap() > 1024 {
			t.Fatalf("超過 MaxCap 的緩衝區被放回 pool 了（cap %d）", b.Cap())
		}
	}

	var zero BufferPool // 零值可以直接用，MaxCap 0 表示不限制
	zero.Put(zero.Get())

	// 好幾個 goroutine 同時借還（用 -race 跑）。
	var wg sync.WaitGroup
	for i := range 20 {
		wg.Go(func() {
			for range 50 {
				b := p.Get()
				if b == nil {
					return
				}
				b.WriteString(strings.Repeat("飯", i))
				p.Put(b)
			}
		})
	}
	wg.Wait()
}

func TestReceipt(t *testing.T) {
	items := []Item{{"排骨飯", 2, 90}, {"滷肉飯", 1, 50}, {"雞腿飯", 10, 100}}
	want := "排骨飯 x2 180\n滷肉飯 x1 50\n雞腿飯 x10 1000\n合計 1230\n"
	if got := Receipt(items); got != want {
		t.Fatalf("Receipt =\n%s\nwant\n%s", got, want)
	}
	if got := Receipt(nil); got != "合計 0\n" {
		t.Errorf("Receipt(nil) = %q，want %q", got, "合計 0\n")
	}
	if n := testing.AllocsPerRun(100, func() { Receipt(items) }); n > 1 {
		t.Errorf("Receipt 配置了 %v 次，最多 1 次", n)
	}
}

func TestFirstLine(t *testing.T) {
	data := make([]byte, 0, 1<<20)
	data = append(data, "訂單紀錄 2026-11-20\n12:05,排骨飯,2,90\n"...)
	got := FirstLine(data)
	if string(got) != "訂單紀錄 2026-11-20" {
		t.Fatalf("FirstLine = %q", got)
	}
	if cap(got) > 64 {
		t.Errorf("結果的 cap 是 %d：整個 1 MB 的 data 都被留在記憶體裡", cap(got))
	}
	copy(data, "XXXXXXXX") // data 被拿去讀下一個檔案
	if string(got) != "訂單紀錄 2026-11-20" {
		t.Errorf("data 被重複使用之後，第一行變成 %q", got)
	}
	if got := FirstLine([]byte("只有一行")); string(got) != "只有一行" {
		t.Errorf("沒有換行時 = %q", got)
	}
}
