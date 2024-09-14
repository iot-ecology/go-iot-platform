import { scriptWaringHistory } from '@/services/ant-design-pro/api';
import { type ProColumns, ProFormDateRangePicker, ProTable } from '@ant-design/pro-components';
import { Form, Modal } from 'antd';
import dayjs from 'dayjs';
import moment from 'moment';
import React, { useState } from 'react';

export type Props = {
  updateModalOpen: boolean;
  onCancel: (flag?: boolean, formVals?: API.ScriptWaringListItem) => void;
  values: API.ScriptWaringListItem;
};

const ScriptWaringHistoryList: React.FC<Props> = (props) => {
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
      up_time_start: Math.floor(start / 1000),
      up_time_end: Math.floor(end / 1000),
      ID: props.values.ID,
    };
    let c = await scriptWaringHistory(parm);
    setParam(c.data);
  };

  const pp: ProColumns[] = [
    {
      key: 'up_time',
      title: '上报时间',
      hideInSearch: true,
      valueType: 'dateTime',
      dataIndex: 'up_time',
      render: (_, record) => {
        return dayjs(record.up_time * 1000).format('YYYY-MM-DD HH:mm:ss');
      },
    },
    {
      key: 'param',
      title: '参数',
      hideInSearch: true,
      dataIndex: 'param',
      render: (_, record) => {
        return (
          <>
            <div>{JSON.stringify(record)}</div>
          </>
        );
      },
    },
    {
      key: 'script',
      title: '脚本',
      hideInSearch: true,
      valueType: 'code',
      dataIndex: 'script',
    },

    {
      key: 'value',
      title: '是否报警',
      hideInSearch: true,
      dataIndex: 'value',
    },
    {
      key: 'insert_time',
      title: '处理时间',
      hideInSearch: true,
      valueType: 'dateTime',
      dataIndex: 'insert_time',
      render: (_, record) => {
        return dayjs(record.insert_time * 1000).format('YYYY-MM-DD HH:mm:ss');
      },
    },
  ];

  return (
    <Modal
      form={form}
      destroyOnClose
      title={'报警记录'}
      width={'75%'}
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

      <ProTable dataSource={param} search={false} options={false} columns={pp} />
    </Modal>
  );
};

export default ScriptWaringHistoryList;
