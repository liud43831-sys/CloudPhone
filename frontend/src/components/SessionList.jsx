import React, { useState, useEffect } from 'react';
import { 
  Table, 
  Button, 
  Space, 
  Popconfirm,
  Tag,
  Modal,
  Form,
  Input,
  message,
  Empty
} from 'antd';
import { 
  DeleteOutlined, 
  PoweroffOutlined,
  PlusOutlined 
} from '@ant-design/icons';
import dayjs from 'dayjs';
import { sessionAPI } from '../api/client';
import './SessionList.css';

export default function SessionList({ deviceId, onClose }) {
  const [sessions, setSessions] = useState([]);
  const [loading, setLoading] = useState(false);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchSessions();
    const interval = setInterval(fetchSessions, 5000); // Refresh every 5 seconds
    return () => clearInterval(interval);
  }, [deviceId]);

  const fetchSessions = async () => {
    try {
      setLoading(true);
      const response = await sessionAPI.list({ 
        device_id: deviceId,
        status: 'active'
      });
      setSessions(response.data.data || []);
    } catch (error) {
      console.error('Failed to fetch sessions:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateSession = async (values) => {
    try {
      await sessionAPI.create({
        device_id: deviceId,
        user_id: values.user_id
      });
      message.success('Session created successfully');
      setIsModalVisible(false);
      form.resetFields();
      fetchSessions();
    } catch (error) {
      message.error('Failed to create session');
    }
  };

  const handleTerminateSession = async (sessionId) => {
    try {
      await sessionAPI.terminate(sessionId);
      message.success('Session terminated');
      fetchSessions();
    } catch (error) {
      message.error('Failed to terminate session');
    }
  };

  const columns = [
    {
      title: 'Session ID',
      dataIndex: 'session_id',
      key: 'session_id',
      render: (text) => <code style={{ fontSize: '12px' }}>{text?.substring(0, 12)}...</code>
    },
    {
      title: 'User ID',
      dataIndex: 'user_id',
      key: 'user_id'
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (status) => (
        <Tag color={status === 'active' ? 'green' : 'red'}>
          {status.toUpperCase()}
        </Tag>
      )
    },
    {
      title: 'Start Time',
      dataIndex: 'start_time',
      key: 'start_time',
      render: (time) => dayjs(time).format('YYYY-MM-DD HH:mm:ss')
    },
    {
      title: 'Duration',
      dataIndex: 'duration',
      key: 'duration',
      render: (duration) => duration ? `${Math.floor(duration / 60)}m ${duration % 60}s` : '-'
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_, record) => (
        <Space size="small">
          <Popconfirm
            title="Terminate Session"
            description="Are you sure you want to terminate this session?"
            onConfirm={() => handleTerminateSession(record.id)}
            okText="Yes"
            cancelText="No"
          >
            <Button 
              danger
              size="small"
              icon={<PoweroffOutlined />}
            >
              Terminate
            </Button>
          </Popconfirm>
        </Space>
      )
    }
  ];

  return (
    <div className="session-list">
      <div className="session-header">
        <h3>Active Sessions</h3>
        <Button 
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => setIsModalVisible(true)}
        >
          New Session
        </Button>
      </div>

      {sessions.length === 0 ? (
        <Empty description="No active sessions" />
      ) : (
        <Table
          columns={columns}
          dataSource={sessions}
          rowKey="id"
          loading={loading}
          pagination={{ pageSize: 10 }}
          size="small"
        />
      )}

      <Modal
        title="Create New Session"
        visible={isModalVisible}
        onOk={() => form.submit()}
        onCancel={() => setIsModalVisible(false)}
      >
        <Form form={form} onFinish={handleCreateSession}>
          <Form.Item
            label="User ID"
            name="user_id"
            rules={[{ required: true, message: 'Please input user ID' }]}
          >
            <Input placeholder="e.g., user@example.com" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
