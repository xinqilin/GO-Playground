package parking

import (
	"errors"
	"testing"
)

func TestNormalizePlate(t *testing.T) {
	tests := map[string]struct {
		in, want string
		err      error
	}{
		"已經是標準寫法": {"ABC-1234", "ABC-1234", nil},
		"小寫加空白":   {" abc 1234 ", "ABC-1234", nil},
		"兩個字母":    {"ab1234", "AB-1234", nil},
		"數字在前":    {"1234-ab", "1234-AB", nil},
		"空字串":     {"", "", ErrBadPlate},
		"數字太少":    {"ABC-123", "", ErrBadPlate},
		"字母太多":    {"ABCD-1234", "", ErrBadPlate},
		"全形字母":    {"ＡＢＣ-1234", "", ErrBadPlate},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := NormalizePlate(tt.in)
			if got != tt.want || !errors.Is(err, tt.err) {
				t.Errorf("NormalizePlate(%q) = %q, %v；want %q, %v", tt.in, got, err, tt.want, tt.err)
			}
		})
	}
}

// FuzzNormalizePlate：任何輸入都不能 panic；整理成功的結果，再整理一次要一模一樣。
func FuzzNormalizePlate(f *testing.F) {
	for _, s := range []string{"ABC-1234", " abc 1234", "1234ab", ""} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		p, err := NormalizePlate(s)
		if err != nil {
			return
		}
		again, err := NormalizePlate(p)
		if err != nil || again != p {
			t.Errorf("NormalizePlate(%q) = %q，再整理一次變成 %q, %v", s, p, again, err)
		}
	})
}

func BenchmarkNormalizePlate(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		NormalizePlate(" abc 1234")
	}
}
