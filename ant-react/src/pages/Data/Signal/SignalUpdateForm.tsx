import { mqttList } from '@/services/ant-design-pro/api';
import { FormattedMessage } from '@@/exports';
import { ProFormDigit, ProFormSelect, ProFormText } from '@ant-design/pro-components';
import { Form, Modal } from 'antd';
import React, { useEffect } from 'react';

export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: API.SignalListItem) => void;
  onSubmit: (values: API.SignalListItem) => Promise<void>;
  updateModalOpen: boolean;

  values: API.SignalListItem;
};
const SignalUpdateForm: React.FC<UpdateFormProps> = (props) => {
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
          disabled={true}
          valueEnum={{
            MQTT: { text: 'MQTT', status: 'success' },
            HTTP: { text: 'HTTP', status: 'success' },
            WebSocket: { text: 'WebSocket', status: 'success' },
            TCP: { text: 'TCP', status: 'success' },
            COAP: { text: 'COAP', status: 'success' },
          }}
          key={'protocol'}
          label={<FormattedMessage id="pages.signal.protocol" />}
          name="protocol"
        />
        <ProFormSelect
          multiple={false}
          disabled={true}
          dependencies={['protocol']}
          request={async (params) => {
            console.log(params);
            if (params.protocol === 'MQTT') {
              let res = await mqttList();
              return res.data;
            }
          }}
          fieldProps={{
            showSearch: true,
            allowClear: false,
            fieldNames: {
              label: 'client_id',
              value: 'ID',
            },
          }}
          key={'device_uid'}
          label={<FormattedMessage id="pages.signal.device_uid" />}
          name="device_uid"
        />
        <ProFormText
          key={'identification_code'}
          label={<FormattedMessage id="pages.signal.identification_code" />}
          name="identification_code"
        />

        <ProFormText key={'name'} label={<FormattedMessage id="pages.signal.name" />} name="name" />
        <ProFormText
          key={'alias'}
          label={<FormattedMessage id="pages.signal.alias" />}
          name="alias"
        />
        <ProFormSelect
          key={'type'}
          valueEnum={{
            数字: { status: 'success', text: '数字' },
            文本: { status: 'default', text: '文本' },
          }}
          label={<FormattedMessage id="pages.signal.type" />}
          name="type"
        />
        <ProFormText key={'unit'} label={<FormattedMessage id="pages.signal.unit" />} name="unit" />
        <ProFormDigit
          fieldProps={{ precision: 0 }}
          key={'cache_size'}
          label={<FormattedMessage id="pages.signal.cache_size" />}
          name="cache_size"
        />
      </Form>
    </Modal>
  );
};

export default SignalUpdateForm;
