## 11.1 将查询函数和修改函数分离（Separate Query from Modifier）

1. 名称

2. 一个简单的速写

```javascript
function getTotalOutstandAndSendBill(){
    const result = customer.invoices.reduce((total.each) => eache.amount + total, 0);
    sendBill();
    return result;
}
```

重构为：

```javascript
function totalOutstanding(){
    return customer.invoices.reduce((total.each) => eache.amount + total, 0);
}
function sendBill(){
    emailGeteway.send(formatBill(customer));
}
```

3. 动机

任何有返回值的函数，都不应该有看得到的副作用——命令与查询分离

4. 做法

- 复制整个函数，将其作为一个查询来命名
- 从新建的查询函数中去掉所有造成副作用的语句
- 执行静态检查
- 查询所有调用原函数的地方。如果调用处用到了该函数的返回值，就将其改为调用新建的查询函数，并在下面马上再调用一次原函数。每次修改之后都要测试
- 从原函数中去掉返回值
- 测试

5. 范例

```javascript
function alertForMiscreant(people){
    for (const o of people) {
        if (p == "Don") {
            setOffAlarms();
            return "Don";
        }
        if (p === "John") {
            setOffAlarms();
            return "John";
        }
    }
    return "";
}
const found = alertForMiscreant(people);
```

首先复制整个函数，用它的查询部分为其命名

```javascript
function findMiscreant(people){
   for (const o of people) {
        if (p == "Don") {
            setOffAlarms();
            return "Don";
        }
        if (p === "John") {
            setOffAlarms();
            return "John";
        }
    }
    return ""; 
}
```

返回在新建的查询函数中去掉副作用

```javascript
function findMiscreant(people){
   for (const o of people) {
        if (p == "Don") {
            return "Don";
        }
        if (p === "John") {
            return "John";
        }
    }
    return ""; 
}
```

然后找到所有原函数的调用者，将其改为调用新建的查询函数，并在其后调用一次原函数

```javascript
const found = findMiscreant(people);
alertForMiscreant(people);
```

修改原函数，去掉返回值

```javascript
function alertForMiscreant(people){
    for (const o of people) {
        if (p == "Don") {
            setOffAlarms();
        }
        if (p === "John") {
            setOffAlarms();
        }
    }
    return;
}
```

现在两个函数有大量的重复代码。使用替换算法，修该函数

```javascript
function alertForMiscreant(people){
	if (findMiscreant(people) != "") setOffAlarms();
}
```

















































