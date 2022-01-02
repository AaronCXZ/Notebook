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

## 11.4 保持对象完整（Preserve Whole Object）

1. 名称

2. 一个简单的速写

```javascript
const low = aRoom.daysTempRange.low;
const high = aRoom.daysTempRange.high;
if (aPlan.withinRange(low, high))
```

重构为：

```javascript
if (aPlan.withinRange(aRomm.days.TempRange))
```

3. 动机

如果将来被调的函数需要从记录中导出更多的数据，就不用为此修改参数列表，并且传递整个记录也能缩短参数列表，让函数调用更容易看懂。

4. 做法

- 新建一个空函数，给它以其王总的参数列表（即传入完成对象作为参数）
- 在新函数体内调用旧函数，并把新的参数映射到旧的参数列表
- 执行静态检查
- 逐一修改旧函数的调用者，令其使用新函数，每次修改之后执行测试
- 所有调用处都修改过来之后，使用内联函数把旧函数内联带新函数体内
- 给新函数改名，从重构开始时的容易搜索的临时名字，改为使用旧函数的名字，同时修改所有调用处

5. 范例

- 范例一：

```javascript
const low = aRoom.daysTempRange.low;
const high = aRoom.daysTempRange.high;
if (!aPlan.withinRange(low, high))
    alerts.push("room temoerature went outside range");
class HeatingPlan {
    withinRange(bottm, top){
        return (bottm >= this.,_temperatureRange.low) && (top <= this._temperatureRange.high);
    }
}
```

首先在HeatingPlan类中添加新的函数

```javascript
xxNEWWithinRange(aNumberRange){}
```

然后新函数调用现有函数

```javascript
xxNEWWithinRange(aNumberRange){
    return this.withinRange(aNumberRange.low, aNumberRange.high);
}
```

修改调用方代码，调用新函数

```javascript
if (!aPlan.xxNEWWithinRange(aRoom.daysTempRange))
    alerts.push("room temoerature went outside range");
```

内联函数

```javascript
class HeatingPlan {
    xxNEWWithinRange(aNumberRange){
    	return (aNumberRange.low >= this.,_temperatureRange.low) && (aNumberRange.high <= this._temperatureRange.high);
	}
}
```

新函数改名为旧函数的名称

```javascript
class HeatingPlan {
    withinRange(aNumberRange){
    	return (aNumberRange.low >= this.,_temperatureRange.low) && (aNumberRange.high <= this._temperatureRange.high);
	}
}
```

修改调用方代码

```javascript
if (!aPlan.withinRange(aRoom.daysTempRange))
    alerts.push("room temoerature went outside range");
```

- 范例二：换个方式创建新函数

调用方代码：

```javascript
const low = aRoom.daysTempRange.low;
const high = aRoom.daysTempRange.high;
if (!aPlan.withinRange(low, high))
    alerts.push("room temoerature went outside range");
```

首先把对旧函数的调用从条件判断中解放出来

```javascript
const low = aRoom.daysTempRange.low;
const high = aRoom.daysTempRange.high;
const isWithinRange = aPlan.withinRange(low, high);
if (!isWithinRange)
    alerts.push("room temoerature went outside range");
```

然后把输入参数也提炼出来

```javascript
const tempRange = aRoom.daysTempRange
const low = tempRange.low;
const high = tempRange.high;
const isWithinRange = aPlan.withinRange(low, high);
if (!isWithinRange)
    alerts.push("room temoerature went outside range");
```

使用提炼函数来创建新函数

```javascript
const tempRange = aRoom.daysTempRange
const isWithinRange = xxNEWWithinRange(aPlan, tempRange);
if (!isWithinRange)
    alerts.push("room temoerature went outside range");
function xxNEWWithinRange(aPlan, tempRange){
	const low = tempRange.low;
	const high = tempRange.high;
	const isWithinRange = aPlan.withinRange(low, high);
    return isWithinRange;
}
```

把新函数搬移到HeatingPlan类中

```javascript
class HeatingPlan {
	function xxNEWWithinRange(tempRange){
		const low = tempRange.low;
		const high = tempRange.high;
		const isWithinRange = this.withinRange(low, high);
    	return isWithinRange;
	} 
}
```

剩下的步骤与范例一一样，替换其他调用者，然后把旧函数内联到新函数中，然后改名即可

## 11.5 以查询取代参数（Replace Parameter with Query）

1. 名称

2. 一个简单的速写

```javascript
availableVacation(anEmployee, anEmployee.grade);
function availableVacation(anEmployee, grade){...}
```

重构为：

```javascript
availableVacation(anEmployee)
function availableVacation(anEmployee){
    const grade = enEmployee.grade;
    ...
}
```

