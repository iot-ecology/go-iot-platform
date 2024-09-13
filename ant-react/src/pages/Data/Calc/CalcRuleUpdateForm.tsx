import { FormattedMessage } from '@@/exports';
import { ProFormDigit, ProFormText } from '@ant-design/pro-components';
import { Form, Modal } from 'antd';
import React, { useEffect } from 'react';

export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: API.CalcRuleListItem) => void;
  onSubmit: (values: API.CalcRuleListItem) => Promise<void>;
  updateModalOpen: boolean;

  values: API.CalcRuleListItem;
};
const CalcRuleUpdateForm: React.FC<UpdateFormProps> = (props) => {
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
          label={<FormattedMessage id="pages.calc-rule.name" />}
          name="name"
        />
        <ProFormText
          key={'cron'}
          label={<FormattedMessage id="pages.calc-rule.cron" />}
          name="cron"
        />
        <ProFormText
          key={'script'}
          label={<FormattedMessage id="pages.calc-rule.script" />}
          name="script"
        />
        <ProFormDigit
          key={'offset'}
          label={<FormattedMessage id="pages.calc-rule.offset" />}
          name="offset"
        />

        <ProFormText
          key={'mock_value'}
          label={<FormattedMessage id="pages.calc-rule.mock_value" />}
          name="mock_value"
        />
      </Form>
    </Modal>
  );
};

export default CalcRuleUpdateForm;
