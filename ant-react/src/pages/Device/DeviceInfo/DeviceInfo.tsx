import DeviceInfoUpdateForm from '@/pages/Device/DeviceInfo/DeviceInfoUpdateForm';
import {
  addDevice,
  deleteDeviceInfo,
  devicePage,
  productList,
  updateDeviceInfo,
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
  ProFormDatePicker,
  ProFormDigit,
  ProFormSelect,
  ProFormText,
  ProTable,
} from '@ant-design/pro-components';
import { useIntl } from '@umijs/max';
import { Button, Drawer, Form, message, Popconfirm } from 'antd';
import dayjs from 'dayjs';
import React, { useRef, useState } from 'react';

async function handleAdd(fields: API.DeviceInfoItem) {
  const hide = message.loading('正在添加');
  try {
    let req = { ...fields };
    req.source = Number(req.source);
    let newVar = await addDevice(req);
    if (newVar.code === 20000) {
      hide();
      message.success('Added successfully');
      return true;
    } else {
      message.error(newVar.message);
    }
  } catch (error) {
    hide();
    message.error('Adding failed, please try again!');
    return false;
  }
}

async function handleUpdate(fields: API.DeviceInfoItem) {
  const hide = message.loading('正在更新');
  try {
    let req = { ...fields };
    req.manufacturing_date = dayjs(req.manufacturing_date).format('YYYY-MM-DDTHH:mm:ssZ');
    req.procurement_date = dayjs(req.procurement_date).format('YYYY-MM-DDTHH:mm:ssZ');
    req.warranty_expiry = dayjs(req.warranty_expiry).format('YYYY-MM-DDTHH:mm:ssZ');
    req.source = Number(req.source);

    await updateDeviceInfo(req);
    hide();
    message.success('更新成功');
    return true;
  } catch (error) {
    hide();
    message.error('更新 failed, please try again!');
    return false;
  }
}

async function handleRemove(ID: number | undefined) {
  const hide = message.loading('正在删除');
  try {
    await deleteDeviceInfo(ID);
    hide();
    message.success('删除 successfully');
    return true;
  } catch (error) {
    hide();
    message.error('Delete failed, please try again!');
    return false;
  }
}

