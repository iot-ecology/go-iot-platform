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
      <a-button style="margin: 10px 0" type="primary" @click="modalVisible = true">{{ $t('message.addition') }}</a-button>

      <a-table :columns="columns" :data-source="list" bordered :pagination="pagination" @change="handleTableChange">
        <template #bodyCell="{ column, text, record }">
          <template v-if="column.dataIndex === 'start_date'">
            <div>{{formatDate(text)}}</div>
          </template>
          <template v-if="column.dataIndex === 'end_date'">
            <div>{{formatDate(text)}}</div>
          </template>
          <template v-if="column.dataIndex === 'operation'">
            <div class="editable-row-operations">
              <span v-if="editableData[record.key]">
                <a-button type="primary" size="small" style="margin-right: 10px" @click="save(record.key)">{{$t('message.save')}}</a-button>
                <a-popconfirm :title="$t('message.sureEdit')" @confirm="cancel(record.key)">
                  <a-button type="primary" size="small">{{$t('message.cancel')}}</a-button>
                </a-popconfirm>
              </span>
              <span v-else>
                <a-button type="primary" size="small" @click="edit(record)">{{$t('message.edit')}}</a-button>
                <a-popconfirm :title="$t('message.sureDelete')" :okText="$t('message.yes')" :cancelText="$t('message.no')" @confirm="confirm(record.id)">
                  <a-button type="primary" size="small" danger style="margin-left: 10px;">{{$t('message.delete')}}</a-button>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>
      <!--新增-->
      <a-modal style="width:700px" :okText="$t('message.confirm')" :cancelText="$t('message.cancel')" v-model:open="modalVisible" :destroy-on-close="true" :title="$t('message.addition')" @ok="onAddData()">
        <a-form ref="formRef" :label-col="{ style: { width: '150px' } }" :labelWrap="true" :rules="rules" :model="form">
          <a-form-item label="生产计划名称" name="name">
            <a-input v-model:value="form.name" style="width: 400px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
          <a-form-item label="生产计划开始日期" name="start_date">
            <a-date-picker v-model:value="form.start_date" style="width: 400px;" />
          </a-form-item>
          <a-form-item label="生产计划结束日期" name="end_date">
            <a-date-picker v-model:value="form.end_date" style="width: 400px;" />
          </a-form-item>
          <a-form-item label="状态" name="name">
            <a-select v-model:value="form.status" style="width: 400px;">
              <a-select-option value="1">准备中</a-select-option>
              <a-select-option value="2">进行中</a-select-option>
              <a-select-option value="3">已完成</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="生产计划描述" name="description">
            <a-textarea v-model:value="form.description" style="width: 400px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
          <a-form-item label="生产产品" name="">
            <div v-for="(item, index) in form.product_plans" :key="index" style="display: flex;align-items: center;margin-bottom: 10px">
                <Product v-model="item.product_id" style="width: 200px" />
                <a-input-number placeholder="请输入数量" v-model:value="item.quantity" style="width: 100px;margin: 0 10px" />
                <a-button shape="circle" @click="onUpdateProductFun(1,index)" style="margin-right: 8px">+</a-button>
                <a-button v-if="index!==0" shape="circle" @click="onUpdateProductFun(-1,index)">-</a-button>
            </div>
          </a-form-item>
        </a-form>
      </a-modal>
    </a-card>
  </div>
</template>
<script lang="ts" setup>
import {onMounted, reactive, ref, UnwrapRef, watch} from 'vue'
import {
  ProductionPlanCreate, ProductionPlanDelete, ProductionPlanDetail,
  ProductionPlanPage, ProductionPlanUpdate,
} from "@/api";
import {useI18n} from "vue-i18n";
import {Rule} from "ant-design-vue/es/form";
import {message} from "ant-design-vue";
import dayjs from "dayjs";
import  Product  from '@/components/select/Product.vue'
import {cloneDeep} from "lodash-es";

