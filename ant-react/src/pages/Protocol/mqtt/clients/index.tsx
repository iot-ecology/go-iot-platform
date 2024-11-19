import MqttUpdateForm from '@/pages/Protocol/mqtt/clients/MqttUpdateForm';
import {
  addMQTT,
  deleteMQTT,
  mqttPage,
  mqttScriptCheck,
  sendMqttClientMessage,
  setMqttScriptCheck,
  startMqttClient,
  stopMqttClient,
  updateMQTT,
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
  ProFormDigit,
  ProFormRadio,
  ProFormSelect,
  ProFormText,
  ProFormTextArea,
  ProTable,
  TableDropdown,
} from '@ant-design/pro-components';
import { useIntl } from '@umijs/max';
import { Button, Drawer, message, Popconfirm, Tag } from 'antd';
import React, { useRef, useState } from 'react';
import { history } from 'umi';

function generateUUID() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
    let r = (Math.random() * 16) | 0,
      v = c == 'x' ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

const handleAdd = async (fields: API.MqttListItem) => {
  const hide = message.loading('正在添加');
  try {
    await addMQTT({ ...fields });
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
    await deleteMQTT(id);
    hide();
    message.success('删除 successfully');
    return true;
  } catch (error) {
    hide();
    message.error('Delete failed, please try again!');
    return false;
  }
};

