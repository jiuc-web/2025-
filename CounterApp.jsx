import React,{useState,useEffect,useCallback} from "react";
import './CounterApp.css';

//自定义HOOK：封装计数器逻辑
const useCounter = (initialValue = 0) => {
    const [count , setCount] = useState(initialValue);


    //增加计数
    const increment = useCallback(() => {
        setCount(prev => prev + 1);
    },[]);

    //减少计数
    const decrement = useCallback(() => {
        setCount(prev => prev -1);
    },[])

    //重置计数
    const reset = useCallback(() => {
        setCount(initialValue);
    },[initialValue]);

    //数值变化时打印日志
    useEffect(() => {
        console.log(`当前计数值：${count}`);
    },[count]);
    return {count,increment,decrement,reset};
}

//计数器组件
const CounterApp = () => {
    const {count,increment,decrement,reset} = useCounter(0);
    return (
        <div className="counter-container">
            <h1>计数器:{count}</h1>
            <div className="button-group">
                <button className="counter-button" onClick={increment}>+ 增加</button>
                <button className="counter-button" onClick={decrement}>- 减少</button>
                <button className="counter-button" onClick={reset}>重置</button>

            </div>
        </div>
    );
};

export default CounterApp;