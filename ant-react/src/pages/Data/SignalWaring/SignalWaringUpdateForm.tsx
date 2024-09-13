import { mqttList, signalList } from '@/services/ant-design-pro/api';
import { FormattedMessage } from '@@/exports';
import { ProFormDigit, ProFormSelect, ProFormText } from '@ant-design/pro-components';
import { Form, Modal } from 'antd';
import React, { useEffect } from 'react';

export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: API.SignalWaringItem) => void;
  onSubmit: (values: API.SignalWaringItem) => Promise<void>;
  updateModalOpen: boolean;

  values: API.SignalWaringItem;
};
const SignalWaringUpdateForm: React.FC<UpdateFormProps> = (props) => {
  const [form] = Form.useForm();
  useEffect(() => {
    props.values.in_or_out = String(props.values.in_or_out);
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

        <ProFormSelect
          valueEnum={{
            MQTT: { text: 'MQTT', status: 'success' },
            HTTP: { text: 'HTTP', status: 'success' },
            WebSocket: { text: 'WebSocket', status: 'success' },
            'TCP/IP': { text: 'TCP/IP', status: 'success' },
            COAP: { text: 'COAP', status: 'success' },
          }}
          key={'protocol'}
          label={<FormattedMessage id="pages.signal.waring.protocol" />}
          name="protocol"
        />
        <ProFormSelect
          multiple={false}
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
          label={<FormattedMessage id="pages.signal.waring.device_uid" />}
          name="device_uid"
        />

        <ProFormText
          key={'identification_code'}
          label={<FormattedMessage id="pages.signal.waring.identification_code" />}
          name="identification_code"
        />

        <ProFormSelect
          multiple={false}
          dependencies={['device_uid', 'protocol']}
          request={async (params) => {
            console.log(params);
            if (params.device_uid && params.protocol) {
              params.ty = '数字';
              let newVar = await signalList(params);
              return newVar.data;
            }
            return [];
          }}
          fieldProps={{
            showSearch: true,
            allowClear: false,
            fieldNames: {
              label: 'alias',
              value: 'ID',
            },
          }}
          key={'signal_id'}
          label={<FormattedMessage id="pages.signal.waring.signal_id" />}
          name="signal_id"
        />
        <ProFormDigit
          key={'min'}
          label={<FormattedMessage id="pages.signal.waring.min" />}
          name="min"
        />
        <ProFormDigit
          key={'max'}
          label={<FormattedMessage id="pages.signal.waring.max" />}
          name="max"
        />
        <ProFormSelect
          key={'in_or_out'}
          label={<FormattedMessage id="pages.signal.waring.in_or_out" />}
          name="in_or_out"
          valueEnum={{
            1: { text: '范围内报警', status: 'success' },
            0: { text: '范围外报警', status: 'success' },
          }}
        />
      </Form>
    </Modal>
  );
};

export default SignalWaringUpdateForm;
