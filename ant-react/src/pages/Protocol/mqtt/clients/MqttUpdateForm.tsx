import { FormattedMessage } from '@@/exports';
import { ProFormDigit, ProFormText } from '@ant-design/pro-components';
import { Form, Modal } from 'antd';
import React, { useEffect } from 'react';

export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: API.MqttListItem) => void;
  onSubmit: (values: API.MqttListItem) => Promise<void>;
  updateModalOpen: boolean;

  values: API.MqttListItem;
};
const MqttUpdateForm: React.FC<UpdateFormProps> = (props) => {
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
          key={'ID'}
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
          disabled={true}
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
      </Form>
    </Modal>
  );
};

export default MqttUpdateForm;
