import { deptList } from '@/services/ant-design-pro/api';
import { FormattedMessage } from '@@/exports';
import { ProFormSelect, ProFormText } from '@ant-design/pro-components';
import { Form, message, Modal } from 'antd';
import React, { useEffect, useState } from 'react';

export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: API.DeptListItem) => void;
  onSubmit: (values: API.DeptListItem) => Promise<void>;
  updateModalOpen: boolean;

  values: API.DeptListItem;
};
const DeptUpdateForm: React.FC<UpdateFormProps> = (props) => {
  const [options, setOptions] = useState<{ label: string; value: string | number }[]>([]);

  const [form] = Form.useForm();
  useEffect(() => {
    form.resetFields();
    form.setFieldsValue(props.values);
  });

  const fetchData = async () => {
    try {
      const result = await deptList(); // 调用API获取数据
      if (result && result.data) {
        if (Array.isArray(result.data)) {
          const formattedOptions = result.data.map((item) => ({
            label: item.name, // 假设每项数据都有一个name属性
            value: item.ID, // 假设每项数据都有一个id属性
          }));

          // @ts-ignore
          setOptions(formattedOptions);
        }
      }
    } catch (error) {
      message.error('请求数据失败，请稍后再试！');
    }
  };

  return (
    <Modal
      afterOpenChange={async (flag: boolean) => {
        await fetchData();
      }}
      key="userupdateform"
      destroyOnClose
      forceRender={true}
      open={props.updateModalOpen}
      onCancel={(vvv) => {
        props.onCancel();
      }}
      onOk={async () => {
        let success = await form.validateFields();
        if (success) {
          props.onSubmit(form.getFieldsValue());
        }
      }}
      onClose={() => {
        props.onCancel();
      }}
    >
      <Form key={'userupdateform'} form={form} style={{ padding: '32px 40px 48px' }}>
        <ProFormText
          disabled={true}
          key={'id'}
          label={<FormattedMessage id="pages.id" />}
          name="ID"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
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
      </Form>
    </Modal>
  );
};

export default DeptUpdateForm;
