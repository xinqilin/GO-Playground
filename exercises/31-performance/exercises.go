//go:build !solution

package rush

import "bytes"

// Summarize 整理同一個 benchmark 跑好幾次的結果（每次的 ns/op）：
// median 是中位數（個數是偶數時取中間兩個的平均），
// spread 是離中位數最遠的那一個差了中位數的百分之幾。
// 例如 {100, 104, 98} 的中位數是 100，spread 是 4。
// samples 是空的時候兩個都回傳 0；中位數是 0 時算不出百分比，spread 回傳 0。
// 不能改動 samples 本身的順序。
func Summarize(samples []float64) (median, spread float64) {
	// TODO: 實作
	return 0, 0
}

// AppendLine 把一行收據「排骨飯 x2 180」加上換行，接在 dst 後面並回傳，180 是 qty*price。
// dst 的容量夠時不能配置任何記憶體：不要用 fmt，數字用 strconv.AppendInt。
func AppendLine(dst []byte, dish string, qty, price int) []byte {
	// TODO: 實作
	return dst
}

// ParseOrder 解析「12:05,排骨飯,2,90」這種一行訂單紀錄。必須剛好四個欄位：
// 時間是 5 個字元的 HH:MM（兩個數字、冒號、兩個數字），菜名不能是空的，
// 數量是 1 到 99，單價大於 0。格式不對時回傳 ErrBadLine。
// 解析正確的一行時不能配置任何記憶體：Time 和 Dish 直接用 line 的子字串。
func ParseOrder(line string) (Order, error) {
	// TODO: 實作
	return Order{}, nil
}

// Get 從 pool 拿一個緩衝區，裡面一定是空的（Len() 是 0）。pool 裡沒有時建一個新的。
func (p *BufferPool) Get() *bytes.Buffer {
	// TODO: 實作
	return nil
}

// Put 把用完的緩衝區還回去。MaxCap 不是 0 而且 b.Cap() 超過 MaxCap 時直接丟掉，
// 免得偶爾一張特別大的收據，讓 pool 一直佔著大塊記憶體。
func (p *BufferPool) Put(b *bytes.Buffer) {
	// TODO: 實作
}

// Receipt 把 items 排成收據，每一項一行（格式和 AppendLine 一樣），最後一行是「合計 230」：
//
//	排骨飯 x2 180
//	滷肉飯 x1 50
//	合計 230
//
// 整個函式最多只能配置一次記憶體：先算出總長度，一次配置好。
func Receipt(items []Item) string {
	// TODO: 實作
	return ""
}

// FirstLine 回傳 data 的第一行（不含換行）。data 通常是剛讀進來、好幾 MB 的整個檔案，
// 呼叫的人只想留下第一行，之後 data 會被丟掉或重複使用。
//
// 找錯：留下第一行，結果整個檔案都跟著留在記憶體裡；而且 data 被重複使用之後，第一行也跟著變了。
func FirstLine(data []byte) []byte {
	if before, _, ok := bytes.Cut(data, []byte{'\n'}); ok {
		return before
	}
	return data
}
