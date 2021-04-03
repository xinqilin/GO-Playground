# GO-practice
- ctrl-shift-p -> go -> update -> select all 

### 基本語法/操作
```go

package main
import "fmt"
func main()  {
	fmt.Println("Hello world")
}

```

```go build 檔名```
```go run 檔名   也可以 ./檔名```

- int int8 int16 int32 int64
- float32 float64
- string
- bool 
- rune
- byte 
- complex64 complex128
- pointer 的init value 是nil

```go

	var x int
	x = 10
	fmt.Println(x)
	var isTrue bool = true
	fmt.Println(isTrue)
	var c rune = 'b'
	fmt.Println(c)
	var s string = "Test Test Test"
	fmt.Println(s)
	var f float64 = 3.1415
	fmt.Println(f)

    // 也可以跟js一樣
    fmt.Println('a',"字串",2,3.14)

    //another assign 簡短格式
    //須給init value
    //不能提供数据类型。
    //只能用在函数内部
     i:=1

     //匿名變數 _ 
     //不站內存 

```

### input
```go
fmt.Scanln(&圈圈叉叉變數)  

// & 是address
var x int
fmt.Println(&x) // print address
```

### if

```go
if true {
    //判斷條件不用 括號 XD  其他跟java一樣
    
}

```

### for

```go 
for int  {
    //判斷條件不用 括號 XD  其他跟java一樣
    
}

```

### pointer

```go
	var x int = 5
	fmt.Println(x) //5

	fmt.Println(&x) //0xc0000160c8

	var pointer *int = &x
	fmt.Println(pointer) // 0xc0000160c8

	fmt.Println(*pointer) //5

```

### struct

example:

```go

type Point struct{
    x int
    y int
}

// 實體化 (java 的constructer)
func main(){
    var p1 Point=Point{3,4};
    var p2 Point=Point{x:1,y:2}
}

```

### array

```go

var arr [4]int
var arr2 [4]int{0,1,2,3}

fmt.Println(len(arr)) // get array length

a := [2]int{1, 2}
b := [...]int{1, 2}
c := [2]int{1, 3}
fmt.Println(a == b, a == c, b == c) // "true false false"
d := [3]int{1, 2}
fmt.Println(a == d) // compiler error, length not equals


//get data
for k, v := range a {
    fmt.Println(k, v)
}

```




