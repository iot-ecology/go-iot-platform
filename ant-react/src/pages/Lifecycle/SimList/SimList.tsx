import SimUpdateForm from '@/pages/Lifecycle/SimList/SimUpdateForm';
import { addSim, deleteSimCard, simCardPage, updateSimCard } from '@/services/ant-design-pro/api';
import { FormattedMessage } from '@@/exports';
import { PlusOutlined } from '@ant-design/icons';
import {
  type ActionType,
  ModalForm,
  PageContainer,
  type ProColumns,
  ProDescriptions,
  type ProDescriptionsItemProps,
  ProFormDatePicker,
  ProFormText,
  ProTable,
} from '@ant-design/pro-components';
import { useIntl } from '@umijs/max';
import { Button, Drawer, message, Popconfirm } from 'antd';
import dayjs from 'dayjs';
import React, { useRef, useState } from 'react';

const handleAdd = async (fields: API.SimListItem) => {
  const hide = message.loading('正在添加');
  try {
    await addSim({ ...fields });
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
    await deleteSimCard(id);
    hide();
    message.success('删除 successfully');
    return true;
  } catch (error) {
    hide();
    message.error('Delete failed, please try again!');
    return false;
  }
};

const handlerUpdate = async (fields: API.SimListItem) => {
  const hide = message.loading('正在更新');
  try {
    await updateSimCard({ ...fields });
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
  const [currentRow, setCurrentRow] = useState<API.SimListItem>();

  const columns: ProColumns<API.SimListItem>[] = [
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
      key: 'access_number',
      title: <FormattedMessage id="pages.sim.access_number" />,
      hideInSearch: false,
      dataIndex: 'access_number',
    },
    {
      key: 'iccid',
      title: <FormattedMessage id="pages.sim.iccid" />,
      hideInSearch: false,
      dataIndex: 'iccid',
    },
    {
      key: 'imsi',
      title: <FormattedMessage id="pages.sim.imsi" />,
      hideInSearch: false,
      dataIndex: 'imsi',
    },
    {
      key: 'operator',
      title: <FormattedMessage id="pages.sim.operator" />,
      hideInSearch: false,
      dataIndex: 'operator',
    },
    {
      key: 'expiration',
      title: <FormattedMessage id="pages.sim.expiration" />,
      hideInSearch: false,
      dataIndex: 'expiration',
      render: (_, record) => {
        return dayjs(record.expiration).format('YYYY-MM-DD');
      },
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
        <Popconfirm
          key="delete"
          title={<FormattedMessage id="pages.deleteConfirm" defaultMessage="确定要删除吗？" />}
          onConfirm={async () => {
            // todo: 删除接口
            const success = await handleRemove(record.ID);
            if (success) {
              if (actionRef.current) {
                await actionRef.current.reload();
              }
            }
          }}
          okText={<FormattedMessage id="pages.yes" defaultMessage="确定" />}
          cancelText={<FormattedMessage id="pages.no" defaultMessage="取消" />}
        >
          <Button danger>
            <FormattedMessage id="pages.delete" defaultMessage="删除" />
          </Button>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<API.SimListItem, API.PageParams>
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
        request={simCardPage}
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
          const success = await handleAdd(value as API.SimListItem);
          if (success) {
            handleModalOpen(false);
            if (actionRef.current) {
              await actionRef.current.reload();
            }
          }
        }}
      >
        <ProFormText
          key={'access_number'}
          label={<FormattedMessage id="pages.sim.access_number" />}
          name="access_number"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText
          key={'iccid'}
          label={<FormattedMessage id="pages.sim.iccid" />}
          name="iccid"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText
          key={'imsi'}
          label={<FormattedMessage id="pages.sim.imsi" />}
          name="imsi"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText
          key={'operator'}
          label={<FormattedMessage id="pages.sim.operator" />}
          name="operator"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormDatePicker
          transform={(value) => {
            return dayjs(value).format('YYYY-MM-DDTHH:mm:ssZ');
          }}
          key={'expiration'}
          label={<FormattedMessage id="pages.sim.expiration" />}
          name="expiration"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />
      </ModalForm>

      <SimUpdateForm
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
        {currentRow?.access_number && (
          <ProDescriptions<API.SimListItem>
            column={2}
            title={currentRow?.access_number}
            request={async () => ({
              data: currentRow || {},
            })}
            params={{
              id: currentRow?.access_number,
            }}
            columns={columns as ProDescriptionsItemProps<API.SimListItem>[]}
          />
        )}
      </Drawer>
    </PageContainer>
  );
};

export default Admin;
