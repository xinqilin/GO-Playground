package parking

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"
)

// Payments 是收費的對外介面。正式環境接金流公司的 API，測試時換成假的。
type Payments interface {
	Charge(ctx context.Context, plate string, amount int) error
}

// ChargeTimeout 是等金流回應的上限。
const ChargeTimeout = 5 * time.Second

var (
	ErrInside    = errors.New("這台車已經在場內")
	ErrNotInside = errors.New("場內沒有這台車")
)

// Record 是一筆出場紀錄。
type Record struct {
	Plate   string
	In, Out time.Time
	Fee     int
}

// Gate 是停車場的柵欄機。現在幾點由 now 決定，測試時可以換成假的時鐘。
type Gate struct {
	now func() time.Time
	pay Payments

	mu     sync.Mutex
	inside map[string]time.Time
	done   []Record
}

// NewGate 建立柵欄機。正式環境傳 time.Now。
func NewGate(now func() time.Time, pay Payments) *Gate {
	return &Gate{now: now, pay: pay, inside: make(map[string]time.Time)}
}

// Enter 記下進場時間。
func (g *Gate) Enter(plate string) error {
	p, err := NormalizePlate(plate)
	if err != nil {
		return err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.inside[p]; ok {
		return fmt.Errorf("%s：%w", p, ErrInside)
	}
	g.inside[p] = g.now()
	return nil
}

// Exit 計費、收費，成功才放行。收費失敗時車留在場內，可以再試一次。
func (g *Gate) Exit(ctx context.Context, plate string) (Record, error) {
	p, err := NormalizePlate(plate)
	if err != nil {
		return Record{}, err
	}
	g.mu.Lock()
	in, ok := g.inside[p]
	delete(g.inside, p) // 先拿出來：同一台車同時出場兩次，也只會收一次錢
	g.mu.Unlock()
	if !ok {
		return Record{}, fmt.Errorf("%s：%w", p, ErrNotInside)
	}

	rec := Record{Plate: p, In: in, Out: g.now()}
	rec.Fee, err = Fee(rec.In, rec.Out)
	if err == nil && rec.Fee > 0 {
		ctx, cancel := context.WithTimeout(ctx, ChargeTimeout)
		err = g.pay.Charge(ctx, p, rec.Fee)
		cancel()
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	if err != nil {
		if _, again := g.inside[p]; !again { // 等金流的時候同一個車牌又進場了，以新的進場時間為準
			g.inside[p] = in // 放回場內
		}
		return Record{}, fmt.Errorf("%s 出場失敗：%w", p, err)
	}
	g.done = append(g.done, rec)
	return rec, nil
}

// Records 回傳出場紀錄的複本，依出場的先後。
func (g *Gate) Records() []Record {
	g.mu.Lock()
	defer g.mu.Unlock()
	return slices.Clone(g.done)
}
