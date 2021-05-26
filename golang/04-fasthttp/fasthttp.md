##### fasthttp的优点

1. net/http的实现是一个连接新建一个goroutine；fasthttp是利用一个worker复用goroutine，减轻runtime调度goroutine的压力。
2. net/http解析的请求数据很多放在map[string]string(http.Header)或map[string][]string(http.Request.From)，有不必要的[]byte到string的转换，是可以规避的。
3. net/http解析HTTP请求每次生成新的*http.Request和http.ResponseWriter；fasthttp解析HTTP数据到*fasthttp.RequestCtx，然后使用sync.Pool复用结构实例，减少了对象的数量。
4. fasthttp会延迟解析HTTP请求中的数据，尤其是Body部分，这样节省了很多不直接操作Body的情况的消耗。

##### web开发的要点

1. 处理tcp通信
2. 对URL/URI进行处理
3. HTTP数据处理

##### fasthttp的性能优化思路

1. 重写了tcp之上进行HTTP握手、连接、通讯的goroutine pool实现。
2. 对http数据基本按传输时的二进制进行延迟处理，交由开发者按需决定。
3. 对二进制的数据进行了缓冲池处理，需要开发者手工处理已达到零内存分配。

##### RequestCtx操作

*RequestCtx综合http.Request和http.ResponseWriter的操作，可以更方便的读取和返回数据。

```go
func httpHandle(ctx *fasthttp.RequestCtx) {
    ctx.SetContentType("text/html")
    fmt.Fprintf(ctx, "Method:%s <br/>", ctx.Method())
    fmt.Fprintf(ctx, "URI:%s <br/>", ctx.)
    fmt.Fprintf(ctx, "RemoteAddr:%s <br/>", ctx.RemoteAddr())
	fmt.Fprintf(ctx, "UserAgent:%s <br/>", ctx.UserAgent())
	fmt.Fprintf(ctx, "Header.Accept:%s <br/>",ctx.Request.Header.Peek("Accept"))
    fmt.Fprintf(ctx, "IP:%s <br/>", ctx.RemoteIP())
	fmt.Fprintf(ctx, "Host:%s <br/>", ctx.Host())
	fmt.Fprintf(ctx, "ConnectTime:%s <br/>", ctx.ConnTime())
	fmt.Fprintf(ctx, "IsGET:%v <br/>", ctx.IsGet())
}
```

##### 表单数据

```go
func httpHandle(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html")

	// GET ?abc=abc&abc=123
	getValues := ctx.QueryArgs()
	fmt.Fprintf(ctx, "GET abc=%s <br/>",
		getValues.Peek("abc")) // Peek 只获取第一个值
	fmt.Fprintf(ctx, "GET abc=%s <br/>",
		bytes.Join(getValues.PeekMulti("abc"), []byte(","))) // PeekMulti 获取所有值

	// POST xyz=xyz&xyz=123
	postValues := ctx.PostArgs()
	fmt.Fprintf(ctx, "POST xyz=%s <br/>",
		postValues.Peek("xyz"))
	fmt.Fprintf(ctx, "POST xyz=%s <br/>",
		bytes.Join(postValues.PeekMulti("xyz"), []byte(",")))
}
```

##### Body消息体

```go
func httpHandle(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody() // 获取到的是 []byte
	fmt.Fprintf(ctx, "Body:%s", body)

	// 因为是 []byte，解析 JSON 很简单
	var v interface{}
	json.Unmarshal(body,&v)
}

func httpHandle2(ctx *fasthttp.RequestCtx) {
	ungzipBody, err := ctx.Request.BodyGunzip()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusServiceUnavailable)
		return
	}
	fmt.Fprintf(ctx, "Ungzip Body:%s", ungzipBody)
}
```

##### 上传文件

```go
func httpHandle(ctx *fasthttp.RequestCtx) {
	// 这里直接获取到 multipart.FileHeader, 需要手动打开文件句柄
	f, err := ctx.FormFile("file")
	if err != nil {
		ctx.SetStatusCode(500)
		fmt.Println("get upload file error:", err)
		return
	}
	fh, err := f.Open()
	if err != nil {
		fmt.Println("open upload file error:", err)
		ctx.SetStatusCode(500)
		return
	}
	defer fh.Close() // 记得要关

	// 打开保存文件句柄
	fp, err := os.OpenFile("saveto.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("open saving file error:", err)
		ctx.SetStatusCode(500)
		return
	}
	defer fp.Close() // 记得要关

	if _, err = io.Copy(fp, fh); err != nil {
		fmt.Println("save upload file error:", err)
		ctx.SetStatusCode(500)
		return
	}
	ctx.Write([]byte("save file successfully!"))
}
```

##### 返回内容

```go
func httpHandle(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("hello,fasthttp")
	// 因为实现不同，fasthttp 的返回内容不是即刻返回的
	// 不同于标准库，添加返回内容后设置状态码，也是有效的
	ctx.SetStatusCode(404)

	// 返回的内容也是可以获取的，不需要标准库的用法，需要自己扩展 http.ResponseWriter
	fmt.Println(string(ctx.Response.Body()))
}
// 下载文件
func httpHandle(ctx *fasthttp.RequestCtx) {
	ctx.SendFile("abc.txt")
}
```

##### RequestCtx复用引发数据竞争

RequestCtx在fasthttp中使用sync.Pool复用。在执行完RequestHandler后当前使用的RequestCtx就返回池中等待下次使用，如果业务逻辑有跨goroutine使用RequestCtx，可能会遇到：同一个RequestCtx在RequestHandler结束时放回池中，立刻被另一个连接使用，业务goroutine还在使用这个RequestCtx，读取的数据发生变化。

两种处理方式：

1. 给这次请求处理设置timeout，保证RequestCtx的使用时RequestHandler没有结束。提供了fasthttp.TimeoutHandler

```go
func httpHandle(ctx *fasthttp.RequestCtx) {
	resCh := make(chan string, 1)
	go func() {
		// 这里使用 ctx 参与到耗时的逻辑中
		time.Sleep(5 * time.Second)
		resCh <- string(ctx.FormValue("abc"))
	}()

	// RequestHandler 阻塞，等着 ctx 用完或者超时
	select {
	case <-time.After(1 * time.Second):
		ctx.TimeoutError("timeout")
	case r := <-resCh:
		ctx.WriteString("get: abc = " + r)
	}
}
```

2. fasthttp 不推荐复制 RequestCtx。但是根据业务思考，如果只是收到请求数据立即返回，后续处理数据的情况，复制 RequestCtx.Request 是可以的，因此也可以使用：

```go
func httpHandle(ctx *fasthttp.RequestCtx) {
	var req fasthttp.Request
	ctx.Request.CopyTo(&req)
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("GET abc=" + string(req.URI().QueryArgs().Peek("abc")))
	}()
	ctx.WriteString("hello fasthttp")
}
```

##### BytesBuffer

```go
func httpHandle(ctx *fasthttp.RequestCtx) {
	b := fasthttp.AcquireByteBuffer()
	b.B = append(b.B, "Hello "...)
	// 这里是编码过的 HTML 文本了，&gt;strong 等
	b.B = fasthttp.AppendHTMLEscape(b.B, "<strong>World</strong>")
	defer fasthttp.ReleaseByteBuffer(b) // 记得释放

	ctx.Write(b.B)
}
```

