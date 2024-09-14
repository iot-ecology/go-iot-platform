import { FormattedMessage } from '@@/exports';
import { ProFormText } from '@ant-design/pro-components';
import { Form, Modal } from 'antd';
import React, { useEffect } from 'react';

export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: API.DingDingListItem) => void;
  onSubmit: (values: API.DingDingListItem) => Promise<void>;
  updateModalOpen: boolean;

  values: API.DingDingListItem;
};
const DingDingUpdateForm: React.FC<UpdateFormProps> = (props) => {
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
          key={'id'}
          label={<FormattedMessage id="pages.id" />}
          disabled={true}
          name="ID"
        />
        <ProFormText
          key={'name'}
          label={<FormattedMessage id="pages.dingding.name" />}
          name="name"
        />
        <ProFormText
          key={'access_token'}
          label={<FormattedMessage id="pages.dingding.access_token" />}
          name="access_token"
        />
        <ProFormText
          key={'secret'}
          label={<FormattedMessage id="pages.dingding.secret" />}
          name="secret"
        />
        <ProFormText
          key={'content'}
          label={<FormattedMessage id="pages.dingding.content" />}
          name="content"
        />
      </Form>
    </Modal>
  );
};

export default DingDingUpdateForm;