3. 动机

参数列表应该尽量避免重复，并且参数列表越短就越容易理解

4. 做法

- 如果有必要，使用提炼函数将参数的计算过程提炼到一个独立的函数中
- 将函数体内引用该参数的地方改为调用新建的函数。每次修改后执行测试
- 全部替换完成后，使用改变函数声明将该参数去掉

5. 范例

```javascript
class Order {
    get finalPrice(){
        const basePrice = this.quantity * this.itemPrice;
        let discountLevel;
        if (this.quantity > 100) discountLevel = 2;
        else discountLevel = 1;
        return this.discountedPrice(basePrice, discountLevel);
    }
    discountedPrice(basePrice, discountLevel){
        switch (discountLevel){
            case 1: return basePrice * 0.05;
            case 2: return basePrice * 0.9;
        }
    }
}
```

使用一查询取代临时变量

```javascript
class Order {
    get finalPrice(){
        const basePrice = this.quantity * this.itemPrice;
        return this.discountedPrice(basePrice, this.discountLevel);
    }
    get discountLevel(){
        return (this.quantity > 100) ? 2 : 1;
    }
    discountedPrice(basePrice, discountLevel){
        switch (discountLevel){
            case 1: return basePrice * 0.05;
            case 2: return basePrice * 0.9;
        }
    }
}
```

将discountedPrice函数中调用discountLevel参数改为调用查询函数

```javascript
class Order {
    get finalPrice(){
        const basePrice = this.quantity * this.itemPrice;
        return this.discountedPrice(basePrice, this.discountLevel);
    }
    get discountLevel(){
        return (this.quantity > 100) ? 2 : 1;
    }
    discountedPrice(basePrice, discountLevel){
        switch (this.discountLevel){
            case 1: return basePrice * 0.05;
            case 2: return basePrice * 0.9;
        }
    }
}
```

然后改变函数声明，移除参数

```javascript
class Order {
    get finalPrice(){
        const basePrice = this.quantity * this.itemPrice;
        return this.discountedPrice(basePrice);
    }
    get discountLevel(){
        return (this.quantity > 100) ? 2 : 1;
    }
    discountedPrice(basePrice){
        switch (this.discountLevel){
            case 1: return basePrice * 0.05;
            case 2: return basePrice * 0.9;
        }
    }
}
```

## 11.6 以参数取代查询（Replace Query with Parameter）

1. 名称

2. 一个简单的速写

```javascript
targetTemperature(aPlan)
function targetTemperature(aPlan){
    currentTemperature = thermostat.currentTemperature;
    ...
}
```

重构为：

```javascript
targetTemperature(aPlan, thermostat.currentTemperature)
function targetTemperature(aPlan){...}
```

3. 动机

如果把所有依赖关系都变成参数，会导致参数列表冗长重复，如果作用域之间的共享太多，会导致函数间依赖过度，需要权衡

4. 做法

- 对执行查询操作的代码使用提炼变量，将其从函数体中分离出来
- 选择函数体代码已经不再执行查询操作，对这部分代码使用提炼函数
- 使用内联变量，消除刚才提炼出来的变量
- 对原来的函数使用内联函数
- 对新函数改名，该会原来函数的名字

5. 范例

```javascript
class HeatingPlan {
    get targetTemperature(){
        if (thermostat.selectedTemperature > this._max) return this._max;
        else if (thermostat.selectedTemperature < this._min) return this._min;
        else return thermostat.selectedTemperature;
    }
}

if (thermostat.targetTemperature > thermostat.currentTemperature) setToHeat();
else if (thermostat.targetTemperature < thermostat.currentTemperature) setToCool();
else setOff();
```

首先提炼变量

```javascript
class HeatingPlan {
    get targetTemperature(){
        const selectedTemperature = thermostat.selectedTemperature;
        if (selectedTemperature > this._max) return this._max;
        else if (selectedTemperature < this._min) return this._min;
        else return selectedTemperature;
    }
}
```

然后提炼函数

```javascript
class HeatingPlan {
    get targetTemperature(){
        const selectedTemperature = thermostat.selectedTemperature;
        return xxNEWtargetTemperature(selectedTemperature);
    }
    xxNEWtargetTemperature(selectedTemperature){
        if (selectedTemperature > this._max) return this._max;
        else if (selectedTemperature < this._min) return this._min;
        else return selectedTemperature;
    }
}
```

把刚才提炼的变量内联回去

