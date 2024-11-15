import { ClickhouseTransmitList, deviceList, mqttList } from '@/services/ant-design-pro/api';
import { FormattedMessage } from '@@/exports';
import {
  ProFormRadio,
  ProFormSelect,
  ProFormText,
  ProFormTextArea,
} from '@ant-design/pro-components';
import { Form, Modal } from 'antd';
import React, { useEffect } from 'react';

export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: API.ClickhouseTransmitBindListItem) => void;
  onSubmit: (values: API.ClickhouseTransmitBindListItem) => Promise<void>;
  updateModalOpen: boolean;

  values: API.ClickhouseTransmitBindListItem;
};
const ClickhouseTransmitBindUpdateForm: React.FC<UpdateFormProps> = (props) => {
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
          key={'id'}
          disabled={true}
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
          valueEnum={{
            MQTT: { text: 'MQTT', status: 'success' },
            HTTP: { text: 'HTTP', status: 'success' },
            WebSocket: { text: 'WebSocket', status: 'success' },
            TCP: { text: 'TCP', status: 'success' },
            COAP: { text: 'COAP', status: 'success' },
          }}
          key={'protocol'}
          label={<FormattedMessage id="pages.CassandraTransmitBind.protocol" />}
          name="protocol"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />

        <ProFormSelect
          multiple={false}
          dependencies={['protocol']}
          onChange={async (value) => {
            form.setFieldValue('device_uid', value);
          }}
          request={async (params) => {
            console.log(params);
            if (params.protocol === 'MQTT') {
              let res = await mqttList();
              return res.data;
            } else {
              let c = await deviceList();
              let r = [];
              c.data.forEach((e) => {
                if (e.protocol === params.protocol) {
                  r.push({
                    client_id: e.sn,
                    ID: e.ID,
                  });
                }
              });
              return r;
            }
          }}
          fieldProps={{
            onClick: (v) => {
              console.log(v);
            },
            showSearch: true,
            allowClear: false,
            fieldNames: {
              label: 'client_id',
              value: 'ID',
            },
          }}
          key={'device_uid'}
          label={<FormattedMessage id="pages.CassandraTransmitBind.device_uid" />}
          name="device_uid"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />

        <ProFormText
          key={'identification_code'}
          label={<FormattedMessage id="pages.CassandraTransmitBind.identification_code" />}
          name="identification_code"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormSelect
          request={async () => {
            let r = await ClickhouseTransmitList();
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
          key={'clickhouse_transmit_id'}
          label={<FormattedMessage id="pages.CassandraTransmitBind.cassandra_transmit_id" />}
          name="clickhouse_transmit_id"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />
        <ProFormText
          key={'database'}
          label={<FormattedMessage id="pages.CassandraTransmitBind.database" />}
          name="database"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText
          key={'table'}
          label={<FormattedMessage id="pages.CassandraTransmitBind.table" />}
          name="table"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormTextArea
          tooltip={
            'function main(jsonData) {\n' +
            '    var c = []\n' +
            '    for (var jsonDatum of jsonData) {\n' +
            '        var time = jsonDatum.Time;\n' +
            '        var arr = []\n' +
            '        var timeField = {\n' +
            '            "FieldName": "time",\n' +
            '            "Value": time\n' +
            '        }\n' +
            '\t\t\n' +
            '        arr.push(timeField)\n' +
            '        for (var e of jsonDatum.DataRows) {\n' +
            '            if (e.Name == "a") {\n' +
            '                var aField = {\n' +
            '                    "FieldName": "name",\n' +
            '                    "Value": e.Value\n' +
            '                }\n' +
            '                arr.push(aField)\n' +
            '            }\n' +
            '        }\n' +
            '        c.push(arr)\n' +
            '    }\n' +
            '    return c;\n' +
            '}'
          }
          key={'script'}
          label={<FormattedMessage id="pages.CassandraTransmitBind.script" />}
          name="script"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />

        <ProFormRadio.Group
          name="enable"
          label={<FormattedMessage id="pages.CassandraTransmitBind.enable" />}
          options={[
            {
              label: '启用',
              value: true,
            },
            {
              label: '停用',
              value: false,
            },
          ]}
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.select" />,
            },
          ]}
        />
      </Form>
    </Modal>
  );
};

export default ClickhouseTransmitBindUpdateForm;
