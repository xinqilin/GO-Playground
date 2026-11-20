package fuzzbug

import "testing"

// 這些案例都通過了，看起來沒問題。
func TestNormalizePlate(t *testing.T) {
	for in, want := range map[string]string{"abc 1234": "ABC-1234", "AB-1234": "AB-1234", "abc1234": "ABC-1234"} {
		if got, err := NormalizePlate(in); got != want || err != nil {
			t.Errorf("NormalizePlate(%q) = %q, %v；want %q", in, got, err, want)
		}
	}
}

func FuzzNormalizePlate(f *testing.F) {
	f.Add("abc 1234")
	f.Add("AB-1234")
	f.Fuzz(func(t *testing.T, s string) {
		p, err := NormalizePlate(s)
		if err != nil {
			return
		}
		if again, err := NormalizePlate(p); err != nil || again != p {
			t.Errorf("NormalizePlate(%q) = %q，再整理一次變成 %q, %v", s, p, again, err)
		}
	})
}
