import ProductUpdateForm from '@/pages/Lifecycle/ProductList/ProductUpdateForm';
import RenderTags from '@/pages/Lifecycle/ProductList/TagsRow';
import {
  addProduct,
  deleteProduct,
  productPage,
  updateProduct,
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
  ProForm,
  ProFormDigit,
  ProFormText,
  ProFormUploadDragger,
  ProTable,
} from '@ant-design/pro-components';
import { useIntl } from '@umijs/max';
import { Button, Drawer, Form, Image, message } from 'antd';
import React, { useRef, useState } from 'react';

const handleAdd = async (fields: API.ProductItem) => {
  const hide = message.loading('正在添加');
  try {
    let image_string = '';
    if (typeof fields.image_url === 'string') {
    } else if (Array.isArray(fields.image_url)) {
      fields.image_url.forEach((file) => {
        image_string += file.response.file_path;
      });
    }
    // 拷贝
    const req = { ...fields };
    req.image_url = image_string;
    await addProduct(req);
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
    await deleteProduct(id);
    hide();
    message.success('删除 successfully');
    return true;
  } catch (error) {
    hide();
    message.error('Delete failed, please try again!');
    return false;
  }
};

const handlerUpdate = async (fields: API.ProductItem) => {
  const hide = message.loading('正在更新');
  try {
    console.log(fields);
    await updateProduct({ ...fields });
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
  const [currentRow, setCurrentRow] = useState<API.ProductItem>();

  const columns: ProColumns<API.ProductItem>[] = [
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
      key: 'pages.product.name',
      title: <FormattedMessage id="pages.product.name" />,
      hideInSearch: false,
      dataIndex: 'name',
    },
    {
      key: 'pages.product.description',
      title: <FormattedMessage id="pages.product.description" />,
      hideInSearch: false,
      dataIndex: 'description',
    },
    {
      key: 'pages.product.sku',
      title: <FormattedMessage id="pages.product.sku" />,
      hideInSearch: false,
      dataIndex: 'sku',
    },
    {
      key: 'pages.product.price',
      title: <FormattedMessage id="pages.product.price" />,
      hideInSearch: false,
      dataIndex: 'price',
    },
    {
      key: 'pages.product.cost',
      title: <FormattedMessage id="pages.product.cost" />,
      hideInSearch: false,
      dataIndex: 'cost',
    },
    {
      key: 'pages.product.quantity',
      title: <FormattedMessage id="pages.product.quantity" />,
      hideInSearch: false,
      dataIndex: 'quantity',
    },
    {
      key: 'pages.product.minimum_stock',
      title: <FormattedMessage id="pages.product.minimum_stock" />,
      hideInSearch: false,
      dataIndex: 'minimum_stock',
    },
    {
      key: 'pages.product.warranty_period',
      title: <FormattedMessage id="pages.product.warranty_period" />,
      hideInSearch: false,
      dataIndex: 'warranty_period',
    },
    {
      key: 'pages.product.status',
      title: <FormattedMessage id="pages.product.status" />,
      hideInSearch: false,
      dataIndex: 'status',
    },
    {
      key: 'pages.product.tags',
      title: <FormattedMessage id="pages.product.tags" />,
      hideInSearch: false,
      dataIndex: 'tags',
      render: (e, r) => {
        if (r.tags) {
          return RenderTags(r.tags);
        } else {
          return <div>-</div>;
        }
      },
    },
    {
      key: 'pages.product.image_url',
      title: <FormattedMessage id="pages.product.image_url" />,
      hideInSearch: false,
      dataIndex: 'image_url',
      render: (e, r) => {
        return <Image src={'http://localhost:8080' + '/file/download?path=' + r.image_url}></Image>;
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
                await actionRef.current.reload();
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
      <ProTable<API.ProductItem, API.PageParams>
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
        request={productPage}
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
          const success = await handleAdd(value as API.ProductItem);
          if (success) {
            handleModalOpen(false);
            if (actionRef.current) {
              await actionRef.current.reload();
            }
          }
        }}
      >
        <ProForm.Group>
          <ProFormText
            key={'name'}
            label={<FormattedMessage id="pages.product.name" />}
            name="name"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.input" />,
              },
            ]}
          />
          <ProFormText
            key={'description'}
            label={<FormattedMessage id="pages.product.description" />}
            name="description"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.input" />,
              },
            ]}
          />
          <ProFormText
            key={'sku'}
            label={<FormattedMessage id="pages.product.sku" />}
            name="sku"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.input" />,
              },
            ]}
          />
        </ProForm.Group>
        <ProForm.Group>
          <ProFormDigit
            key={'price'}
            label={<FormattedMessage id="pages.product.price" />}
            name="price"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.input" />,
              },
            ]}
          />
          <ProFormDigit
            key={'cost'}
            label={<FormattedMessage id="pages.product.cost" />}
            name="cost"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.input" />,
              },
            ]}
          />
          <ProFormDigit
            key={'quantity'}
            label={<FormattedMessage id="pages.product.quantity" />}
            name="quantity"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.input" />,
              },
            ]}
          />
          <ProFormDigit
            key={'minimum_stock'}
            label={<FormattedMessage id="pages.product.minimum_stock" />}
            name="minimum_stock"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.input" />,
              },
            ]}
          />
          <ProFormDigit
            key={'warranty_period'}
            label={<FormattedMessage id="pages.product.warranty_period" />}
            name="warranty_period"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.input" />,
              },
            ]}
          />
        </ProForm.Group>
        <ProFormText
          key={'status'}
          label={<FormattedMessage id="pages.product.status" />}
          name="status"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText
          key={'tags'}
          label={<FormattedMessage id="pages.product.tags" />}
          placeholder={'请输入,如果存在多个请用逗号分割'}
          name="tags"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormUploadDragger
          name="image_url"
          label={<FormattedMessage id="pages.product.image_url" />}
          fieldProps={{
            accept: '.jpg,.png', // 设置接受的文件类型
            name: 'file',
            action: '/api/file/update',
            multiple: false,
          }}
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />
      </ModalForm>

      <ProductUpdateForm
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
        {currentRow?.name && (
          <ProDescriptions<API.ProductItem>
            column={2}
            title={currentRow?.name}
            request={async () => ({
              data: currentRow || {},
            })}
            params={{
              id: currentRow?.name,
            }}
            columns={columns as ProDescriptionsItemProps<API.ProductItem>[]}
          />
        )}
      </Drawer>
    </PageContainer>
  );
};

export default Admin;