```javascript
class HeatingPlan {
    get targetTemperature(){
        return xxNEWtargetTemperature(thermostat.selectedTemperature);
    }
    xxNEWtargetTemperature(selectedTemperature){
        if (selectedTemperature > this._max) return this._max;
        else if (selectedTemperature < this._min) return this._min;
        else return selectedTemperature;
    }
}
```

修改调用方代码

```javascript
if (thePlan.xxNEWtargetTemperature(thermostat.targetTemperature) > thermostat.currentTemperature) setToHeat();
else if (thePlan.xxNEWtargetTemperature(thermostat.targetTemperature) < thermostat.currentTemperature) setToCool();
else setOff();
```

把新函数改名为旧函数

```javascript
class HeatingPlan {
    targetTemperature(selectedTemperature){
        if (selectedTemperature > this._max) return this._max;
        else if (selectedTemperature < this._min) return this._min;
        else return selectedTemperature;
    }
}
if (thePlan.targetTemperature(thermostat.targetTemperature) > thermostat.currentTemperature) setToHeat();
else if (thePlan.targetTemperature(thermostat.targetTemperature) < thermostat.currentTemperature) setToCool();
else setOff();
```

## 11.7 移除设值函数（Remove Setting Method）

1. 名称

2. 一个简单的速写

```javascript
class Person{
    get name(){...}
    set name(aString){...}
}
```

重构为：

```javascript
class Person{
    get name(){...}
}
```

3. 动机

4. 做法

- 如果构造函数尚无法得到想要设入字段的值，就使用改变函数声明将这个值以参数的形式传入构造函数。在构造函数中调用设值函数，对字段设值
- 移除所有在构造函数之外对设值函数的调用，改为使用新的构造函数，每次修改之后都要测试
- 使用内联函数消去设值函数，如果可能的话，把字段声明为不可变
- 参数

5. 范例

```javascript
class Person {
    get name(){return this._name;}
    set name(arg){this._name = arg;}
    get id(){return this._id;}
    set id(arg){this._id = arg;}
}
const martin = new Person();
martin.name = "martin";
martin.id = "1234";
```

首先使用改变函数声明在构造函数中添加对应的参数

```javascript
class Person {
    constructor(id){
        this._id = id;
    }
    get name(){return this._name;}
    set name(arg){this._name = arg;}
    get id(){return this._id;}
    set id(arg){this._id = arg;}
}
const martin = new Person("1234");
martin.name = "martin";
martin.id = "1234";
```

使用内联函数消去设值函数

```javascript
class Person {
    constructor(id){
        this._id = id;
    }
    get name(){return this._name;}
    set name(arg){this._name = arg;}
    get id(){return this._id;}
}
const martin = new Person("1234");
martin.name = "martin";

```

## 11.8 以工厂函数取代构造函数（Replace Constructor with Factory Function）

1. 名称

2. 一个简单的速写

```javascript
leadEngineer = new Employee(document.leadEngineer, 'E');
```

重构为：

```javascript
leadEngineer = createEngineer(document.leadEngineer);
```

3. 动机

构造函数一般有一些丑陋的局限性。工厂函数不受限制。

4. 做法

- 新建一个工厂函数，让它调用现有的构造函数
- 将调用构造函数的代码改为调用工厂函数
- 每修改一处，就执行测试
- 尽量缩小构造函数的可见范围

5. 范例

```javascript
class Employee {
    constructor(name, typeCode){
        this._name = name;
        this._typeCode = typeCode;
    }
    get name() {return this._name;}
    get type() {
        return Employee.legalTyeCodes[this._typeCode];
    }
    static get leadlTypeCodes(){
        return {"E": "Engineer", "M": "Manager", "S": "Salesman"};
    }
}
candidata = new Employee(document.name, document.empType);
const leadEngineer = new Employee(document.leadEngineer, 'E');
```

创建工厂函数

```javascript
function createEmployee(name, typeCode){
    return new Employee(name, typeCode);
}
```

调用方修改为使用工厂函数

```javascript
candidata = createEmployee(document.name, document.empType);
const leadEngineer = createEmployee(document.leadEngineer, 'E');
```

但是不够优雅，再新建一个工厂函数

```javascript
function createEngineer(name){
    return new Employee(name, 'E');
}
const leadEngineer = createEngineer(document.leadEngineer);
```

## 11.9 以命令取代函数（Replace Function with Command）

1. 名称

2. 一个简单的速写

```javascript
function score(candidate, medicalExam, scoringGuide){
    let result = 0;
    let healthLevel = 0;
    ...
}
```

重构为：

```javascript
class Scorer {
    constructor(candidate, medicalExam, scoringGuide){
        this._candidate = candidate;
        this._medicalExam = medicalExam;
		this._scoringGuide = scoringGuide;
    }
    execute(){
        thsi._result = 0;
        this._healthLevel = 0;
        ...
    }
}
```