const handlerUpdate = async (fields: API.MqttListItem) => {
  const hide = message.loading('正在更新');
  try {
    await updateMQTT({ ...fields });
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
  const [sendMqttMessage, setSendMqttMessage] = useState<boolean>(false);

  const [setScriptModalOpen, handleSetScriptModalOpen] = useState<boolean>(false);

  const actionRef = useRef<ActionType>();
  const [currentRow, setCurrentRow] = useState<API.MqttListItem>();

  const [curMqttClientId, setCurMqttClientId] = useState<number>(-1);
  const [canSubmmitScript, setSubmmitScript] = useState<boolean>(false);

  const handleSignalConfig = (record: any) => {
    history.push({
      pathname: '/data/signal?id=' + record.ID + '&protocol=MQTT',
    });
  };

  const handleMockSend = (record: any) => {
    setCurMqttClientId(record.ID);
    setSendMqttMessage(true);
    setCurrentRow(record);
  };

  const handleSetScript = (record: any) => {
    setCurMqttClientId(record.ID);
    handleSetScriptModalOpen(true);
  };

  const handleStart = async (record: any) => {
    let newVar = await startMqttClient(record.ID);
    if (newVar.code === 20000) {
      await actionRef.current?.reload();
      message.success('启动成功');
    } else {
      message.error('启动失败 ' + newVar.message);
    }
  };

  const handleStop = async (record: any) => {
    let newVar = await stopMqttClient(record.ID);
    if (newVar.code == 20000) {
      await actionRef.current?.reload();
      message.success('停止成功');
    } else {
      message.error('停止失败 ' + newVar.message);
    }
  };

  const columns: ProColumns<API.MqttListItem>[] = [
    {
      key: 'ID',
      title: <FormattedMessage id="pages.id" defaultMessage="唯一码" />,
      hideInSearch: true,
      dataIndex: 'ID', // @ts-ignore
      width: '100px',
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
      key: 'host',
      title: <FormattedMessage id="pages.mqtt.host" />,
      hideInSearch: true,
      dataIndex: 'host',
    },

    {
      key: 'port',
      title: <FormattedMessage id="pages.mqtt.port" />,
      hideInSearch: true,
      dataIndex: 'port',
    },

    {
      key: 'client_id',
      title: <FormattedMessage id="pages.mqtt.client_id" />,
      hideInSearch: false,
      dataIndex: 'client_id',
    },

    {
      key: 'username',
      title: <FormattedMessage id="pages.mqtt.username" />,
      hideInSearch: true,
      dataIndex: 'username',
    },

    {
      key: 'password',
      title: <FormattedMessage id="pages.mqtt.password" />,
      hideInSearch: true,
      valueType: 'password',
      dataIndex: 'password',
    },

    {
      key: 'subtopic',
      title: <FormattedMessage id="pages.mqtt.subtopic" />,
      hideInSearch: true,
      dataIndex: 'subtopic',
    },

    {
      key: 'start',
      title: <FormattedMessage id="pages.mqtt.start" />,
      hideInSearch: true,
      dataIndex: 'start',
      width: '100px',
      render: (_, record) => {
        if (record.start) {
          return <Tag color="green">启动</Tag>;
        } else {
          return <Tag color="red">停止</Tag>;
        }
      },
    },

    {
      key: 'script',
      title: <FormattedMessage id="pages.mqtt.script" />,
      hideInSearch: true,
      hideInTable: false,
      dataIndex: 'script',
      valueType: 'code',
    },

    {
      title: <FormattedMessage id="pages.searchTable.titleOption" defaultMessage="Operating" />,
      dataIndex: 'option',
      valueType: 'option',
      width: '150px',
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
        record.start === false && (
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
          </Popconfirm>
        ),
        <TableDropdown
          key="actionGroup"
          onSelect={(e) => {
            if (e === 'signal-config') {
              handleSignalConfig(record);
            }
            if (e === 'mock-send') {
              handleMockSend(record);
            }
            if (e === 'set-script') {
              handleSetScript(record);
            }
            if (e === 'start') {
              handleStart(record);
            }
            if (e === 'end') {
              handleStop(record);
            }
          }}
          menus={[
            {
              key: 'signal-config',
              name: <FormattedMessage id="pages.config-signal" defaultMessage="信号配置" />,
            },
            {
              key: 'mock-send',
              name: <FormattedMessage id="pages.mock-send" defaultMessage="模拟发送" />,
            },
            {
              key: 'set-script',
              name: record.start === false && (
                <FormattedMessage id="pages.set-script" defaultMessage="设置脚本" />
              ),
            },
            {
              key: 'start',
              name: record.start === false && (
                <FormattedMessage id="pages.start" defaultMessage="启动" />
              ),
            },
            {
              key: 'stop',
              name: record.start === true && (
                <FormattedMessage id="pages.stop" defaultMessage="停止" />
              ),
            },
          ].filter(({ name }) => !!name)}
        />,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<API.MqttListItem, API.PageParams>
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
        request={mqttPage}
        columns={columns}
      />

      <ModalForm
        title={'模拟发送'}
        open={sendMqttMessage}
        onOpenChange={setSendMqttMessage}
        modalProps={{
          destroyOnClose: true,
        }}
        onFinish={async (value) => {
          console.log('设置脚本参数');
        }}
        submitter={{
          render: (props, defaultDoms) => {
            return [
              <Button
                key="submmit"
                type={'primary'}
                onClick={async () => {
                  console.log(props.form?.getFieldsValue());
                  let fieldsValue = props.form?.getFieldsValue();
                  fieldsValue.qos = Number(fieldsValue.qos);
                  let newVar = await sendMqttClientMessage(currentRow?.client_id, fieldsValue);
                  if (newVar.code === 20000) {
                    setSendMqttMessage(false);
                    message.success('模拟发送成功');
                  }
                }}
              >
                发送
              </Button>,
            ];
          },
        }}
      >
        <ProFormText
          label={<FormattedMessage id={'pages.mqtt.payload'} />}
          name={'payload'}
          key={'payload'}
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />{' '}
        <ProFormSelect
          valueEnum={{
            0: { text: '最多发送一次，可能会丢失', status: 'success' },
            1: { text: '至少发送一次，可能会重复', status: 'success' },
            2: { text: '发送一次，确保消息不会重复', status: 'success' },
          }}
          label={<FormattedMessage id={'pages.mqtt.qos'} />}
          name={'qos'}
          key={'qos'}
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />{' '}
        <ProFormText
          label={<FormattedMessage id={'pages.mqtt.topic'} />}
          name={'topic'}
          key={'topic'}
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />{' '}
        <ProFormRadio.Group
          initialValue={false}
          options={[
            {
              label: '是',
              value: true,
            },
            { label: '否', value: false },
          ]}
          label={<FormattedMessage id={'pages.mqtt.retained'} />}
          name={'retained'}
          key={'retained'}
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />
      </ModalForm>

      <ModalForm
        title={'设置数据解析脚本'}
        open={setScriptModalOpen}
        onOpenChange={handleSetScriptModalOpen}
        modalProps={{
          destroyOnClose: true,
        }}
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
                    props.form?.getFieldValue('script'),
                  );
                  if (newVar.code === 20000) {
                    props.form?.setFieldValue('mock-result', JSON.stringify(newVar.data));
                    setSubmmitScript(true);
                    message.success('模拟执行成功');
                  } else {
                    message.error('模拟执行失败，请检查脚本');
                  }
                }}
              >
                验证
              </Button>,
              <Button
                disabled={!canSubmmitScript}
                key="submmit"
                type={'primary'}
                onClick={async () => {
                  let newVar = await setMqttScriptCheck(
                    curMqttClientId,
                    props.form?.getFieldValue('script'),
                  );
                  if (newVar.code === 20000) {
                    setSubmmitScript(false);
                    await actionRef.current?.reload();
                    handleSetScriptModalOpen(false);
                    message.success('数据解析脚本设置成功');
                  }
                }}
              >
                提交
              </Button>,
            ];
          },
        }}
      >
        <div>{curMqttClientId}</div>
        <ProFormTextArea
          initialValue={
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
            curMqttClientId +
            '",\n' +
            '        "DeviceUid": "' +
            curMqttClientId +
            '",\n' +
            '        "Nc": nc\n' +
            '    };\n' +
            '    return [result];\n' +
            '}'
          }
          label={<FormattedMessage id={'pages.script'} />}
          name={'script'}
          key={'script'}
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        ></ProFormTextArea>
        <ProFormTextArea
          label={<FormattedMessage id="pages.mock-param" />}
          name={'mock-param'}
          key={'mock-param'}
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        ></ProFormTextArea>
        <ProFormTextArea
          label={<FormattedMessage id={'pages.mock-result'} />}
          name={'mock-result'}
          key={'mock-result'}
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        ></ProFormTextArea>
      </ModalForm>
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
          const success = await handleAdd(value as API.MqttListItem);
          if (success) {
            handleModalOpen(false);
            if (actionRef.current) {
              await actionRef.current.reload();
            }
          }
        }}
      >
        <ProFormText
          initialValue={'127.0.0.1'}
          key={'host'}
          label={<FormattedMessage id="pages.mqtt.host" />}
          name="host"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormDigit
          initialValue={1883}
          key={'port'}
          label={<FormattedMessage id="pages.mqtt.port" />}
          name="port"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText
          initialValue={generateUUID()}
          key={'client_id'}
          label={<FormattedMessage id="pages.mqtt.client_id" />}
          name="client_id"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText
          initialValue={'admin'}
          key={'username'}
          label={<FormattedMessage id="pages.mqtt.username" />}
          name="username"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText.Password
          initialValue={'admin'}
          key={'password'}
          label={<FormattedMessage id="pages.mqtt.password" />}
          name="password"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText
          initialValue={'/test_topic/0'}
          key={'subtopic'}
          label={<FormattedMessage id="pages.mqtt.subtopic" />}
          name="subtopic"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
      </ModalForm>

      <MqttUpdateForm
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
          <ProDescriptions<API.MqttListItem>
            column={2}
            title={currentRow?.client_id}
            request={async () => ({
              data: currentRow || {},
            })}
            columns={columns as ProDescriptionsItemProps<API.MqttListItem>[]}
          />
        )}
      </Drawer>
    </PageContainer>
  );
};

export default Admin;
