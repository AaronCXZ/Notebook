###### byte类型

官方对byte的定义：

```go
type byte = uint8
```

可以看到byte就是uint8的别名，它是用来区分字节值和8位无符号整数值。其实可以白byte当作一个ASCII码的一个字符。

示例：

```go
var ch byte = 65
var ch byte = '\x41'
var ch byte = 'A'
```

###### []byte类型

[]byte就是一个byte类型的切片，切片本质也是一个结构体，定义如下：

```go
type slice struct {
    array 	unsafe.Pointer
    len 	int
    cap 	int
}
```

array代表底层数组的指针，len代表切片的长度，cap代表切片的容量。

###### string类型

官方定义如下:

```go
type string string
```

string是一个8位字节的集合，通常但不一定代表UTF-8编码的的文本。string可以为空，但不能为nil，string的值是不可能改变的。

string类型本质也是一个结构体，定义如下：

```go
type stringStruct struct {
    str 	unsafe.Pointer
    len 	int
}
```

str指针指向的是某个数组的首地址，len代表数组的长度

```go
func gostringnocapy(str *byte) string {
    ss := stringStruct{str: unsafe.Pointer(str),len: findnull(str)}
    s := *(*string)(unsafe.Pointer(&ss))
    return s
}
```

入参是一个byte类型的指针，从这里可以看出string类型底层是一个byte类型的数组。

###### []byte和string有什么区别

因为在go语言中string类型被设计为不可变的，这样的好处就是：在并发场景下，我们可以在不加锁的控制下，多次使用同已字符串，在保证高效共享的情况下而不担心安全问题。所以将byte类型的数组的基础上封装了string类型。

string类型虽然是不能更改的，但是可以被替换，因为stringStruct中的str指针是可以改变的，只是指针的内容是不可以改变的。示例：

```go
func main()  {
 str := "song"
 fmt.Printf("%p\n",[]byte(str))
 str = "asong"
 fmt.Printf("%p\n",[]byte(str))
}
// 运行结果
0xc00001a090
0xc00001a098
```

可以看出，指针指向的位置发生了变化，也就是说每一次更改字符串，就需要重新分配一次内存，之前分配的空间会被gc回收。

###### string和[]byte标准转换

go语言提供课标准方式对string和[]byte进行转换，示例：

```go
func main() {
    str := "song"
    by := []byte(str)
    
    str1 := string(by)
    fmt.Println(str1)
}
```

+ string类型转换到[]byte类型

对上面的代码执行：go tool compile -N -l -S ./string.go，可以看到调用的值runtime.stringtoslicebyte：

```go
const tmpStringBufSize = 32

type tmpBuf [tmpStringBufSize]byte

func stringtoslicebyte(buf *tmpBuf, s string) []byte {
	var b []byte
	if buf != nil && len(s) <= len(buf) {
		*buf = tmpBuf{}
		b = buf[:len(s)]
	} else {
		b = rawbyteslice(len(s))
	}
	copy(b, s)
	return b
}

// rawbyteslice allocates a new byte slice. The byte slice is not zeroed.
func rawbyteslice(size int) (b []byte) {
	cap := roundupsize(uintptr(size))
	p := mallocgc(cap, nil, false)
	if cap != uintptr(size) {
		memclrNoHeapPointers(add(p, uintptr(size)), cap-uintptr(size))
	}

	*(*slice)(unsafe.Pointer(&b)) = slice{p, size, int(cap)}
	return
}
```

先预定义一个长度为32的数组tmpBuf，字符串长度超过了这个数组的长度，说明[]byte不够用了，需要重新分配一块内存。最后通过调用copy方法实现string到[]byte的拷贝，核心思路就是：将string的底层数组从头复制n个到[]byte对应的底层数组中去。

+ []byte类型转换到string类型

[]byte类型转换到string类型就是调用runtime.slicebytetostring：

```go
// Buf is a fixed-size buffer for the result,
// it is not nil if the result does not escape.
func slicebytetostring(buf *tmpBuf, b []byte) (str string) {
	l := len(b)
	if l == 0 {
		// Turns out to be a relatively common case.
		// Consider that you want to parse out data between parens in "foo()bar",
		// you find the indices and convert the subslice to string.
		return ""
	}
	if raceenabled {
		racereadrangepc(unsafe.Pointer(&b[0]),
			uintptr(l),
			getcallerpc(),
			funcPC(slicebytetostring))
	}
	if msanenabled {
		msanread(unsafe.Pointer(&b[0]), uintptr(l))
	}
	if l == 1 {
		stringStructOf(&str).str = unsafe.Pointer(&staticbytes[b[0]])
		stringStructOf(&str).len = 1
		return
	}

	var p unsafe.Pointer
	if buf != nil && len(b) <= len(buf) {
		p = unsafe.Pointer(buf)
	} else {
		p = mallocgc(uintptr(len(b)), nil, false)
	}
	stringStructOf(&str).str = p
	stringStructOf(&str).len = len(b)
	memmove(p, (*(*slice)(unsafe.Pointer(&b))).array, uintptr(len(b)))
	return
}
```

根据[]byte的长度来决定是否重新分配内存，最后通过memmove可以拷贝数组到字符串。

###### string和[]byte强转换

标准的转换方法都会发生内存拷贝，所以为了减少内存拷贝和内存申请我们可以使用强转换的方式对两者进行转换，标准库中的实现方式：

```go
func slicebytetostringtmp(ptr *byte, n int) (str string) {
 stringStructOf(&str).str = unsafe.Pointer(ptr)
 stringStructOf(&str).len = n
 return
}

func stringtoslicebytetmp(s string) []byte {
    str := (*stringStruct)(unsafe.Pointer(&s))
    ret := slice{array: unsafe.Pointer(str.str), len: str.len, cap: str.len}
    return *(*[]byte)(unsafe.Pointer(&ret))
}
```

可以看出只要使用的就是unsafe.Pointer进行指针替换，因为string和slice的结构字段是相似的。唯一不同的是cap字段，array和str是一致的，len是一致的，所以他们的内存布局上是对齐的，这样就可以直接通过unsafe.Pointer进行指针替换。

###### 两种转换的取舍

标准转换是安全的，强转换是不安全的。强转换的方式性能肯定要比标准转换要好，根据实际场景来取舍。

