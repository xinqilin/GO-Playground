package hello

import "testing"

func TestGreeting(t *testing.T) {
	tests := []struct {
		name, want string
	}{
		{"小明", "哈囉，小明！"},
		{"Go", "哈囉，Go！"},
		{"", "哈囉，陌生人！"},
	}
	for _, tt := range tests {
		if got := Greeting(tt.name); got != tt.want {
			t.Errorf("Greeting(%q) = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestBentoTotal(t *testing.T) {
	tests := []struct {
		qty, want int
	}{
		{-1, 0}, {0, 0},
		{1, 85}, {3, 255}, {9, 765},
		{10, 800}, {12, 960},
	}
	for _, tt := range tests {
		if got := BentoTotal(tt.qty); got != tt.want {
			t.Errorf("BentoTotal(%d) = %d, want %d", tt.qty, got, tt.want)
		}
	}
}

func TestClock(t *testing.T) {
	tests := []struct {
		seconds int
		want    string
	}{
		{0, "00:00:00"},
		{-5, "00:00:00"},
		{59, "00:00:59"},
		{60, "00:01:00"},
		{3725, "01:02:05"},
		{86399, "23:59:59"},
		{360000, "100:00:00"},
	}
	for _, tt := range tests {
		if got := Clock(tt.seconds); got != tt.want {
			t.Errorf("Clock(%d) = %q, want %q", tt.seconds, got, tt.want)
		}
	}
}

func TestBox(t *testing.T) {
	tests := []struct {
		text, want string
	}{
		{"Go", "+----+\n| Go |\n+----+"},
		{"hello", "+-------+\n| hello |\n+-------+"},
		{"", "+--+\n|  |\n+--+"},
	}
	for _, tt := range tests {
		if got := Box(tt.text); got != tt.want {
			t.Errorf("Box(%q) =\n%s\nwant\n%s", tt.text, got, tt.want)
		}
	}
}

func TestGoVersion(t *testing.T) {
	tests := []struct {
		line         string
		major, minor int
		ok           bool
	}{
		{"go 1.27", 1, 27, true},
		{"go 1.27.4", 1, 27, true},
		{"go 1.21.0", 1, 21, true},
		{"go 1", 0, 0, false},
		{"go 1.x", 0, 0, false},
		{"toolchain go1.27.4", 0, 0, false},
		{"golang 1.27", 0, 0, false},
		{"", 0, 0, false},
	}
	for _, tt := range tests {
		major, minor, ok := GoVersion(tt.line)
		if major != tt.major || minor != tt.minor || ok != tt.ok {
			t.Errorf("GoVersion(%q) = %d, %d, %v; want %d, %d, %v",
				tt.line, major, minor, ok, tt.major, tt.minor, tt.ok)
		}
	}
}
