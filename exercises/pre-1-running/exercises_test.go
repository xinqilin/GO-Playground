package running

import (
	"slices"
	"testing"
)

func TestMotto(t *testing.T) {
	if got, want := Motto(), "我會寫 Go 了"; got != want {
		t.Errorf("Motto() = %q, want %q", got, want)
	}
}

func TestReceipt(t *testing.T) {
	want := "排骨便當 85\n雞腿便當 90\n合計 175"
	if got := Receipt(); got != want {
		t.Errorf("Receipt() = %q\nwant %q", got, want)
	}
}

func TestTotal(t *testing.T) {
	if got := Total(); got != 360 {
		t.Errorf("Total() = %d, want 360", got)
	}
}

func TestAverage(t *testing.T) {
	if got := Average(); got != 82 {
		t.Errorf("Average() = %d, want 82", got)
	}
}

func TestSteps(t *testing.T) {
	want := []string{"燒開水", "放茶葉", "倒熱水", "等三分鐘"}
	if got := Steps(); !slices.Equal(got, want) {
		t.Errorf("Steps() = %v\nwant %v", got, want)
	}
}
