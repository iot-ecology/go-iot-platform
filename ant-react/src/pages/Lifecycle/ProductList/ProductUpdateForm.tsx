import { FormattedMessage } from '@@/exports';
import {
  ProForm,
  ProFormDigit,
  ProFormText,
  ProFormUploadDragger,
} from '@ant-design/pro-components';
import { Form, GetProp, message, Modal, UploadFile, UploadProps } from 'antd';
import { RcFile } from 'antd/lib/upload';
import React, { useEffect, useState } from 'react';
type FileType = Parameters<GetProp<UploadProps, 'beforeUpload'>>[0];

export type UpdateFormProps = {
  onCancel: (flag?: boolean, formVals?: API.ProductItem) => void;
  onSubmit: (values: API.ProductItem) => Promise<void>;
  updateModalOpen: boolean;

  values: API.ProductItem;
};

function generateUUID() {
  return 'xxxx-xxxx-4xxx-xxxx'.replace(/[xy]/g, function (c) {
    var r = (Math.random() * 16) | 0,
      v = c === 'x' ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

const ProductUpdateForm: React.FC<UpdateFormProps> = (props) => {
  const [form] = Form.useForm();
  useEffect(() => {
    form.resetFields();
    form.setFieldsValue(props.values);
  });
  const [fileList, setFileList] = useState<UploadFile[]>([]);
  function beforeUpload(form: any, file: RcFile, fileList: RcFile[]) {
    const formData = new FormData();
    fileList.forEach((file) => {
      formData.append('file', file);
    });
    fetch('/api/file/update', {
      method: 'POST',
      body: formData,
    })
      .then((res) => res.json())
      .then((vvv) => {
        setFileList([
          {
            name: vvv.file_path,
            uid: generateUUID(),
            fileName: vvv.file_path,
          },
        ]);
        props.values.image_url = vvv.file_path;
        message.success('upload successfully.');
      })
      .catch((e) => {
        message.error('upload failed.');
      })
      .finally(() => {});
    return false;
  }

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
        <ProForm.Group>
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
            key={'name'}
            label={<FormattedMessage id="pages.product.name" />}
            name="name"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.input" />,
              },
            ]}
          />
          <ProFormText
            key={'description'}
            label={<FormattedMessage id="pages.product.description" />}
            name="description"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.input" />,
              },
            ]}
          />
          <ProFormText
            key={'sku'}
            label={<FormattedMessage id="pages.product.sku" />}
            name="sku"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.input" />,
              },
            ]}
          />
        </ProForm.Group>
        <ProForm.Group>
          <ProFormDigit
            key={'price'}
            label={<FormattedMessage id="pages.product.price" />}
            name="price"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.input" />,
              },
            ]}
          />
          <ProFormDigit
            key={'cost'}
            label={<FormattedMessage id="pages.product.cost" />}
            name="cost"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.input" />,
              },
            ]}
          />
          <ProFormDigit
            key={'quantity'}
            label={<FormattedMessage id="pages.product.quantity" />}
            name="quantity"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.input" />,
              },
            ]}
          />
          <ProFormDigit
            key={'minimum_stock'}
            label={<FormattedMessage id="pages.product.minimum_stock" />}
            name="minimum_stock"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.input" />,
              },
            ]}
          />
          <ProFormDigit
            key={'warranty_period'}
            label={<FormattedMessage id="pages.product.warranty_period" />}
            name="warranty_period"
            rules={[
              {
                required: true,
                message: <FormattedMessage id="pages.rules.input" />,
              },
            ]}
          />
        </ProForm.Group>
        <ProFormText
          key={'status'}
          label={<FormattedMessage id="pages.product.status" />}
          name="status"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormText
          key={'tags'}
          label={<FormattedMessage id="pages.product.tags" />}
          placeholder={'请输入,如果存在多个请用逗号分割'}
          name="tags"
          rules={[
            {
              required: true,
              message: <FormattedMessage id="pages.rules.input" />,
            },
          ]}
        />
        <ProFormUploadDragger
          name="image_url"
          label={<FormattedMessage id="pages.product.image_url" />}
          fieldProps={{
            accept: '.jpg,.png', // 设置接受的文件类型
            name: 'file',
            action: '/api/file/update',
            multiple: false,

            beforeUpload: (file, fileList) => beforeUpload(form, file, fileList),
            fileList: fileList,
          }}
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

export default ProductUpdateForm;
