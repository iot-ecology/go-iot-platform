import { genScriptParam, mockScriptWaringRun } from '@/services/ant-design-pro/api';
import { Divider, Form, Modal } from 'antd';
import React, { useEffect, useState } from 'react';

export type Props = {
  updateModalOpen: boolean;
  onCancel: (flag?: boolean, formVals?: API.ScriptWaringListItem) => void;
  values: API.ScriptWaringListItem;
};

const DebugScript: React.FC<Props> = (props) => {
  const [form] = Form.useForm();
  const [param, setParam] = useState<any>(null);
  const [result, setResult] = useState<any>(null);

  useEffect(() => {
    const fetchData = async () => {
      const response = await genScriptParam(props.values.ID);
      const r = await mockScriptWaringRun(props.values.ID);
      setParam(response);
      setResult(r);
    };

    fetchData();
  }, [props.values]);
  return (
    <Modal
      form={form}
      destroyOnClose
      title={'调试脚本'}
      forceRender={true}
      open={props.updateModalOpen}
      onCancel={() => {
        props.onCancel();
      }}
      onClose={() => {
        props.onCancel();
      }}
      footer={false}
    >
      {/* 数据表格 */}
      <div>{JSON.stringify(param?.data)}</div>
      <Divider />
      <div>{JSON.stringify(result)}</div>
    </Modal>
  );
};

export default DebugScript;
