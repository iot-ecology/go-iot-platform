import { deviceList } from '@/services/ant-design-pro/api';
import { FormattedMessage } from '@@/exports';
import { ProFormSelect, ProFormText, ProFormTextArea } from '@ant-design/pro-components';
import { Form, Modal } from 'antd';
import React, { useEffect } from 'react';

export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: API.TcpHandlerListItem) => void;
  onSubmit: (values: API.TcpHandlerListItem) => Promise<void>;
  updateModalOpen: boolean;

  values: API.TcpHandlerListItem;
};
const TcpHandlerUpdateForm: React.FC<UpdateFormProps> = (props) => {
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
        <ProFormSelect
          key={'device_info_id'}
          disabled={true}
          label={<FormattedMessage id="pages.http.device_info_id" />}
          name="device_info_id"
          request={async () => {
            let a = await deviceList();
            return a.data;
          }}
          fieldProps={{
            showSearch: true,
            allowClear: false,
            fieldNames: {
              label: 'sn',
              value: 'ID',
            },
          }}
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />
        <ProFormText
          key={'name'}
          label={<FormattedMessage id="pages.http.name" />}
          name="name"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText
          key={'username'}
          label={<FormattedMessage id="pages.http.username" />}
          name="username"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText.Password
          key={'password'}
          label={<FormattedMessage id="pages.http.password" />}
          name="password"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormTextArea
          key={'script'}
          label={<FormattedMessage id="pages.http.script" />}
          name="script"
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

export default TcpHandlerUpdateForm;
