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
      <a-button style="margin: 10px 0" type="primary" @click="title =$t('message.addition') ,modalVisible = true">{{ $t('message.addition') }}</a-button>

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
      <a-modal style="width:700px" :okText="$t('message.confirm')" :cancelText="$t('message.cancel')" v-model:open="modalVisible" :destroy-on-close="true" @cancel="onClose()" :title="title" @ok="onAddData()">
        <a-form ref="formRef" :label-col="{ style: { width: '150px' } }" :labelWrap="true" :rules="rules" :model="form">
          <a-form-item :label="$t('message.productionPlanName')" name="name">
            <a-input :disabled="disabled" v-model:value="form.name" style="width: 400px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
          <a-form-item :label="$t('message.timeframe')" name="date">
            <a-range-picker :placeholder="[$t('message.startTime'), $t('message.endTime')]" :disabled="disabled" valueFormat="YYYY-MM-DD" @change="onDateFun" v-model:value="form.date" style="width: 400px;" />
          </a-form-item>
          <a-form-item :label="$t('message.status')" name="status">
            <a-select :disabled="disabled" v-model:value="form.status" style="width: 400px;">
              <a-select-option value="1">{{ $t('message.preparing') }}</a-select-option>
              <a-select-option value="2">{{ $t('message.ongoing') }}</a-select-option>
              <a-select-option value="3">{{ $t('message.finished') }}</a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item :label="$t('message.describe')" name="description">
            <a-textarea :disabled="disabled" v-model:value="form.description" style="width: 400px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
          <a-form-item :label="$t('message.productionName')" name="">
            <div v-for="(item, index) in form.product_plans" :key="index" style="display: flex;align-items: center;margin-bottom: 10px">
                <Product @dataUnit="dataUnit(index,$event)" :disabled="disabled" v-model="item.product_id" style="width: 200px" />
                <a-input-number :disabled="disabled" :placeholder="$t('message.pleaseEnterTheQuantity')" v-model:value="item.quantity" style="width: 100px;margin: 0 10px" />
                <div style="min-width: 40px" v-show="item.unit">（{{item.unit}}）</div>
                <a-button shape="circle" @click="onUpdateProductFun(1,index)" style="margin: 0 8px">+</a-button>
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
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';
dayjs.extend(utc);
dayjs.extend(timezone);
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
const disabled = ref(false);
const nameStr = ref('');
const modalVisible = ref(false);
const form = reactive({
  name: '',
  start_date: '',
  end_date: '',
  status:'',
  description: '',
  date: '',
  product_plans: [{product_id: '',quantity: '',unit: '' }],
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
    title: t('message.productionPlanStartDate'),
    dataIndex: "start_date",
  },
  {
    title: t('message.endDateOfProductionPlan'),
    dataIndex: "end_date",
  },
  {
    title: t('message.describe'),
    dataIndex: "description",
  },
  {
    title: t('message.status'),
    dataIndex: "status_name",
  },
  {
    title: t('message.operation'),
    dataIndex: "operation",
  }
]);
const title = ref('');


let statusName = {'1': t('message.preparing'),'2': t('message.ongoing'),'3': t('message.finished')}

let rules: Record<string, Rule[]> = {
  name: [{ required: true, message: t('message.pleaseName'), trigger: "change" }],
  date: [{ required: true, message: t('message.pleaseSelectADate'), trigger: "change" }],
  status: [{ required: true, message: t('message.pleaseSelectStatus'), trigger: "change" }],
  description: [{ required: true, message: t('message.pleaseEnterDescription'), trigger: "change" }],
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
      title: t('message.productionPlanStartDate'),
      dataIndex: "start_date",
    },
    {
      title: t('message.endDateOfProductionPlan'),
      dataIndex: "end_date",
    },
    {
      title: t('message.describe'),
      dataIndex: "description",
    },
    {
      title: t('message.status'),
      dataIndex: "status_name",
    },
    {
      title: t('message.operation'),
      dataIndex: "operation",
    }
  ]
  rules = {
    name: [{ required: true, message: t('message.pleaseName'), trigger: "change" }],
    date: [{ required: true, message: t('message.pleaseSelectADate'), trigger: "change" }],
    status: [{ required: true, message: t('message.pleaseSelectStatus'), trigger: "change" }],
    description: [{ required: true, message: t('message.pleaseEnterDescription'), trigger: "change" }],
  }
});
const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD');

const onDateFun = ()=> {
  form.start_date = dayjs(form.date[0]).format('YYYY-MM-DDTHH:mm:ss[Z]')
  form.end_date = dayjs(form.date[1]).format('YYYY-MM-DDTHH:mm:ss[Z]')
}

const onClose = ()=>{
  Object.assign(form,{
    name: '',
    start_date: '',
    end_date: '',
    status:'',
    description: '',
    date: '',
    product_plans: [{product_id: '',quantity: '',unit: '' }]})
}

const dataUnit = (index:number,unit:string)=>{
  form.product_plans[index].unit = unit
}

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
  const { data } = await ProductionPlanDetail(record.id);
  Object.assign(form, cloneDeep(data.data),{
    start_date:dayjs(data.data.start_date).format('YYYY-MM-DDTHH:mm:ss[Z]'),
    end_date:dayjs(data.data.end_date).format('YYYY-MM-DDTHH:mm:ss[Z]')
  })
  disabled.value = dayjs(data.data.start_date).isBefore(dayjs());
  form.date = [form.start_date,form.end_date]
  form.product_plans = cloneDeep(data.data.product_plans)
  modalVisible.value = true
  title.value = t('message.edit')
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