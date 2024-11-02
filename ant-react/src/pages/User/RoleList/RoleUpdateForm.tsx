import { FormattedMessage } from '@@/exports';
import { ProFormSelect, ProFormText } from '@ant-design/pro-components';
import { Form, message, Modal } from 'antd';
import React, { useEffect, useState } from 'react';
import { deptList } from '@/services/ant-design-pro/api';

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
        <ProFormText key={'name'} label={<FormattedMessage id="pages.name" />} name="name" />
        <ProFormText label={<FormattedMessage id="pages.desc" />} name="description" />

      </Form>
    </Modal>
  );
};

export default RoleUpdateForm;
