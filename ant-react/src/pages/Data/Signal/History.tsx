import { influxdbQuery } from '@/services/ant-design-pro/api';
import { ProFormDateRangePicker } from '@ant-design/pro-components';
import { Chart } from '@antv/g2';
import { Form, Modal } from 'antd';
import dayjs from 'dayjs';
import moment from 'moment';
import React, { useRef, useState } from 'react';

export type Props = {
  updateModalOpen: boolean;
  onCancel: (flag?: boolean, formVals?: API.SignalListItem) => void;
  values: API.SignalListItem;
};

const HistoryList: React.FC<Props> = (props) => {
  const [form] = Form.useForm();
  const [startTime, setStartTime] = useState<number | null>(null);
  const [endTime, setEndTime] = useState<number | null>(null);
  const [param, setParam] = useState<string | null>(null);
  const [chartData, setChartData] = useState([]); // State to store chart data
  const chartRef = useRef(null);
  const handleDateChange = async (dates, dateStrings) => {
    const start = dates[0] ? dates[0].valueOf() : null;
    const end = dates[1] ? dates[1].valueOf() : null;
    setStartTime(start);
    setEndTime(end);
    console.log('开始时间戳:', start, '结束时间戳:', end);

    let parm = {
      device_uid: Number(props.values.device_uid),
      protocol: String(props.values.protocol),

      measurement:
        props.values.protocol +
        '_' +
        props.values.device_uid +
        '_' +
        props.values.identification_code,
      fields: [String(props.values.ID), 'storage_time', 'push_time'],
      start_time: Math.floor(start / 1000),
      end_time: Math.floor(end / 1000),
      aggregation: {
        every: 1,
        function: props.values.type === '数字' ? 'mean' : 'first',
        create_empty: false,
      },
    };
    setParam(JSON.stringify(parm));

    let response = await influxdbQuery(parm);
    // setChartData(response.data[props.values.ID]);

    const chart = new Chart({
      container: chartRef.current,
      autoFit: true,
      height: 400,
    });

    chart
      .line()
      .data(response.data[props.values.ID])
      .encode('x', (d) => {
        return dayjs(d._time).format('YYYY-MM-DD HH:mm:ss');
      })
      .encode('y', '_value')

      .axis('y', {
        title: props.values.alias + '(' + props.values.unit + ')',
      });
    chart.render();
  };

  return (
    <Modal
      form={form}
      width={'75%'}
      destroyOnClose
      title={'历史数据'}
      forceRender={true}
      open={props.updateModalOpen}
      onCancel={() => {
        props.onCancel();
      }}
      onClose={() => {
        props.onCancel();
      }}
      footer={false}
    >
      <ProFormDateRangePicker
        fieldProps={{
          showTime: true,
          format: 'YYYY-MM-DD HH:mm:ss',
          onChange: handleDateChange,
          presets: [
            {
              label: '最近一小时',
              value: () => [moment().subtract(1, 'hours'), moment()],
            },
            {
              label: '最近一周',
              value: () => [moment().subtract(1, 'weeks'), moment()],
            },
            {
              label: '最近一个月',
              value: () => [moment().subtract(1, 'months'), moment()],
            },
            {
              label: '最近三个月',
              value: () => [moment().subtract(3, 'months'), moment()],
            },
            {
              label: '最近半年',
              value: () => [moment().subtract(6, 'months'), moment()],
            },
            {
              label: '最近一年',
              value: () => [moment().subtract(1, 'years'), moment()],
            },
            {
              label: '今年至今',
              value: () => [moment().startOf('year'), moment()],
            },
          ],
        }}
        name="contractTime"
        label="上报时间"
      />

      {/* 数据表格 */}
      <div id="ccc" ref={chartRef}></div>
    </Modal>
  );
};

export default HistoryList;
