<template>
  <div class="comp-preview">
    <a-card :bordered="true">
      <a-button style="margin: 10px 0" type="primary" @click="modalVisible = true">{{ $t('message.addition') }}</a-button>

      <a-table :columns="columns" :data-source="list" bordered :pagination="pagination" @change="handleTableChange">
        <template #bodyCell="{ column, text, record }">
          <template v-if="['min', 'source', 'in_or_out'].includes(String(column.dataIndex))">
            <div>
              <a-switch
                  v-if="editableData[record.key] && column.dataIndex === 'in_or_out'"
                  v-model:checked="editableData[record.key][column.dataIndex]"
                  :checked-children="$t('message.internalAlarm')"
                  :un-checked-children="$t('message.externalAlarm')"
              />
              <a-input-number
                  v-else-if="editableData[record.key] && column.dataIndex == 'min'"
                  v-model:value="editableData[record.key][column.dataIndex]"
                  :max="editableData[record.key]['max']"
                  style="margin: -5px 0"
              />
              <a-input-number
                  v-else-if="editableData[record.key] && column.dataIndex == 'max'"
                  v-model:value="editableData[record.key][column.dataIndex]"
                  :min="editableData[record.key]['min']"
                  style="margin: -5px 0"
              />
              <template v-else>
                <div v-if="column.dataIndex !== 'source'">{{ text }}</div>
                <div v-else>{{ text===1 ? '内部' : '外部' }}</div>
              </template>
            </div>
          </template>
          <template v-else-if="column.dataIndex === 'operation'">
            <div class="editable-row-operations">
              <span v-if="editableData[record.key]">
                <a-button type="primary" size="small" style="margin-right: 10px" @click="save(record.key)">{{$t('message.save')}}</a-button>
                <a-popconfirm :title="$t('message.sureEdit')" @confirm="cancel(record.key)">
                  <a-button type="primary" size="small">{{$t('message.cancel')}}</a-button>
                </a-popconfirm>
              </span>
              <span v-else>
                <a-button type="primary" size="small" @click="edit(record.key)">{{$t('message.edit')}}</a-button>
                <a-button type="primary" size="small" style="margin-left: 10px" @click="onWaringHistory(record)">{{ $t('message.alarmHistory') }}</a-button>
                <a-popconfirm :title="$t('message.sureDelete')" :okText="$t('message.yes')" :cancelText="$t('message.no')" @confirm="confirm(record.ID)">
                  <a-button type="primary" size="small" danger style="margin-left: 10px;">{{$t('message.delete')}}</a-button>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>
      <a-modal :okText="$t('message.confirm')" :cancelText="$t('message.cancel')" v-model:open="modalVisible" :destroy-on-close="true" :title="$t('message.addition')" @ok="onAddData()">
        <a-form ref="formRef" :label-col="{ style: { width: '120px' } }" :labelWrap="true" :rules="rules" :model="form">

          <a-form-item :label="$t('message.productId')" name="product_id">
<!--            <a-input-number v-model:value="form.product_id" style="width: 200px;" />-->
            <Product v-model="form.product_id"  style="width: 200px;" />
          </a-form-item>
          <a-form-item :label="$t('message.deviceSn')" name="sn">
            <a-input v-model:value="form.sn" style="width: 200px;" />
          </a-form-item>
          <a-form-item :label="$t('message.deviceSource')" name="source">
            <a-radio-group name="radioGroup" v-model:value="form.source">
              <a-radio :value="1">{{ $t('message.inside') }}</a-radio>
              <a-radio :value="2">{{ $t('message.outside') }}</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item v-if="form.source===1"  :label="$t('message.manufacturing_date')" name="manufacturing_date">
            <a-date-picker show-time v-model:value="form.manufacturing_date" style="width: 200px;" />
          </a-form-item>
          <a-form-item v-if="form.source===2" :label="$t('message.procurement_date')" name="procurement_date">
            <a-date-picker  show-time v-model:value="form.procurement_date" style="width: 200px;" />
          </a-form-item>
          <a-form-item :label="$t('message.warrantyExpiry')" name="warranty_expiry">
            <a-date-picker  show-time v-model:value="form.warranty_expiry" style="width: 200px;" />
          </a-form-item>

        </a-form>
      </a-modal>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import {onMounted, reactive, ref, watch,UnwrapRef} from "vue";
import { message } from "ant-design-vue";
import { type Rule } from "ant-design-vue/es/form";
import { cloneDeep } from "lodash-es";
import {
  DeviceInfoCreate,
  DeviceInfoDelete,
  DeviceInfoPage,
  DeviceInfoUpdate,
} from "@/api";
import {useI18n} from "vue-i18n";
import  Product  from '@/components/select/Product.vue'

