import { calcRuleMockRun } from '@/services/ant-design-pro/api';
import { ProFormDateRangePicker } from '@ant-design/pro-components';
import { Form, Modal } from 'antd';
import moment from 'moment';
import React, { useState } from 'react';

export type Props = {
  updateModalOpen: boolean;
  onCancel: (flag?: boolean, formVals?: API.CalcRuleListItem) => void;
  values: API.CalcRuleListItem;
};

const MockRun: React.FC<Props> = (props) => {
  const [form] = Form.useForm();
  const [startTime, setStartTime] = useState<number | null>(null);
  const [endTime, setEndTime] = useState<number | null>(null);
  const [param, setParam] = useState<string | null>(null);
  const handleDateChange = async (dates, dateStrings) => {
    const start = dates[0] ? dates[0].valueOf() : null;
    const end = dates[1] ? dates[1].valueOf() : null;
    setStartTime(start);
    setEndTime(end);
    console.log('开始时间戳:', start, '结束时间戳:', end);

    let parm = {
      start_time: Math.floor(start / 1000),
      end_time: Math.floor(end / 1000),
      id: props.values.ID,
    };
    let response = await calcRuleMockRun(parm);
    setParam(JSON.stringify(response.data));
  };

  return (
    <Modal
      form={form}
      destroyOnClose
      title={'模拟计算'}
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

      <div>{param}</div>
    </Modal>
  );
};

export default MockRun;
