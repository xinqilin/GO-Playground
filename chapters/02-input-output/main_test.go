package main

import (
	"os"
	"strings"
	"testing"
)

func Example_showVerbs() {
	showVerbs(os.Stdout)
	// Output:
	// {1 2} | {X:1 Y:2} | main.point{X:1, Y:2} | main.point
	// 42 |    42 | 42    | 00042
	// 3.141590 | 3.14 | 3.141590e+02
	// Go | "Go" | 476f
	// true | 世 | U+4E16
}

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr string
	}{
		{name: "三個數字", in: "3\n4\n5\n", want: "共 3 個數字，總和 12，平均 4.00\n"},
		{name: "略過空白行", in: "1\n\n 2 \n", want: "共 2 個數字，總和 3，平均 1.50\n"},
		{name: "沒有輸入", in: "", want: "沒有輸入任何數字\n"},
		{name: "不是整數", in: "1\nabc\n", wantErr: `第 2 行 "abc" 不是整數`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out strings.Builder
			err := run(strings.NewReader(tt.in), &out)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("err = %v, want %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if got := out.String(); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
