import React from 'react';
import { Card, Tag, Space, Button, Tooltip, Statistic, Row, Col } from 'antd';
import { 
  DeleteOutlined, 
  EditOutlined, 
  LinkOutlined,
  CloudOutlined
} from '@ant-design/icons';
import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import './DeviceCard.css';

dayjs.extend(relativeTime);

export default function DeviceCard({ device, onConnect, onEdit, onDelete }) {
  const isOnline = device.status === 'online';
  const statusColor = {
    'online': 'green',
    'offline': 'red',
    'busy': 'orange'
  };

  return (
    <Card 
      className={`device-card ${isOnline ? 'online' : 'offline'}`}
      bordered={false}
    >
      <div className="device-header">
        <Space direction="vertical" style={{ width: '100%' }}>
          <div className="device-title">
            <CloudOutlined className="device-icon" />
            <h3>{device.device_name}</h3>
            <Tag color={statusColor[device.status]}>
              {device.status.toUpperCase()}
            </Tag>
          </div>
          
          <div className="device-info">
            <span className="info-label">Device ID:</span>
            <code className="device-id">{device.device_id}</code>
          </div>
        </Space>
      </div>

      <div className="device-stats">
        <Row gutter={16}>
          <Col xs={12} sm={6}>
            <Statistic 
              title="Memory" 
              value={device.memory_total} 
              suffix="MB"
              valueStyle={{ fontSize: '14px' }}
            />
          </Col>
          <Col xs={12} sm={6}>
            <Statistic 
              title="Storage" 
              value={device.storage_total} 
              suffix="MB"
              valueStyle={{ fontSize: '14px' }}
            />
          </Col>
          <Col xs={12} sm={6}>
            <Statistic 
              title="Active" 
              value={device.active_sessions}
              valueStyle={{ fontSize: '14px' }}
            />
          </Col>
          <Col xs={12} sm={6}>
            <Statistic 
              title="Region" 
              value={device.region || 'N/A'}
              valueStyle={{ fontSize: '14px' }}
            />
          </Col>
        </Row>
      </div>

      <div className="device-meta">
        <small>
          {device.last_heartbeat ? `Last heartbeat: ${dayjs(device.last_heartbeat).fromNow()}` : 'Never'}
        </small>
      </div>

      <div className="device-actions">
        <Space>
          <Tooltip title="Connect to this device">
            <Button 
              type="primary"
              icon={<LinkOutlined />}
              size="small"
              disabled={!isOnline}
              onClick={() => onConnect(device)}
            >
              Connect
            </Button>
          </Tooltip>
          <Tooltip title="Edit device">
            <Button 
              icon={<EditOutlined />}
              size="small"
              onClick={() => onEdit(device)}
            >
              Edit
            </Button>
          </Tooltip>
          <Tooltip title="Delete device">
            <Button 
              danger
              icon={<DeleteOutlined />}
              size="small"
              onClick={() => onDelete(device.id)}
            >
              Delete
            </Button>
          </Tooltip>
        </Space>
      </div>
    </Card>
  );
}
