import React from 'react';
import { Layout, Button, Space } from 'antd';
import { LogoutOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../stores';
import './Header.css';

const { Header } = Layout;

export default function AppHeader() {
  const navigate = useNavigate();
  const { logout } = useAuthStore();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <Header className="app-header">
      <div className="header-content">
        <div className="logo">
          <h1>☁️ CloudPhone</h1>
          <span className="subtitle">Control Platform</span>
        </div>
        <Space>
          <Button 
            type="primary" 
            danger 
            icon={<LogoutOutlined />}
            onClick={handleLogout}
          >
            Logout
          </Button>
        </Space>
      </div>
    </Header>
  );
}