3. 动机

命令对象：将函数封装成自己的对象。提供了更大的控制灵活性和更强的表达能力

4. 做法

+ 为想要包装的函数创建一个空的类，根据该函数的名字为其命名
+ 使用搬移函数把函数移动到空的类里
+ 可以考虑给每个参数创建一个字段，并在构造函数中添加对应的参数

5. 范例

```javascript
function score(candidate, medicalExam, scoringGuide){
    let result = 0;
    let healthLevel = 0;
    let highMedicalRiskFlag = false;
    if (medicalExam.isSmoker){
        healthLevel += 10;
        highMedicalRiskFlag = true;
    }
    let certificationGrade = "regular";
    if (scoringGuide.stateWithLowCertification(candidate.originState)){
        certificationGrade = "low";
        result -= 5;
    }
    result -= Math.max(healthLevel - 5, 0);
    return result;
}
```

创建一个空的类，把上述函数搬移到这个类里

```javascript
function score(candidate, medicalExam, scoringGuide){
	return new Scorer().execute(candidate, medicalExam, scoringGuide);
}
class Scorer {
    execute(candidate, medicalExam, scoringGuide){
    	let result = 0;
    	let healthLevel = 0;
    	let highMedicalRiskFlag = false;
    	if (medicalExam.isSmoker){
    	    healthLevel += 10;
    	    highMedicalRiskFlag = true;
    	}
    	let certificationGrade = "regular";
    	if (scoringGuide.stateWithLowCertification(candidate.originState)){
    	    certificationGrade = "low";
    	    result -= 5;
    	}
    	result -= Math.max(healthLevel - 5, 0);
    	return result;
    }
}
```

搬移参数到构造函数

```javascript
function score(candidate, medicalExam, scoringGuide){
	return new Scorer(candidate, medicalExam, scoringGuide).execute();
}
class Scorer {
    constructor(candidate, medicalExam, scoringGuide){
        this._candidate = candidate;
        this._medicalExam = medicalExam;
        this._scoringGuide = scoringGuide;
    }
    execute(){
    	let result = 0;
    	let healthLevel = 0;
    	let highMedicalRiskFlag = false;
    	if (thsi._medicalExam.isSmoker){
    	    healthLevel += 10;
    	    highMedicalRiskFlag = true;
    	}
    	let certificationGrade = "regular";
    	if (this._scoringGuide.stateWithLowCertification(this._candidate.originState)){
    	    certificationGrade = "low";
    	    result -= 5;
    	}
    	result -= Math.max(healthLevel - 5, 0);
    	return result;
    }
}
```

把局部变量变成字段

```javascript
function score(candidate, medicalExam, scoringGuide){
	return new Scorer(candidate, medicalExam, scoringGuide).execute();
}
class Scorer {
    constructor(candidate, medicalExam, scoringGuide){
        this._candidate = candidate;
        this._medicalExam = medicalExam;
        this._scoringGuide = scoringGuide;
    }
    execute(){
    	this._result = 0;
    	this._healthLevel = 0;
    	this._highMedicalRiskFlag = false;
    	if (thsi._medicalExam.isSmoker){
       		this._healthLevel += 10;
        	this._highMedicalRiskFlag = true;
    	}
    	this._certificationGrade = "regular";
    	if (this._scoringGuide.stateWithLowCertification(this._candidate.originState)){
        	this._certificationGrade = "low";
        	this._result -= 5;
    	}
    	this._result -= Math.max(this._healthLevel - 5, 0);
    	return this._result;
    }
}
```

使用提炼函数等重构手法

```javascript
function score(candidate, medicalExam, scoringGuide){
	return new Scorer(candidate, medicalExam, scoringGuide).execute();
}
class Scorer {
    constructor(candidate, medicalExam, scoringGuide){
        this._candidate = candidate;
        this._medicalExam = medicalExam;
        this._scoringGuide = scoringGuide;
    }
    execute(){
    	this._result = 0;
    	this._healthLevel = 0;
    	this._highMedicalRiskFlag = false;
        
        this.scoreSmoking();
    	this._certificationGrade = "regular";
    	if (this._scoringGuide.stateWithLowCertification(this._candidate.originState)){
        	this._certificationGrade = "low";
        	this._result -= 5;
    	}
    	this._result -= Math.max(this._healthLevel - 5, 0);
    	return this._result;
    }
    scoreSmoking(){
        if (thsi._medicalExam.isSmoker){
        	this._healthLevel += 10;
        	this._highMedicalRiskFlag = true;
    	}
    }
}
```



