interface DataItem {
  client_id: string;
  host: string;
  port: number;
  username: string;
  password: string;
  subtopic: string;
}
const { t,locale } = useI18n();
const formRef = ref<HTMLFormElement | null>(null);
const modalVisible = ref(false);
const modalTime = ref(false);
const editableData: UnwrapRef<Record<string, DataItem>> = reactive({});
const form = reactive({
  manufacturing_date: null,
  procurement_date: null,
  sn: "",
  product_id: "",
  warranty_expiry: null,
  source: 1
});
let columns = [
  {
    title: t('message.uniCode'),
    dataIndex: "id",
  },
  {
    title: t('message.productId'),
    dataIndex: "product_id",
  },
  {
    title: t('message.deviceSn'),
    dataIndex: "sn",
  },
  {
    title: t('message.deviceSource'),
    dataIndex: "source",
    render: ({ record }: any) => {
      return record.source===1 ? t('message.inside') : t('message.outside');
    },
  },
  {
    title: t('message.manufacturing_date'),
    dataIndex: "manufacturing_date",
  },
  {
    title: t('message.procurement_date'),
    dataIndex: "procurement_date",
  },
  {
    title: t('message.warrantyExpiry'),
    dataIndex: "warranty_expiry",
  },
  {
    title: t('message.operation'),
    dataIndex: "operation",
  }
];
const list = ref([]);
let rules: Record<string, Rule[]> = {
  manufacturing_date: [{ required: true, message: t('message.manufacturing_date'), trigger: "change" }],
  procurement_date: [{ required: true, message: t('message.procurement_date'), trigger: "change" }],
  warranty_expiry: [{ required: false, message: t('message.warrantyExpiry'), trigger: "change" }],
  source: [{ required: true, message: t('message.deviceSource'), trigger: "change" }],
  sn: [{ required: true, message: t('message.deviceSn'), trigger: "blur" }],
  product_id: [{ required: true, message: t('message.productId'), trigger: "blur" }],
};

const pagination = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
  showSizeChanger: true, // 显示每页显示条目数选择器
});
const formObj = reactive({ ID: "", up_time_start: 0, up_time_end: 0, date: "" });
const columnsResult = ref([
  {
    title: t('message.reportingTime'),
    dataIndex: "up_time",
  },
  {
    title: t('message.dataValue'),
    dataIndex: "value",
  },
  {
    title: t('message.processingTime'),
    dataIndex: "insert_time",
  },
]);
watch(locale, () => {
  columns = [
    {
      title: t('message.uniCode'),
      dataIndex: "ID",
    },
    {
      title: t('message.manufacturing_date'),
      dataIndex: "manufacturing_date",
    },
    {
      title: t('message.procurement_date'),
      dataIndex: "procurement_date",
    },
    {
      title: t('message.productId'),
      dataIndex: "product_name",
    },
    {
      title: t('message.deviceSn'),
      dataIndex: "sn",
    },
    {
      title: t('message.deviceSource'),
      dataIndex: "source",
      render: ({ record }: any) => {
        return record.source===1 ? t('message.inside') : t('message.outside');
      },
    },
    {
      title: t('message.warrantyExpiry'),
      dataIndex: "warranty_expiry",
    },
    {
      title: t('message.operation'),
      dataIndex: "operation",
    }
  ]
  rules = {
    manufacturing_date: [{ required: true, message: t('message.pleaseMinimum'), trigger: "change" }],
    procurement_date: [{ required: true, message: t('message.pleaseMaximum'), trigger: "change" }],
    source: [{ required: true, message: t('message.pleaseAlarm'), trigger: "change" }],
    warranty_expiry: [{ required: true, message: t('message.pleaseTime'), trigger: "change" }],
    sn: [{ required: true, message: t('message.pleaseTime'), trigger: "change" }],
  }
  columnsResult.value = [
    {
      title: t('message.reportingTime'),
      dataIndex: "up_time",
    },
    {
      title: t('message.dataValue'),
      dataIndex: "value",
    },
    {
      title: t('message.processingTime'),
      dataIndex: "insert_time",
    },
  ]
});

const onAddData = () => {
  (formRef.value as HTMLFormElement)
      .validate()
      .then(() => {
        let data = { ...form }
        // if(!data.procurement_date) delete  data.procurement_date;
        // if(!data.manufacturing_date) delete  data.manufacturing_date;
        DeviceInfoCreate({ ...data }).then(async ({ data }) => {
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
};

const edit = (key: string) => {
  editableData[key] = cloneDeep(list.value.filter((item) => key === item.key)[0]);
};
const save = async (key: string) => {
  Object.assign(list.value.filter((item) => key === item.key)[0], editableData[key]);
  const data = list.value.filter((item) => key === item.key)[0];
  // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
  delete editableData[key];
  data.source = data.source ? 1 : 0;
  await DeviceInfoUpdate(data);
  await pageList();
};
const cancel = (key: string) => {
  delete editableData[key];
};
const confirm = async (id: string) => {
  DeviceInfoDelete(id).then(async ({ data }) => {
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

const pageList = async () => {
  const { data } = await DeviceInfoPage({ page: pagination.current, page_size: pagination.pageSize });
  console.log(data)
  pagination.total = data.data?.total || 0;
  list.value = data.data.data?.map((item: any, index: number) => ({
    key: index,
    id: item.ID,
    manufacturing_date: item.manufacturing_date,
    procurement_date: item.procurement_date,
    product_id: item.product_id,
    source: item.source === 1,
    sn: item.sn,
    warranty_expiry: item.warranty_expiry,
  }));
};

const handleTableChange = async (page: any) => {
  pagination.current = page.current;
  pagination.pageSize = page.pageSize;
  await pageList();
};

const onWaringHistory = (record: any) => {
  formObj.ID = record.ID;
  modalTime.value = true;
};
onMounted(async()=>{
  await pageList()
})
</script>

<style lang="less" scoped>
.comp-preview {
  width: 100%;
  height: 100%;
}
</style>
