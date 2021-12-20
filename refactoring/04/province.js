// 生产json数据
function sampleProvinceData() {
    return {
        name: "Asia",
        producers: [
            { name: "Byzantium", cost: 10, production: 9 },
            { name: "Attalia", cost: 12, production: 10 },
            { name: "Sinope", cost: 10, production: 6 },
        ],
        demand: 30,
        price: 20
    };
}

// 省份
class Province {
    constructor(doc) {
        // 名字
        this._name = doc.name;
        // 生产者
        this._producers = [];
        // 总产量
        this._totalProduction = 0;
        // 需求
        this._demand = doc.demand;
        // 价格
        this._price = doc.price;
        doc.producers.forEach(d => this.addProducer(new Producer(this, d)));
    }
    // 添加生产者
    addProducer(arg) {
        this._producers.push(arg);
        this._totalProduction += arg.production;
    }
    get name() { return this._name; }
    get producers() { return this._producers.slice(); }
    get totalProduction() { return this._totalProduction; }
    get demand() { return this._demand; }
    set demand(arg) { this._demand = arg; }
    get price() { return this._price; }
    set price(arg) { this._price = arg; }
}

// 生产者
class Producer {
    constructor(aProvince, data) {
        // 省份
        this._province = aProvince;
        // 费用
        this._cost = data.cost;
        // 名称
        this._name = data.name;
        // 生产
        this._production = data.production || 0;
    }
    get name() { return this._name; }
    get cost() { return this._cost; }
    set cost(arg) { this._cost = parseInt(arg); }

    get production() { return this._production; }
    set production(amountStr) {
        const amount = parseInt(amountStr);
        const newProduction = Number.isNaN(amount) ? 0 : amount;
        this._production.totalProduction += newProduction - this._production;
        this._production = newProduction;
    }
    // 差额
    get shortfall() { return this._demand - this.totalProduction; }
    // 利润
    get prodfit() { return this.demandValue - this.demandCost; }
    //  成本
    get demandCost() {
        let remainingDemand = this.demand;
        let result = 0;
        this.producers.
            sort((a, b) => a.cost - b.cost)
            .forEach(p => {
                const contribution = Math.min(remainingDemand, p.production);
                remainingDemand -= contribution;
                result += contribution * p.cost;
            });
        return result;
    }
    // 总量价值
    get demandValue() {
        return this.satisfiedDemand * this.price;
    }
    // 产量个需求取最小值
    get satisfiedDemand() {
        return Math.min(this._demand, this.totalProduction);
    }
}
