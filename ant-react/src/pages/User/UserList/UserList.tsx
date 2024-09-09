import {
  type ActionType,
  FooterToolbar,
  ModalForm,
  PageContainer,
  type ProColumns,
  ProFormText,
  ProTable
} from '@ant-design/pro-components';
import React, {useRef, useState} from 'react';
import {Button, message} from 'antd';
import {useIntl} from '@umijs/max';
import {FormattedMessage} from "@@/exports";
import {PlusOutlined} from "@ant-design/icons";
import {addUser, userList} from "@/services/ant-design-pro/api";

const handleAdd = async (fields: API.UserListItem) => {
  const hide = message.loading('正在添加');
  try {
    await addUser({...fields});
    hide();
    message.success('Added successfully');
    return true;
  } catch (error) {
    hide();
    message.error('Adding failed, please try again!');
    return false;
  }
};

const handleRemove = async (id: any) => {
  console.log("删除数据")
}
const Admin: React.FC = () => {
  const intl = useIntl();
  /**
   * @en-US Pop-up window of new window
   * @zh-CN 新建窗口的弹窗
   *  */
  const [createModalOpen, handleModalOpen] = useState<boolean>(false);
  /**
   * @en-US The pop-up window of the distribution update window
   * @zh-CN 分布更新窗口的弹窗
   * */
  const [updateModalOpen, handleUpdateModalOpen] = useState<boolean>(false);

  const [showDetail, setShowDetail] = useState<boolean>(false);

  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<API.UserListItem>();
  const [selectedRowsState, setSelectedRows] = useState<API.UserListItem[]>([]);

  const columns: ProColumns<API.UserListItem>[] = [{
    title: (<FormattedMessage
      id="pages.id"
      defaultMessage="唯一码"

    />), hideInSearch: true, dataIndex: 'name', // @ts-ignore
    tip: 'The rule name is the unique key', render: (dom, entity) => {
      return (<a
        onClick={() => {
          setCurrentRow(entity);
          setShowDetail(true);
        }}
      >
        {dom}
      </a>);
    },
  },

    {
      title: (<FormattedMessage
        id="pages.user-list.username"
        defaultMessage="用户名"
      />), hideInSearch: true, dataIndex: 'username',
    }, {
      title: (<FormattedMessage
        id="pages.user-list.password"
        defaultMessage="密码"
      />), hideInSearch: true, dataIndex: 'password',
    }, {
      title: (<FormattedMessage
        id="pages.user-list.mail"
        defaultMessage="邮箱"
      />), hideInSearch: true, dataIndex: 'mail',
    },

    {
      title: <FormattedMessage id="pages.searchTable.titleOption" defaultMessage="Operating"/>,
      dataIndex: 'option',
      valueType: 'option',
      render: (_, record) => [<a
        key="config"
        onClick={() => {
          handleUpdateModalOpen(true);
          setCurrentRow(record);
        }}
      >
        <FormattedMessage id="pages.searchTable.config" defaultMessage="Configuration"/>
      </a>,],
    },];

  return (<PageContainer>
    <ProTable<API.UserListItem, API.PageParams>
      headerTitle={intl.formatMessage({
        id: 'pages.searchTable.title', defaultMessage: 'Enquiry form',
      })}
      actionRef={actionRef}
      rowKey="key"
      search={{
        labelWidth: 120,
      }}
      toolBarRender={() => [<Button
        type="primary"
        key="primary"
        onClick={() => {
          handleModalOpen(true);
        }}
      >
        <PlusOutlined/> <FormattedMessage id="pages.searchTable.new" defaultMessage="New"/>
      </Button>,]}
      request={userList}
      columns={columns}
      rowSelection={{
        onChange: (_, selectedRows) => {
          setSelectedRows(selectedRows);
        },
      }}
    />
    {selectedRowsState?.length > 0 && (<FooterToolbar
      extra={<div>
        <FormattedMessage id="pages.searchTable.chosen" defaultMessage="Chosen"/>{' '}
        <a style={{fontWeight: 600}}>{selectedRowsState.length}</a>{' '}
        <FormattedMessage id="pages.searchTable.item" defaultMessage="项"/>
        &nbsp;&nbsp;
        <span>
                <FormattedMessage
                  id="pages.searchTable.totalServiceCalls"
                  defaultMessage="Total number of service calls"
                />{' '}
          {selectedRowsState.reduce((pre, item) => pre + item.callNo!, 0)}{' '}
          <FormattedMessage id="pages.searchTable.tenThousand" defaultMessage="万"/>
              </span>
      </div>}
    >
      <Button
        onClick={async () => {
          await handleRemove(selectedRowsState);
          setSelectedRows([]);
          actionRef.current?.reloadAndRest?.();
        }}
      >
        <FormattedMessage
          id="pages.searchTable.batchDeletion"
          defaultMessage="Batch deletion"
        />
      </Button>
      <Button type="primary">
        <FormattedMessage
          id="pages.searchTable.batchApproval"
          defaultMessage="Batch approval"
        />
      </Button>
    </FooterToolbar>)}
    <ModalForm
      title={intl.formatMessage({
        id: 'pages.new', defaultMessage: '新增',
      })}
      width="75%"
      open={createModalOpen}
      onOpenChange={handleModalOpen}
      onFinish={async (value) => {
        console.log(value)
        debugger
        const success = await handleAdd(value as API.UserListItem);
        if (success) {
          handleModalOpen(false);
          if (actionRef.current) {
            actionRef.current.reload();
          }
        }
      }}
    >
      <ProFormText
        label={<FormattedMessage
          id="pages.user-list.username"
        />}
        name="username"
      />
      <ProFormText
        label={<FormattedMessage
          id="pages.user-list.mail"
        />}
        name="mail"
      />
      <ProFormText
        label={<FormattedMessage
          id="pages.user-list.password"
        />}


        name="password"
      />

    </ModalForm>


  </PageContainer>);

};

export default Admin;
