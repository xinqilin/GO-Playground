package orders

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestParsePickup(t *testing.T) {
	ok := []struct {
		in   string
		want int
	}{
		{"18:30", 1110}, {"11:00", 660}, {"20:00", 1200}, {"12:05", 725},
	}
	for _, tt := range ok {
		got, err := ParsePickup(tt.in)
		if got != tt.want || err != nil {
			t.Errorf("ParsePickup(%q) = %d, %v, want %d, nil", tt.in, got, err, tt.want)
		}
	}

	bad := []struct {
		in   string
		want error
	}{
		{"1830", ErrTimeFormat}, {"18:3x", ErrTimeFormat}, {"24:00", ErrTimeFormat},
		{"18:60", ErrTimeFormat}, {"-1:30", ErrTimeFormat}, {"", ErrTimeFormat},
		{"20:01", ErrClosed}, {"08:30", ErrClosed}, {"10:59", ErrClosed},
	}
	for _, tt := range bad {
		_, err := ParsePickup(tt.in)
		if !errors.Is(err, tt.want) {
			t.Errorf("ParsePickup(%q) error = %v, want 包著 %v", tt.in, err, tt.want)
		}
	}

	if _, err := ParsePickup("1830"); err == nil || err.Error() != `取餐時間 "1830"：格式要是 HH:MM` {
		t.Errorf("ParsePickup(\"1830\") 的訊息是 %v", err)
	}
	if _, err := ParsePickup("21:00"); err == nil || err.Error() != "取餐時間 21:00：不在營業時間" {
		t.Errorf("ParsePickup(\"21:00\") 的訊息是 %v", err)
	}
}

