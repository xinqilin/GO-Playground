package values

import "testing"

func TestElapsed(t *testing.T) {
	tests := []struct{ startH, startM, endH, endM, want int }{
		{9, 50, 11, 5, 75},
		{8, 0, 8, 45, 45},
		{12, 30, 12, 30, 0},
		{0, 0, 23, 59, 1439},
	}
	for _, tt := range tests {
		if got := Elapsed(tt.startH, tt.startM, tt.endH, tt.endM); got != tt.want {
			t.Errorf("Elapsed(%d, %d, %d, %d) = %d, want %d", tt.startH, tt.startM, tt.endH, tt.endM, got, tt.want)
		}
	}
}

func TestWeeksAndDays(t *testing.T) {
	tests := []struct{ days, weeks, rest int }{
		{17, 2, 3},
		{14, 2, 0},
		{6, 0, 6},
		{0, 0, 0},
	}
	for _, tt := range tests {
		weeks, rest := WeeksAndDays(tt.days)
		if weeks != tt.weeks || rest != tt.rest {
			t.Errorf("WeeksAndDays(%d) = %d, %d, want %d, %d", tt.days, weeks, rest, tt.weeks, tt.rest)
		}
	}
}

func TestIsHalfFare(t *testing.T) {
	tests := []struct {
		age  int
		want bool
	}{
		{5, true},
		{11, true},
		{12, false},
		{40, false},
		{64, false},
		{65, true},
		{80, true},
	}
	for _, tt := range tests {
		if got := IsHalfFare(tt.age); got != tt.want {
			t.Errorf("IsHalfFare(%d) = %v, want %v", tt.age, got, tt.want)
		}
	}
}

func TestPredict(t *testing.T) {
	a, b, c := Predict()
	if a != 17 || b != 5 || c != -8 {
		t.Errorf("Predict() = %d, %d, %d；再追蹤一次，特別注意每一行用的是當時最新的值", a, b, c)
	}
}

func TestCoins(t *testing.T) {
	tests := []struct{ amount, c50, c10, c5, c1 int }{
		{78, 1, 2, 1, 3},
		{100, 2, 0, 0, 0},
		{49, 0, 4, 1, 4},
		{6, 0, 0, 1, 1},
		{0, 0, 0, 0, 0},
	}
	for _, tt := range tests {
		c50, c10, c5, c1 := Coins(tt.amount)
		if c50 != tt.c50 || c10 != tt.c10 || c5 != tt.c5 || c1 != tt.c1 {
			t.Errorf("Coins(%d) = %d, %d, %d, %d, want %d, %d, %d, %d",
				tt.amount, c50, c10, c5, c1, tt.c50, tt.c10, tt.c5, tt.c1)
		}
	}
}

func TestRotate(t *testing.T) {
	a, b, c := Rotate("小明", "小華", "小美")
	if a != "小華" || b != "小美" || c != "小明" {
		t.Errorf("Rotate(小明, 小華, 小美) = %s, %s, %s, want 小華, 小美, 小明", a, b, c)
	}
}

func TestAttendance(t *testing.T) {
	tests := []struct {
		part, total int
		want        float64
	}{
		{18, 20, 90},
		{3, 8, 37.5},
		{1, 4, 25},
		{20, 20, 100},
	}
	for _, tt := range tests {
		if got := Attendance(tt.part, tt.total); got != tt.want {
			t.Errorf("Attendance(%d, %d) = %v, want %v", tt.part, tt.total, got, tt.want)
		}
	}
}
