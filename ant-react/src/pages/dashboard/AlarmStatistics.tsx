import React from 'react';
import { Card, Row, Col, Statistic } from 'antd';
import { ExclamationCircleOutlined } from '@ant-design/icons';

const AlarmStatistics: React.FC = () => {
  const alarmData = [
    { title: '高危报警', value: 32, icon: <ExclamationCircleOutlined style={{ color: '#f5222d' }} /> },
    { title: '中危报警', value: 15, icon: <ExclamationCircleOutlined style={{ color: '#faad14' }} /> },
    { title: '低危报警', value: 5, icon: <ExclamationCircleOutlined style={{ color: '#52c41a' }} /> },
    { title: '已处理报警', value: 50, icon: <ExclamationCircleOutlined style={{ color: '#1890ff' }} /> },
  ];

  return (
    <Card title="报警数据统计" bordered={false} style={{ borderRadius: '12px', boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)' }}>
      <Row gutter={16}>
        {alarmData.map((data) => (
          <Col span={6} key={data.title}>
            <Statistic
              title={data.title}
              value={data.value}
              prefix={data.icon}
              style={{ textAlign: 'center' }}
            />
          </Col>
        ))}
      </Row>
    </Card>
  );
};

export default AlarmStatistics;
