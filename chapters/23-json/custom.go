package main

import (
	"encoding/json/v2"
	"fmt"
	"slices"
	"time"
	"uuid"
)

// Spice 是辣度。程式裡用數字比大小，JSON 裡用中文字。
type Spice int

const (
	NotSpicy Spice = iota
	Mild
	Medium
	Hot
)

var spiceNames = [...]string{"不辣", "小辣", "中辣", "大辣"}

// MarshalText 實作 encoding.TextMarshaler：JSON 會把回傳的文字編成字串。
func (s Spice) MarshalText() ([]byte, error) {
	if s < NotSpicy || s > Hot {
		return nil, fmt.Errorf("辣度 %d 超出範圍", int(s))
	}
	return []byte(spiceNames[s]), nil
}

// UnmarshalText 實作 encoding.TextUnmarshaler。要改到呼叫者的值，所以 receiver 是指標。
func (s *Spice) UnmarshalText(b []byte) error {
	i := slices.Index(spiceNames[:], string(b))
	if i < 0 {
		return fmt.Errorf("不認得的辣度 %q", b)
	}
	*s = Spice(i)
	return nil
}

// Bowl 是一碗湯麵。
type Bowl struct {
	Dish  string `json:"dish"`
	Spice Spice  `json:"spice"`
}

// textDemo Spice 在 JSON 裡是中文，讀回來是可以比大小的數字；當 map 的 key 也可以。
func textDemo() {
	b, err := json.Marshal([]Bowl{{"麻辣鴨血", Hot}, {"清燉牛肉麵", NotSpicy}})
	fmt.Println(string(b), err)

	var bowl Bowl
	err = json.Unmarshal([]byte(`{"dish":"酸辣湯","spice":"中辣"}`), &bowl)
	fmt.Println(bowl.Dish, int(bowl.Spice), bowl.Spice > Mild, err)

	err = json.Unmarshal([]byte(`{"dish":"酸辣湯","spice":"超辣"}`), &bowl)
	fmt.Println(explain(err))

	sold := map[Spice]int{Hot: 3, Mild: 5}
	b, err = json.Marshal(sold, json.Deterministic(true))
	fmt.Println(string(b), err)
}

// Ticket 是叫號單。uuid.UUID 實作了 TextMarshaler，在 JSON 裡就是標準格式的字串。
type Ticket struct {
	ID    uuid.UUID `json:"id"`
	Table int       `json:"table"`
}

// uuidDemo Go 1.27 的 uuid 套件。範例用固定的 UUID，輸出才會每次都一樣。
func uuidDemo() {
	t := Ticket{ID: uuid.MustParse("019a9b1c-3f00-7d2e-8a41-5c6d7e8f9012"), Table: 7}
	b, err := json.Marshal(t)
	fmt.Println(string(b), err)

	var back Ticket
	err = json.Unmarshal(b, &back)
	fmt.Println(back.ID == t.ID, err)

	err = json.Unmarshal([]byte(`{"id":"A-1024","table":7}`), &back)
	fmt.Println(explain(err))

	fmt.Println(len(uuid.New().String())) // 每次都不一樣，這裡只印長度
}

// Recipe 是食譜，Prep 是準備時間。
type Recipe struct {
	Dish string        `json:"dish"`
	Prep time.Duration `json:"prep"`
}

// durationDemo v2 的 time.Duration 沒有預設的 JSON 格式；呼叫的一方可以用 WithMarshalers 指定。
func durationDemo() {
	r := Recipe{Dish: "紅燒牛肉", Prep: 3*time.Hour + 30*time.Minute}
	_, err := json.Marshal(r)
	fmt.Println(explain(err))

	asText := json.MarshalFunc(func(d time.Duration) ([]byte, error) {
		return json.Marshal(d.String())
	})
	b, err := json.Marshal(r, json.WithMarshalers(asText))
	fmt.Println(string(b), err)
}
