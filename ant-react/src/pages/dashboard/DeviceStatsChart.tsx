import React, { useEffect, useRef, useState } from 'react';
import {Card, Row} from 'antd';
import { Chart } from '@antv/g2';
import axios from 'axios';

interface DeviceStats {
  creation_date: string;
  total_devices: number;
  protocol: string;
}

const DeviceStatsChart: React.FC = () => {
  const [data, setData] = useState<DeviceStats[]>([]);
  const chartRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get('http://localhost:8080/device-stats');
        setData(response.data.data);
      } catch (err) {
        console.error('Failed to fetch data', err);
      }
    };

    fetchData();
  }, []);

  useEffect(() => {
    if (data.length > 0 && chartRef.current) {
      const chart = new Chart({
        container: chartRef.current,
        autoFit: true,
      });

      const processedData = data.map(item => ({
        creation_date: item.creation_date,
        total_devices: item.total_devices,
        protocol: item.protocol,
      }));

      chart
        .data(processedData)
        .scale('creation_date', {
          type: 'time',
          nice: true,
        })
        .scale('total_devices', {
          nice: true,
        })
        .encode('x', 'creation_date')
        .encode('y', 'total_devices')
        .encode('color', 'protocol')
        .axis('xDim', {
          title: {
            fontSize: '12', // 文本大小
            textAlign: 'center', // 文本对齐方式
            fill: '#999', // 文本颜色
            // ...
          }
        })
        .axis('y', { labelFormatter: (d) => `${d}个` });

      chart.line().encode('shape', 'smooth');
      chart.point().encode('shape', 'point').tooltip(false);

      chart.render();

      return () => chart.destroy(); // 清理图表实例
    }
  }, [data]);

  return (
    <Card title="Device Statistics Chart" >
      <div ref={chartRef} style={{ height: '400px' }} />
    </Card>
  );
};

export default DeviceStatsChart;
