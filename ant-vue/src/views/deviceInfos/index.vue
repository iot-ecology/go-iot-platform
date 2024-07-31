<template>
  <div class="comp-preview">
    <a-card :bordered="true">
      <a-button style="margin: 10px 0" type="primary" @click="modalVisible = true">{{ $t("message.addition") }}</a-button>

      <a-table :columns="columns" :data-source="list" bordered :pagination="pagination" @change="handleTableChange">
        <template #bodyCell="{ column, text, record }">
          <template v-if="['product_id', 'source', 'manufacturing_date', 'procurement_date', 'error_rate', 'push_interval'].includes(String(column.dataIndex))">
            <div>
              <template v-if="column.dataIndex === 'product_id'">
                <template v-if="editableData[record.key]">
                  <Product v-model="editableData[record.key][column.dataIndex]" style="width: 120px" />
                </template>
                <template v-else>
                  {{ record.product_name }}
                </template>
              </template>
              <!-- <template v-else-if="editableData[record.key] && column.dataIndex === 'sn'">
                <a-input v-model:value="editableData[record.key][column.dataIndex]" style="width: 80px" />
              </template> -->
              <template v-else-if="column.dataIndex === 'source'">
                <template v-if="editableData[record.key]">
                  <a-radio-group name="radioGroup" v-model:value="editableData[record.key][column.dataIndex]">
                    <a-radio :value="1">{{ $t("message.inside") }}</a-radio>
                    <a-radio :value="2">{{ $t("message.outside") }}</a-radio>
                  </a-radio-group>
                </template>
                <template v-else>
                  {{ record.source === 1 ? $t("message.inside") : $t("message.outside") }}
                </template>
              </template>
              <template v-else-if="editableData[record.key] && column.dataIndex === 'manufacturing_date' && editableData[record.key]['source'] == 1">
                <a-date-picker :valueFormat="dataFormat" show-time v-model:value="editableData[record.key][column.dataIndex]" style="width: 200px" />
              </template>
              <template v-else-if="editableData[record.key] && column.dataIndex === 'procurement_date' && editableData[record.key]['source'] == 2">
                <a-date-picker :valueFormat="dataFormat" show-time v-model:value="editableData[record.key][column.dataIndex]" style="width: 200px" />
              </template>
              <!-- <template v-else-if="editableData[record.key] && column.dataIndex === 'warranty_expiry'">
                <a-date-picker :valueFormat="dataFormat" v-model:value="editableData[record.key][column.dataIndex]" show-time style="width: 200px; margin: -5px 0" />
              </template> -->
              <template v-else-if="editableData[record.key] && column.dataIndex === 'error_rate'">
                <a-input-number v-model:value="editableData[record.key][column.dataIndex]" style="width: 80px" />
              </template>
              <template v-else-if="editableData[record.key] && column.dataIndex === 'push_interval'">
                <a-input-number v-model:value="editableData[record.key][column.dataIndex]" style="width: 80px" />
              </template>
              <template v-else>
                {{ text }}
              </template>
            </div>
          </template>
          <template v-else-if="column.dataIndex === 'operation'">
            <div class="editable-row-operations">
              <span v-if="editableData[record.key]">
                <a-button type="primary" size="small" style="margin-right: 10px" @click="save(record.key)">{{ $t("message.save") }}</a-button>
                <a-popconfirm :title="$t('message.sureEdit')" @confirm="cancel(record.key)">
                  <a-button type="primary" size="small">{{ $t("message.cancel") }}</a-button>
                </a-popconfirm>
              </span>
              <span v-else>
                <a-button type="primary" size="small" @click="edit(record.key)">{{ $t("message.edit") }}</a-button>
                <a-button type="primary" size="small" style="margin-left: 10px" @click="onWaringHistory(record)">{{ $t("message.alarmHistory") }}</a-button>
                <a-popconfirm :title="$t('message.sureDelete')" :okText="$t('message.yes')" :cancelText="$t('message.no')" @confirm="deleteById(record.ID)">
                  <a-button type="primary" size="small" danger style="margin-left: 10px">{{ $t("message.delete") }}</a-button>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>
      <a-modal
        :okText="$t('message.confirm')"
        :cancelText="$t('message.cancel')"
        v-model:open="modalVisible"
        :destroyOnClose="true"
        :title="$t('message.addition')"
        @ok="onAddData()"
        @cancel="onCancelData()"
      >
        <a-form ref="formRef" :label-col="{ style: { width: '140px' } }" :labelWrap="true" :rules="rules" :model="form">
          <a-form-item :label="$t('message.productId')" name="product_id">
            <!--            <a-input-number v-model:value="form.product_id" style="width: 200px;" />-->
            <Product v-model="form.product_id" style="width: 200px" />
          </a-form-item>
          <a-form-item :label="$t('message.deviceSn')" name="sn">
            <a-input v-model:value="form.sn" style="width: 200px" />
          </a-form-item>
          <a-form-item :label="$t('message.deviceSource')" name="source">
            <a-radio-group name="radioGroup" v-model:value="form.source">
              <a-radio :value="1">{{ $t("message.inside") }}</a-radio>
              <a-radio :value="2">{{ $t("message.outside") }}</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item v-if="form.source === 1" :label="$t('message.manufacturing_date')" name="manufacturing_date">
            <a-date-picker show-time v-model:value="form.manufacturing_date" style="width: 200px" />
          </a-form-item>
          <a-form-item v-if="form.source === 2" :label="$t('message.procurement_date')" name="procurement_date">
            <a-date-picker show-time v-model:value="form.procurement_date" style="width: 200px" />
          </a-form-item>
          <!-- <a-form-item :label="$t('message.warrantyExpiry')" name="warranty_expiry">
            <a-date-picker show-time v-model:value="form.warranty_expiry" style="width: 200px" />
          </a-form-item> -->
          <a-form-item :label="$t('message.errorRate')" name="error_rate">
            <a-input-number v-model:value="form.error_rate" style="width: 200px" />
          </a-form-item>
          <a-form-item :label="$t('message.pushInterval')" name="push_interval">
            <a-input-number v-model:value="form.push_interval" style="width: 200px" />
          </a-form-item>
        </a-form>
      </a-modal>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import dayjs from "dayjs";
