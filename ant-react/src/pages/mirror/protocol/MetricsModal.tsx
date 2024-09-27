// PodInfoModal.tsx
import React, {useEffect, useRef} from 'react';
import {Modal, Card, Row, Col} from 'antd';
import {adminMetrics, podMetrics} from "@/services/ant-design-pro/api";
import {Chart} from '@antv/g2';
import dayjs from "dayjs";

const MetricsChart = (props) => {
  const chartRef = useRef(null);
  const chartInstance = useRef(null);
  const a = [];

  const fetchData = async () => {
    try {
      const response = await podMetrics(props.c);
      const parsedData = parseMetricsData(response.metrics, props.values.ID);
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



const PodInfoModal = ({visible, pod, onClose, type}) => {
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
    // 可以添加更多的指标
  ];

  let c = {}
  c.node_info=pod.node_info
  c.name = pod.node_info.name
  c.type = type
  return (
    <Modal
      width={"75%"}
      title={pod.name}
      visible={visible}
      onCancel={onClose}
      footer={null}
    >
      <p>Host: {pod.node_info.host}</p>
      <p>Port: {pod.node_info.port}</p>
      <p>Node Name: {pod.node_info.name}</p>
      <Row gutter={16}>
        {metricsProps.map((metric, index) => (
          <Col span={12} key={index} style={{marginBottom: 20}}>
            <MetricsChart values={metric} c={c}/>
          </Col>
        ))}
      </Row>
    </Modal>
  );
};

export default PodInfoModal;
