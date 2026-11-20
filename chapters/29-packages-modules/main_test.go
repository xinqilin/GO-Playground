package main

func Example_exportDemo() {
	exportDemo()
	// Output:
	// <nil> <nil>
	// 價格要大於 0
	// 90 true [排骨飯 滷肉飯]
	// 排骨飯 x2 180
	// 滷肉飯 x1 50
	// 合計 230
}

func Example_aliasDemo() {
	aliasDemo()
	// Output:
	// 備註：<b>不要辣</b>
	// 備註：&lt;b&gt;不要辣&lt;/b&gt;
}

func Example_initOrderDemo() {
	initOrderDemo()
	// Output:
	// steps 的 init → price → qty → total → main 的 init
	// total = 180
}

func Example_goVersionDemo() {
	goVersionDemo()
	// Output:
	// go1.21.0 比 go1.22rc1：-1
	// go1.22rc1 比 go1.22.0：-1
	// go1.9 比 go1.10：-1
	// go1.27 false true
}
