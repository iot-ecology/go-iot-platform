import { deviceList } from '@/services/ant-design-pro/api';
import { FormattedMessage } from '@@/exports';
import { ProFormSelect, ProFormText, ProFormTextArea } from '@ant-design/pro-components';
import { Form, Modal } from 'antd';
import React, { useEffect } from 'react';

export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: API.WsHandlerListItem) => void;
  onSubmit: (values: API.WsHandlerListItem) => Promise<void>;
  updateModalOpen: boolean;

  values: API.WsHandlerListItem;
};
const WsHandlerUpdateForm: React.FC<UpdateFormProps> = (props) => {
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
        />
        <ProFormText key={'name'} label={<FormattedMessage id="pages.http.name" />} name="name" />
        <ProFormText
          key={'username'}
          label={<FormattedMessage id="pages.http.username" />}
          name="username"
        />
        <ProFormText.Password
          key={'password'}
          label={<FormattedMessage id="pages.http.password" />}
          name="password"
        />
        <ProFormTextArea
          key={'script'}
          label={<FormattedMessage id="pages.http.script" />}
          name="script"
        />
      </Form>
    </Modal>
  );
};

export default WsHandlerUpdateForm;
