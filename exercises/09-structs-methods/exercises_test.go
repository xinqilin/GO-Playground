package shop

import (
	"errors"
	"slices"
	"testing"
)

func TestDoubleBooked(t *testing.T) {
	a5 := Seat{"A", 5}
	b3 := Seat{"B", 3}
	tests := []struct {
		name   string
		orders []Seat
		want   Seat
		wantOK bool
	}{
		{"沒有訂位", nil, Seat{}, false},
		{"沒有重複", []Seat{a5, b3, {"A", 6}}, Seat{}, false},
		{"同排不同號", []Seat{a5, {"A", 3}, {"B", 5}}, Seat{}, false},
		{"一個重複", []Seat{a5, b3, a5}, a5, true},
		{"先被訂第二次的", []Seat{a5, b3, b3, a5}, b3, true},
	}
	for _, tt := range tests {
		got, ok := DoubleBooked(tt.orders)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("%s: DoubleBooked(%v) = %v, %v, want %v, %v", tt.name, tt.orders, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestCart(t *testing.T) {
	var c Cart
	c.Add("排骨飯", 85, 2)
	c.Add("紅茶", 25, 1)
	c.Add("排骨飯", 85, 1)
	want := []Line{{"排骨飯", 85, 3}, {"紅茶", 25, 1}}
	if !slices.Equal(c.Lines, want) {
		t.Errorf("Lines = %v, want %v", c.Lines, want)
	}
	if got := c.Total(); got != 280 {
		t.Errorf("Total() = %d, want 280", got)
	}
}

func TestCard(t *testing.T) {
	var c Card
	if err := c.TopUp(150); !errors.Is(err, ErrBadAmount) {
		t.Errorf("TopUp(150) = %v, want ErrBadAmount", err)
	}
	if err := c.TopUp(0); !errors.Is(err, ErrBadAmount) {
		t.Errorf("TopUp(0) = %v, want ErrBadAmount", err)
	}
	if err := c.TopUp(500); err != nil || c.Balance != 500 {
		t.Errorf("TopUp(500) = %v, Balance = %d, want nil, 500", err, c.Balance)
	}
	if err := c.Pay(20); err != nil || c.Balance != 480 {
		t.Errorf("Pay(20) = %v, Balance = %d, want nil, 480", err, c.Balance)
	}
	if err := c.Pay(1000); !errors.Is(err, ErrInsufficient) || c.Balance != 480 {
		t.Errorf("Pay(1000) = %v, Balance = %d, want ErrInsufficient, 480", err, c.Balance)
	}
	if err := c.Pay(480); err != nil || c.Balance != 0 {
		t.Errorf("Pay(480) = %v, Balance = %d, want nil, 0", err, c.Balance)
	}
}

func TestInventory(t *testing.T) {
	var inv Inventory
	if got := inv.Count("滷蛋"); got != 0 {
		t.Errorf("空庫存 Count = %d, want 0", got)
	}
	inv.Add("滷蛋", 10)
	inv.Add("滷蛋", 5)
	if got := inv.Count("滷蛋"); got != 15 {
		t.Errorf("Count(滷蛋) = %d, want 15", got)
	}
	if err := inv.Remove("滷蛋", 4); err != nil || inv.Count("滷蛋") != 11 {
		t.Errorf("Remove(滷蛋, 4) = %v, Count = %d, want nil, 11", err, inv.Count("滷蛋"))
	}
	if err := inv.Remove("滷蛋", 20); !errors.Is(err, ErrInsufficient) || inv.Count("滷蛋") != 11 {
		t.Errorf("Remove(滷蛋, 20) = %v, Count = %d, want ErrInsufficient, 11", err, inv.Count("滷蛋"))
	}
	var empty Inventory
	if err := empty.Remove("豆干", 1); !errors.Is(err, ErrInsufficient) {
		t.Errorf("從空庫存 Remove = %v, want ErrInsufficient", err)
	}
}

func TestLoggedCart(t *testing.T) {
	var c LoggedCart
	c.Add("排骨飯", 85, 2)
	c.Add("排骨飯", 85, 1)
	if want := []Line{{"排骨飯", 85, 3}}; !slices.Equal(c.Lines, want) {
		t.Errorf("Lines = %v, want %v", c.Lines, want)
	}
	if want := []string{"加入 排骨飯 x2", "加入 排骨飯 x1"}; !slices.Equal(c.Log, want) {
		t.Errorf("Log = %q, want %q", c.Log, want)
	}
	if got := c.Total(); got != 255 {
		t.Errorf("Total() = %d, want 255（Total 是從 Cart 提升上來的）", got)
	}
}

func TestCounter(t *testing.T) {
	var c Counter
	c.Inc()
	c.Inc()
	c.Inc()
	if got := c.Value(); got != 3 {
		t.Errorf("Inc 三次之後 Value() = %d, want 3", got)
	}
}

func TestRaiseAll(t *testing.T) {
	staff := []Employee{{"小明", 30000}, {"阿美", 32500}}
	RaiseAll(staff, 10)
	want := []Employee{{"小明", 33000}, {"阿美", 35750}}
	if !slices.Equal(staff, want) {
		t.Errorf("RaiseAll 之後 = %v, want %v", staff, want)
	}
}
