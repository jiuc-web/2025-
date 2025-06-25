function createCounter(n) {
    let count = n; // 使用闭包保存当前计数值
    return function counter() {
        const current = count; // 保存当前值
        count += 1; // 每次调用后递增
        return current; // 返回递增前的值
    };
}

function testCounter(n, calls) {
    const counter = createCounter(n);
    const result = [];
    for (const call of calls) {
        if (call === "call") {
            result.push(counter());
        }
    }
    return result;
}

console.log(testCounter(10, ["call", "call", "call"])); // [10, 11, 12]

console.log(testCounter(-2, ["call", "call", "call", "call", "call"])); // [-2, -1, 0, 1, 2]