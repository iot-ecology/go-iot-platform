import { FormattedMessage } from '@@/exports';
import { ProFormText } from '@ant-design/pro-components';
import { Form, Modal } from 'antd';
import React, { useEffect } from 'react';

export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: API.UserListItem) => void;
  onSubmit: (values: API.UserListItem) => Promise<void>;
  updateModalOpen: boolean;

  values: API.UserListItem;
};
const UserUpdateForm: React.FC<UpdateFormProps> = (props) => {
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
      onOk={() => {
        props.onSubmit(form.getFieldsValue());
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
        />
        <ProFormText
          key={'username'}
          label={<FormattedMessage id="pages.user-list.username" />}
          name="username"
        />
        <ProFormText
          key={'email'}
          label={<FormattedMessage id="pages.user-list.mail" />}
          name="email"
        />
        <ProFormText.Password
          key={'password'}
          label={<FormattedMessage id="pages.user-list.password" />}
          name="password"
        />
      </Form>
    </Modal>
  );
};

export default UserUpdateForm;
