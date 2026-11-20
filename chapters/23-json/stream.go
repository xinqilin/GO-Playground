package main

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strings"
)

// Sale 是一筆營業額。
type Sale struct {
	Shop   string `json:"shop"`
	Amount int    `json:"amount"`
}

// sumJSONLines 讀 JSON Lines（每一行一筆 JSON）：一次只解一筆，檔案再大也只佔一筆的記憶體。
func sumJSONLines(r io.Reader) (map[string]int, error) {
	dec := jsontext.NewDecoder(r)
	totals := map[string]int{}
	for n := 1; ; n++ {
		var s Sale
		err := json.UnmarshalDecode(dec, &s)
		if errors.Is(err, io.EOF) {
			return totals, nil
		}
		if err != nil {
			return nil, fmt.Errorf("第 %d 筆：%w", n, err)
		}
		totals[s.Shop] += s.Amount
	}
}

func linesDemo() {
	const input = `{"shop":"台北店","amount":1280}
{"shop":"台中店","amount":960}

{"shop":"台北店","amount":450}
`
	totals, err := sumJSONLines(strings.NewReader(input))
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, shop := range slices.Sorted(maps.Keys(totals)) {
		fmt.Println(shop, totals[shop])
	}

	_, err = sumJSONLines(strings.NewReader(input + `{"shop":"高雄店","amount":"三百"}`))
	fmt.Println(explain(err))
}

// countExpensive 從一個很大的 JSON 陣列裡一筆一筆讀，數價格超過 limit 的品項，不會把整個陣列放進記憶體。
func countExpensive(r io.Reader, limit int) (int, error) {
	dec := jsontext.NewDecoder(r)
	tok, err := dec.ReadToken() // 先讀掉開頭的 [
	if err != nil {
		return 0, err
	}
	if tok.Kind() != '[' {
		return 0, fmt.Errorf("要的是陣列，開頭卻是 %v", tok.Kind())
	}
	n := 0
	for dec.PeekKind() != ']' {
		var item MenuItem
		if err := json.UnmarshalDecode(dec, &item); err != nil {
			return 0, err
		}
		if item.Price > limit {
			n++
		}
	}
	return n, nil
}

func arrayDemo() {
	menu := `[{"name":"滷肉飯","price":35},{"name":"雞腿飯","price":100},{"name":"牛肉麵","price":180}]`
	fmt.Println(countExpensive(strings.NewReader(menu), 90))
	fmt.Println(countExpensive(strings.NewReader(`{"name":"滷肉飯"}`), 90))
}

// encodeDemo jsontext.Encoder 每寫完一個值就換行，剛好是 JSON Lines 的格式。
func encodeDemo() {
	enc := jsontext.NewEncoder(os.Stdout)
	for _, s := range []Sale{{"台北店", 1730}, {"台中店", 960}} {
		if err := json.MarshalEncode(enc, s); err != nil {
			fmt.Println(err)
			return
		}
	}
}
