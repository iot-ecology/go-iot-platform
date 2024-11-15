import DeviceUidShow from '@/pages/Data/Signal/DeviceUidShow';
import {
  addInfluxdbTransmitBind,
  deleteInfluxdbTransmitBind,
  deviceList,
  InfluxdbTransmitBindPage,
  InfluxdbTransmitList,
  mqttList,
  updateInfluxdbTransmitBind,
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
  ProFormRadio,
  ProFormSelect,
  ProFormText,
  ProTable,
} from '@ant-design/pro-components';
import { useIntl } from '@umijs/max';
import { Button, Drawer, Form, message } from 'antd';
import React, { useRef, useState } from 'react';
import InfluxdbTransmitBindUpdateForm from './InfluxdbTransmitBindUpdateForm';

const handleAdd = async (fields: API.InfluxdbTransmitBindListItem) => {
  const hide = message.loading('正在添加');
  try {
    await addInfluxdbTransmitBind({ ...fields });
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
    await deleteInfluxdbTransmitBind(id);
    hide();
    message.success('删除 successfully');
    return true;
  } catch (error) {
    hide();
    message.error('Delete failed, please try again!');
    return false;
  }
};

const handlerUpdate = async (fields: API.InfluxdbTransmitBindListItem) => {
  const hide = message.loading('正在更新');
  try {
    await updateInfluxdbTransmitBind({ ...fields });
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
  const [currentRow, setCurrentRow] = useState<API.InfluxdbTransmitBindListItem>();

  const columns: ProColumns<API.InfluxdbTransmitBindListItem>[] = [
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
      key: 'device_uid',
      title: <FormattedMessage id="pages.CassandraTransmitBind.device_uid" />,
      hideInSearch: true,
      dataIndex: 'device_uid',
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
      key: 'protocol',
      title: <FormattedMessage id="pages.CassandraTransmitBind.protocol" />,
      hideInSearch: true,
      dataIndex: 'protocol',
      valueType: 'select',
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
      title: <FormattedMessage id="pages.CassandraTransmitBind.identification_code" />,
      hideInSearch: true,
      dataIndex: 'identification_code',
    },
    {
      key: 'influxdb_transmit_id',
      title: <FormattedMessage id="pages.CassandraTransmitBind.cassandra_transmit_id" />,
      hideInSearch: true,
      dataIndex: 'influxdb_transmit_id',
    },
    {
      key: 'org',
      title: <FormattedMessage id="pages.CassandraTransmitBind.org" />,
      hideInSearch: true,
      dataIndex: 'org',
    },
    {
      key: 'bucket',
      title: <FormattedMessage id="pages.CassandraTransmitBind.bucket" />,
      hideInSearch: true,
      dataIndex: 'bucket',
    },

    {
      key: 'measurement',
      title: <FormattedMessage id="pages.CassandraTransmitBind.measurement" />,
      hideInSearch: true,
      dataIndex: 'measurement',
    },

    {
      key: 'script',
      title: <FormattedMessage id="pages.CassandraTransmitBind.script" />,
      hideInSearch: true,
      dataIndex: 'script',
      valueType: 'code',
    },
    {
      key: 'enable',
      title: <FormattedMessage id="pages.CassandraTransmitBind.enable" />,
      hideInSearch: true,
      dataIndex: 'enable',
      render: (dom, entity) => {
        return <>{entity.enable ? '启用' : '禁用'}</>;
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
  const [form] = Form.useForm();

  return (
    <PageContainer>
      <ProTable<API.InfluxdbTransmitBindListItem, API.PageParams>
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
        request={InfluxdbTransmitBindPage}
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
        modalProps={{
          destroyOnClose: true,
        }}
        onFinish={async (value) => {
          const success = await handleAdd(value as API.InfluxdbTransmitBindListItem);
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
          label={<FormattedMessage id="pages.CassandraTransmitBind.protocol" />}
          name="protocol"
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
          onChange={async (value) => {
            form.setFieldValue('device_uid', value);
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
            },
            showSearch: true,
            allowClear: false,
            fieldNames: {
              label: 'client_id',
              value: 'ID',
            },
          }}
          key={'device_uid'}
          label={<FormattedMessage id="pages.CassandraTransmitBind.device_uid" />}
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
          label={<FormattedMessage id="pages.CassandraTransmitBind.identification_code" />}
          name="identification_code"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormSelect
          request={async () => {
            let r = await InfluxdbTransmitList();
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
          key={'influxdb_transmit_id'}
          label={<FormattedMessage id="pages.CassandraTransmitBind.cassandra_transmit_id" />}
          name="influxdb_transmit_id"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />
        <ProFormText
          key={'org'}
          label={<FormattedMessage id="pages.CassandraTransmitBind.org" />}
          name="org"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText
          key={'bucket'}
          label={<FormattedMessage id="pages.CassandraTransmitBind.bucket" />}
          name="bucket"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText
          key={'measurement'}
          label={<FormattedMessage id="pages.CassandraTransmitBind.measurement" />}
          name="measurement"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText
          key={'script'}
          label={<FormattedMessage id="pages.CassandraTransmitBind.script" />}
          name="script"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />

        <ProFormRadio.Group
          name="enable"
          label={<FormattedMessage id="pages.CassandraTransmitBind.enable" />}
          options={[
            {
              label: '启用',
              value: true,
            },
            {
              label: '停用',
              value: false,
            },
          ]}
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />
      </ModalForm>

      <InfluxdbTransmitBindUpdateForm
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
          <ProDescriptions<API.InfluxdbTransmitBindListItem>
            column={2}
            title={currentRow?.database}
            request={async () => ({
              data: currentRow || {},
            })}
            columns={columns as ProDescriptionsItemProps<API.InfluxdbTransmitBindListItem>[]}
          />
        )}
      </Drawer>
    </PageContainer>
  );
};

export default Admin;
