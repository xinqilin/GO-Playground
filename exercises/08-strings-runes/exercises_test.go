package text

import (
	"maps"
	"testing"
)

func TestMaskName(t *testing.T) {
	tests := []struct{ name, want string }{
		{"王小明", "王*明"},
		{"王明", "王*"},
		{"歐陽小明", "歐**明"},
		{"李", "李"},
		{"", ""},
		{"Amy", "A*y"},
	}
	for _, tt := range tests {
		if got := MaskName(tt.name); got != tt.want {
			t.Errorf("MaskName(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestNormalizePhone(t *testing.T) {
	tests := []struct {
		in   string
		want string
		ok   bool
	}{
		{"0912-345-678", "0912345678", true},
		{"0912 345 678", "0912345678", true},
		{"+886 912 345 678", "0912345678", true},
		{"+886-912-345-678", "0912345678", true},
		{"0912345678", "0912345678", true},
		{"02-1234-5678", "", false},
		{"0912-345-67a", "", false},
		{"091234567", "", false},
		{"09123456789", "", false},
		{"", "", false},
	}
	for _, tt := range tests {
		got, ok := NormalizePhone(tt.in)
		if got != tt.want || ok != tt.ok {
			t.Errorf("NormalizePhone(%q) = %q, %v; want %q, %v", tt.in, got, ok, tt.want, tt.ok)
		}
	}
}

func TestParseOrder(t *testing.T) {
	got, err := ParseOrder("排骨x2, 雞腿x1,排骨x1 ")
	if want := map[string]int{"排骨": 3, "雞腿": 1}; err != nil || !maps.Equal(got, want) {
		t.Errorf("ParseOrder = %v, %v; want %v, nil", got, err, want)
	}
	if got, err := ParseOrder(""); err != nil || got == nil || len(got) != 0 {
		t.Errorf("ParseOrder(\"\") = %v, %v; want 空的 map, nil", got, err)
	}
	for _, bad := range []string{"排骨", "排骨xa", "排骨x0", "排骨x-1", "x2", "排骨x2,,雞腿x1"} {
		if got, err := ParseOrder(bad); err == nil || got != nil {
			t.Errorf("ParseOrder(%q) = %v, %v; want nil 和一個錯誤", bad, got, err)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		s    string
		n    int
		want string
	}{
		{"雞腿便當加滷蛋", 4, "雞腿便當…"},
		{"排骨飯", 5, "排骨飯"},
		{"排骨飯", 3, "排骨飯"},
		{"排骨飯", 2, "排骨…"},
		{"Coffee", 3, "Cof…"},
		{"", 2, ""},
	}
	for _, tt := range tests {
		if got := Truncate(tt.s, tt.n); got != tt.want {
			t.Errorf("Truncate(%q, %d) = %q, want %q", tt.s, tt.n, got, tt.want)
		}
	}
}

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"上海自來水來自海上", true},
		{"人人", true},
		{"level", true},
		{"Go語", false},
		{"abca", false},
		{"a", true},
		{"", true},
	}
	for _, tt := range tests {
		if got := IsPalindrome(tt.s); got != tt.want {
			t.Errorf("IsPalindrome(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestFormatMenu(t *testing.T) {
	got := FormatMenu([]MenuItem{
		{"排骨飯", 85},
		{"Coffee", 120},
		{"雞腿便當", 100},
	})
	want := "排骨飯     85\n" +
		"Coffee    120\n" +
		"雞腿便當  100"
	if got != want {
		t.Errorf("FormatMenu =\n%s\nwant\n%s", got, want)
	}
	if got := FormatMenu([]MenuItem{{"湯", 0}}); got != "湯  0" {
		t.Errorf("FormatMenu(一個品項) = %q, want %q", got, "湯  0")
	}
	if got := FormatMenu(nil); got != "" {
		t.Errorf("FormatMenu(nil) = %q, want \"\"", got)
	}
}
