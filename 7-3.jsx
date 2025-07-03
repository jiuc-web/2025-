import React, { useState } from 'react';
import {
  BrowserRouter as Router,
  Routes,
  Route,
  Link,
  Navigate,
  useParams,
  Outlet
} from 'react-router-dom';

// 模拟登录状态
const useAuth = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  
  const login = () => setIsLoggedIn(true);
  const logout = () => setIsLoggedIn(false);
  
  return { isLoggedIn, login, logout };
};

// 登录页面
const LoginPage = ({ onLogin }) => {
  return (
    <div>
      <h2>登录页面</h2>
      <button onClick={onLogin}>点击登录</button>
    </div>
  );
};

// 首页
const HomePage = () => <h2>首页</h2>;

// 商品列表
const ProductsPage = () => {
  const products = [
    { id: 1, name: '商品1' },
    { id: 2, name: '商品2' },
    { id: 3, name: '商品3' },
  ];
  
  return (
    <div>
      <h2>商品列表</h2>
      <ul>
        {products.map(product => (
          <li key={product.id}>
            <Link to={`/products/${product.id}`}>{product.name}</Link>
          </li>
        ))}
      </ul>
    </div>
  );
};

// 商品详情
const ProductDetailPage = () => {
  const { id } = useParams();
  return <h2>商品详情 - ID: {id}</h2>;
};

// 用户中心布局 (嵌套路由)
const ProfileLayout = () => {
  return (
    <div>
      <h2>用户中心</h2>
      <nav>
        <Link to="orders">我的订单</Link> | 
        <Link to="settings">设置</Link>
      </nav>
      <Outlet />
    </div>
  );
};

// 订单页面
const OrdersPage = () => <h3>我的订单</h3>;

// 设置页面
const SettingsPage = () => <h3>设置</h3>;

// 404页面
const NotFoundPage = () => <h2>404 - 页面未找到</h2>;

// 路由守卫
const PrivateRoute = ({ children, isLoggedIn }) => {
  return isLoggedIn ? children : <Navigate to="/login" />;
};

// 主应用组件
const App = () => {
  const { isLoggedIn, login } = useAuth();
  
  return (
    <Router>
      <nav>
        <Link to="/">首页</Link> | 
        <Link to="/products">商品列表</Link> | 
        <Link to="/profile">用户中心</Link> | 
        {!isLoggedIn ? (
          <Link to="/login">登录</Link>
        ) : (
          <button onClick={logout}>退出</button>
        )}
      </nav>
      
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/products" element={<ProductsPage />} />
        <Route path="/products/:id" element={<ProductDetailPage />} />
        
        <Route path="/login" element={<LoginPage onLogin={login} />} />
        
        <Route 
          path="/profile" 
          element={
            <PrivateRoute isLoggedIn={isLoggedIn}>
              <ProfileLayout />
            </PrivateRoute>
          }
        >
          <Route index element={<Navigate to="orders" replace />} />
          <Route path="orders" element={<OrdersPage />} />
          <Route path="settings" element={<SettingsPage />} />
        </Route>
        
        <Route path="*" element={<NotFoundPage />} />
      </Routes>
    </Router>
  );
};

export default App;