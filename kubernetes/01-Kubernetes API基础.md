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



