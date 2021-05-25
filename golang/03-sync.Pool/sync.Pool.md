### 使用方法

`sysc.Pool`是一个内存池。带GC功能的语言都存在STW问题，需要回收的内存块越多，STW持续时间就越长如果能让new出来的变量，一直不被回收，得到重复利用，就减轻了GC的压力。

示例如下：

```go
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
    c:= engine.pool.Get().(*Context)
    c.writemem.reset(w)
    c.Request = req
    c.reset()
    
    engine.handleHTTPRequest(c)
    
    engine.pool.Put(c)
}
```

先Get获取内存空间，基于这个内存做相关的处理，然后再将这个内存还回(Put)到sync.Pool。

### Pool结构

![sync.Pool全景图](./img/pool.png)

### 源码图解

![Pool.Get](./img/Get source.png)

![Pool.Put](./img/Put source.png)

![Pool.Get流程](./img/Pool.Get.png)

![Pool.Put流程](./img/Pool.Put.png)

![Pool.GC流程](./img/Pool.GC.png)

### Sync.Pool梳理

1. Pool的内容会清理？清理会造成数据丢失吗？

Go 会在每个`GC`周期内定期清理`sync.Pool`内的数据。

要分几个方面来说：

* 已经从`sync.Pool Get`的值，在`poolClean`时虽说将`pool.local`置成了`nil`，`Get`到的值依然是有效的，是被`GC`标记为黑色的，不会被`GC`回收，当`Put`后又重新加入到`sync.Pool`中

* 在第一个`GC`周期内`Put`到`sync.Pool`的数值。在第二个`GC`周期没有被`Get`使用，就会被放在`local.victim`中。如果第三个`GC`周期仍然没有被使用就会被`GC`回收。

2. runtime.GOMAXPROCS与Pool之间的关系

```go
s := p.localSize
l := p.local
if uintptr(pid) < s {
    return indexLocal(l, pid), pid
}

if p.local == nil {
    allPools = append(allPools, p)
}
// If GOMAXPROCS changes between GCs, we re-allocate the array and lose the old one.
size := runtime.GOMAXPROCS(0)
local := make([]poolLocal, size)
atomic.StorePointer(&p.local, unsafe.Pointer(&local[0])) // store-release
runtime_StoreReluintptr(&p.localSize, uintptr(size))     // store-release
```

`runtime.GOMAXPROCS(0)`是获取当前最大的P的数量，`sync.Pool`的`poolLocal`数量受P的数量的影响，会开辟`runtime.GOMAXPROCS(0)`个`poolLocal`。某些场景夏我们会使用`runtime.GOMAXPROCS(N)`来改变P的数量，会使`sync.Pool`的`pool.poolLocal`释放重新开辟新的空间。

3. 为什么要开辟runtime.GOMAXPROCS个local

`pool.local`是个`poolLocal`结构，这个结构体是private + shared链表组成，在多goroutine的Get/Put下是有数据竞争的，如果只有一个local就需要加锁来操作。每个P的local就能减少加锁造成的数据竞争问题。

4. New()的作用，假如没有New会出现什么情况

从上面的pool.Get流程图可以看出，从sync.Pool获取一个内存会尝试从当前private，shared，其它的P的shared获取或者victim获取，如果获取不到时才会调用New函数来获取，也就是New()函数才是真正开辟内存空间的，New()开辟出来的内存空间使用完毕后，调用pool.Put函数放入到sync.Pool中被重复利用。如果没有New()函数则无法开辟内存空间。

5. 先Put，再Get会出现什么情况

```go
func main(){
    pool:= sync.Pool{
        New: func() interface{} {
            return item{}
        },
    }
    pool.Put(item{value:1})
    data := pool.Get()
    fmt.Println(data)
}
```

如果直接跑这个例子，能得到想象的结果，但是在某些情况下就不是这个结果了。

不能把值Pool.Put到sync.Pool中，再使用Pool.Get取出来，因为sync.Pool不是map或者slice，放入的值有可能拿不到，sync.Pool的数据结构不支持做这个事情。因为

* sync.Pool的poolCleanup函数在系统GC时会被调用，Put到sync.Pool的值，由于可能一直得不到利用，被在某个GC周期释放掉
* 不同的goroutine绑定的P有可能不一样，当前P对应的goroutine放入到sync.Pool的值有可能被其它P对应的goroutine取到，导致房钱goroutine再也取不到这个值
* 使用runtime.GOMAXPROCS(N)来改变P的数量，会使sync.Pool的pool.poolLocal释放重新开辟新的空间，导致sync.Pool被释放掉
* 以及其它情况......

6. 只Get不Put会内存泄露吗

Pool.Get的时候会尝试从当前private、shared、其它P的shared获取或者victim获取，如果实在取不到才会调用New函数来获取，New出来的内容本身还受系统GC来控制，所以如果我们提供的New实现不存在内存泄露的话，那么sycn.Pool是不会内存泄露的，当New出来的变量如果不被使用，就会被GC给回收。

如果不Put回sync.Pool，会造成Get的时候每次都调用New来从堆栈申请空间，达不到减轻GC压力的效果。

### 使用场景

一般情况下，当线上高并发业务出现GC问题需要优化是，才需要使用sync.Pool。

### 使用注意点

sync.Pool不能被复制

从pool.Get出来的值进行数据的清空(reset)，防止垃圾数据污染。
