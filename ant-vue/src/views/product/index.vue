<template>
  <div class="role-preview">
    <a-card :bordered="true">
      <a-form layout="inline">
        <a-form-item :label="$t('message.name')">
          <a-input v-model:value="nameStr" style="width: 300px" :placeholder="$t('message.pleaseEnter')" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="pageList()">{{ $t('message.search') }}</a-button>
        </a-form-item>
      </a-form>
      <a-button style="margin: 10px 0" type="primary" @click="modalVisible = true,title = $t('message.addition')">{{ $t('message.addition') }}</a-button>

      <a-table :columns="columns" :data-source="list" bordered :pagination="pagination" @change="handleTableChange">
        <template #bodyCell="{ column, text, record }">
          <template v-if="column.dataIndex === 'image_url'">
            <div>
              <a-image :width="100"
                  :src="`http://192.168.3.100:8080/file/download?path=${text}`"
              />
            </div>
          </template>
          <template v-if="column.dataIndex === 'tags'">
            <div>
              <template v-for="tag in text?.split(',')" :key="tag">
                <a-tooltip v-if="tag.length > 10" :title="tag">
                  <a-tag>
                    {{ `${tag.slice(0, 10)}...` }}
                  </a-tag>
                </a-tooltip>
                <a-tag v-else>
                  {{ tag }}
                </a-tag>
              </template>
            </div>
          </template>
          <template v-if="column.dataIndex === 'operation'">
            <div class="editable-row-operations">
              <span>
                <a-button type="primary" size="small" @click="edit(record)">{{$t('message.edit')}}</a-button>
                <a-popconfirm :title="$t('message.sureDelete')" :okText="$t('message.yes')" :cancelText="$t('message.no')" @confirm="confirm(record.id)">
                  <a-button type="primary" size="small" danger style="margin-left: 10px;">{{$t('message.delete')}}</a-button>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>
      <!--新增/编辑-->
      <a-modal @cancel="handleModalClose" :okText="$t('message.confirm')" :cancelText="$t('message.cancel')" v-model:open="modalVisible" :destroy-on-close="true" :title="title" @ok="onAddData()">
        <a-form ref="formRef" :label-col="{ style: { width: '120px' } }" :labelWrap="true" :rules="rules" :model="form">
          <a-form-item :label="$t('message.name')" name="name">
            <a-input v-model:value="form.name" style="width: 350px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
          <a-form-item :label="$t('message.productDes')" name="description">
            <a-textarea v-model:value="form.description" style="width: 350px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
          <a-form-item :label="$t('message.imageUrl')" name="image_url">
            <Upload :imgStr="imageUrl" v-model="form.image_url" />
          </a-form-item>
          <a-form-item :label="$t('message.minimum_stock')" name="minimum_stock">
            <a-input-number :precision="0" v-model:value="form.minimum_stock" style="width: 350px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
          <a-form-item :label="$t('message.price')" name="price">
            <a-input-number :precision="2" v-model:value="form.price" style="width: 350px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
          <a-form-item :label="$t('message.quantity')" name="quantity">
            <a-input-number :precision="0" v-model:value="form.quantity" style="width: 350px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
          <a-form-item :label="$t('message.sku')" name="sku">
            <a-input v-model:value="form.sku" style="width: 350px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
          <a-form-item :label="$t('message.status')" name="status">
            <a-input v-model:value="form.status" style="width: 350px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
          <a-form-item :label="$t('message.tags')" name="tags">
            <Tag :tags="tags" v-model="form.tags" style="width: 350px" />
          </a-form-item>
          <a-form-item :label="$t('message.warranty_period')" name="warranty_period">
            <a-input-number :precision="0" v-model:value="form.warranty_period" style="width: 350px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
        </a-form>
      </a-modal>
    </a-card>
  </div>
</template>
<script lang="ts" setup>
import {onMounted, reactive, ref, watch} from 'vue'
import {ProductCreate, ProductDelete, ProductPage, ProductUpdate} from "@/api";
import {useI18n} from "vue-i18n";
import {Rule} from "ant-design-vue/es/form";
import {message} from "ant-design-vue";
import { cloneDeep } from "lodash-es";
import Upload from '@/components/upload/index.vue'
import Tag from '@/components/Tag/index.vue'

const { t,locale } = useI18n();
const formRef = ref<HTMLFormElement | null>(null);
const pagination = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
  showSizeChanger: true, // 显示每页显示条目数选择器
});
const nameStr = ref('');
const modalVisible = ref(false);
const form = reactive({
  id: "",
  description: '',
  name: '',
  minimum_stock: '',
  image_url: '',
  price: '',
  quantity:'',
  sku:'',
  status:'',
  tags:'',
  warranty_period:''
});
const list = ref([]);
const imageUrl = ref('');
const tags = ref('')
const columns = ref([
  {
    title: t('message.uniCode'),
    dataIndex: "id",
  },
  {
    title: t('message.name'),
    dataIndex: "name",
  },
  {
    title: t('message.productDes'),
    dataIndex: "description",
  },
  {
    title: t('message.imageUrl'),
    dataIndex: "image_url",
  },
  {
    title: t('message.minimum_stock'),
    dataIndex: "minimum_stock",
  },
  {
    title: t('message.price'),
    dataIndex: "price",
  },
  {
    title: t('message.quantity'),
    dataIndex: "quantity",
  },
  {
    title: t('message.sku'),
    dataIndex: "sku",
  },
  {
    title: t('message.status'),
    dataIndex: "status",
  },
  {
    title: t('message.tags'),
    dataIndex: "tags",
  },
  {
    title: t('message.warranty_period'),
    dataIndex: "warranty_period",
  },
  {
    title: t('message.operation'),
    dataIndex: "operation",
  }
]);
const title = ref(t('message.addition'))

