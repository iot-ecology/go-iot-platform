import {
  addShipmentRecord,
  deleteShipmentRecord,
  FindByShipmentProductDetail,
  productList,
  shipmentPage,
} from '@/services/ant-design-pro/api';
import { FormattedMessage } from '@@/exports';
import { PlusOutlined } from '@ant-design/icons';
import {
  type ActionType,
  EditableProTable,
  ModalForm,
  PageContainer,
  type ProColumns,
  ProDescriptions,
  type ProDescriptionsItemProps,
  ProForm,
  ProFormDatePicker,
  ProFormInstance,
  ProFormSelect,
  ProFormText,
  ProTable,
} from '@ant-design/pro-components';
import { useIntl } from '@umijs/max';
import { Button, Divider, Drawer, message } from 'antd';
import dayjs from 'dayjs';
import React, { useRef, useState } from 'react';

const handleAdd = async (fields: API.ShipmentRecordListItem) => {
  const hide = message.loading('正在添加');
  try {
    let req = { ...fields };
    if (!req.product_plans) {
      req.product_plans = [];
    }
    let newVar = await addShipmentRecord(req);
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
};

const handleRemove = async (id: any) => {
  const hide = message.loading('正在删除');
  try {
    await deleteShipmentRecord(id);
    hide();
    message.success('删除 successfully');
    return true;
  } catch (error) {
    hide();
    message.error('Delete failed, please try again!');
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

  const [showDetail, setShowDetail] = useState<boolean>(false);

  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<API.ShipmentRecordListItem>();
  const [dataSource, setDataSource] = useState<API.RecordInfo>([]);

  const pp: ProColumns<API.ProductPlanCreateParam>[] = [
    {
      key: 'product_id',
      title: '产品',
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
      key: 'quantity',
      title: '发货数量',
      hideInSearch: true,
      valueType: 'digit',
      dataIndex: 'quantity',
    },
  ];
  const columns: ProColumns<API.ShipmentRecordListItem>[] = [
    {
      key: 'ID',
      title: <FormattedMessage id="pages.id" defaultMessage="唯一码" />,
      hideInSearch: true,
      dataIndex: 'ID', // @ts-ignore

      render: (dom, entity) => {
        return (
          <a
            onClick={async () => {
              let newVar = await FindByShipmentProductDetail(entity.ID);
              entity.product_plans = newVar.data;

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
      key: 'shipment_date',
      title: <FormattedMessage id="pages.shipment.shipment_date" />,
      hideInSearch: true,
      dataIndex: 'shipment_date',
      render: (_, record) => {
        return dayjs(record.shipment_date).format('YYYY-MM-DD');
      },
    },
    {
      key: 'technician',
      title: <FormattedMessage id="pages.shipment.technician" />,
      hideInSearch: true,
      dataIndex: 'technician',
    },
    {
      key: 'customer_name',
      title: <FormattedMessage id="pages.shipment.customer_name" />,
      hideInSearch: false,
      dataIndex: 'customer_name',
    },
    {
      key: 'customer_phone',
      title: <FormattedMessage id="pages.shipment.customer_phone" />,
      hideInSearch: true,
      dataIndex: 'customer_phone',
    },
    {
      key: 'customer_address',
      title: <FormattedMessage id="pages.shipment.customer_address" />,
      hideInSearch: true,
      dataIndex: 'customer_address',
    },
    {
      key: 'tracking_number',
      title: <FormattedMessage id="pages.shipment.tracking_number" />,
      hideInSearch: true,
      dataIndex: 'tracking_number',
    },
    {
      key: 'status',
      title: <FormattedMessage id="pages.shipment.status" />,
      hideInSearch: false,
      dataIndex: 'status',
      valueType: 'select',

      valueEnum: {
        '1': { text: '已发货', status: 'default' },
        '2': { text: '已收货', status: 'success' },
        '3': { text: '未发货', status: 'error' },
      },
    },
    {
      key: 'description',
      title: <FormattedMessage id="pages.shipment.description" />,
      hideInSearch: true,
      dataIndex: 'description',
    },

    {
      title: <FormattedMessage id="pages.searchTable.titleOption" defaultMessage="Operating" />,
      dataIndex: 'option',
      valueType: 'option',
      render: (_, record) => [
        record.status !== '3' ? (
          <>
            <Button
              type={'primary'}
              onClick={async () => {
                let newVar = await FindByShipmentProductDetail(record.ID);
                record.product_plans = newVar.data;

                setCurrentRow(record);
                setShowDetail(true);
              }}
            >
              查看
            </Button>
          </>
        ) : (
          <Button
            key="deleted"
            onClick={async () => {
              let success = await handleRemove(record.ID);
              if (success) {
                if (actionRef.current) {
                  actionRef.current.reload();
                }
              }
            }}
            danger={true}
          >
            <FormattedMessage id="pages.deleted" defaultMessage="删除" />
          </Button>
        ),
      ],
    },
  ];
  const formRef = useRef<ProFormInstance>();
  const [editableKeys, setEditableRowKeys] = useState<React.Key[]>();
  return (
    <PageContainer>
      <ProTable<API.ShipmentRecordListItem, API.PageParams>
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
        request={shipmentPage}
        columns={columns}
      />
      <ModalForm
        formRef={formRef}
        key={'add'}
        title={intl.formatMessage({
          id: 'pages.new',
          defaultMessage: '新增',
        })}
        width="75%"
        open={createModalOpen}
        onOpenChange={handleModalOpen}
        onFinish={async (value) => {
          const success = await handleAdd(value as API.ShipmentRecordListItem);
          if (success) {
            handleModalOpen(false);
            if (actionRef.current) {
              await actionRef.current.reload();
            }
          }
        }}
      >
        <ProFormText
          key={'technician'}
          label={<FormattedMessage id="pages.shipment.technician" />}
          name="technician"
        />
        <ProFormText
          key={'customer_name'}
          label={<FormattedMessage id="pages.shipment.customer_name" />}
          name="customer_name"
        />
        <ProFormText
          key={'customer_phone'}
          label={<FormattedMessage id="pages.shipment.customer_phone" />}
          name="customer_phone"
        />
        <ProFormText
          key={'customer_address'}
          label={<FormattedMessage id="pages.shipment.customer_address" />}
          name="customer_address"
        />

        <ProFormSelect
          valueEnum={{
            1: { text: '已发货', status: 'default' },
            2: { text: '已收货', status: 'success' },
            3: { text: '未发货', status: 'error' },
          }}
          key={'status'}
          label={<FormattedMessage id="pages.shipment.status" />}
          name="status"
        />
        <ProFormText
          key={'description'}
          label={<FormattedMessage id="pages.shipment.description" />}
          name="description"
        />

        <ProFormDatePicker
          transform={(value) => {
            return dayjs(value).format('YYYY-MM-DDTHH:mm:ssZ');
          }}
          initialValue={dayjs().format('YYYY-MM-DDTHH:mm:ssZ')}
          key={'shipment_date'}
          label={<FormattedMessage id="pages.shipment.shipment_date" />}
          name="shipment_date"
        />

        <ProForm.Item
          label="数组数据"
          name="product_plans"
          initialValue={dataSource}
          trigger="onValuesChange"
        >
          <EditableProTable<API.ProductPlanCreateParam>
            rowKey="id"
            toolBarRender={false}
            columns={pp}
            recordCreatorProps={{
              newRecordType: 'dataSource',
              // @ts-ignore
              position: 'bottom',
              // @ts-ignore
              record: () => ({
                id: new Date(),
                product_id: null,
                quantity: null,
              }),
            }}
            editable={{
              type: 'multiple',
              editableKeys,
              onChange: setEditableRowKeys,
              actionRender: (row, _, dom) => {
                return [dom.delete];
              },
            }}
          />
        </ProForm.Item>
      </ModalForm>

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
          <>
            <ProDescriptions<API.ShipmentRecordListItem>
              column={2}
              title={currentRow?.customer_name}
              request={async () => ({
                data: currentRow || {},
              })}
              params={{
                id: currentRow?.customer_name,
              }}
              columns={columns as ProDescriptionsItemProps<API.ShipmentRecordListItem>[]}
            ></ProDescriptions>
            <Divider />
            <h3>{'发货详情'}</h3>
            <ProTable<API.ProductPlanCreateParam>
              dataSource={currentRow.product_plans}
              search={false}
              options={false}
              columns={pp}
            />
          </>
        )}
      </Drawer>
    </PageContainer>
  );
};

export default Admin;
