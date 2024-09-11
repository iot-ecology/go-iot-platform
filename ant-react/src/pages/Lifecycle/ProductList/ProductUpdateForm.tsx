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
        setFileList([]);
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
      onOk={() => {
        props.onSubmit(form.getFieldsValue());
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
          />
          <ProFormText
            key={'name'}
            label={<FormattedMessage id="pages.product.name" />}
            name="name"
          />
          <ProFormText
            key={'description'}
            label={<FormattedMessage id="pages.product.description" />}
            name="description"
          />
          <ProFormText key={'sku'} label={<FormattedMessage id="pages.product.sku" />} name="sku" />
        </ProForm.Group>
        <ProForm.Group>
          <ProFormDigit
            key={'price'}
            label={<FormattedMessage id="pages.product.price" />}
            name="price"
          />
          <ProFormDigit
            key={'cost'}
            label={<FormattedMessage id="pages.product.cost" />}
            name="cost"
          />
          <ProFormDigit
            key={'quantity'}
            label={<FormattedMessage id="pages.product.quantity" />}
            name="quantity"
          />
          <ProFormDigit
            key={'minimum_stock'}
            label={<FormattedMessage id="pages.product.minimum_stock" />}
            name="minimum_stock"
          />
          <ProFormDigit
            key={'warranty_period'}
            label={<FormattedMessage id="pages.product.warranty_period" />}
            name="warranty_period"
          />
        </ProForm.Group>
        <ProFormText
          key={'status'}
          label={<FormattedMessage id="pages.product.status" />}
          name="status"
        />
        <ProFormText
          key={'tags'}
          label={<FormattedMessage id="pages.product.tags" />}
          placeholder={'请输入,如果存在多个请用逗号分割'}
          name="tags"
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
        />
      </Form>
    </Modal>
  );
};

export default ProductUpdateForm;