interface DataItem {
  name: string;
  description: string;
  can_del: boolean;
}
const { t,locale } = useI18n();
const formRef = ref<HTMLFormElement | null>(null);
const editableData: UnwrapRef<Record<string, DataItem>> = reactive({});
const pagination = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
  showSizeChanger: true, // 显示每页显示条目数选择器
});
const nameStr = ref('');
const modalVisible = ref(false);
const form = reactive({
  name: '',
  start_date: '',
  end_date: '',
  status:'',
  description: '',
  product_plans: [{product_id: '',quantity: '' }],
});
const list = ref([]);
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
    title: '生产计划开始日期',
    dataIndex: "start_date",
  },
  {
    title: '生产计划结束日期',
    dataIndex: "end_date",
  },
  {
    title: '描述',
    dataIndex: "description",
  },
  {
    title: '状态',
    dataIndex: "status_name",
  },
  {
    title: t('message.operation'),
    dataIndex: "operation",
  }
]);
const statusName = {'1':'准备中','2':'进行中','3':'已完成'}


let rules: Record<string, Rule[]> = {
  name: [{ required: true, message: t('message.pleaseEnterOperatorName'), trigger: "change" }],
  start_date: [{ required: true, message: t('message.pleaseEnterImsi'), trigger: "change" }],
  end_date: [{ required: true, message: t('message.pleaseEnterIccid'), trigger: "change" }],
  description: [{ required: true, message: t('message.pleaseEnterAccessNumber'), trigger: "change" }],
};

watch(locale, () => {
  columns.value = [
    {
      title: t('message.uniCode'),
      dataIndex: "id",
    },
    {
      title: t('message.operatorName'),
      dataIndex: "operator",
    },
    {
      title: t('message.imsi'),
      dataIndex: "imsi",
    },
    {
      title: t('message.iccid'),
      dataIndex: "iccid",
    },
    {
      title: t('message.accessNumber'),
      dataIndex: "access_number",
    },
    {
      title: t('message.expiration'),
      dataIndex: "expiration",
    },
    {
      title: t('message.operation'),
      dataIndex: "operation",
    }
  ]
  rules = {
    operator: [{ required: true, message: t('message.pleaseEnterOperatorName'), trigger: "change" }],
    access_number: [{ required: true, message: t('message.pleaseEnterImsi'), trigger: "change" }],
    expiration: [{ required: true, message: t('message.pleaseEnterIccid'), trigger: "change" }],
    iccid: [{ required: true, message: t('message.pleaseEnterAccessNumber'), trigger: "change" }],
    imsi: [{ required: true, message: t('message.pleaseEnterExpiration'), trigger: "change" }],
  }
});
const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD');

const pageList = async () => {
  const { data } = await ProductionPlanPage({name:nameStr.value, page: pagination.current, page_size: pagination.pageSize });
  pagination.total = data.data?.total || 0;
  list.value = data.data.data?.map((item: any, index: number) => ({
    key: index,
    id: item.ID,
    description: item.description,
    start_date: item.start_date,
    end_date: item.end_date,
    status: item.status,
    status_name: statusName[item.status],
    name: item.name
  }));
};

const edit = async(record: string) => {
  modalVisible.value = true
  const { data } = await ProductionPlanDetail(record.id);
  console.log(data.data);
  Object.assign(form, cloneDeep(data.data))
  console.log(form);
  // editableData[key] = cloneDeep(list.value.filter((item) => key === item.key)[0]);
};
const save = async (key: string) => {
  Object.assign(list.value.filter((item) => key === item.key)[0], editableData[key]);
  const data = list.value.filter((item) => key === item.key)[0];
  delete editableData[key];
  data.start_date = dayjs(data.start_date).format('YYYY-MM-DDTHH:mm:ssZ')
  data.end_date = dayjs(data.end_date).format('YYYY-MM-DDTHH:mm:ssZ')
  ProductionPlanUpdate(data).then(async ({ data }) => {
    if (data.code === 20000) {
      message.success(data.message);
      await pageList();
    } else {
      message.error(`${t('message.operationFailed')}:${data.data}`);
    }
  }).catch(e=>{
    console.error(e)
  });
};
const cancel = (key: string) => {
  delete editableData[key];
};

const confirm = async (id: string) => {
  ProductionPlanDelete(id).then(async ({ data }) => {
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
        console.log(form)
        const productFun = form.id ? ProductionPlanUpdate: ProductionPlanCreate
        productFun({ ...form }).then(async ({ data }) => {
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

/*新增生产产品*/
const onUpdateProductFun = (number: number, index: number) => {
  if( number > 0 ) {
    form.product_plans.push({product_id: '',quantity:''})
  } else {
    form.product_plans.splice(index,1)
  }
}
onMounted(async () => {
  await pageList();
});

</script>
<style lang="less" scoped>

</style>