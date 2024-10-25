import React, { useEffect, useRef } from 'react';
import { Chart } from '@antv/g2';
import axios from 'axios';
import dayjs from "dayjs";
import { Card, Row, Col } from 'antd';
import {adminMetrics} from "@/services/ant-design-pro/api";

const MetricsChart = (props) => {
  const chartRef = useRef(null);
  const chartInstance = useRef(null);
  const a = [];

  const fetchData = async () => {
    try {
      const response = await adminMetrics();
      const dataString = response;

      const parsedData = parseMetricsData(dataString, props.values.ID);
      a.push(...parsedData);
      updateChart(a);
    } catch (error) {
      console.error('Error fetching data:', error);
    }
  };

  const parseMetricsData = (dataString, id) => {
    const lines = dataString.split('\n');
    const result = [];

    lines.forEach(line => {
      if (line.startsWith('#')) return;

      const parts = line.split(' ');
      const metric = parts[0];
      const value = parseFloat(parts[1]);

      if (metric.includes(id)) {
        const time = Date.now();
        result.push({ _time: time, _value: value });
      }
    });

    return result;
  };

  const drawChart = (data) => {
    const chart = new Chart({
      container: chartRef.current,
      autoFit: true,
      height: 400,
    });

    chart
      .line()
      .data(data)
      .encode('x', (d) => dayjs(d._time).format('YYYY-MM-DD HH:mm:ss'))
      .encode('y', '_value')
      .axis('y', {
        title: `${props.values.alias} (${props.values.unit})`,
      });

    chart.render();
    chartInstance.current = chart;
  };

  const updateChart = (newData) => {
    if (chartInstance.current) {
      chartInstance.current.changeData(newData);
    }
  };

  useEffect(() => {
    fetchData();
    drawChart([]);
    const interval = setInterval(fetchData, 5000);

    return () => clearInterval(interval);
  }, [props.values.ID]);

  return (
    <Card title={props.values.alias} bordered={false}>
      <div ref={chartRef}></div>
    </Card>
  );
};

const App = () => {
  const metricsProps = [
    {
      ID: 'go_gc_duration_seconds',
      alias: '垃圾回收时长',
      unit: '秒',
    },
    {
      ID: 'go_goroutines',
      alias: '当前 Goroutines 数量',
      unit: '个',
    },
    {
      ID: 'go_memstats_alloc_bytes',
      alias: '分配的内存字节',
      unit: '字节',
    },
    {
      ID: 'go_threads',
      alias: '创建的 OS 线程数',
      unit: '个',
    },
    {
      ID: 'promhttp_metric_handler_requests_total{code="200"}',
      alias: '响应200',
      unit: '次数',
    },
    {
      ID: 'promhttp_metric_handler_requests_total{code="500"}',
      alias: '响应500',
      unit: '次数',
    },
    {
      ID: 'promhttp_metric_handler_requests_total{code="503"}',
      alias: '响应503',
      unit: '次数',
    },
  ];

  return (
    <div style={{ padding: 20 }}>
      <h1>监控数据图表</h1>
      <Row gutter={16}>
        {metricsProps.map((metric, index) => (
          <Col span={12} key={index} style={{ marginBottom: 20 }}> {/* 添加下间距 */}
            <MetricsChart values={metric} />
          </Col>
        ))}
      </Row>
    </div>
  );
};

export default App;
