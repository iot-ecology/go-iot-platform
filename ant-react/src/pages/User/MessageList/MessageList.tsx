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
      valueType: 'select',
      request: async () => {
        return cc();
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
