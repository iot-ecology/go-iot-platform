import DeptUpdateForm from '@/pages/User/DeptList/DeptUpdateForm';
import { addDept, deleteDept, deptList, deptPage, updateDept } from '@/services/ant-design-pro/api';
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
import React, { useEffect, useRef, useState } from 'react';

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
  const [currentRow, setCurrentRow] = useState<API.DeptListItem>();

  const handleRemove = async (id: any) => {
    const hide = message.loading('正在删除');
    try {
      await deleteDept(id);
      hide();
      message.success('删除 successfully');
      return true;
    } catch (error) {
      hide();
      message.error('Delete failed, please try again!');
      return false;
    }
  };

  const columns: ProColumns<API.DeptListItem>[] = [
    {
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
      title: <FormattedMessage id="pages.name" defaultMessage="名称" />,
      hideInSearch: false,
      dataIndex: 'name',
    },
    {
      title: <FormattedMessage id="pages.dept-list.pname" defaultMessage="上级部门" />,
      hideInSearch: true,
      dataIndex: 'parent_name',
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
          title={<FormattedMessage id="pages.deleteConfirm" defaultMessage="确定要删除此项吗？" />}
          onConfirm={async () => {
            const success = await handleRemove(record.ID);
            if (success) {
              actionRef.current?.reload();
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

  async function handleAdd(value1: API.DeptListItem) {
    const hide = message.loading('正在添加');
    try {
      await addDept({ ...value1 });
      hide();
      message.success('Added successfully');
      return true;
    } catch (error) {
      hide();
      message.error('Adding failed, please try again!');
      return false;
    }
  }

  const [options, setOptions] = useState<{ label: string; value: string | number }[]>([]);

  const loadDeptList = async () => {
    try {
      const result = await deptList(); // 调用API获取数据
      if (result && result.data && Array.isArray(result.data)) {
        const formattedOptions = result.data.map((item) => ({
          label: item.name, // 假设每项数据都有一个name属性
          value: item.ID, // 假设每项数据都有一个id属性
        }));

        // @ts-ignore
        setOptions(formattedOptions);
      }
    } catch (error) {
      message.error('请求数据失败，请稍后再试！');
    }
  };

  useEffect(() => {
    loadDeptList();
  }, []);

  async function handlerUpdate(value: API.DeptListItem) {
    const hide = message.loading('正在更新');
    try {
      await updateDept({ ...value });
      hide();
      message.success('更新成功');
      return true;
    } catch (error) {
      hide();
      message.error('更新 failed, please try again!');
      return false;
    }
  }

  return (
    <PageContainer>
      <ProTable<API.DeptListItem, API.PageParams>
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
        request={deptPage}
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
          const success = await handleAdd(value as API.DeptListItem);
          if (success) {
            await loadDeptList();
            handleModalOpen(false);
            if (actionRef.current) {
              actionRef.current.reload();
            }
          }
        }}
      >
        <ProFormText
          key={'name'}
          label={<FormattedMessage id="pages.name" />}
          name="name"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />

        <ProFormSelect
          key={'parent_id'}
          label={<FormattedMessage id="pages.dept-top" />}
          name="parent_id"
          options={options}
          rules={[
            {
              required: false,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />
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
        {currentRow?.name && (
          <ProDescriptions<API.DeptListItem>
            column={2}
            title={currentRow?.name}
            request={async () => ({
              data: currentRow || {},
            })}
            params={{
              id: currentRow?.ID,
            }}
            columns={columns as ProDescriptionsItemProps<API.DeptListItem>[]}
          />
        )}
      </Drawer>

      <DeptUpdateForm
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
              actionRef.current.reload();
            }
          }
        }}
      />
    </PageContainer>
  );
};

export default Admin;
