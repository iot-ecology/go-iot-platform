import ScriptNameShow from '@/pages/Data/ScriptParam/ScriptNameShow';
import ScriptParamUpdateForm from '@/pages/Data/ScriptParam/ScriptParamUpdateForm';
import { initSearchDeviceUidForMqtt } from '@/pages/Data/Signal';
import DeviceUidShow from '@/pages/Data/Signal/DeviceUidShow';
import SignalNameShow from '@/pages/Data/Signal/SignalNameShow';
import {
  addScriptWaringParam,
  deleteScriptWaringParam,
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
import { Button, Drawer, message } from 'antd';
import React, { useRef, useState } from 'react';

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

  const [showDetail, setShowDetail] = useState<boolean>(false);

  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<API.ScriptWaringParamListItem>();
  const [historyModalOpen, handleHistoryModalOpen] = useState<boolean>(false);
  const [calcRuleResultModalOpen, handleCalcRuleResultModalOpen] = useState<boolean>(false);
  const [searchProtocol, setSearchProtocol] = useState<string>('MQTT');
  const [searchDeviceUid, setSearchDeviceUid] = useState<number | string>('');
  const [opDeviceUid, setOpDeviceUid] = useState<any>();

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
      key: 'protocol',
      title: <FormattedMessage id="pages.waring-param.protocol" />,
      hideInSearch: false,
      dataIndex: 'protocol',
      initialValue: 'MQTT',
      onChange: (value) => {
        setSearchProtocol(value);
      },
      fieldProps: (form, config: any) => ({
        onChange: async (value) => {
          if (value === 'MQTT') {
            await initSearchDeviceUidForMqtt(setSearchDeviceUid, setOpDeviceUid);
          } else {
            setSearchDeviceUid('');
            setOpDeviceUid([
              {
                client_id: 'ccc',
                ID: '1',
              },
            ]);
          }
          setSearchProtocol(value);
        },
      }),
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
      title: <FormattedMessage id="pages.waring-param.device_uid" />,
      hideInSearch: false,
      dataIndex: 'device_uid',
      valueType: 'select',
      fieldProps: (form, config: any) => ({
        showSearch: true,
        allowClear: false,
        fieldNames: {
          label: 'client_id',
          value: 'ID',
        },
      }),
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
      hideInSearch: true,
      dataIndex: 'name',
    },

    {
      key: 'signal_id',
      title: <FormattedMessage id="pages.waring-param.signal_id" />,
      hideInSearch: false,
      dataIndex: 'signal_id',
      valueType: 'select',
      order: 1,
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
        required: true,
        showSearch: true,
        allowClear: false,
        fieldNames: {
          label: 'alias',
          value: 'ID',
        },
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
              return [];
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
        />

        <ProFormText
          key={'identification_code'}
          label={<FormattedMessage id="pages.waring-param.identification_code" />}
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
          label={<FormattedMessage id="pages.waring-param.signal_id" />}
          name="signal_id"
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
        />

        <ProFormText
          key={'name'}
          label={<FormattedMessage id="pages.waring-param.name" />}
          name="name"
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
