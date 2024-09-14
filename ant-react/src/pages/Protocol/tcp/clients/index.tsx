import TcpHandlerUpdateForm from '@/pages/Protocol/tcp/clients/TcpHandlerUpdateForm';
import {
  addTcpHandler,
  deltedTcpHandler,
  deviceList,
  mqttScriptCheck,
  TcpHandlerPage,
  updateTcpHandler,
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
  ProFormTextArea,
  ProTable,
} from '@ant-design/pro-components';
import { useIntl } from '@umijs/max';
import { Button, Drawer, Form, message } from 'antd';
import React, { useRef, useState } from 'react';

const handleAdd = async (fields: API.TcpHandlerListItem) => {
  const hide = message.loading('正在添加');
  try {
    await addTcpHandler({ ...fields });
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
    await deltedTcpHandler(id);
    hide();
    message.success('删除 successfully');
    return true;
  } catch (error) {
    hide();
    message.error('Delete failed, please try again!');
    return false;
  }
};

const handlerUpdate = async (fields: API.TcpHandlerListItem) => {
  const hide = message.loading('正在更新');
  try {
    await updateTcpHandler({ ...fields });
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
  const [currentRow, setCurrentRow] = useState<API.TcpHandlerListItem>();
  const [setScriptModalOpen, handleSetScriptModalOpen] = useState<boolean>(false);
  const [createScript, handlerCreateScript] = useState<string>(false);

  const columns: ProColumns<API.TcpHandlerListItem>[] = [
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
      key: 'device_info_id',
      title: <FormattedMessage id="pages.http.device_info_id" />,
      hideInSearch: false,
      dataIndex: 'device_info_id',
      valueType: 'select',
      request: async () => {
        let a = await deviceList();
        return a.data.filter((e) => {
          return e.protocol == 'TCP';
        });
      },
      fieldProps: {
        showSearch: true,
        fieldNames: {
          label: 'sn',
          value: 'ID',
        },
      },
    },
    {
      key: 'name',
      title: <FormattedMessage id="pages.http.name" />,
      hideInSearch: true,
      dataIndex: 'name',
    },
    {
      key: 'username',
      title: <FormattedMessage id="pages.http.username" />,
      hideInSearch: true,
      dataIndex: 'username',
    },
    {
      key: 'password',
      title: <FormattedMessage id="pages.http.password" />,
      hideInSearch: true,
      valueType: 'password',
      dataIndex: 'password',
    },
    {
      key: 'script',
      title: <FormattedMessage id="pages.http.script" />,
      hideInSearch: true,
      valueType: 'code',
      dataIndex: 'script',
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
        <Button
          key="check-script"
          onClick={async () => {
            setCurrentRow(record);
            handleSetScriptModalOpen(true);
          }}
        >
          <FormattedMessage id="pages.check-script" defaultMessage="检查脚本" />
        </Button>,
      ],
    },
  ];
  const [form] = Form.useForm();

  return (
    <PageContainer>
      <ProTable<API.TcpHandlerListItem, API.PageParams>
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
        request={TcpHandlerPage}
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
          const success = await handleAdd(value as API.TcpHandlerListItem);
          if (success) {
            handleModalOpen(false);
            if (actionRef.current) {
              await actionRef.current.reload();
            }
          }
        }}
      >
        <ProFormSelect
          key={'device_info_id'}
          label={<FormattedMessage id="pages.http.device_info_id" />}
          name="device_info_id"
          request={async () => {
            let a = await deviceList();
            return a.data.filter((e) => {
              return e.protocol == 'TCP';
            });
          }}
          onChange={(v) => {
            var s =
              'function main(nc) {\n' +
              '    var dataRows = [\n' +
              '        { "Name": "Temperature", "Value": "23" },\n' +
              '        { "Name": "Humidity", "Value": "30" },\n' +
              '        { "Name": "A", "Value": nc },\n' +
              '    ];\n' +
              '    var result = {\n' +
              '        "Time":  Math.floor(Date.now() / 1000),\n' +
              '        "DataRows": dataRows,\n' +
              '        "IdentificationCode": "' +
              v +
              '",\n' +
              '        "DeviceUid": "' +
              v +
              '",\n' +
              '        "Nc": nc\n' +
              '    };\n' +
              '    return [result];\n' +
              '}';
            handlerCreateScript(s);

            form.setFieldValue('script', s);
          }}
          fieldProps={{
            showSearch: true,
            allowClear: false,
            fieldNames: {
              label: 'sn',
              value: 'ID',
            },
          }}
        />
        <ProFormText key={'name'} label={<FormattedMessage id="pages.http.name" />} name="name" />
        <ProFormText
          key={'username'}
          label={<FormattedMessage id="pages.http.username" />}
          name="username"
        />
        <ProFormText.Password
          key={'password'}
          label={<FormattedMessage id="pages.http.password" />}
          name="password"
        />
        <ProFormTextArea
          initialValue={createScript}
          fieldProps={{
            value: createScript,
          }}
          key={'script'}
          label={<FormattedMessage id="pages.http.script" />}
          name="script"
        />
      </ModalForm>

      <ModalForm
        title={'验证脚本'}
        open={setScriptModalOpen}
        onOpenChange={handleSetScriptModalOpen}
        onFinish={async (value) => {
          console.log('设置脚本参数');
        }}
        submitter={{
          render: (props, defaultDoms) => {
            return [
              <Button key="ok">取消</Button>,
              <Button
                key="check"
                onClick={async () => {
                  let newVar = await mqttScriptCheck(
                    props.form?.getFieldValue('mock-param'),
                    currentRow?.script,
                  );
                  if (newVar.code === 20000) {
                    props.form?.setFieldValue('mock-result', JSON.stringify(newVar.data));
                    message.success('模拟执行成功');
                  } else {
                    message.error('模拟执行失败，请检查脚本');
                  }
                }}
              >
                验证
              </Button>,
            ];
          },
        }}
      >
        <ProFormTextArea
          label={<FormattedMessage id="pages.mock-param" />}
          name={'mock-param'}
          key={'mock-param'}
        ></ProFormTextArea>
        <ProFormTextArea
          label={<FormattedMessage id={'pages.mock-result'} />}
          name={'mock-result'}
          key={'mock-result'}
        ></ProFormTextArea>
      </ModalForm>
      <TcpHandlerUpdateForm
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
          <ProDescriptions<API.TcpHandlerListItem>
            column={2}
            title={currentRow?.ID}
            request={async () => ({
              data: currentRow || {},
            })}
            columns={columns as ProDescriptionsItemProps<API.TcpHandlerListItem>[]}
          />
        )}
      </Drawer>
    </PageContainer>
  );
};

export default Admin;
