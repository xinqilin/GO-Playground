package main

import "fmt"

// Loader 和 Saver 各只有一個方法；Store 把兩個組合起來。
type Loader interface {
	Load(key string) (string, bool)
}

type Saver interface {
	Save(key, value string)
}

type Store interface {
	Loader
	Saver
}

// memStore 把資料放在 map 裡，同時滿足 Loader、Saver 和 Store。
type memStore struct {
	data map[string]string
}

func (m *memStore) Load(key string) (string, bool) {
	v, ok := m.data[key]
	return v, ok
}

func (m *memStore) Save(key, value string) {
	m.data[key] = value
}

// lastOrder 只需要讀，所以參數只要求 Loader。
func lastOrder(l Loader, member string) string {
	if v, ok := l.Load(member); ok {
		return v
	}
	return "第一次來"
}

func embedInterfaceDemo() {
	s := &memStore{data: map[string]string{}}
	var st Store = s
	st.Save("小明", "排骨飯")
	fmt.Println(lastOrder(st, "小明"))
	fmt.Println(lastOrder(s, "阿美"))
}

// remindAll 提醒每一位還沒取餐的客人。它只依賴 Notifier，測試時可以換成假的。
func remindAll(n Notifier, waiting []int) {
	for _, no := range waiting {
		n.Notify(fmt.Sprintf("%d 號請取餐", no))
	}
}

func remindDemo() {
	remindAll(SMS{Phone: "0912-345-678"}, []int{17, 21})
}
