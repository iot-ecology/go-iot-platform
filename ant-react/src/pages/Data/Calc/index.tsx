import CalcRuleResult from '@/pages/Data/Calc/CalcRuleResult';
import CalcRuleUpdateForm from '@/pages/Data/Calc/CalcRuleUpdateForm';
import MockRun from '@/pages/Data/Calc/MockRun';
import {
  addCalcRule,
  calcRulePage,
  calcRuleStart,
  calcRuleStop,
  deleteCalcRule,
  updateCalcRule,
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
  ProFormText,
  ProTable,
} from '@ant-design/pro-components';
import { useIntl } from '@umijs/max';
import { Button, Drawer, message, Tag } from 'antd';
import React, { useRef, useState } from 'react';
import { history } from 'umi';

const handleAdd = async (fields: API.CalcRuleListItem) => {
  const hide = message.loading('正在添加');
  try {
    var options = { ...fields };
    options.offset = Number(options.offset);
    await addCalcRule(options);
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
    await deleteCalcRule(id);
    hide();
    message.success('删除 successfully');
    return true;
  } catch (error) {
    hide();
    message.error('Delete failed, please try again!');
    return false;
  }
};

const handlerUpdate = async (fields: API.CalcRuleListItem) => {
  const hide = message.loading('正在更新');
  try {
    var dt = { ...fields };
    dt.offset = Number(dt.offset);
    await updateCalcRule(dt);
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
  const [currentRow, setCurrentRow] = useState<API.CalcRuleListItem>();
  const [historyModalOpen, handleHistoryModalOpen] = useState<boolean>(false);
  const [calcRuleResultModalOpen, handleCalcRuleResultModalOpen] = useState<boolean>(false);

  const columns: ProColumns<API.CalcRuleListItem>[] = [
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
      key: 'name',
      title: <FormattedMessage id="pages.calc-rule.name" />,
      hideInSearch: false,
      dataIndex: 'name',
    },
    {
      key: 'cron',
      title: <FormattedMessage id="pages.calc-rule.cron" />,
      hideInSearch: true,
      valueType: 'code',
      dataIndex: 'cron',
    },
    {
      key: 'script',
      title: <FormattedMessage id="pages.calc-rule.script" />,
      hideInSearch: true,
      dataIndex: 'script',
      valueType: 'code',
    },
    {
      key: 'offset',
      title: <FormattedMessage id="pages.calc-rule.offset" />,
      hideInSearch: true,
      dataIndex: 'offset',
    },
    {
      key: 'start',
      title: <FormattedMessage id="pages.calc-rule.start" />,
      hideInSearch: true,
      dataIndex: 'start',
      render: (dom, entity) => {
        return (
          <>
            {entity.start ? (
              <Tag color={'blue'}>
                <FormattedMessage id="pages.calc-rule.start.true" defaultMessage="已启动" />
              </Tag>
            ) : (
              <Tag color={'red'}>
                <FormattedMessage id="pages.calc-rule.start.false" defaultMessage="已停止" />
              </Tag>
            )}
          </>
        );
      },
    },
    {
      key: 'mock_value',
      title: <FormattedMessage id="pages.calc-rule.mock_value" />,
      hideInSearch: true,
      dataIndex: 'mock_value',
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
        record.start ? (
          <>
            <Button
              key="start"
              onClick={async () => {
                // todo: 修改
                setCurrentRow(record);
                await calcRuleStop(record.ID);
                await actionRef.current?.reload();
              }}
            >
              <FormattedMessage id="pages.stop" defaultMessage="暂停" />
            </Button>
          </>
        ) : (
          <>
            <Button
              key="start"
              onClick={async () => {
                // todo: 修改
                setCurrentRow(record);
                await calcRuleStart(record.ID);
                await actionRef.current?.reload();
              }}
            >
              <FormattedMessage id="pages.start" defaultMessage="启动" />
            </Button>
          </>
        ),
        <Button
          key="mock-start"
          onClick={() => {
            // todo: 修改
            handleHistoryModalOpen(true);
            setCurrentRow(record);
          }}
        >
          <FormattedMessage id="pages.mock.start" defaultMessage="模拟执行" />
        </Button>,
        <Button
          key="set-param"
          onClick={async () => {
            setCurrentRow(record);
            handleCalcRuleResultModalOpen(true);
          }}
        >
          <FormattedMessage id="pages.history" defaultMessage="历史数据" />
        </Button>,
        <Button
          key="set-param"
          onClick={async () => {
            history.push({
              pathname: '/data/calc-param?id=' + record.ID,
            });
          }}
        >
          <FormattedMessage id="pages.set-param" defaultMessage="参数设置" />
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
      <ProTable<API.CalcRuleListItem, API.PageParams>
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
        request={calcRulePage}
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
          const success = await handleAdd(value as API.CalcRuleListItem);
          if (success) {
            handleModalOpen(false);
            if (actionRef.current) {
              await actionRef.current.reload();
            }
          }
        }}
      >
        <ProFormText
          key={'name'}
          label={<FormattedMessage id="pages.calc-rule.name" />}
          name="name"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText
          key={'cron'}
          label={<FormattedMessage id="pages.calc-rule.cron" />}
          name="cron"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText
          key={'script'}
          label={<FormattedMessage id="pages.calc-rule.script" />}
          name="script"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormDigit
          key={'offset'}
          label={<FormattedMessage id="pages.calc-rule.offset" />}
          name="offset"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
      </ModalForm>
      <CalcRuleUpdateForm
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
      <MockRun
        updateModalOpen={historyModalOpen}
        onCancel={async () => {
          handleHistoryModalOpen(false);

          await actionRef.current?.reload();
        }}
        values={currentRow || {}}
      />{' '}
      <CalcRuleResult
        updateModalOpen={calcRuleResultModalOpen}
        onCancel={async () => {
          handleCalcRuleResultModalOpen(false);

          await actionRef.current?.reload();
        }}
        values={currentRow || {}}
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
          <ProDescriptions<API.CalcRuleListItem>
            column={2}
            title={currentRow?.name}
            request={async () => ({
              data: currentRow || {},
            })}
            params={{
              id: currentRow?.name,
            }}
            columns={columns as ProDescriptionsItemProps<API.CalcRuleListItem>[]}
          />
        )}
      </Drawer>
    </PageContainer>
  );
};

export default Admin;