import { onMounted, reactive, ref, UnwrapRef } from "vue";
import { message } from "ant-design-vue";
import { cloneDeep } from "lodash-es";
import { DEVICEINFO_CREATE, DEVICEINFO_DELETE, DEVICEINFO_PAGE, DeviceInfoUpdate } from "@/api";
import { useI18n } from "vue-i18n";
import Product from "@/components/select/Product.vue";
import { useColumns } from "@/hooks/useColumns";
import { useRoles } from "@/hooks/useRoles";
import { usePagination } from "@/hooks/usePagination";
import { useAddPage } from "@/hooks/useAddPage";
import { useDeletePage } from "@/hooks/useDeletePage";
import { useQueryPage } from "@/hooks/useQueryPage";

interface DataItem {
  client_id: string;
  host: string;
  port: number;
  username: string;
  password: string;
  subtopic: string;
}
const { t } = useI18n();
const formRef = ref<HTMLFormElement | null>(null);
const modalVisible = ref(false);
const modalTime = ref(false);
const editableData: UnwrapRef<Record<string, DataItem>> = reactive({});
const dataFormat = "YYYY-MM-DD HH:mm:ss";
const form = reactive({
  manufacturing_date: null,
  procurement_date: null,
  sn: "",
  product_id: "",
  // warranty_expiry: null,
  error_rate: 0,
  push_interval: 0,
  source: 1,
});

const { columns } = useColumns([
  {
    title: "message.uniCode",
    dataIndex: "ID",
  },
  {
    title: "message.productId",
    dataIndex: "product_id",
  },
  {
    title: "message.deviceSn",
    dataIndex: "sn",
  },
  {
    title: "message.deviceSource",
    dataIndex: "source",
  },
  {
    title: "message.manufacturing_date",
    dataIndex: "manufacturing_date",
  },
  {
    title: "message.procurement_date",
    dataIndex: "procurement_date",
  },
  {
    title: "message.warrantyExpiry",
    dataIndex: "warranty_expiry",
  },
  {
    title: "message.errorRate",
    dataIndex: "error_rate",
  },
  {
    title: "message.pushInterval",
    dataIndex: "push_interval",
  },
  {
    title: "message.operation",
    dataIndex: "operation",
  },
]);
const list = ref([]);

const { rules } = useRoles({
  manufacturing_date: [{ required: true, message: "message.manufacturing_date", trigger: "change" }],
  procurement_date: [{ required: true, message: "message.procurement_date", trigger: "change" }],
  // warranty_expiry: [{ required: false, message: "message.warrantyExpiry", trigger: "change" }],
  error_rate: [{ required: true, message: "message.errorRate", trigger: "change" }],
  push_interval: [{ required: true, message: "message.pushInterval", trigger: "change" }],
  source: [{ required: true, message: "message.deviceSource", trigger: "change" }],
  sn: [{ required: true, message: "message.deviceSn", trigger: "blur" }],
  product_id: [{ required: true, message: "message.productId", trigger: "blur" }],
});

const { pagination } = usePagination();
const formObj = reactive({ ID: "", up_time_start: 0, up_time_end: 0, date: "" });

const edit = (key: string) => {
  editableData[key] = cloneDeep(list.value.filter((item) => key === item.key)[0]);
};
const save = async (key: string) => {
  Object.assign(list.value.filter((item) => key === item.key)[0], editableData[key]);
  const data = list.value.filter((item) => key === item.key)[0];
  // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
  delete editableData[key];
  // await DeviceInfoUpdate(data);
  let params = {
    ...data,
    manufacturing_date: data.manufacturing_date ? dayjs(data.manufacturing_date).format("YYYY-MM-DDTHH:mm:ssZ") : undefined,
    procurement_date: data.procurement_date ? dayjs(data.procurement_date).format("YYYY-MM-DDTHH:mm:ssZ") : undefined,
    warranty_expiry: data.warranty_expiry ? dayjs(data.warranty_expiry).format("YYYY-MM-DDTHH:mm:ssZ") : undefined,
  };
  DeviceInfoUpdate(params).then(async ({ data }) => {
    if (data.code === 20000) {
      message.success(data.message);
      await pageList();
    } else {
      message.error(`${t("message.operationFailed")}:${data.data}`);
    }
  });
  await pageList();
};
const cancel = (key: string) => {
  delete editableData[key];
};

const { pageList } = useQueryPage({ path: DEVICEINFO_PAGE, type: "get" }, pagination, list);

const { onAddData, onCancelData } = useAddPage({ path: DEVICEINFO_CREATE, type: "post" }, formRef, form, modalVisible, pageList);

const { deleteById } = useDeletePage({ path: DEVICEINFO_DELETE, type: "delete" }, pageList);
const handleTableChange = async (page: any) => {
  pagination.current = page.current;
  pagination.pageSize = page.pageSize;
  await pageList();
};

const onWaringHistory = (record: any) => {
  formObj.ID = record.id;
  modalTime.value = true;
};
onMounted(async () => {
  await pageList();
});
</script>

<style lang="less" scoped>
.comp-preview {
  width: 100%;
  height: 100%;
}
</style>
