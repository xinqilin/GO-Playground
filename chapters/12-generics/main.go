// 第 12 章：泛型。
//
// 泛型是一個模具：先把形狀刻好，要用的時候再決定倒進去的是 int、string 還是別的型別。
// 執行：go run ./chapters/12-generics
package main

import (
	"cmp"
	"fmt"
	"strconv"
)

// Number 是自訂的型別限制。~int 代表「底層型別是 int 的都可以」，所以自訂的 Celsius 也能用。
type Number interface {
	~int | ~int64 | ~float64
}

// Sum 對任何 Number 都能加總。
func Sum[T Number](nums []T) T {
	var total T // T 的零值
	for _, n := range nums {
		total += n
	}
	return total
}

// Map 把 []T 轉成 []U；T、U 由呼叫時的參數推斷出來。
func Map[T, U any](s []T, f func(T) U) []U {
	out := make([]U, 0, len(s))
	for _, v := range s {
		out = append(out, f(v))
	}
	return out
}

// MaxOf 回傳 slice 裡最大的值。cmp.Ordered 包含所有能用 < 比較的型別。
func MaxOf[T cmp.Ordered](s []T) (T, bool) {
	var best T
	if len(s) == 0 {
		return best, false
	}
	best = s[0]
	for _, v := range s[1:] {
		if v > best {
			best = v
		}
	}
	return best, true
}

// Stack 是泛型型別：同一份程式碼，可以是 Stack[int]，也可以是 Stack[string]。
type Stack[T any] struct {
	items []T
}

// Push 放一個值到最上面。
func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

// Pop 拿走最上面的值；空的時候回傳零值和 false。
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	v := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return v, true
}

// Convert 是 generic method（Go 1.27 起支援）：方法自己也能宣告型別參數 U。
func (s *Stack[T]) Convert[U any](f func(T) U) *Stack[U] {
	return &Stack[U]{items: Map(s.items, f)}
}

// Celsius 的底層型別是 float64，所以符合 Number。
type Celsius float64

func main() {
	fmt.Println("== 泛型函式 ==")
	fmt.Println(Sum([]int{1, 2, 3}))
	fmt.Println(Sum([]Celsius{20.5, 21.5}))
	fmt.Println(Map([]int{1, 2, 3}, strconv.Itoa))
	fmt.Println(MaxOf([]string{"pear", "apple", "zoo"}))
	fmt.Println(MaxOf([]int{}))

	fmt.Println("== 泛型型別 ==")
	var s Stack[int]
	s.Push(1)
	s.Push(2)
	top, _ := s.Pop()
	fmt.Println(top, s.items)

	fmt.Println("== generic method ==")
	s.Push(3)
	labels := s.Convert(func(n int) string { return "#" + strconv.Itoa(n) })
	fmt.Println(labels.items)
}
