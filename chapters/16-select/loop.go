package main

import (
	"fmt"
	"slices"
)

// from 在另一個 goroutine 裡依序送出 items，送完就 close。
func from(items ...string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for _, it := range items {
			out <- it
		}
	}()
	return out
}

// combine 用 for-select 同時收內用和外帶兩個出餐口，兩邊都 close 才結束。
// 某一邊 close 之後把它設成 nil：nil channel 永遠收不到東西，select 就不會再選那個 case。
func combine(dineIn, takeout <-chan string) []string {
	var all []string
	for dineIn != nil || takeout != nil {
		select {
		case o, ok := <-dineIn:
			if !ok {
				dineIn = nil
				continue
			}
			all = append(all, "內用 "+o)
		case o, ok := <-takeout:
			if !ok {
				takeout = nil
				continue
			}
			all = append(all, "外帶 "+o)
		}
	}
	return all
}

// loopDemo：兩邊交錯的順序每次都可能不同，所以排序之後再印。
func loopDemo() {
	all := combine(from("排骨飯", "雞腿飯"), from("魚排飯"))
	slices.Sort(all)
	for _, o := range all {
		fmt.Println(o)
	}
}
