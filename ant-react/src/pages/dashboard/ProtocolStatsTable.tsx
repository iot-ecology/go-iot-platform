import React from 'react';
import { Card, Col, Row } from 'antd';
import { UpOutlined, DownOutlined } from '@ant-design/icons';

interface ProtocolStats {
  protocol: string;
  today_processed: number;
  yesterday_processed: number;
  total_processed: number;
}

const mockData: ProtocolStats[] = [
  { protocol: 'MQTT', today_processed: 10, yesterday_processed: 8, total_processed: 100 },
  { protocol: 'HTTP', today_processed: 5, yesterday_processed: 6, total_processed: 50 },
  { protocol: 'WebSocket', today_processed: 8, yesterday_processed: 10, total_processed: 80 },
  { protocol: 'CoAP', today_processed: 2, yesterday_processed: 2, total_processed: 20 },
];

const ProtocolStatsCards: React.FC = () => {
  return (
    <Row gutter={16}>
      {mockData.map((protocol) => (
        <Col span={6} key={protocol.protocol}>
          <Card
            bordered={false}
            style={{
              marginBottom: 24,
              borderRadius: '12px',
              boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)',
              backgroundColor: '#f0f4f8',
              padding: '30px',
              transition: 'transform 0.2s, box-shadow 0.2s',
            }}
            hoverable
          >
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <div>
                <h4 style={{ margin: '0 0 12px', fontSize: '1.5rem', color: '#333' }}>{protocol.protocol}</h4>
                <p style={{ margin: '0', color: '#666', fontSize: '1rem' }}>Today: <strong>{protocol.today_processed}</strong></p>
                <p style={{ margin: '0', color: '#666', fontSize: '1rem' }}>Yesterday: <strong>{protocol.yesterday_processed}</strong></p>
                <p style={{ margin: '0', color: '#666', fontSize: '1rem' }}>Total: <strong>{protocol.total_processed}</strong></p>
              </div>
              <div>
                {protocol.today_processed > protocol.yesterday_processed ? (
                  <UpOutlined style={{ color: '#4caf50', fontSize: '28px' }} />
                ) : (
                  <DownOutlined style={{ color: '#f44336', fontSize: '28px' }} />
                )}
              </div>
            </div>
          </Card>
        </Col>
      ))}
    </Row>
  );
};

export default ProtocolStatsCards;
