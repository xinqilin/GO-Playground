package main

import (
	"context"
	"fmt"
	"time"
)

// treeDemo：總店底下有東區和西區，東區底下有一張訂單。
// 取消會往下傳給所有衍生出來的 ctx，但不會往上、也不會傳給旁邊的兄弟。
func treeDemo() {
	shop, closeShop := context.WithCancel(context.Background())
	defer closeShop()
	east, closeEast := context.WithCancel(shop)
	defer closeEast()
	west, closeWest := context.WithCancel(shop)
	defer closeWest()
	order, closeOrder := context.WithTimeout(east, time.Hour)
	defer closeOrder()

	closeEast()
	<-order.Done()
	fmt.Println("總店", shop.Err())
	fmt.Println("東區", east.Err())
	fmt.Println("東區的訂單", order.Err())
	fmt.Println("西區", west.Err())

	closeShop()
	<-west.Done()
	fmt.Println("總店關門之後，西區", west.Err())
}