func TestBuy(t *testing.T) {
	limits := map[string]int{"排骨飯": 2}
	if err := Buy("排骨飯", 2, limits); err != nil {
		t.Errorf("Buy(排骨飯, 2) = %v (%T), want nil", err, err)
	}
	if err := Buy("滷蛋", 99, limits); err != nil {
		t.Errorf("Buy(滷蛋, 99) = %v (%T), want nil（不限購）", err, err)
	}
	err := Buy("排骨飯", 3, limits)
	le, ok := errors.AsType[*LimitError](err)
	if !ok || le.Item != "排骨飯" || le.Max != 2 {
		t.Fatalf("Buy(排骨飯, 3) = %v, want *LimitError{排骨飯, 2}", err)
	}
	if got, want := err.Error(), "排骨飯 每人限購 2 份"; got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestCheckCart(t *testing.T) {
	limits := map[string]int{"排骨飯": 2, "雞腿飯": 1}
	if err := CheckCart(map[string]int{"排骨飯": 2, "滷蛋": 5}, limits); err != nil {
		t.Errorf("沒有超過限購時 CheckCart = %v (%T), want nil", err, err)
	}

	cart := map[string]int{"雞腿飯": 2, "排骨飯": 3, "滷蛋": 5}
	err := CheckCart(cart, limits)
	if err == nil {
		t.Fatal("CheckCart 沒有回報超過限購")
	}
	if got, want := err.Error(), "排骨飯 每人限購 2 份\n雞腿飯 每人限購 1 份"; got != want {
		t.Errorf("訊息 = %q, want %q", got, want)
	}
	if le, ok := errors.AsType[*LimitError](err); !ok || le.Item != "排骨飯" {
		t.Errorf("errors.AsType 找不到第一個 *LimitError：%v", err)
	}
}

func TestRetry(t *testing.T) {
	errTimeout := errors.New("連線逾時")

	calls := 0
	err := Retry(5, func() error {
		calls++
		if calls < 3 {
			return errTimeout
		}
		return nil
	})
	if err != nil || calls != 3 {
		t.Errorf("第 3 次成功：Retry = %v，呼叫 %d 次, want nil，3 次", err, calls)
	}

	calls = 0
	err = Retry(3, func() error { calls++; return errTimeout })
	if calls != 3 || !errors.Is(err, errTimeout) || err.Error() != "試了 3 次：連線逾時" {
		t.Errorf("一直失敗：Retry = %v，呼叫 %d 次, want 「試了 3 次：連線逾時」，3 次", err, calls)
	}

	calls = 0
	err = Retry(5, func() error { calls++; return fmt.Errorf("付款：%w", ErrPermanent) })
	if calls != 1 || !errors.Is(err, ErrPermanent) {
		t.Errorf("不能重試的錯誤：Retry = %v，呼叫 %d 次, want 包著 ErrPermanent，1 次", err, calls)
	}

	calls = 0
	Retry(0, func() error { calls++; return errTimeout })
	if calls != 1 {
		t.Errorf("times 是 0 時呼叫了 %d 次, want 1", calls)
	}
}

// stepError 是測試用的自訂錯誤型別，有 Unwrap() error。
type stepError struct {
	step string
	err  error
}

func (e *stepError) Error() string { return e.step + "失敗：" + e.err.Error() }
func (e *stepError) Unwrap() error { return e.err }

func TestCauses(t *testing.T) {
	a := errors.New("a")
	d := errors.New("d")
	tests := []struct {
		name string
		err  error
		want []string
	}{
		{"nil", nil, nil},
		{"單一", a, []string{"a"}},
		{"包一層", fmt.Errorf("b：%w", a), []string{"b：a", "a"}},
		{"Join", errors.Join(a, fmt.Errorf("c：%w", d)), []string{"a\nc：d", "a", "c：d", "d"}},
		{"兩個 %w", fmt.Errorf("%w；%w", a, d), []string{"a；d", "a", "d"}},
		{"自訂型別", &stepError{"付款", fmt.Errorf("x：%w", d)}, []string{"付款失敗：x：d", "x：d", "d"}},
	}
	for _, tt := range tests {
		if got := Causes(tt.err); !slices.Equal(got, tt.want) {
			t.Errorf("%s: Causes = %q, want %q", tt.name, got, tt.want)
		}
	}
}

func TestIsSoldOut(t *testing.T) {
	tests := []struct {
		err  error
		want bool
	}{
		{ErrSoldOut, true},
		{fmt.Errorf("訂單 17：%w", ErrSoldOut), true},
		{errors.Join(errors.New("x"), ErrSoldOut), true},
		{errors.New("賣完了"), false},
		{errors.New("沒貨"), false},
		{nil, false},
	}
	for _, tt := range tests {
		if got := IsSoldOut(tt.err); got != tt.want {
			t.Errorf("IsSoldOut(%v) = %v, want %v", tt.err, got, tt.want)
		}
	}
}

// fakeFile 是測試用的 io.WriteCloser。
type fakeFile struct {
	writeErr error
	closeErr error
	closed   int
	buf      strings.Builder
}

func (f *fakeFile) Write(p []byte) (int, error) {
	if f.writeErr != nil {
		return 0, f.writeErr
	}
	return f.buf.Write(p)
}

func (f *fakeFile) Close() error {
	f.closed++
	return f.closeErr
}

func TestSaveLines(t *testing.T) {
	var ok fakeFile
	if err := SaveLines(&ok, []string{"排骨飯 x2", "滷蛋 x1"}); err != nil || ok.closed != 1 || ok.buf.String() != "排骨飯 x2\n滷蛋 x1\n" {
		t.Errorf("正常情況：err = %v，關閉 %d 次，內容 %q", err, ok.closed, ok.buf.String())
	}

	errDisk := errors.New("磁碟滿了")
	closeFail := fakeFile{closeErr: errDisk}
	if err := SaveLines(&closeFail, []string{"x"}); !errors.Is(err, errDisk) {
		t.Errorf("關閉失敗時 SaveLines = %v, want 包著 %v", err, errDisk)
	}

	errWrite := errors.New("寫入失敗")
	writeFail := fakeFile{writeErr: errWrite}
	if err := SaveLines(&writeFail, []string{"x"}); !errors.Is(err, errWrite) || writeFail.closed != 1 {
		t.Errorf("寫入失敗時 SaveLines = %v，關閉 %d 次, want 包著 %v，1 次", err, writeFail.closed, errWrite)
	}

	both := fakeFile{writeErr: errWrite, closeErr: errDisk}
	if err := SaveLines(&both, []string{"x"}); !errors.Is(err, errWrite) || !errors.Is(err, errDisk) {
		t.Errorf("兩個都失敗時 SaveLines = %v, want 兩個錯誤都在", err)
	}
}
