import { FormattedMessage } from '@@/exports';
import { ProFormDigit, ProFormText } from '@ant-design/pro-components';
import { Form, Modal } from 'antd';
import React, { useEffect } from 'react';

export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: API.MongoTransmitListItem) => void;
  onSubmit: (values: API.MongoTransmitListItem) => Promise<void>;
  updateModalOpen: boolean;

  values: API.MongoTransmitListItem;
};
const MongoTransmitUpdateForm: React.FC<UpdateFormProps> = (props) => {
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
          key={'ID'}
          label={<FormattedMessage id="pages.id" />}
          name="ID"
        />
        <ProFormText
          key={'name'}
          label={<FormattedMessage id="pages.CassandraTransmit.name" />}
          name="name"
        />
        <ProFormText
          key={'host'}
          label={<FormattedMessage id="pages.CassandraTransmit.host" />}
          name="host"
        />
        <ProFormDigit
          key={'port'}
          label={<FormattedMessage id="pages.CassandraTransmit.port" />}
          name="port"
        />
        <ProFormText
          key={'username'}
          label={<FormattedMessage id="pages.CassandraTransmit.username" />}
          name="username"
        />
        <ProFormText
          key={'password'}
          label={<FormattedMessage id="pages.CassandraTransmit.password" />}
          name="password"
        />
      </Form>
    </Modal>
  );
};

export default MongoTransmitUpdateForm;
