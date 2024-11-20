import ScriptParamUpdateForm from '@/pages/Data/ScriptParam/ScriptParamUpdateForm';
import DeviceUidShow from '@/pages/Data/Signal/DeviceUidShow';
import SignalNameShow from '@/pages/Data/Signal/SignalNameShow';
import { initSearchSignalId } from '@/pages/Data/SignalWaring';
import {
  addScriptWaringParam,
  deleteScriptWaringParam,
  deviceList,
  mqttList,
  scriptWaringList,
  scritpParamPage,
  signalList,
  updateScriptWaringParam,
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
  ProFormSelect,
  ProFormText,
  ProTable,
} from '@ant-design/pro-components';
import { useIntl } from '@umijs/max';
import { Button, Drawer, message, Popconfirm } from 'antd';
import React, { useRef, useState } from 'react';
import ScriptNameShow from './ScriptNameShow';

const handleAdd = async (fields: API.ScriptWaringParamListItem) => {
  const hide = message.loading('正在添加');
  try {
    let options = { ...fields };
    options.device_uid = Number(options.device_uid);
    options.signal_id = Number(options.signal_id);
    options.signal_delay_waring_id = Number(options.signal_delay_waring_id);
    await addScriptWaringParam(options);
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
    await deleteScriptWaringParam(id);
    hide();
    message.success('删除 successfully');
    return true;
  } catch (error) {
    hide();
    message.error('Delete failed, please try again!');
    return false;
  }
};

