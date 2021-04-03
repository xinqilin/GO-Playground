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

- int
- float64
- string
- bool 
- rune

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

    //another assign
     i:=1

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




