package failing

import (
	"os"
	"strings"
	"testing"
)

// t.Errorf 記下失敗，繼續檢查下一個案例。
func TestMemberPrice(t *testing.T) {
	tests := []struct{ price, want int }{
		{100, 90},
		{95, 86}, // 85.5 四捨五入
		{55, 50}, // 49.5 四捨五入
	}
	for _, tt := range tests {
		if got := MemberPrice(tt.price); got != tt.want {
			t.Errorf("MemberPrice(%d) = %d，want %d", tt.price, got, tt.want)
		}
	}
}

// t.Fatalf 記下失敗，並且立刻結束這個測試：後面都用到 data，讀不到就沒有必要再檢查。
func TestMenuPrices(t *testing.T) {
	data, err := os.ReadFile("testdata/menu.txt")
	if err != nil {
		t.Fatalf("讀不到菜單：%v", err)
	}
	if !strings.Contains(string(data), "排骨飯") {
		t.Error("菜單裡沒有排骨飯")
	}
}
