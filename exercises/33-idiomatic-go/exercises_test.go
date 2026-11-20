package idiom

import (
	"slices"
	"testing"
)

func TestStutters(t *testing.T) {
	tests := []struct {
		pkg, name string
		want      bool
	}{
		{"ticket", "TicketMachine", true},
		{"ticket", "Machine", false},
		{"ticket", "Ticket", false},
		{"ticket", "Tickets", false},
		{"http", "HTTPServer", true},
		{"http", "Server", false},
		{"user", "UserID", true},
		{"io", "IOReader", true},
		{"strings", "Builder", false},
		{"bento", "Be", false},
	}
	for _, tt := range tests {
		if got := Stutters(tt.pkg, tt.name); got != tt.want {
			t.Errorf("Stutters(%q, %q) = %v，想要 %v", tt.pkg, tt.name, got, tt.want)
		}
	}
}

func TestDocProblems(t *testing.T) {
	tests := []struct {
		name, doc string
		want      []string
	}{
		{"Fee", "Fee 計算停車費。", nil},
		{"Fee", "Fee computes the parking fee.", nil},
		{"Fee", "  Fee 計算停車費。\n", nil},
		{"Fee", "", []string{ProblemMissing}},
		{"Fee", "   ", []string{ProblemMissing}},
		{"Fee", "計算停車費。", []string{ProblemName}},
		{"Fee", "FeeTable 是費率表。", []string{ProblemName}},
		{"Fee", "Fee 計算停車費", []string{ProblemPeriod}},
		{"Fee", "停車費", []string{ProblemName, ProblemPeriod}},
	}
	for _, tt := range tests {
		if got := DocProblems(tt.name, tt.doc); !slices.Equal(got, tt.want) {
			t.Errorf("DocProblems(%q, %q) = %q，想要 %q", tt.name, tt.doc, got, tt.want)
		}
	}
}

func TestTally(t *testing.T) {
	var tally Tally
	if got := tally.Count("排骨飯"); got != 0 {
		t.Errorf("零值的 Count = %d，想要 0", got)
	}
	if got := tally.Top(3); len(got) != 0 {
		t.Errorf("零值的 Top(3) = %q，想要空的", got)
	}
	tally.Add("排骨飯", 2)
	tally.Add("滷肉飯", 3)
	tally.Add("雞腿飯", 1)
	tally.Add("排骨飯", 1)
	tally.Add("魚排飯", 3)
	if got := tally.Count("排骨飯"); got != 3 {
		t.Errorf("Count(排骨飯) = %d，想要 3", got)
	}
	tests := []struct {
		n    int
		want []string
	}{
		{1, []string{"排骨飯"}},
		{3, []string{"排骨飯", "滷肉飯", "魚排飯"}},
		{10, []string{"排骨飯", "滷肉飯", "魚排飯", "雞腿飯"}},
		{0, []string{}},
	}
	for _, tt := range tests {
		if got := tally.Top(tt.n); !slices.Equal(got, tt.want) {
			t.Errorf("Top(%d) = %q，想要 %q", tt.n, got, tt.want)
		}
	}
}

func TestOptions(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
		want []string
	}{
		{"預設", nil, []string{"A001", "A002", "A003"}},
		{"前綴", []Option{WithPrefix("外帶")}, []string{"外帶001", "外帶002", "外帶003"}},
		{"起始與位數", []Option{WithStart(998), WithWidth(3)}, []string{"A998", "A999", "A1000"}},
		{"全部", []Option{WithPrefix("B"), WithWidth(2), WithStart(7)}, []string{"B07", "B08", "B09"}},
		{"後面的選項蓋過前面的", []Option{WithPrefix("B"), WithPrefix("C")}, []string{"C001", "C002", "C003"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(tt.opts...)
			var got []string
			for range len(tt.want) {
				got = append(got, m.Next())
			}
			if !slices.Equal(got, tt.want) {
				t.Errorf("號碼 = %q，想要 %q", got, tt.want)
			}
		})
	}
	a, b := New(), New()
	a.Next()
	if got := b.Next(); got != "A001" {
		t.Errorf("兩台叫號機不能共用號碼：第二台的第一號 = %q", got)
	}
}

func TestMixedCaps(t *testing.T) {
	tests := []struct {
		snake    string
		exported bool
		want     string
	}{
		{"order_count", true, "OrderCount"},
		{"order_count", false, "orderCount"},
		{"user_id", true, "UserID"},
		{"user_id", false, "userID"},
		{"api_base_url", true, "APIBaseURL"},
		{"api_base_url", false, "apiBaseURL"},
		{"id", true, "ID"},
		{"id", false, "id"},
		{"line2_total", true, "Line2Total"},
		{"http__server_", true, "HTTPServer"},
		{"_private_key", false, "privateKey"},
		{"identity", true, "Identity"},
		{"", true, ""},
	}
	for _, tt := range tests {
		if got := MixedCaps(tt.snake, tt.exported); got != tt.want {
			t.Errorf("MixedCaps(%q, %v) = %q，想要 %q", tt.snake, tt.exported, got, tt.want)
		}
	}
}

func TestQueueCopies(t *testing.T) {
	nums := make([]int, 3, 10)
	copy(nums, []int{11, 12, 13})
	q := QueueOf(nums)

	nums[0] = 99
	if got := q.Waiting(); !slices.Equal(got, []int{11, 12, 13}) {
		t.Fatalf("改了傳進去的 slice 之後，Waiting() = %v，想要 [11 12 13]", got)
	}

	got := q.Waiting()
	got[1] = 99
	_ = append(got, 98)
	q.Join(14)
	if got := q.Waiting(); !slices.Equal(got, []int{11, 12, 13, 14}) {
		t.Fatalf("改了 Waiting 回傳的 slice 之後，Waiting() = %v，想要 [11 12 13 14]", got)
	}
	if got := nums[:4]; got[3] == 14 {
		t.Errorf("Join 寫進了呼叫端的陣列：nums[:4] = %v", got)
	}

	var empty Queue
	empty.Join(1)
	if got := empty.Waiting(); !slices.Equal(got, []int{1}) {
		t.Errorf("零值的 Queue：Waiting() = %v，想要 [1]", got)
	}
}
