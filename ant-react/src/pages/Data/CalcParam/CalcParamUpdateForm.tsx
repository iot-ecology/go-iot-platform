import { calcRuleList, mqttList, signalList } from '@/services/ant-design-pro/api';
import { FormattedMessage } from '@@/exports';
import { ProFormSelect, ProFormText } from '@ant-design/pro-components';
import { Form, Modal } from 'antd';
import React, { useEffect, useState } from 'react';

export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: API.CalcParamListItem) => void;
  onSubmit: (values: API.CalcParamListItem) => Promise<void>;
  updateModalOpen: boolean;

  values: API.CalcParamListItem;
};
const CalcParamUpdateForm: React.FC<UpdateFormProps> = (props) => {
  const [form] = Form.useForm();
  const [createDeviceUid, setCreateDeviceUid] = useState<string>();
  const [createProp, setCreateProp] = useState<string>();
  const [opSignal, setOpSignal] = useState<any>();
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
          key={'calc_rule_id'}
          label={<FormattedMessage id="pages.calc-param.calc_rule_id" />}
          name="calc_rule_id"
          disabled={true}
          request={async () => {
            let r = await calcRuleList();
            return r.data;
          }}
          fieldProps={{
            showSearch: true,
            allowClear: false,
            fieldNames: {
              label: 'name',
              value: 'ID',
            },
          }}
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
          label={<FormattedMessage id="pages.calc-param.protocol" />}
          name="protocol"
          onChange={(v) => {
            setCreateProp(v);
          }}
          fieldProps={{
            value: createProp,
            onClick: (v) => {
              setCreateDeviceUid('');
            },
          }}
        />

        <ProFormSelect
          multiple={false}
          dependencies={['protocol']}
          request={async (params) => {
            console.log(params);
            if (params.protocol === 'MQTT') {
              let res = await mqttList();
              return res.data;
            } else {
              return [];
            }
          }}
          onChange={async (v) => {
            setCreateDeviceUid(v);

            if (createDeviceUid && createProp) {
              let params = {
                ty: '数字',
                device_uid: v,
                protocol: createProp,
              };
              let newVar = await signalList(params);
              setOpSignal(newVar.data);
              debugger;
            }
          }}
          fieldProps={{
            value: createDeviceUid,
            showSearch: true,
            allowClear: false,
            fieldNames: {
              label: 'client_id',
              value: 'ID',
            },
          }}
          key={'device_uid'}
          label={<FormattedMessage id="pages.calc-param.device_uid" />}
          name="device_uid"
        />

        <ProFormText
          key={'identification_code'}
          label={<FormattedMessage id="pages.calc-param.identification_code" />}
          name="identification_code"
        />

        <ProFormSelect
          multiple={false}
          dependencies={['device_uid', 'protocol']}
          request={async (params) => {
            console.log(params);
            debugger;
            if (params.device_uid && params.protocol) {
              params.ty = '数字';
              let newVar = await signalList(params);
              return newVar.data;
            }
            return [];
          }}
          fieldProps={{
            options: opSignal,
            showSearch: true,
            allowClear: false,
            fieldNames: {
              label: 'alias',
              value: 'ID',
            },
          }}
          key={'signal_id'}
          label={<FormattedMessage id="pages.calc-param.signal_id" />}
          name="signal_id"
        />

        <ProFormText
          key={'name'}
          label={<FormattedMessage id="pages.calc-param.name" />}
          name="name"
        />

        <ProFormSelect
          valueEnum={{
            mean: { status: 'success', text: 'mean' },
            sum: { status: 'success', text: 'sum' },
            max: { status: 'success', text: 'max' },
            min: { status: 'success', text: 'min' },
            原始: { status: 'success', text: '原始' },
          }}
          key={'reduce'}
          label={<FormattedMessage id="pages.calc-param.reduce" />}
          name="reduce"
        />
      </Form>
    </Modal>
  );
};

export default CalcParamUpdateForm;