let rules: Record<string, Rule[]> = {
  name: [{ required: true, message: t('message.pleaseName'), trigger: "blur" }],
  description: [{ required: true, message: t('message.pleaseEnterDescription'), trigger: "blur" }],
  image_url: [{ required: true, message: t('message.pleaseUploadPictures'), trigger: "change" }],
  minimum_stock: [{ required: true, message: t('message.pleaseEnterMinimumStock'), trigger: "change" }],
  price: [{ required: true, message: t('message.pleaseEnterPrice'), trigger: "change" }],
  quantity: [{ required: true, message: t('message.pleaseEnterQuantity'), trigger: "change" }],
  sku: [{ required: true, message: t('message.pleaseEnterSku'), trigger: "change" }],
  status: [{ required: true, message: t('message.pleaseEnterStatus'), trigger: "change" }],
  tags: [{ required: true, message: t('message.pleaseEnterTags'), trigger: "change" }],
  warranty_period: [{ required: true, message: t('message.pleaseEnterWarrantyPeriod'), trigger: "change" }],
};
watch(locale, () => {
  columns.value = [
    {
      title: t('message.uniCode'),
      dataIndex: "id",
    },
    {
      title: t('message.name'),
      dataIndex: "name",
    },
    {
      title: t('message.productDes'),
      dataIndex: "description",
    },
    {
      title: t('message.imageUrl'),
      dataIndex: "image_url",
    },
    {
      title: t('message.minimum_stock'),
      dataIndex: "minimum_stock",
    },
    {
      title: t('message.price'),
      dataIndex: "price",
    },
    {
      title: t('message.quantity'),
      dataIndex: "quantity",
    },
    {
      title: t('message.sku'),
      dataIndex: "sku",
    },
    {
      title: t('message.status'),
      dataIndex: "status",
    },
    {
      title: t('message.tags'),
      dataIndex: "tags",
    },
    {
      title: t('message.warranty_period'),
      dataIndex: "warranty_period",
    },
    {
      title: t('message.operation'),
      dataIndex: "operation",
    }
  ]
  rules = {
    name: [{ required: true, message: t('message.pleaseName'), trigger: "blur" }],
    description: [{ required: true, message: t('message.pleaseEnterDescription'), trigger: "blur" }],
    image_url: [{ required: true, message: t('message.pleaseUploadPictures'), trigger: "change" }],
    minimum_stock: [{ required: true, message: t('message.pleaseEnterMinimumStock'), trigger: "change" }],
    price: [{ required: true, message: t('message.pleaseEnterPrice'), trigger: "change" }],
    quantity: [{ required: true, message: t('message.pleaseEnterQuantity'), trigger: "change" }],
    sku: [{ required: true, message: t('message.pleaseEnterSku'), trigger: "change" }],
    status: [{ required: true, message: t('message.pleaseEnterStatus'), trigger: "change" }],
    tags: [{ required: true, message: t('message.pleaseEnterTags'), trigger: "change" }],
    warranty_period: [{ required: true, message: t('message.pleaseEnterWarrantyPeriod'), trigger: "change" }],
  }
});
watch(
    () => form.tags,
    () => {
      (formRef.value as HTMLFormElement)?.clearValidate("tags");
    },
);
watch(
    () => form.image_url,
    () => {
      (formRef.value as HTMLFormElement)?.clearValidate("image_url");
    },
);
const pageList = async () => {
  const { data } = await ProductPage({name: nameStr.value, page: pagination.current, page_size: pagination.pageSize });
  pagination.total = data.data?.total || 0;
  list.value = data.data.data?.map((item: any, index: number) => ({
    key: index,
    id: item.ID,
    name: item.name,
    description: item.description,
    image_url: item.image_url,
    minimum_stock: item.minimum_stock,
    price: item.price,
    quantity: item.quantity,
    sku: item.sku,
    status: item.status,
    tags: item.tags,
    warranty_period: item.warranty_period
  }));
};

const edit = (record: any) => {
  Object.assign(form,record)
  title.value = t('message.edit');
  imageUrl.value = form.image_url;
  tags.value = form.tags;
  modalVisible.value = true;
};

const confirm = async (id: string) => {
  ProductDelete(id).then(async ({ data }) => {
    if (data.code === 20000) {
      message.success(data.message);
      await pageList();
    } else {
      message.success(data.message);
    }
  }).catch(e=>{
    console.error(e)
  });
};

const onAddData = async() => {
  (formRef.value as HTMLFormElement)
      .validate()
      .then(() => {
        const operateFun = form.id ? ProductUpdate : ProductCreate
        if(!form.id) delete form.id;
        operateFun({ ...form }).then(async ({ data }) => {
          if (data.code === 20000) {
            message.success(data.message);
            modalVisible.value = false;
            formRef.value?.resetFields();
            await pageList();
          } else {
            message.error(`${t('message.operationFailed')}:${data.data}`);
          }
        }).catch(e=>{
          console.error(e)
        });
      })
      .catch(e => {
        console.error(e)
      });
}

const handleTableChange = async (page: any) => {
  pagination.current = page.current;
  pagination.pageSize = page.pageSize;
  await pageList();
};

// 处理模态框关闭事件
const handleModalClose = () => {
  formRef.value?.resetFields();
  imageUrl.value = '';
  tags.value = '';
  Object.assign(form,{
    id: "",
    description: '',
    name: '',
    minimum_stock: '',
    image_url: '',
    price: '',
    quantity:'',
    sku:'',
    status:'',
    tags:'',
    warranty_period:''})
};

onMounted(async () => {
  await pageList();
});

</script>
<style lang="less">

</style>