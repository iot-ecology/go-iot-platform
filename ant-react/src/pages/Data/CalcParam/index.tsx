import CalcParamUpdateForm from '@/pages/Data/CalcParam/CalcParamUpdateForm';
import { initSearchDeviceUidForMqtt } from '@/pages/Data/Signal';
import DeviceUidShow from '@/pages/Data/Signal/DeviceUidShow';
import SignalNameShow from '@/pages/Data/Signal/SignalNameShow';
import {
  addCalcParam,
  calcParamPage,
  calcRuleList,
  deleteCalcParam,
  deviceList,
  mqttList,
  signalList,
  updateCalcParam,
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

const handleAdd = async (fields: API.CalcParamListItem) => {
  const hide = message.loading('正在添加');
  try {
    await addCalcParam({ ...fields });
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
    await deleteCalcParam(id);
    hide();
    message.success('删除 successfully');
    return true;
  } catch (error) {
    hide();
    message.error('Delete failed, please try again!');
    return false;
  }
};

const handlerUpdate = async (fields: API.CalcParamListItem) => {
  const hide = message.loading('正在更新');
  try {
    await updateCalcParam({ ...fields });
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
  const [currentRow, setCurrentRow] = useState<API.CalcParamListItem>();
  const [createDeviceUid, setCreateDeviceUid] = useState<string>();
  const [searchProtocol, setSearchProtocol] = useState<string>('MQTT');
  const [opDeviceUid, setOpDeviceUid] = useState<any>();
  const [searchDeviceUid, setSearchDeviceUid] = useState<number | string>('');

  const columns: ProColumns<API.CalcParamListItem>[] = [
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
      title: <FormattedMessage id="pages.calc-param.protocol" />,
      order:99,
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
      title: <FormattedMessage id="pages.calc-param.identification_code" />,
      hideInSearch: true,
      dataIndex: 'identification_code',
    },

    {
      key: 'device_uid',
      dependencies:["protocol"],
      title: <FormattedMessage id="pages.calc-param.device_uid" />,
      order:98,
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
      title: <FormattedMessage id="pages.calc-param.name" />,
      hideInSearch: true,
      dataIndex: 'name',
    },

    {
      key: 'signal_id',
      title: <FormattedMessage id="pages.calc-param.signal_id" />,
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
      key: 'reduce',
      title: <FormattedMessage id="pages.calc-param.reduce" />,
      hideInSearch: true,
      dataIndex: 'reduce',
    },
    {
      key: 'calc_rule_id',
      title: <FormattedMessage id="pages.calc-param.calc_rule_id" />,
      order:100,
      hideInSearch: false,
      dataIndex: 'calc_rule_id',
      request: async () => {
        let r = await calcRuleList();
        return r.data;
      },
      fieldProps: {
        showSearch: true,
        allowClear: false,
        fieldNames: {
          label: 'name',
          value: 'ID',
        },
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
      <ProTable<API.CalcParamListItem, API.PageParams>
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
        request={async (params, sort, filter) => {
          return await calcParamPage(params);
        }}
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
          const success = await handleAdd(value as API.CalcParamListItem);
          if (success) {
            handleModalOpen(false);
            if (actionRef.current) {
              await actionRef.current.reload();
            }
          }
        }}
      >
        <ProFormSelect
          key={'calc_rule_id'}
          label={<FormattedMessage id="pages.calc-param.calc_rule_id" />}
          name="calc_rule_id"
          request={async () => {
            let r = await calcRuleList();
            return r.data;
          }}
          fieldProps={{
            showSearch: true,
            allowClear: false,
            fieldNames: {
              label: 'name',
              value: 'ID',
            },
          }}
        />
        <ProFormSelect
          valueEnum={{
            MQTT: { text: 'MQTT', status: 'success' },
            HTTP: { text: 'HTTP', status: 'success' },
            WebSocket: { text: 'WebSocket', status: 'success' },
            TCP: { text: 'TCP', status: 'success' },
            COAP: { text: 'COAP', status: 'success' },
          }}
          key={'protocol'}
          label={<FormattedMessage id="pages.calc-param.protocol" />}
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
          label={<FormattedMessage id="pages.calc-param.device_uid" />}
          name="device_uid"
        />

        <ProFormText
          key={'identification_code'}
          label={<FormattedMessage id="pages.calc-param.identification_code" />}
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
          label={<FormattedMessage id="pages.calc-param.signal_id" />}
          name="signal_id"
        />

        <ProFormText
          key={'name'}
          label={<FormattedMessage id="pages.calc-param.name" />}
          name="name"
        />

        <ProFormSelect
          valueEnum={{
            mean: { status: 'success', text: 'mean' },
            sum: { status: 'success', text: 'sum' },
            max: { status: 'success', text: 'max' },
            min: { status: 'success', text: 'min' },
            原始: { status: 'success', text: '原始' },
          }}
          key={'reduce'}
          label={<FormattedMessage id="pages.calc-param.reduce" />}
          name="reduce"
        />
      </ModalForm>

      <CalcParamUpdateForm
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
          <ProDescriptions<API.CalcParamListItem>
            column={2}
            title={currentRow?.name}
            request={async () => ({
              data: currentRow || {},
            })}
            params={{
              id: currentRow?.name,
            }}
            columns={columns as ProDescriptionsItemProps<API.CalcParamListItem>[]}
          />
        )}
      </Drawer>
    </PageContainer>
  );
};

export default Admin;
