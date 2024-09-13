import SignalUpdateForm from '@/pages/Data/Signal/SignalUpdateForm';
import {
  addSignal,
  deleteSignal,
  mqttList,
  signalPage,
  updateSignal,
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
  ProFormSelect,
  ProFormText,
  ProTable,
} from '@ant-design/pro-components';
import { useIntl } from '@umijs/max';
import { Button, Drawer, Form, message } from 'antd';
import React, { useRef, useState } from 'react';

const handleAdd = async (fields: API.SignalListItem) => {
  const hide = message.loading('正在添加');
  try {
    await addSignal({ ...fields });
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
    await deleteSignal(id);
    hide();
    message.success('删除 successfully');
    return true;
  } catch (error) {
    hide();
    message.error('Delete failed, please try again!');
    return false;
  }
};

const handlerUpdate = async (fields: API.SignalListItem) => {
  const hide = message.loading('正在更新');
  try {
    await updateSignal({ ...fields });
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
  const [prop, handlerProp] = useState<string>('MQTT');

  const [showDetail, setShowDetail] = useState<boolean>(false);

  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<API.SignalListItem>();

  const columns: ProColumns<API.SignalListItem>[] = [
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
      key: 'protocol',
      title: <FormattedMessage id="pages.signal.protocol" />,
      hideInSearch: false,
      dataIndex: 'protocol',
      initialValue: 'MQTT',
      valueEnum: {
        MQTT: { text: 'MQTT', status: 'success' },
        HTTP: { text: 'HTTP', status: 'success' },
        WebSocket: { text: 'WebSocket', status: 'success' },
        'TCP/IP': { text: 'TCP/IP', status: 'success' },
        COAP: { text: 'COAP', status: 'success' },
      },
    },
    {
      key: 'identification_code',
      title: <FormattedMessage id="pages.signal.identification_code" />,
      hideInSearch: true,
      dataIndex: 'identification_code',
    },
    {
      key: 'device_uid',
      title: <FormattedMessage id="pages.signal.device_uid" />,
      hideInSearch: false,
      dataIndex: 'device_uid',
    },
    {
      key: 'name',
      title: <FormattedMessage id="pages.signal.name" />,
      hideInSearch: true,
      dataIndex: 'name',
    },
    {
      key: 'alias',
      title: <FormattedMessage id="pages.signal.alias" />,
      hideInSearch: true,
      dataIndex: 'alias',
    },
    {
      key: 'type',
      title: <FormattedMessage id="pages.signal.type" />,
      hideInSearch: false,
      dataIndex: 'type',
      valueType: 'select',
      initialValue: '数字',
      valueEnum: {
        数字: { text: '数字', status: 'success' },
        文本: { text: '文本', status: 'success' },
      },
    },
    {
      key: 'unit',
      title: <FormattedMessage id="pages.signal.unit" />,
      hideInSearch: true,
      valueType: 'text',
      dataIndex: 'unit',
    },
    {
      key: 'cache_size',
      title: <FormattedMessage id="pages.signal.cache_size" />,
      hideInSearch: true,
      valueType: 'digit',
      dataIndex: 'cache_size',
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

  const [form] = Form.useForm();

  return (
    <PageContainer>
      <ProTable<API.SignalListItem, API.PageParams>
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
        request={signalPage}
        columns={columns}
      />
      <ModalForm
        form={form}
        key={'add'}
        title={intl.formatMessage({
          id: 'pages.new',
          defaultMessage: '新增',
        })}
        width="75%"
        open={createModalOpen}
        onOpenChange={handleModalOpen}
        onFinish={async (value) => {
          const success = await handleAdd(value as API.SignalListItem);
          if (success) {
            handleModalOpen(false);
            if (actionRef.current) {
              await actionRef.current.reload();
            }
          }
        }}
      >
        <ProFormSelect
          valueEnum={{
            MQTT: { text: 'MQTT', status: 'success' },
            HTTP: { text: 'HTTP', status: 'success' },
            WebSocket: { text: 'WebSocket', status: 'success' },
            'TCP/IP': { text: 'TCP/IP', status: 'success' },
            COAP: { text: 'COAP', status: 'success' },
          }}
          key={'protocol'}
          label={<FormattedMessage id="pages.signal.protocol" />}
          name="protocol"
        />
        <ProFormSelect
          multiple={false}
          dependencies={['protocol']}
          request={async (params) => {
            console.log(params);
            if (params.protocol === 'MQTT') {
              let res = await mqttList();
              return res.data;
            }
          }}
          fieldProps={{
            showSearch: true,
            allowClear: false,
            fieldNames: {
              label: 'client_id',
              value: 'ID',
            },
          }}
          key={'device_uid'}
          label={<FormattedMessage id="pages.signal.device_uid" />}
          name="device_uid"
        />
        <ProFormText
          key={'identification_code'}
          label={<FormattedMessage id="pages.signal.identification_code" />}
          name="identification_code"
        />

        <ProFormText key={'name'} label={<FormattedMessage id="pages.signal.name" />} name="name" />
        <ProFormText
          key={'alias'}
          label={<FormattedMessage id="pages.signal.alias" />}
          name="alias"
        />
        <ProFormSelect
          key={'type'}
          valueEnum={{
            数字: { status: 'success', text: '数字' },
            文本: { status: 'default', text: '文本' },
          }}
          label={<FormattedMessage id="pages.signal.type" />}
          name="type"
        />
        <ProFormText key={'unit'} label={<FormattedMessage id="pages.signal.unit" />} name="unit" />
        <ProFormDigit
          fieldProps={{ precision: 0 }}
          key={'cache_size'}
          label={<FormattedMessage id="pages.signal.cache_size" />}
          name="cache_size"
        />
      </ModalForm>

      <SignalUpdateForm
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
          <ProDescriptions<API.SignalListItem>
            column={2}
            title={currentRow?.access_number}
            request={async () => ({
              data: currentRow || {},
            })}
            params={{
              id: currentRow?.access_number,
            }}
            columns={columns as ProDescriptionsItemProps<API.SignalListItem>[]}
          />
        )}
      </Drawer>
    </PageContainer>
  );
};

export default Admin;
