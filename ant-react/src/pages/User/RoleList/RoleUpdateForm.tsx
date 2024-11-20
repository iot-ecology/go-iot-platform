import { FormattedMessage } from '@@/exports';
import { ProFormText } from '@ant-design/pro-components';
import { Form, Modal } from 'antd';
import React, { useEffect } from 'react';

export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: API.RoleListItem) => void;
  onSubmit: (values: API.RoleListItem) => Promise<void>;
  updateModalOpen: boolean;

  values: API.RoleListItem;
};
const RoleUpdateForm: React.FC<UpdateFormProps> = (props) => {
  const [form] = Form.useForm();
  useEffect(() => {
    form.resetFields();
    form.setFieldsValue(props.values);
  });

  return (
    <Modal
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
        <ProFormText
          label={<FormattedMessage id="pages.desc" />}
          name="description"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
      </Form>
    </Modal>
  );
};

export default RoleUpdateForm;
