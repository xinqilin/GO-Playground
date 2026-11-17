package variables

import (
	"math"
	"testing"
)

func TestSplitBill(t *testing.T) {
	tests := []struct {
		total, people, each, rest int
	}{
		{1000, 3, 333, 1},
		{900, 3, 300, 0},
		{85, 2, 42, 1},
		{5, 10, 0, 5},
		{100, 0, 0, 0},
		{100, -2, 0, 0},
	}
	for _, tt := range tests {
		each, rest := SplitBill(tt.total, tt.people)
		if each != tt.each || rest != tt.rest {
			t.Errorf("SplitBill(%d, %d) = %d, %d; want %d, %d",
				tt.total, tt.people, each, rest, tt.each, tt.rest)
		}
	}
}

func TestTables(t *testing.T) {
	tests := []struct {
		guests, perTable, want int
	}{
		{10, 4, 3},
		{8, 4, 2},
		{1, 4, 1},
		{4, 1, 4},
		{0, 4, 0},
		{10, 0, 0},
		{-3, 4, 0},
	}
	for _, tt := range tests {
		if got := Tables(tt.guests, tt.perTable); got != tt.want {
			t.Errorf("Tables(%d, %d) = %d, want %d", tt.guests, tt.perTable, got, tt.want)
		}
	}
}

func TestAverage(t *testing.T) {
	tests := []struct {
		a, b, c int
		want    float64
	}{
		{80, 85, 91, 256.0 / 3},
		{1, 2, 2, 5.0 / 3},
		{90, 90, 90, 90},
		{-1, 0, 0, -1.0 / 3},
	}
	for _, tt := range tests {
		if got := Average(tt.a, tt.b, tt.c); math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("Average(%d, %d, %d) = %v, want %v", tt.a, tt.b, tt.c, got, tt.want)
		}
	}
}

func TestFormatSize(t *testing.T) {
	if KB != 1024 || MB != 1024*1024 || GB != 1024*1024*1024 {
		t.Fatalf("KB, MB, GB = %d, %d, %d; want 1024, 1048576, 1073741824", KB, MB, GB)
	}
	tests := []struct {
		bytes int64
		want  string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1023, "1023 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048575, "1024.0 KB"},
		{3145728, "3.0 MB"},
		{1073741824, "1.0 GB"},
		{5 << 40, "5120.0 GB"},
		{-1, "0 B"},
	}
	for _, tt := range tests {
		if got := FormatSize(tt.bytes); got != tt.want {
			t.Errorf("FormatSize(%d) = %q, want %q", tt.bytes, got, tt.want)
		}
	}
}

func TestPerm(t *testing.T) {
	cashier := Checkout
	manager := Checkout | Refund | EditPrice | CloseDay

	if got := Grant(cashier, Refund); got != Checkout|Refund {
		t.Errorf("Grant(Checkout, Refund) = %04b, want %04b", got, Checkout|Refund)
	}
	if got := Grant(cashier, Checkout); got != Checkout {
		t.Errorf("Grant(Checkout, Checkout) = %04b, want %04b", got, Checkout)
	}
	if got := Revoke(manager, EditPrice|CloseDay); got != Checkout|Refund {
		t.Errorf("Revoke(全部, EditPrice|CloseDay) = %04b, want %04b", got, Checkout|Refund)
	}
	if got := Revoke(cashier, Refund); got != Checkout {
		t.Errorf("Revoke(Checkout, Refund) = %04b, want %04b", got, Checkout)
	}

	hasTests := []struct {
		p, q Perm
		want bool
	}{
		{manager, Refund, true},
		{manager, Refund | CloseDay, true},
		{cashier, Checkout, true},
		{cashier, Refund, false},
		{cashier, Checkout | Refund, false},
		{cashier, 0, true},
	}
	for _, tt := range hasTests {
		if got := Has(tt.p, tt.q); got != tt.want {
			t.Errorf("Has(%04b, %04b) = %v, want %v", tt.p, tt.q, got, tt.want)
		}
	}
}

func TestAddStock(t *testing.T) {
	tests := []struct {
		stock, n uint8
		want     uint8
		ok       bool
	}{
		{10, 5, 15, true},
		{0, 0, 0, true},
		{250, 5, 255, true},
		{250, 6, 255, false},
		{255, 1, 255, false},
		{200, 200, 255, false},
	}
	for _, tt := range tests {
		got, ok := AddStock(tt.stock, tt.n)
		if got != tt.want || ok != tt.ok {
			t.Errorf("AddStock(%d, %d) = %d, %v; want %d, %v", tt.stock, tt.n, got, ok, tt.want, tt.ok)
		}
	}
}

func TestTicket(t *testing.T) {
	tests := []struct {
		current, delta, want int
	}{
		{1, 1, 2},
		{998, 1, 999},
		{999, 1, 1},
		{999, 0, 999},
		{1, -1, 999},
		{10, -9, 1},
		{10, -10, 999},
		{500, 999, 500},
		{500, 1000, 501},
		{5, -2000, 3},
	}
	for _, tt := range tests {
		if got := Ticket(tt.current, tt.delta); got != tt.want {
			t.Errorf("Ticket(%d, %d) = %d, want %d", tt.current, tt.delta, got, tt.want)
		}
	}
}
