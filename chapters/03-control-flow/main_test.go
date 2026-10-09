package main

import "testing"

func Example() {
	main()
	// Output:
	// == if 可以先做一件事 ==
	// 42 的兩倍是 84
	// abc 不是數字
	// == for 的四種寫法 ==
	// [0 1 2]
	// 243
	// [0 1 2]
	// 5
	// == 九九乘法表（前三列） ==
	// 1×1=1  2×1=2  3×1=3
	// 1×2=2  2×2=4  3×2=6
	// 1×3=3  2×3=6  3×3=9
	// == 有標籤的 break ==
	// 找到 3 6
	// == switch ==
	// 95 優
	// 72 及格
	// 40 不及格
	// 週末 平日
	// == fallthrough ==
	// ans 是 a
	// 進到 b 的區塊
}

func TestGrade(t *testing.T) {
	tests := []struct {
		score int
		want  string
	}{
		{100, "優"},
		{90, "優"},
		{89, "及格"},
		{60, "及格"},
		{59, "不及格"},
	}
	for _, tt := range tests {
		if got := grade(tt.score); got != tt.want {
			t.Errorf("grade(%d) = %q, want %q", tt.score, got, tt.want)
		}
	}
}
