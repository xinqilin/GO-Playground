package main

import "fmt"

// finished 是一份做好的餐，seq 是它在點餐順序裡的位置。
type finished struct {
	seq  int
	dish string
}

// inOrder 讓先做好的餐先在一旁等著，照 seq 的順序一份一份送出；in 收完時 close。
// pending 裡放已經做好、但前面還有餐沒好的；next 是下一個該送出的 seq。
func inOrder(in <-chan finished) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		pending := make(map[int]string)
		next := 0
		for f := range in {
			pending[f.seq] = f.dish
			for {
				dish, ok := pending[next]
				if !ok {
					break
				}
				out <- dish
				delete(pending, next)
				next++
			}
		}
	}()
	return out
}

// orderDemo：餐做好的順序是 2、0、3、1，送出去的順序還是 0、1、2、3。
func orderDemo() {
	in := make(chan finished, 4)
	in <- finished{2, "魚排飯"}
	in <- finished{0, "排骨飯"}
	in <- finished{3, "滷肉飯"}
	in <- finished{1, "雞腿飯"}
	close(in)

	for dish := range inOrder(in) {
		fmt.Println(dish)
	}
}
