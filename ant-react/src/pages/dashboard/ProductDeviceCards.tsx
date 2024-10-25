import React, { useEffect } from 'react';
import { Chart } from '@antv/g2';
import {Card} from "antd";

const ProductDeviceChart: React.FC = () => {
  const data = [
    { item: '产品 1', count: 120, percent: 0.4 },
    { item: '产品 2', count: 85, percent: 0.285 },
    { item: '产品 3', count: 50, percent: 0.1667 },
    { item: '产品 4', count: 200, percent: 0.6667 },
  ];

  useEffect(() => {
    const chart = new Chart({
      container: 'container',
      autoFit: true,
      height: 400,
    });

    chart.options({
      type: 'view',
      coordinate: { type: 'theta', outerRadius: 0.8, innerRadius: 0.5 },
      data,
      children: [
        {
          type: 'interval',
          encode: { y: 'percent', color: 'item' },
          transform: [{ type: 'stackY' }],
          labels: [
            {
              position: 'outside',
              text: (data) => `${data.item}: ${Math.round(data.percent * 100)}%`,
            },
          ],
          tooltip: {
            items: [
              (data) => ({
                name: data.item,
                value: `${Math.round(data.percent * 100)}%`,
              }),
            ],
          },
        },
        {
          type: 'text',
          style: {
            text: `${data.reduce((acc, item) => acc + item.count, 0)}`,
            x: '50%',
            y: '50%',
            dy: -25,
            fontSize: 44,
            fill: '#8c8c8c',
            textAlign: 'center',
          },
        },
        {
          type: 'text',
          style: {
            text: '',
            x: '50%',
            y: '50%',
            dx: 30,
            dy: 20,
            fontSize: 34,
            fill: '#8c8c8c',
            textAlign: 'center',
          },
        },
      ],
    });

    chart.render();

    return () => chart.destroy(); // 清理
  }, [data]);

  return (
    <Card >
    <div
      id="container"

    />
    </Card>
  );
};

export default ProductDeviceChart;
