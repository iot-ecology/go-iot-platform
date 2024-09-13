import DeviceUidShow from '@/pages/Data/Signal/DeviceUidShow';
import SignalWaringUpdateForm from '@/pages/Data/SignalWaring/SignalWaringUpdateForm';
import {
  addSignalWaring,
  deleteSignalWaring,
  mqttList,
  signalList,
  signalWaringPage,
  updateSignalWaring,
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
import { Button, Drawer, message } from 'antd';
import React, { useRef, useState } from 'react';

const handleAdd = async (fields: API.SignalWaringItem) => {
  const hide = message.loading('正在添加');
  try {
    var dt = { ...fields };
    dt.device_uid = Number(dt.device_uid);
    dt.in_or_out = Number(dt.in_or_out);
    await addSignalWaring(dt);
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
    await deleteSignalWaring(id);
    hide();
    message.success('删除 successfully');
    return true;
  } catch (error) {
    hide();
    message.error('Delete failed, please try again!');
    return false;
  }
};

const handlerUpdate = async (fields: API.SignalWaringItem) => {
  const hide = message.loading('正在更新');
  try {
    var dt = { ...fields };
    dt.device_uid = Number(dt.device_uid);
    dt.in_or_out = Number(dt.in_or_out);
    await updateSignalWaring(dt);
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
  const [currentRow, setCurrentRow] = useState<API.SignalWaringItem>();
  const [searchProtocol, setSearchProtocol] = useState<string>('MQTT');
  const [searchDeviceUid, setSearchDeviceUid] = useState<number>();

  const columns: ProColumns<API.SignalWaringItem>[] = [
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
      title: <FormattedMessage id="pages.signal.waring.protocol" />,
      hideInSearch: false,
      dataIndex: 'protocol',
      initialValue: 'MQTT',
      order: 3,
      fieldProps: {
        onChange: (value) => {
          setSearchProtocol(value);
        },
      },
      formItemProps: {
        rules: [
          {
            required: true,
            message: '此项为必填项',
          },
        ],
      },
      valueEnum: {
        MQTT: { text: 'MQTT', status: 'success' },
        HTTP: { text: 'HTTP', status: 'success' },
        WebSocket: { text: 'WebSocket', status: 'success' },
        'TCP/IP': { text: 'TCP/IP', status: 'success' },
        COAP: { text: 'COAP', status: 'success' },
      },
    },
    {
      key: 'device_uid',
      title: <FormattedMessage id="pages.signal.waring.device_uid" />,
      order: 2,
      hideInSearch: false,
      dataIndex: 'device_uid',
      valueType: 'select',
      formItemProps: {
        rules: [
          {
            required: true,
            message: '此项为必填项',
          },
        ],
      },
      fieldProps: {
        onChange: (value) => {
          setSearchDeviceUid(value);
        },
        showSearch: true,
        allowClear: false,
        fieldNames: {
          label: 'client_id',
          value: 'ID',
        },
      },
      request: async (params, props) => {
        console.log(searchProtocol);

        if (searchProtocol === 'MQTT') {
          let res = await mqttList();
          return res.data;
        } else {
          return [];
        }
      },
      render: (dom, entity) => {
        return (
          <>
            <DeviceUidShow
              protocol={entity.protocol}
              device_uid={entity.device_uid}
            ></DeviceUidShow>
          </>
        );
      },
    },
    {
      key: 'identification_code',
      title: <FormattedMessage id="pages.signal.waring.identification_code" />,
      hideInSearch: true,
      dataIndex: 'identification_code',
    },
    {
      key: 'signal_id',
      title: <FormattedMessage id="pages.signal.waring.signal_id" />,
      hideInSearch: false,
      dataIndex: 'signal_id',
      valueType: 'select',
      order: 1,
      formItemProps: {
        rules: [
          {
            required: true,
            message: '此项为必填项',
          },
        ],
      },
      fieldProps: {
        required: true,
        showSearch: true,
        allowClear: false,
        fieldNames: {
          label: 'alias',
          value: 'ID',
        },
      },
      request: async (params, props) => {
        let res = await signalList({
          protocol: searchProtocol,
          device_uid: searchDeviceUid,
          type: '数字',
        });
        return res.data;
      },
    },
    {
      key: 'min',
      title: <FormattedMessage id="pages.signal.waring.min" />,
      hideInSearch: true,
      dataIndex: 'min',
    },
    {
      key: 'max',
      title: <FormattedMessage id="pages.signal.waring.max" />,
      hideInSearch: true,
      dataIndex: 'max',
    },
    {
      key: 'in_or_out',
      title: <FormattedMessage id="pages.signal.waring.in_or_out" />,
      hideInSearch: true,
      dataIndex: 'in_or_out',
      valueType: 'select',
      valueEnum: {
        1: { text: '范围内报警', status: 'success' },
        0: { text: '范围外报警', status: 'success' },
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
      <ProTable<API.SignalWaringItem, API.PageParams>
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
        request={signalWaringPage}
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
        onFinish={async (value) => {
          const success = await handleAdd(value as API.SignalWaringItem);
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
          label={<FormattedMessage id="pages.signal.waring.protocol" />}
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
          label={<FormattedMessage id="pages.signal.waring.device_uid" />}
          name="device_uid"
        />

        <ProFormText
          key={'identification_code'}
          label={<FormattedMessage id="pages.signal.waring.identification_code" />}
          name="identification_code"
        />

        <ProFormSelect
          multiple={false}
          dependencies={['device_uid', 'protocol']}
          request={async (params) => {
            console.log(params);
            if (params.device_uid && params.protocol) {
              params.ty = '数字';
              let newVar = await signalList(params);
              return newVar.data;
            }
            return [];
          }}
          fieldProps={{
            showSearch: true,
            allowClear: false,
            fieldNames: {
              label: 'alias',
              value: 'ID',
            },
          }}
          key={'signal_id'}
          label={<FormattedMessage id="pages.signal.waring.signal_id" />}
          name="signal_id"
        />
        <ProFormDigit
          key={'min'}
          label={<FormattedMessage id="pages.signal.waring.min" />}
          name="min"
        />
        <ProFormDigit
          key={'max'}
          label={<FormattedMessage id="pages.signal.waring.max" />}
          name="max"
        />
        <ProFormSelect
          key={'in_or_out'}
          label={<FormattedMessage id="pages.signal.waring.in_or_out" />}
          name="in_or_out"
          valueEnum={{
            1: { text: '范围内报警', status: 'success' },
            0: { text: '范围外报警', status: 'success' },
          }}
        />
      </ModalForm>

      <SignalWaringUpdateForm
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
          <ProDescriptions<API.SignalWaringItem>
            column={2}
            request={async () => ({
              data: currentRow || {},
            })}
            columns={columns as ProDescriptionsItemProps<API.SignalWaringItem>[]}
          />
        )}
      </Drawer>
    </PageContainer>
  );
};

export default Admin;
