import { calcRuleResultQuery } from '@/services/ant-design-pro/api';
import { type ProColumns, ProFormDateRangePicker, ProTable } from '@ant-design/pro-components';
import { Form, Modal } from 'antd';
import dayjs from 'dayjs';
import moment from 'moment';
import React, { useState } from 'react';

export type Props = {
  updateModalOpen: boolean;
  onCancel: (flag?: boolean, formVals?: API.CalcRuleListItem) => void;
  values: API.CalcRuleListItem;
};

const CalcRuleResult: React.FC<Props> = (props) => {
  const [form] = Form.useForm();
  const [startTime, setStartTime] = useState<number | null>(null);
  const [endTime, setEndTime] = useState<number | null>(null);
  const [param, setParam] = useState<any>(null);
  const handleDateChange = async (dates, dateStrings) => {
    const start = dates[0] ? dates[0].valueOf() : null;
    const end = dates[1] ? dates[1].valueOf() : null;
    setStartTime(start);
    setEndTime(end);
    console.log('开始时间戳:', start, '结束时间戳:', end);

    let parm = {
      start_time: Math.floor(start / 1000),
      end_time: Math.floor(end / 1000),
      rule_id: props.values.ID,
    };
    let response = await calcRuleResultQuery(parm);
    setParam(response.data);
  };

  const pp: ProColumns[] = [
    {
      key: 'ex_time',
      title: '执行时间',
      hideInSearch: true,
      valueType: 'dateTime',
      dataIndex: 'ex_time',
      render: (_, record) => {
        return dayjs(record.ex_time * 1000).format('YYYY-MM-DD HH:mm:ss');
      },
    },
    {
      key: 'start_time',
      title: '计算时间范围(开始)',
      hideInSearch: true,
      valueType: 'dateTime',
      dataIndex: 'start_time',
      render: (_, record) => {
        return dayjs(record.start_time * 1000).format('YYYY-MM-DD HH:mm:ss');
      },
    },
    {
      key: 'end_time',
      title: '计算时间范围(结束)',
      hideInSearch: true,
      valueType: 'dateTime',
      dataIndex: 'end_time',
      render: (_, record) => {
        return dayjs(record.end_time * 1000).format('YYYY-MM-DD HH:mm:ss');
      },
    },
    {
      key: 'param',
      title: '参数',
      hideInSearch: true,
      valueType: 'jsonCode',
      dataIndex: 'param',
      render: (_, record) => {
        return (
          <>
            <div>{JSON.stringify(record.param)}</div>
          </>
        );
      },
    },
    {
      key: 'result',
      title: '参数',
      hideInSearch: true,
      valueType: 'jsonCode',
      dataIndex: 'result',
      render: (_, record) => {
        return (
          <>
            <div>{JSON.stringify(record.result)}</div>
          </>
        );
      },
    },
  ];

  return (
    <Modal
      width={'75%'}
      form={form}
      destroyOnClose
      title={'历史计算结果'}
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

      <ProTable dataSource={param} search={false} options={false} columns={pp} />
    </Modal>
  );
};

export default CalcRuleResult;
