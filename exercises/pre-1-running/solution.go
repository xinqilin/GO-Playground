//go:build solution

package running

func Motto() string {
	return "我會寫 Go 了"
}

func Receipt() string {
	return "排骨便當 85\n雞腿便當 90\n合計 175"
}

func Total() int {
	return 4 * 90
}

// Average 要用括號讓加法先算，再除以 3。
func Average() int {
	return (72 + 85 + 91) / 3
}

// Steps 要先燒開水，才有熱水可以倒。
func Steps() []string {
	return []string{
		"燒開水",
		"放茶葉",
		"倒熱水",
		"等三分鐘",
	}
}
