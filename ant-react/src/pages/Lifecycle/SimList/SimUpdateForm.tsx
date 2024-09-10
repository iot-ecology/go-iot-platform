import { FormattedMessage } from '@@/exports';
import { ProFormDatePicker, ProFormText } from '@ant-design/pro-components';
import { Form, Modal } from 'antd';
import dayjs from 'dayjs';
import React, { useEffect } from 'react';

export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: API.UserListItem) => void;
  onSubmit: (values: API.UserListItem) => Promise<void>;
  updateModalOpen: boolean;

  values: API.UserListItem;
};
const SimUpdateForm: React.FC<UpdateFormProps> = (props) => {
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
          key={'access_number'}
          label={<FormattedMessage id="pages.sim.access_number" />}
          name="access_number"
        />
        <ProFormText key={'iccid'} label={<FormattedMessage id="pages.sim.iccid" />} name="iccid" />
        <ProFormText key={'imsi'} label={<FormattedMessage id="pages.sim.imsi" />} name="imsi" />
        <ProFormText
          key={'operator'}
          label={<FormattedMessage id="pages.sim.operator" />}
          name="operator"
        />
        <ProFormDatePicker
          transform={(value) => {
            return dayjs(value).format('YYYY-MM-DDTHH:mm:ssZ');
          }}
          key={'expiration'}
          label={<FormattedMessage id="pages.sim.expiration" />}
          name="expiration"
        />
      </Form>
    </Modal>
  );
};

export default SimUpdateForm;
