import {
  type ActionType,
  FooterToolbar,
  ModalForm,
  PageContainer,
  type ProColumns,
  ProDescriptions,
  type ProDescriptionsItemProps,
  ProFormText,
  ProTable
} from '@ant-design/pro-components';
import React, {useRef, useState} from 'react';
import {Button, Drawer} from 'antd';
import {useIntl} from '@umijs/max';
import {FormattedMessage} from "@@/exports";
import {PlusOutlined} from "@ant-design/icons";
import {rule} from "@/services/ant-design-pro/api";
import UpdateForm from "@/pages/TableList/components/UpdateForm";


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
  const [currentRow, setCurrentRow] = useState<API.RuleListItem>();
  const [selectedRowsState, setSelectedRows] = useState<API.RuleListItem[]>([]);

  const columns: ProColumns<API.RuleListItem>[] = [{
    title: (<FormattedMessage
      id="pages.id"
      defaultMessage="唯一码"

    />),
    hideInSearch:true,
    dataIndex: 'id', // @ts-ignore
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
        id="pages.name"
        defaultMessage="名称"
      />), hideInSearch: false, dataIndex: 'name',
    },
    {
      title: (<FormattedMessage
        id="pages.desc"
        defaultMessage="描述"
      />), hideInSearch: true, dataIndex: 'desc',
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
    <ProTable<API.RuleListItem, API.PageParams>
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
      request={rule}
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
        id: 'pages.searchTable.createForm.newRule', defaultMessage: 'New rule',
      })}
      width="75%"
      open={createModalOpen}
      onOpenChange={handleModalOpen}
      onFinish={async (value) => {
        const success = await handleAdd(value as API.RuleListItem);
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
          id="pages.name"
        />}
        name="name"
      />

    </ModalForm>


  </PageContainer>);

};

export default Admin;
