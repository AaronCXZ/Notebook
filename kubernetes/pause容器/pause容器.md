#### 概念

Pause容器是基础容器，作为init pod存在，pod中的其它容器都是从pause容器fork出来的。

每个Pod里运行这一个特殊的被称之为Pause的容器，其它容器则为业务容器，这些业务容器共享Pause容器的网络栈和Volume挂载卷，因此他们之间的通信和数据交换更为高效，同一个Pod里的容器之间仅需通过localhost就能互相通信。

pause容器主要为每个业务容器提供以下功能：

1. PID命名空间：Pod中不同应用程序可以看到其它应用程序的进程ID。

2. 网络命名空间：Pod中的多个容器能够访问同一个IP和端口范围。

3. IPC命名空间：Pod中的多个容器能够使用SystemV IPC或者POSIX消息队列进行通信。

4. UTS命名空间：Pod中的多个容器共享一个主机名和Volumes。

5. Pod中的每个容器可以访问在Pod级别定义的Vilumes。

*[pause提供的功能](./img/pause功能.png)*