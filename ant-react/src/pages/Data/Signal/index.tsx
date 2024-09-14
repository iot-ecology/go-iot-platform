import DeviceUidShow from '@/pages/Data/Signal/DeviceUidShow';
import HistoryList from '@/pages/Data/Signal/History';
import SignalUpdateForm from '@/pages/Data/Signal/SignalUpdateForm';
import {
  addSignal,
  deleteSignal,
  deviceList,
  mqttList,
  signalPage,
  updateSignal,
} from '@/services/ant-design-pro/api';
import { history, useLocation } from 'umi';

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
import React, { useEffect, useRef, useState } from 'react';

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

export async function initSearchDeviceUidForMqtt(
  setSearchDeviceUid: (
    value: ((prevState: number | string) => number | string) | number | string,
  ) => void,
  setOpDeviceUid: (value: any) => void,
) {
  let res = await mqttList();
  let data = res.data;
  var value1 = data[0].ID;
  setSearchDeviceUid(value1);
  setOpDeviceUid(data);
  return value1;
}

const Admin: React.FC = () => {
  const intl = useIntl();
  const location = useLocation();
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
  const [historyModalOpen, handleHistoryModalOpen] = useState<boolean>(false);

  const [showDetail, setShowDetail] = useState<boolean>(false);
  const [searchProtocol, setSearchProtocol] = useState<string>('MQTT');
  const [searchDeviceUid, setSearchDeviceUid] = useState<number | string>('');

  const [opDeviceUid, setOpDeviceUid] = useState<any>();
  useEffect(() => {
    const queryParams = new URLSearchParams(location.search);

    if (queryParams.get('protocol')) {
      setSearchProtocol(queryParams.get('protocol'));
    }
    if (queryParams.get('id')) {
      setSearchDeviceUid(Number(queryParams.get('id')));
    } else {
      initSearchDeviceUidForMqtt(setSearchDeviceUid, setOpDeviceUid);
    }
  }, []);
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
      onChange: (value) => {
        setSearchProtocol(value);
      },
      fieldProps: (form, config: any) => ({
        onChange: async (value) => {
          if (value === 'MQTT') {
            await initSearchDeviceUidForMqtt(setSearchDeviceUid, setOpDeviceUid);
          } else {
            setSearchDeviceUid('');

            let c = await deviceList();
            let r = [];
            c.data.forEach((e) => {
              if (e.protocol === value) {
                r.push({
                  client_id: e.sn,
                  ID: e.ID,
                });
              }
            });
            setOpDeviceUid(r);
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
      title: <FormattedMessage id="pages.signal.identification_code" />,
      hideInSearch: true,
      dataIndex: 'identification_code',
    },
    {
      key: 'device_uid',
      title: <FormattedMessage id="pages.signal.device_uid" />,
      hideInSearch: false,
      dataIndex: 'device_uid',
      valueType: 'select',
      fieldProps: (form, config: any) => ({
        onChange: (value) => {
          setSearchDeviceUid(Number(value));
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
          key="history"
          onClick={async () => {
            setCurrentRow(record);
            handleHistoryModalOpen(true);
          }}
        >
          <FormattedMessage id="pages.history" defaultMessage="历史数据" />
        </Button>,
        <Button
          key="waring-setting"
          onClick={async () => {
            history.push({
              pathname:
                '/data/signal-waring?id=' +
                record.ID +
                '&protocol=' +
                record.protocol +
                '&client_id=' +
                record.device_uid,
            });
          }}
        >
          <FormattedMessage id="pages.waring-setting" defaultMessage="报警配置" />
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
        onReset={() => {
          setSearchProtocol('MQTT');
          setSearchDeviceUid('');
        }}
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
        request={async (params, sorter, filter) => {
          if (searchProtocol) {
            params.protocol = searchProtocol;
          }
          if (searchDeviceUid) {
            params.device_uid = Number(searchDeviceUid);
          }
          return signalPage(params);
        }}
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
            TCP: { text: 'TCP', status: 'success' },
            COAP: { text: 'COAP', status: 'success' },
          }}
          key={'protocol'}
          label={<FormattedMessage id="pages.signal.protocol" />}
          name="protocol"
        />
        <ProFormSelect
          multiple={false}
          dependencies={['protocol']}
          onChange={async (value) => {
            form.setFieldValue('device_uid', value);
            debugger;
          }}
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
          fieldProps={{
            onClick: (v) => {
              console.log(v);
              debugger;
            },
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

      <HistoryList
        updateModalOpen={historyModalOpen}
        onCancel={() => {
          handleHistoryModalOpen(false);
        }}
        values={currentRow || {}}
      />
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
        {currentRow?.ID && (
          <ProDescriptions<API.SignalListItem>
            column={2}
            title={currentRow?.name}
            request={async () => ({
              data: currentRow || {},
            })}
            columns={columns as ProDescriptionsItemProps<API.SignalListItem>[]}
          />
        )}
      </Drawer>
    </PageContainer>
  );
};

export default Admin;