const handlerUpdate = async (fields: API.ScriptWaringParamListItem) => {
  const hide = message.loading('正在更新');
  try {
    var dt = { ...fields };
    dt.device_uid = Number(dt.device_uid);
    dt.signal_id = Number(dt.signal_id);
    dt.signal_delay_waring_id = Number(dt.signal_delay_waring_id);
    await updateScriptWaringParam(dt);
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
  const [createDeviceUid, setCreateDeviceUid] = useState<string>();
  /**
   * @en-US The pop-up window of the distribution update window
   * @zh-CN 分布更新窗口的弹窗
   * */
  const [updateModalOpen, handleUpdateModalOpen] = useState<boolean>(false);
  const [searchSignalId, setSearchSignalId] = useState<number | string>();

  const [showDetail, setShowDetail] = useState<boolean>(false);

  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<API.ScriptWaringParamListItem>();
  const [historyModalOpen, handleHistoryModalOpen] = useState<boolean>(false);
  const [calcRuleResultModalOpen, handleCalcRuleResultModalOpen] = useState<boolean>(false);
  const [searchProtocol, setSearchProtocol] = useState<string>('MQTT');
  const [searchDeviceUid, setSearchDeviceUid] = useState<number | string>('');
  const [opDeviceUid, setOpDeviceUid] = useState<any>();
  const [opSignal, setOpSignal] = useState<any>();

  const columns: ProColumns<API.ScriptWaringParamListItem>[] = [
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
      key: 'signal_delay_waring_id',
      title: <FormattedMessage id="pages.waring-param.signal_delay_waring_id" />,
      hideInSearch: true,
      dataIndex: 'signal_delay_waring_id',
      fieldProps: {
        showSearch: true,
        allowClear: false,
        fieldNames: {
          label: 'name',
          value: 'ID',
        },
      },
      render: (dom, entity) => {
        return (
          <>
            <ScriptNameShow id={entity.signal_delay_waring_id} />
          </>
        );
      },
    },
    {
      key: 'protocol',
      title: <FormattedMessage id="pages.waring-param.protocol" />,
      hideInSearch: false,
      dataIndex: 'protocol',
      initialValue: 'MQTT',
      fieldProps: {
        onChange: async (value) => {
          setSearchDeviceUid(Number(value));
          await initSearchSignalId(searchProtocol, value, setOpSignal, setSearchSignalId);
        },
        value: searchDeviceUid,

        showSearch: true,
        allowClear: false,
        fieldNames: {
          label: 'client_id',
          value: 'ID',
        },
        options: opDeviceUid,
        placeholder: '请选择',
      },
      order: 99,
      valueEnum: {
        MQTT: { text: 'MQTT', status: 'success' },
        HTTP: { text: 'HTTP', status: 'success' },
        WebSocket: { text: 'WebSocket', status: 'success' },
        TCP: { text: 'TCP', status: 'success' },
        COAP: { text: 'COAP', status: 'success' },
      },
    },
    {
      key: 'identification_code',
      title: <FormattedMessage id="pages.waring-param.identification_code" />,
      hideInSearch: true,
      dataIndex: 'identification_code',
    },

    {
      key: 'device_uid',
      order: 98,
      title: <FormattedMessage id="pages.waring-param.device_uid" />,
      hideInSearch: false,
      dataIndex: 'device_uid',
      valueType: 'select',
      fieldProps: {
        onChange: async (value) => {
          setSearchDeviceUid(Number(value));
          await initSearchSignalId(searchProtocol, value, setOpSignal, setSearchSignalId);
        },
        value: searchDeviceUid,

        showSearch: true,
        allowClear: false,
        fieldNames: {
          label: 'client_id',
          value: 'ID',
        },
        options: opDeviceUid,
        placeholder: '请选择',
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
      key: 'name',
      title: <FormattedMessage id="pages.waring-param.name" />,
      hideInSearch: false,
      dataIndex: 'name',
    },

    {
      key: 'signal_id',
      title: <FormattedMessage id="pages.waring-param.signal_id" />,
      order: 97,
      hideInSearch: false,
      dataIndex: 'signal_id',
      valueType: 'select',
      render: (dom, entity) => {
        return (
          <>
            <SignalNameShow id={entity.signal_id}></SignalNameShow>
          </>
        );
      },
      formItemProps: {
        rules: [
          {
            required: true,
            message: '此项为必填项',
          },
        ],
      },
      fieldProps: {
        value: searchSignalId,
        onChange: (value) => {
          setSearchSignalId(Number(value));
        },
        required: true,
        showSearch: true,
        allowClear: false,
        fieldNames: {
          label: 'alias',
          value: 'ID',
        },
        options: opSignal,
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
      <ProTable<API.ScriptWaringParamListItem, API.PageParams>
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
        request={scritpParamPage}
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
          const success = await handleAdd(value as API.ScriptWaringParamListItem);
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
            TCP: { text: 'TCP', status: 'success' },
            COAP: { text: 'COAP', status: 'success' },
          }}
          key={'protocol'}
          label={<FormattedMessage id="pages.waring-param.protocol" />}
          name="protocol"
          fieldProps={{
            onClick: (v) => {
              setCreateDeviceUid('');
            },
          }}
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />

        <ProFormSelect
          multiple={false}
          dependencies={['protocol']}
          request={async (params) => {
            console.log(params);
            if (params.protocol === 'MQTT') {
              let res = await mqttList();
              return res.data;
            } else {
              let c = await deviceList();
              let r = [];
              c.data.forEach((e) => {
                if (e.protocol === params.protocol) {
                  r.push({
                    client_id: e.sn,
                    ID: e.ID,
                  });
                }
              });
              return r;
            }
          }}
          onChange={(v) => {
            setCreateDeviceUid(v);
          }}
          fieldProps={{
            value: createDeviceUid,
            showSearch: true,
            allowClear: false,
            fieldNames: {
              label: 'client_id',
              value: 'ID',
            },
          }}
          key={'device_uid'}
          label={<FormattedMessage id="pages.waring-param.device_uid" />}
          name="device_uid"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />

        <ProFormText
          key={'identification_code'}
          label={<FormattedMessage id="pages.waring-param.identification_code" />}
          name="identification_code"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
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
          label={<FormattedMessage id="pages.waring-param.signal_id" />}
          name="signal_id"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />

        <ProFormSelect
          multiple={false}
          request={async (params) => {
            let newVar = await scriptWaringList();
            return newVar.data;
          }}
          fieldProps={{
            showSearch: true,
            allowClear: false,
            fieldNames: {
              label: 'name',
              value: 'ID',
            },
          }}
          key={'signal_delay_waring_id'}
          label={<FormattedMessage id="pages.waring-param.signal_delay_waring_id" />}
          name="signal_delay_waring_id"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />

        <ProFormText
          key={'name'}
          label={<FormattedMessage id="pages.waring-param.name" />}
          name="name"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
      </ModalForm>
      <ScriptParamUpdateForm
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
          <ProDescriptions<API.ScriptWaringParamListItem>
            column={2}
            title={currentRow?.name}
            request={async () => ({
              data: currentRow || {},
            })}
            params={{
              id: currentRow?.name,
            }}
            columns={columns as ProDescriptionsItemProps<API.ScriptWaringParamListItem>[]}
          />
        )}
      </Drawer>
    </PageContainer>
  );
};

export default Admin;
