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

## 11.2 函数参数化（Parameterize Function）

1. 名称

2. 一个简单的速写

```javascript
function tenPercentRaise(aPerson){
    aPerson.salary = aPerson.salary.multiply(1.1);   
}
function fivePercentRaise(aPerson){
    aPerson.salary = aPerson.salary.multiply(1.05);
}
```

重构为：

```javascript
function raise(aPerson, factor){
    aPerson.salary = aPerson.salary.multiply(1 + factor);
}
```

3. 动机

4. 做法

- 从一组相似的函数中选择一个
- 运用改变函数声明，把需要作为参数传入的字面量添加到参数列表中
- 修改该函数所有的调用处，使其在调用时传入该字面量
- 测试
- 修改函数体，令其使用新传入的参数。每使用一个新参数都要测试
- 对于其他与之相似的函数，逐一将其调用处改为调用已经参数化的函数。每次修改后都要测试

5. 范例

```javascript
function baseCharge(usage){
    if (usage < 0) return usd(0);
    const amount = 
          bottomBand(usage) * 0.03 
    	  + middleBand(usage) * 0.05 
          + topBand(usage) * 0.07;
    return usd(amount);
}
function bottomBand(usage){
    return Math.min(usage, 100);
}
function middleBand(usage){
    return usage > 100 ? Math.min(usage, 200) - 100 : 0;
}
function topBand(usage){
    return usage > 200 ? usage - 200 : 0;
}
```

middleBand函数添加参数

```javascript
function withinBand(usage, bottom, top){
    return usage > 100 ? Math.min(usage, top) - bottom : 0;
}
function baseCharge(usage){
    if (usage < 0) return usd(0);
    const amount = 
          bottomBand(usage) * 0.03 
    	  + withinBand(usage， 100， 200) * 0.05 
          + topBand(usage) * 0.07;
    return usd(amount);
}
```

bottomBand函数用withinBand函数代替

```javascript
function baseCharge(usage){
    if (usage < 0) return usd(0);
    const amount = 
          withinBand(usage, 0, 100) * 0.03 
    	  + withinBand(usage， 100， 200) * 0.05 
          + topBand(usage) * 0.07;
    return usd(amount);
}
```

topBand函数也用withinBand函数代替，Infinity代表无穷大。

```javascript
function baseCharge(usage){
    if (usage < 0) return usd(0);
    const amount = 
          withinBand(usage, 0, 100) * 0.03 
    	  + withinBand(usage， 100， 200) * 0.05 
          + withinBand(usage, 200, Infinity) * 0.07;
    return usd(amount);
}
```

## 11.3 移除标记参数（Remove Flag Argument）

1. 名称

2. 一个简单的速写

```javascript
function setDimension(name, value){
    if (name === "height"){
        this._height = value;
        return;
    }
    if (name === "width"){
        this._width = value;
        return;
    }
}
```

重构为：

```javascript
function setHeight(value) {this._height = value;}
function setWidth(value) {this._width = value;}
```

3. 动机

标记参数：调用者用它来指示被调用函数应该执行哪一部分逻辑

4. 做法

- 针对参数的每一种可能值，新建一个明确函数
- 对于“用字面量值作为参数”的函数调用者，将其改为调用新建的明确函数

5. 范例

两组调用代码

```javascript
aShipment.deliveryDate = deliveryDate(anOrder, true);
aShipment.deliveryDate = deliveryDate(anOrder, false);
```

deliveryDate函数的主体

```javascript
function deliveryDate(anOrder, isRush){
    if (isRush){
        let deliveryTime;
        if (["MA", "CT"].includes(anOrder.deliveryState)) deliveryTime = 1;
        else if (["NY", "NH"].includes(anOrder.deliveryState)) deliveryTime = 2;
        else deliveryTime = 3;
        return anOrder.placedOn.plusDays(1 + deliveryTime);
    }
    else {
        let deliveryTime;
        if (["MA", "CT", "NY"].includes(anOrder.deliveryState)) deliveryTime = 2;
        else if (["ME", "NH"].includes(anOrder.deliveryState)) deliveryTime = 3;
        else deliveryTime = 4;
        return anOrder.placedOn.plusDays(2 + deliveryTime);
    }
}
```

使用分解条件表达式

```javascript
function deliveryDate(anOrder, isRush){
	if (isPush) return rushDeliveryDate(anOrder);
    else return regularDeliveryDate(anOrder);
}
function rushDeliveryDate(anOrder){
    let deliveryTime;
    if (["MA", "CT"].includes(anOrder.deliveryState)) deliveryTime = 1;
    else if (["NY", "NH"].includes(anOrder.deliveryState)) deliveryTime = 2;
    else deliveryTime = 3;
    return anOrder.placedOn.plusDays(1 + deliveryTime);
}
function regularDeliveryDate(anOrder){
    let deliveryTime;
    if (["MA", "CT", "NY"].includes(anOrder.deliveryState)) deliveryTime = 2;
    else if (["ME", "NH"].includes(anOrder.deliveryState)) deliveryTime = 3;
    else deliveryTime = 4;
    return anOrder.placedOn.plusDays(2 + deliveryTime);
}
```

修改调用方代码

```javascript
aShipment.deliveryDate = rushDeliveryDate(anOrder);
aShipment.deliveryDate = regularDeliveryDate(anOrder);
```





































