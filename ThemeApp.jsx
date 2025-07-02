import React, { createContext, useContext, useState, useEffect } from 'react';
import './ThemeApp.css';

// 创建主题上下文
const ThemeContext = createContext();

// 自定义Hook：封装主题逻辑
const useTheme = () => {
  const [theme, setTheme] = useState('light');

  const toggleTheme = () => {
    setTheme(prev => (prev === 'light' ? 'dark' : 'light'));
  };

  // 同步body背景色
  useEffect(() => {
    document.body.className = theme;
  }, [theme]);

  return { theme, toggleTheme };
};

// 主题卡片组件
const ThemeCard = () => {
  const { theme } = useContext(ThemeContext);
  
  return (
    <div className={`theme-card ${theme}`}>
      当前主题: {theme === 'light' ? '亮色模式' : '暗色模式'}
    </div>
  );
};

// 主题切换按钮组件
const ThemeButton = () => {
  const { toggleTheme } = useContext(ThemeContext);
  
  return (
    <button 
      onClick={toggleTheme}
      className="theme-button"
    >
      切换主题
    </button>
  );
};

// 主应用组件
const ThemeApp = () => {
  const theme = useTheme();

  return (
    <ThemeContext.Provider value={theme}>
      <div className={`theme-app ${theme.theme}`}>
        <h1>主题切换应用</h1>
        <ThemeCard />
        <ThemeButton />
      </div>
    </ThemeContext.Provider>
  );
};

export default ThemeApp;