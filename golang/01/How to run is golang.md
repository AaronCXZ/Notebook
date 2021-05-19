1. 理解可执行文件
```go
package main

func main() {
    println("hello world!")
}
```
* 编译
`go build -x hello.go`
![go build.png](01.go build.png)

* 不同操作系统上的规范不一样

|linux|windows|MacOS|
|:---|:---:|---:|
EFL|PE|Mach-O|
已linux的可执行文件ELF为例，ELF由几部分构成：

+ ELF header
+ Section header
+ Sections

解析ELF header-->加载文件内容值内存-->从entry point开始执行代码
+ 通过entry point找到go进程的入口，使用readelf
`readelf -h ./hello`
![readelf.png](02.readelf.png)
`dlv exec ./hello`  
![03.dlv exec.png](03.dlv exec.png)
2. GMP调度
