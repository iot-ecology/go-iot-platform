import InfluxdbTransmitUpdateForm from '@/pages/forward/influxdb2/InfluxdbTransmitUpdateForm';
import {
  addInfluxdbTransmit,
  deleteInfluxdbTransmit,
  InfluxdbTransmitUpdateFormPage,
  updateInfluxdbTransmit,
} from '@/services/ant-design-pro/api';
import { FormattedMessage } from '@@/exports';
import { PlusOutlined } from '@ant-design/icons';
import {
  type ActionType,
  ModalForm,
  PageContainer,
  type ProColumns,
  ProDescriptions,
  type ProDescriptionsItemProps,
  ProFormDigit,
  ProFormText,
  ProTable,
} from '@ant-design/pro-components';
import { useIntl } from '@umijs/max';
import { Button, Drawer, message } from 'antd';
import React, { useRef, useState } from 'react';

const handleAdd = async (fields: API.InfluxdbTransmitListItem) => {
  const hide = message.loading('正在添加');
  try {
    await addInfluxdbTransmit({ ...fields });
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
  const hide = message.loading('正在删除');
  try {
    await deleteInfluxdbTransmit(id);
    hide();
    message.success('删除 successfully');
    return true;
  } catch (error) {
    hide();
    message.error('Delete failed, please try again!');
    return false;
  }
};

const handlerUpdate = async (fields: API.InfluxdbTransmitListItem) => {
  const hide = message.loading('正在更新');
  try {
    await updateInfluxdbTransmit({ ...fields });
    hide();
    message.success('更新成功');
    return true;
  } catch (error) {
    hide();
    message.error('更新 failed, please try again!');
    return false;
  }
};
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
  const [currentRow, setCurrentRow] = useState<API.InfluxdbTransmitListItem>();

  const columns: ProColumns<API.InfluxdbTransmitListItem>[] = [
    {
      key: 'ID',
      title: <FormattedMessage id="pages.id" defaultMessage="唯一码" />,
      hideInSearch: true,
      dataIndex: 'ID', // @ts-ignore

      render: (dom, entity) => {
        return (
          <a
            onClick={() => {
              setCurrentRow(entity);
              setShowDetail(true);
            }}
          >
            {dom}
          </a>
        );
      },
    },

    {
      key: 'name',
      title: <FormattedMessage id="pages.InfluxdbTransmit.name" />,
      hideInSearch: false,
      dataIndex: 'name',
    },

    {
      key: 'host',
      title: <FormattedMessage id="pages.InfluxdbTransmit.host" />,
      hideInSearch: true,
      dataIndex: 'host',
    },

    {
      key: 'port',
      title: <FormattedMessage id="pages.InfluxdbTransmit.port" />,
      hideInSearch: true,
      dataIndex: 'port',
    },

    {
      key: 'token',
      title: <FormattedMessage id="pages.InfluxdbTransmit.token" />,
      hideInSearch: true,
      dataIndex: 'token',
    },

    {
      title: <FormattedMessage id="pages.searchTable.titleOption" defaultMessage="Operating" />,
      dataIndex: 'option',
      valueType: 'option',
      render: (_, record) => [
        <Button
          key="update"
          onClick={() => {
            // todo: 修改
            handleUpdateModalOpen(true);
            setCurrentRow(record);
          }}
        >
          <FormattedMessage id="pages.update" defaultMessage="修改" />
        </Button>,

        <Button
          key="delete"
          onClick={async () => {
            // todo: 删除接口
            const success = await handleRemove(record.ID);
            if (success) {
              if (actionRef.current) {
                actionRef.current.reload();
              }
            }
          }}
          danger={true}
        >
          <FormattedMessage id="pages.deleted" defaultMessage="删除" />
        </Button>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<API.InfluxdbTransmitListItem, API.PageParams>
        headerTitle={intl.formatMessage({
          id: 'pages.searchTable.title',
          defaultMessage: 'Enquiry form',
        })}
        actionRef={actionRef}
        rowKey="key"
        search={{
          labelWidth: 120,
        }}
        toolBarRender={() => [
          <Button
            type="primary"
            key="primary"
            onClick={() => {
              handleModalOpen(true);
            }}
          >
            <PlusOutlined /> <FormattedMessage id="pages.searchTable.new" defaultMessage="New" />
          </Button>,
        ]}
        request={InfluxdbTransmitUpdateFormPage}
        columns={columns}
      />
      <ModalForm
        key={'add'}
        title={intl.formatMessage({
          id: 'pages.new',
          defaultMessage: '新增',
        })}
        width="75%"
        open={createModalOpen}
        onOpenChange={handleModalOpen}
        modalProps={{
          destroyOnClose: true,
        }}
        onFinish={async (value) => {
          const success = await handleAdd(value as API.InfluxdbTransmitListItem);
          if (success) {
            handleModalOpen(false);
            if (actionRef.current) {
              await actionRef.current.reload();
            }
          }
        }}
      >
        <ProFormText
          key={'name'}
          label={<FormattedMessage id="pages.InfluxdbTransmit.name" />}
          name="name"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText
          key={'host'}
          label={<FormattedMessage id="pages.InfluxdbTransmit.host" />}
          name="host"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormDigit
          key={'port'}
          label={<FormattedMessage id="pages.InfluxdbTransmit.port" />}
          name="port"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText
          key={'token'}
          label={<FormattedMessage id="pages.InfluxdbTransmit.token" />}
          name="token"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
      </ModalForm>

      <InfluxdbTransmitUpdateForm
        key={'update'}
        updateModalOpen={updateModalOpen}
        values={currentRow || {}}
        onCancel={() => {
          handleUpdateModalOpen(false);
          if (!showDetail) {
            setCurrentRow(undefined);
            setShowDetail(false);
          }
        }}
        onSubmit={async (value) => {
          console.log(value);
          const success = await handlerUpdate(value);
          if (success) {
            handleUpdateModalOpen(false);
            if (actionRef.current) {
              await actionRef.current.reload();
            }
          }
        }}
      />

      <Drawer
        key={'detail'}
        width={600}
        open={showDetail}
        onClose={() => {
          setCurrentRow(undefined);
          setShowDetail(false);
        }}
        closable={false}
      >
        {currentRow?.ID && (
          <ProDescriptions<API.InfluxdbTransmitListItem>
            column={2}
            title={currentRow?.name}
            request={async () => ({
              data: currentRow || {},
            })}
            columns={columns as ProDescriptionsItemProps<API.InfluxdbTransmitListItem>[]}
          />
        )}
      </Drawer>
    </PageContainer>
  );
};

export default Admin;
