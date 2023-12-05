### API服务器

主节点上的控制平面由API服务器，控制器管理器和调度器组成

API服务器时完成中心管理的实体，也是唯一一个直接与etcd交互的组件

API服务器的职责：

- 提供kubernetes API。这个API会在集群内被主节点，工作节点组件和原生kubernetes应用直接调用，也可以被集群外的一些客户端程序调用，比如kubectl
- 为其他集群组件提供代理，比如kubernetes管理大盘，日志串流，服务端口，以及提供kubectl exec 回话

API服务器提供的功能：

- 读取状态：获取单个对象，列出对象，获取变更消息
- 操作状态：创建，更新或删除对象

所有的状态都会持久化在etcd中

#### API服务器的HTTP接口

API服务器提供了一个RESTful HTTP API，提供JSON或Protocol Buffers和数的数据内容。出于性能考虑，集群内的通信主要使用Protobuf。

- HTTP GET方法用于获取特定资源的相关数据，或者列出一组资源
- HTTP POST方法用于创建资源
- HTTP PUT方法用于更新已有的资源
- HTTP PATCH方法用于部分更新已有资源
- HTTP DELETE方法用于销毁一个资源，这个操作时不可逆的

#### API术语

##### 类别：Kind

一个类型实体，每个对象都有一个名为Kind的字段用于告诉客户端它所代表的东西，比如一个Pod

- Object用于表示系统中的初九话的实体对象，比如Pod，Endpoint，这类对象需要名字，并且大部分属于某个命名空间
- List是一种或多种实体的列表集。这些列表对象包含少量的公共信息，比如PodList，NodeList
- 其他特殊用途的型别回用于某些对象上的特定动作或用于一些非持久化的实体，比如/binding或/scale，使用APIGroup和APIResource来提供服务发现的结果，使用Status来返回出错结果

##### API组

一组逻辑上相关的的Kind

##### 版本

每个API组都允许多个版本共存，并且他们大部分也确实如此

##### 资源

小写复数形式的单词会出现在HTTP端点中，用于暴露对于系统中某种对象的CRUD操作，常见的路径包括：

- 根路径：比如.../pods，用于列出对应类型的所有实例
- 提供资源的某个特定实力的路径，比如 .../pods/nginx

资源总是从属于某个API组的某个版本，一个例子：在default命名空间中，又一个/apis/batch/v1/namespaces/default/jobs的资源，其中batch为组，v1为版本，jobs为资源。还有一些集群范围的资源，比如节点/api/v1/nodes

#### API版本

- Alpha级通常默认时被禁止的。这个级别中的功能可能会随时不经通知就被舍弃，所以只适合在一些短期的测试集群中使用
- Beta级通常默认时启用的，，表示这里的功能都经过了必要的测试。不过对象的语意在后续的其他Beta、或正式版本中可能发生不兼容的变化
- Stable出现在正式发布的软件版本中，并会在后续很多版本中继续得到支持

#### 声明式状态管理

大部分API对象都区分资源的期望状态和当前状态。所谓规格（Spec）是对某种资源的期望状态的完整描述，会持久化道etcd中。

Spec用于描述对某种资源所期望的状态

当前状态用于描述对象的观测状态或实际状态，它由控制平面管理

### 通过命令行使用API

```shell
kubectl get deploy/cordns -n kube-system -o yaml
```

可以看到Deployment的spec段是用于定义相关参数的，status段可以看到当前的状态

```shell
kubectl proxy --port=8080
```

这个命令把API服务代理到了本地，并处理了有关身份认证和授权相关的逻辑，这样就可以直接使用HTTP来发送请求，并接受返回的JSON数据

```shell
curl http://127.0.0.1:8080/apis/batch/v1
```

想知道集群能提供哪些API资源，可以通过命令`kubectl api-resources`查看，命令`kubectl api-versions`查看不同资源版本

### API服务器是如何处理请求的

1. HTTP请求由一些在DefaultBuildHandlerChain()中注册并级联在一起的过滤器依次处理。这个处理链在`k8s.io/apiserver/pkg/server/config.go`中定义，它把这个请求通过一系列的过滤器进行处理。不同的过滤器可能会忽略这个请求，页可能在请求上添加一些信息，其实是使用一个`ctx.RequestInfo`结构体维护上下文信息，另一个可能是这个请求没有通过过滤器，则直接返回相应的HTTP返回码，并说明出错的原因。
2. 下一步，`k8s.io/apiserver/pkg/server/handler.go`中定义的复用器会把不同的HTTP路径上的请求路由到相关的处理器上
3. 请求处理器是注册在每一个API组上的，它们对HTTP请求和相关上下文进行处理，并在etcd中存取相关信息

#### 处理器

##### WithPanicRecovery

处理状态恢复并把出错记录到日志中

##### WithRequestInfo

吧RequestInfo添加到请求上下文中

##### WithWaitGroup

把所有非长时间运行的请求放入一个等待组，用于在关闭时更平滑的退出

##### WithTimeoutForNonLongRunningRequest

处理非长时间运行请求的超时，对于watch和proxy这类长时间运行的请求则不需要进行超时处理

##### WithCORS

提供CORS（跨域资源共享）实现，为HTML页面中的JavaScript提供了通过XMLHttpRequest向非本域名下的资源发送请求的能力

##### WithAuthentication

尝试对当前的请求进行身份认证，并把用户信息填入请求上下文。

- 如果认证成功，HTTP头中的Authorization头会被删除。

- 如果认证失败，会返回HTTP 401 状态码

##### WithAudit

为所有的请求加上审计日志信息。

审计日志信息中需要包含请求的源IP地址，用户请求的操作，请求的命名空间等信息

##### WithImpersontion

通过检查那些尝试切换用户的请求，处理切换用户逻辑

##### WithMaxInFlightLimit

限制同时处理的请求数

##### WithAuthorization

通过调用授权模块检查用户的访问权限，如果授权通过，则把所有请求发往复用器，通过复用器把请求分发给对应的处理器。如果用户没有足够的权限，则返回HTTP 403状态码

#### 后续处理请求

- 对/,/version,/apis,/healthz的请求，以及所有非RESTfull API的请求都会被直接处理掉

- 对RESTful资源的请求，会进入请求处理流水线

  - 准入：请求中对象进入准入处理链，这个处理链包含20种不同的准入处理插件，每一个插件都可以是变更阶段的一部分，也可以是验证阶段的一部分，也可以在两个阶段同时出现。在变更阶段，运行请求的内容

    第二准入阶段纯粹是验证，比如验证Pod相关的安全设置，或者在指定命名空间中创建对象之前先验证该命名空间是否存在

  - 验证：对请求带来的每一个系统对象进行一系列复杂的验证。比如对Service名字中出现的字符进行验证，看它是否都由DNS中允许使用的字符组成或者验证Pod中所有的容器名字是否都互不相同

  - 基于etcd的CRUD逻辑：更新逻辑会把对象从etcd中读出，

  - 对于自定义资源

  - 对于Go语言原生资源



