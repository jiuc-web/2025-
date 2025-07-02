import React, { createContext, useContext, useReducer, useEffect, useMemo, useRef, useState } from 'react';
import './TodoApp.css';

// 创建Todo上下文
const TodoContext = createContext();

// 初始状态
const initialState = {
  todos: [],
};

// 定义reducer
const todoReducer = (state, action) => {
  switch (action.type) {
    case 'ADD_TODO':
      return {
        ...state,
        todos: [...state.todos, {
          id: Date.now(),
          text: action.payload,
          completed: false,
        }],
      };
    case 'TOGGLE_TODO':
      return {
        ...state,
        todos: state.todos.map(todo =>
          todo.id === action.payload ? { ...todo, completed: !todo.completed } : todo
        ),
      };
    case 'DELETE_TODO':
      return {
        ...state,
        todos: state.todos.filter(todo => todo.id !== action.payload),
      };
    default:
      return state;
  }
};

// TodoItem组件
const TodoItem = ({ todo }) => {
  const { dispatch } = useContext(TodoContext);

  return (
    <li className="todo-item">
      <input
        type="checkbox"
        checked={todo.completed}
        onChange={() => dispatch({ type: 'TOGGLE_TODO', payload: todo.id })}
        className="todo-checkbox"
      />
      <span className={`todo-text ${todo.completed ? 'completed' : ''}`}>
        {todo.text}
      </span>
      <button
        onClick={() => dispatch({ type: 'DELETE_TODO', payload: todo.id })}
        className="delete-button"
      >
        删除
      </button>
    </li>
  );
};

// TodoList组件
const TodoList = () => {
  const { state } = useContext(TodoContext);
  
  // 使用useMemo优化列表渲染
  const memoizedTodos = useMemo(() => state.todos, [state.todos]);

  return (
    <ul className="todo-list">
      {memoizedTodos.map(todo => (
        <TodoItem key={todo.id} todo={todo} />
      ))}
    </ul>
  );
};

// TodoForm组件
const TodoForm = () => {
  const { dispatch } = useContext(TodoContext);
  const inputRef = useRef(null);
  const [inputValue, setInputValue] = useState('');

  const handleSubmit = e => {
    e.preventDefault();
    if (inputValue.trim()) {
      dispatch({ type: 'ADD_TODO', payload: inputValue });
      setInputValue('');
      inputRef.current.focus();
    }
  };

  return (
    <form onSubmit={handleSubmit} className="todo-form">
      <input
        ref={inputRef}
        type="text"
        value={inputValue}
        onChange={e => setInputValue(e.target.value)}
        placeholder="输入任务内容"
        className="todo-input"
      />
      <button
        type="submit"
        className="add-button"
      >
        添加
      </button>
    </form>
  );
};

// 主应用组件
const TodoApp = () => {
  const [state, dispatch] = useReducer(todoReducer, initialState);

  // 组件挂载时打印消息
  useEffect(() => {
    console.log('Todo List已加载');
  }, []);

  return (
    <TodoContext.Provider value={{ state, dispatch }}>
      <div className="todo-app">
        <h1 className="app-title">简易TodoList</h1>
        <TodoForm />
        <TodoList />
      </div>
    </TodoContext.Provider>
  );
};

export default TodoApp;