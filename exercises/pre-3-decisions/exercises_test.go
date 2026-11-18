package decisions

import "testing"

func TestShippingFee(t *testing.T) {
	tests := []struct {
		total  int
		member bool
		want   int
	}{
		{100, true, 0},
		{100, false, 60},
		{498, false, 60},
		{499, false, 0},
		{1000, false, 0},
	}
	for _, tt := range tests {
		if got := ShippingFee(tt.total, tt.member); got != tt.want {
			t.Errorf("ShippingFee(%d, %v) = %d, want %d", tt.total, tt.member, got, tt.want)
		}
	}
}

func TestSumOdd(t *testing.T) {
	tests := []struct{ n, want int }{
		{7, 16},
		{10, 25},
		{1, 1},
		{0, 0},
		{-3, 0},
	}
	for _, tt := range tests {
		if got := SumOdd(tt.n); got != tt.want {
			t.Errorf("SumOdd(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestCountDigits(t *testing.T) {
	tests := []struct{ n, want int }{
		{2026, 4},
		{0, 1},
		{9, 1},
		{10, 2},
		{100000, 6},
	}
	for _, tt := range tests {
		if got := CountDigits(tt.n); got != tt.want {
			t.Errorf("CountDigits(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

func TestParkingFee(t *testing.T) {
	tests := []struct{ minutes, want int }{
		{0, 0},
		{30, 0},
		{31, 20},
		{60, 20},
		{61, 40},
		{150, 80},
		{300, 180},
		{330, 200},
		{600, 200},
	}
	for _, tt := range tests {
		if got := ParkingFee(tt.minutes); got != tt.want {
			t.Errorf("ParkingFee(%d) = %d, want %d", tt.minutes, got, tt.want)
		}
	}
}

func TestRightTriangle(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{3, "  *\n **\n***"},
		{1, "*"},
		{0, ""},
	}
	for _, tt := range tests {
		if got := RightTriangle(tt.n); got != tt.want {
			t.Errorf("RightTriangle(%d) = %q\nwant %q", tt.n, got, tt.want)
		}
	}
}

func TestSumRange(t *testing.T) {
	tests := []struct{ from, to, want int }{
		{1, 10, 55},
		{5, 5, 5},
		{3, 6, 18},
		{3, 2, 0},
	}
	for _, tt := range tests {
		if got := SumRange(tt.from, tt.to); got != tt.want {
			t.Errorf("SumRange(%d, %d) = %d, want %d", tt.from, tt.to, got, tt.want)
		}
	}
}

func TestDaysToSave(t *testing.T) {
	tests := []struct{ target, daily, want int }{
		{100, 30, 4},
		{100, 50, 2},
		{50, 50, 1},
		{10, 100, 1},
	}
	for _, tt := range tests {
		if got := DaysToSave(tt.target, tt.daily); got != tt.want {
			t.Errorf("DaysToSave(%d, %d) = %d, want %d", tt.target, tt.daily, got, tt.want)
		}
	}
}