const Admin: React.FC = () => {
  const intl = useIntl();
  /**
   * @en-US Pop-up window of new window
   * @zh-CN 新建窗口的弹窗
   *  */
  const [createModalOpen, handleModalOpen] = useState<boolean>(false);
  const [sourceValue, setSourceValue] = React.useState(1);
  /**
   * @en-US The pop-up window of the distribution update window
   * @zh-CN 分布更新窗口的弹窗
   * */
  const [updateModalOpen, handleUpdateModalOpen] = useState<boolean>(false);

  const [showDetail, setShowDetail] = useState<boolean>(false);

  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<API.DeviceInfoItem>();

  const columns: ProColumns<API.DeviceInfoItem>[] = [
    {
      title: <FormattedMessage id="pages.device-info.id" defaultMessage="唯一码" />,
      dataIndex: 'ID', // @ts-ignore
      hideInSearch: true,
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
      title: <FormattedMessage id="pages.device-info.product_id" />,
      hideInSearch: true,
      dataIndex: 'product_id',
      valueType: 'select',
      request: async () => {
        let resp = await productList();
        return resp.data;
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
      title: <FormattedMessage id="pages.device-info.sn" />,
      hideInSearch: false,
      dataIndex: 'sn',
    },
    {
      title: <FormattedMessage id="pages.device-info.manufacturing_date" />,
      hideInSearch: true,
      dataIndex: 'manufacturing_date',
      render: (_, record) => {
        return dayjs(record.manufacturing_date).format('YYYY-MM-DD');
      },
    },
    {
      title: <FormattedMessage id="pages.device-info.procurement_date" />,
      hideInSearch: true,
      dataIndex: 'procurement_date',
      render: (_, record) => {
        if (record.procurement_date) {
          return dayjs(record.procurement_date).format('YYYY-MM-DD');
        } else {
          return '-';
        }
      },
    },
    {
      title: <FormattedMessage id="pages.device-info.source" />,
      hideInSearch: true,
      dataIndex: 'source',
      valueType: 'select',
      valueEnum: {
        1: { text: '自产', status: 'success' },
        2: { text: '采购', status: 'success' },
      },
    },
    {
      title: <FormattedMessage id="pages.device-info.warranty_expiry" />,
      hideInSearch: true,
      dataIndex: 'warranty_expiry',
      render: (_, record) => {
        return dayjs(record.warranty_expiry).format('YYYY-MM-DD');
      },
    },
    {
      title: <FormattedMessage id="pages.device-info.push_interval" />,
      hideInSearch: true,
      dataIndex: 'push_interval',
      valueType: 'digit',
    },
    {
      title: <FormattedMessage id="pages.device-info.error_rate" />,
      hideInSearch: true,
      valueType: 'digit',
      dataIndex: 'error_rate',
    },
    {
      title: <FormattedMessage id="pages.device-info.protocol" />,
      hideInSearch: false,
      valueType: 'select',
      valueEnum: {
        MQTT: { text: 'MQTT', status: 'success' },
        HTTP: { text: 'HTTP', status: 'success' },
        WebSocket: { text: 'WebSocket', status: 'success' },
        TCP: { text: 'TCP', status: 'success' },
        COAP: { text: 'COAP', status: 'success' },
      },
      dataIndex: 'protocol',
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
          title={
            <FormattedMessage
              id="pages.deleteConfirm"
              defaultMessage="Are you sure to delete this record?"
            />
          }
          onConfirm={async () => {
            const success = await handleRemove(record.ID);
            if (success && actionRef.current) {
              actionRef.current.reload();
            }
          }}
          okText={<FormattedMessage id="pages.yes" defaultMessage="Yes" />}
          cancelText={<FormattedMessage id="pages.no" defaultMessage="No" />}
        >
          <Button danger>
            <FormattedMessage id="pages.deleted" defaultMessage="删除" />
          </Button>
        </Popconfirm>,
      ],
    },
  ];
  const [form] = Form.useForm();

  return (
    <PageContainer>
      <ProTable<API.DeviceInfoItem, API.PageParams>
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
        request={devicePage}
        columns={columns}
        rowSelection={{
          onChange: (_, selectedRows) => {
            setSelectedRows(selectedRows);
          },
        }}
      />
      <ModalForm
        form={form}
        title={intl.formatMessage({
          id: 'pages.searchTable.createForm.newRule',
          defaultMessage: 'New rule',
        })}
        width="75%"
        open={createModalOpen}
        onOpenChange={handleModalOpen}
        modalProps={{
          destroyOnClose: true,
          onCancel: () => {
            setSourceValue(0);
          },
        }}
        onFinish={async (value) => {
          const success = await handleAdd(value as API.DeviceInfoItem);
          if (success) {
            handleModalOpen(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
      >
        <ProFormSelect
          request={async () => {
            let resp = await productList();
            return resp.data;
          }}
          fieldProps={{
            showSearch: true,
            allowClear: false,
            fieldNames: {
              label: 'name',
              value: 'ID',
            },
          }}
          label={<FormattedMessage id="pages.device-info.product_id" />}
          name="product_id"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />
        <ProFormText
          label={<FormattedMessage id="pages.device-info.sn" />}
          name="sn"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />

        <ProFormSelect
          valueEnum={{
            1: { text: '自产', status: 'success' },
            2: { text: '采购', status: 'success' },
          }}
          onChange={(value) => {
            console.log('source change', value);
            setSourceValue(value); // 更新 sourceValue 状态
          }}
          label={<FormattedMessage id="pages.device-info.source" />}
          name="source"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />

        {sourceValue === '1' && (
          <ProFormDatePicker
            initialValue={dayjs().format('YYYY-MM-DDTHH:mm:ssZ')}
            transform={(value) => {
              return dayjs(value).format('YYYY-MM-DDTHH:mm:ssZ');
            }}
            label={<FormattedMessage id="pages.device-info.manufacturing_date" />}
            name="manufacturing_date"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.select" />,
              },
            ]}
          />
        )}
        {sourceValue === '2' && (
          <ProFormDatePicker
            initialValue={dayjs().format('YYYY-MM-DDTHH:mm:ssZ')}
            transform={(value) => {
              return dayjs(value).format('YYYY-MM-DDTHH:mm:ssZ');
            }}
            label={<FormattedMessage id="pages.device-info.procurement_date" />}
            name="procurement_date"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.select" />,
              },
            ]}
          />
        )}
        {sourceValue === '1' && (
          <ProFormDatePicker
            initialValue={dayjs().format('YYYY-MM-DDTHH:mm:ssZ')}
            transform={(value) => {
              return dayjs(value).format('YYYY-MM-DDTHH:mm:ssZ');
            }}
            label={<FormattedMessage id="pages.device-info.warranty_expiry" />}
            name="warranty_expiry"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.select" />,
              },
            ]}
          />
        )}
        <ProFormDigit
          label={<FormattedMessage id="pages.device-info.push_interval" />}
          name="push_interval"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormDigit
          label={<FormattedMessage id="pages.device-info.error_rate" />}
          name="error_rate"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormSelect
          valueEnum={{
            MQTT: { text: 'MQTT', status: 'success' },
            HTTP: { text: 'HTTP', status: 'success' },
            WebSocket: { text: 'WebSocket', status: 'success' },
            TCP: { text: 'TCP', status: 'success' },
            COAP: { text: 'COAP', status: 'success' },
          }}
          label={<FormattedMessage id="pages.device-info.protocol" />}
          name="protocol"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />
      </ModalForm>
      <DeviceInfoUpdateForm
        onSubmit={async (value) => {
          const success = await handleUpdate(value);
          if (success) {
            handleUpdateModalOpen(false);
            setCurrentRow(undefined);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
        onCancel={() => {
          handleUpdateModalOpen(false);
          if (!showDetail) {
            setCurrentRow(undefined);
          }
        }}
        updateModalOpen={updateModalOpen}
        values={currentRow || {}}
      />

      <Drawer
        width={600}
        open={showDetail}
        onClose={() => {
          setCurrentRow(undefined);
          setShowDetail(false);
        }}
        closable={false}
      >
        {currentRow?.ID && (
          <ProDescriptions<API.DeviceInfoItem>
            column={2}
            title={currentRow?.sn}
            request={async () => ({
              data: currentRow || {},
            })}
            params={{
              id: currentRow?.sn,
            }}
            columns={columns as ProDescriptionsItemProps<API.DeviceInfoItem>[]}
          />
        )}
      </Drawer>
    </PageContainer>
  );
};

export default Admin;
