package parking

import (
	"errors"
	"testing"
	"time"
)

// taipei 是固定的 UTC+8。用 FixedZone 而不是 LoadLocation，測試就不依賴機器上的時區資料庫。
var taipei = time.FixedZone("UTC+8", 8*60*60)

// at 回傳 2026 年 11 月 day 日 hour:min（UTC+8）。
func at(day, hour, min int) time.Time {
	return time.Date(2026, 11, day, hour, min, 0, 0, taipei)
}

func TestFee(t *testing.T) {
	tests := []struct {
		name    string
		in, out time.Time
		want    int
	}{
		{"剛進場就出來", at(20, 9, 0), at(20, 9, 0), 0},
		{"9 分 59 秒免費", at(20, 9, 0), at(20, 9, 0).Add(Grace - time.Second), 0},
		{"剛好 10 分鐘", at(20, 9, 0), at(20, 9, 10), 20},
		{"剛好 30 分鐘", at(20, 9, 0), at(20, 9, 30), 20},
		{"31 分鐘進位", at(20, 9, 0), at(20, 9, 31), 40},
		{"整個白天到上限", at(20, 8, 0), at(20, 20, 0), DailyCap},
		{"跨夜分兩天算", at(20, 23, 0), at(21, 1, 0), 40 + 40},
		{"跨夜各自進位", at(20, 23, 50), at(21, 0, 20), 20 + 20},
		{"出場在午夜整", at(20, 23, 0), at(21, 0, 0), 40},
		{"停三天", at(20, 12, 0), at(22, 12, 0), 3 * DailyCap},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Fee(tt.in, tt.out)
			if err != nil || got != tt.want {
				t.Errorf("Fee(%s, %s) = %d, %v；want %d",
					tt.in.Format("01-02 15:04"), tt.out.Format("01-02 15:04"), got, err, tt.want)
			}
		})
	}
}

func TestFeeOutBeforeIn(t *testing.T) {
	_, err := Fee(at(20, 9, 0), at(20, 8, 59))
	if !errors.Is(err, ErrOutBeforeIn) {
		t.Errorf("err = %v，want ErrOutBeforeIn", err)
	}
}

// FuzzFee 不檢查某個輸入的答案，而是檢查每個輸入都要成立的性質。
func FuzzFee(f *testing.F) {
	f.Add(int64(0), uint32(30))
	f.Add(int64(23*60+50), uint32(30))
	f.Add(int64(12*60), uint32(3*24*60))
	f.Fuzz(func(t *testing.T, startMin int64, stayMin uint32) {
		in := at(1, 0, 0).Add(time.Duration(startMin%(366*24*60)) * time.Minute)
		out := in.Add(time.Duration(stayMin%(30*24*60)) * time.Minute)
		fee, err := Fee(in, out)
		if err != nil {
			t.Fatalf("Fee(%v, %v)：%v", in, out, err)
		}
		days := int(out.Sub(in)/(24*time.Hour)) + 2 // 停留時間最多跨過幾個日曆日
		if fee < 0 || fee > days*DailyCap {
			t.Errorf("Fee(%v, %v) = %d，超出 0 到 %d", in, out, fee, days*DailyCap)
		}
		if later, _ := Fee(in, out.Add(time.Minute)); later < fee {
			t.Errorf("多停一分鐘反而比較便宜：%d → %d（%v 到 %v）", fee, later, in, out)
		}
	})
}

func BenchmarkFee(b *testing.B) {
	in, out := at(20, 12, 0), at(27, 12, 0) // 停一週
	for b.Loop() {
		Fee(in, out)
	}
}
