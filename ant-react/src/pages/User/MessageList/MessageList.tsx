import { cc, getMessageTypeDescription, messagePage } from '@/services/ant-design-pro/api';
import { FormattedMessage } from '@@/exports';
import {
  type ActionType,
  PageContainer,
  type ProColumns,
  ProTable,
} from '@ant-design/pro-components';
import { useIntl } from '@umijs/max';
import { Tag } from 'antd';
import React, { useRef } from 'react';

const Admin: React.FC = () => {
  const intl = useIntl();

  const actionRef = useRef<ActionType>();

  const ss = cc();

  const columns: ProColumns<API.MessageListItem>[] = [
    {
      title: <FormattedMessage id="pages.id" defaultMessage="唯一码" />,
      hideInSearch: true,
      dataIndex: 'ID', // @ts-ignor
    },
    {
      title: <FormattedMessage id="pages.content" defaultMessage="内容" />,
      hideInSearch: true,
      dataIndex: 'content',
    },
    {
      title: <FormattedMessage id="pages.message.type" defaultMessage="消息类型" />,
      hideInSearch: false,
      dataIndex: 'message_type_id',
      valueEnum: {
        1: { text: '计划开始通知', status: 'Default' },
        2: { text: '计划临期通知', status: 'Default' },
        3: { text: '计划到期通知', status: 'Default' },
        4: { text: '生产开始通知', status: 'Default' },
        5: { text: '生产完成通知', status: 'Default' },
        6: { text: '维修通知', status: 'Default' },
        7: { text: '维修开始通知', status: 'Default' },
        8: { text: '维修结束通知', status: 'Default' },
        9: { text: 'SIM卡超时通知', status: 'Default' },
        10: { text: '设备掉线通知', status: 'Default' },
      },
      render: (e, record) => (
        <Tag color="success">{getMessageTypeDescription(record.message_type_id)}</Tag>
      ),
    },
  ];

  return (
    <PageContainer>
      <ProTable<API.MessageListItem, API.PageParams>
        headerTitle={intl.formatMessage({
          id: 'pages.searchTable.title',
          defaultMessage: 'Enquiry form',
        })}
        actionRef={actionRef}
        rowKey="key"
        search={{
          labelWidth: 120,
        }}
        request={messagePage}
        columns={columns}
      />
    </PageContainer>
  );
};

export default Admin;
