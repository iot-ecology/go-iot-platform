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
          <template v-if="['operator', 'access_number', 'expiration', 'iccid', 'imsi'].includes(String(column.dataIndex))">
            <div>
              <a-input v-if="editableData[record.key] && column.dataIndex!=='expiration'" v-model:value="editableData[record.key][column.dataIndex]" style="margin: -5px 0" />
              <a-date-picker valueFormat="YYYY-MM-DD HH:mm:ss"  v-else-if="editableData[record.key] && column.dataIndex==='expiration'" v-model:value="editableData[record.key][column.dataIndex]" show-time style="margin: -5px 0" />
              <template v-else>
                <div>{{ column.dataIndex === 'expiration'?formatDate(text):text }}</div>
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
                <a-popconfirm :title="$t('message.sureDelete')" :okText="$t('message.yes')" :cancelText="$t('message.no')" @confirm="confirm(record.id)">
                  <a-button type="primary" size="small" danger style="margin-left: 10px;">{{$t('message.delete')}}</a-button>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>
      <!--新增-->
      <a-modal :okText="$t('message.confirm')" :cancelText="$t('message.cancel')" v-model:open="modalVisible" :destroy-on-close="true" :title="$t('message.addition')" @ok="onAddData()">
        <a-form ref="formRef" :label-col="{ style: { width: '150px' } }" :labelWrap="true" :rules="rules" :model="form">
          <a-form-item label="运营商名称" name="operator">
            <a-input v-model:value="form.operator" style="width: 300px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
          <a-form-item label="接入号" name="access_number">
            <a-input v-model:value="form.access_number" style="width: 300px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
          <a-form-item label="到期时间" name="expiration">
            <a-date-picker show-time v-model:value="form.expiration" style="width: 300px;" />
          </a-form-item>
          <a-form-item label="集成电路卡识别码" name="iccid">
            <a-input v-model:value="form.iccid" style="width: 300px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
          <a-form-item label="国际移动用户识别码" name="imsi">
            <a-input v-model:value="form.imsi" style="width: 300px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
        </a-form>
      </a-modal>
    </a-card>
  </div>
</template>
<script lang="ts" setup>
import {onMounted, reactive, ref, UnwrapRef, watch} from 'vue'
import {
  SimCardCreate,
  SimCardDelete,
  SimCardPage,
  SimCardUpdate
} from "@/api";
import {useI18n} from "vue-i18n";
import {Rule} from "ant-design-vue/es/form";
import {message} from "ant-design-vue";
import {cloneDeep} from "lodash-es";
import dayjs from "dayjs";

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
  access_number: '',
  expiration: '',
  iccid: '',
  imsi: '',
  operator: '',
});
const list = ref([]);
const columns = ref([
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
]);
let rules: Record<string, Rule[]> = {
  operator: [{ required: true, message: t('message.pleaseEnterOperatorName'), trigger: "change" }],
  access_number: [{ required: true, message: t('message.pleaseEnterImsi'), trigger: "change" }],
  expiration: [{ required: true, message: t('message.pleaseEnterIccid'), trigger: "change" }],
  iccid: [{ required: true, message: t('message.pleaseEnterAccessNumber'), trigger: "change" }],
  imsi: [{ required: true, message: t('message.pleaseEnterExpiration'), trigger: "change" }],
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
const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss');

const pageList = async () => {
  const { data } = await SimCardPage({AccessNumber: nameStr.value, page: pagination.current, page_size: pagination.pageSize });
  pagination.total = data.data?.total || 0;
  list.value = data.data.data?.map((item: any, index: number) => ({
    key: index,
    id: item.ID,
    operator: item.operator,
    imsi: item.imsi,
    iccid: item.iccid,
    expiration: formatDate(item.expiration),
    access_number: item.access_number
  }));
  console.log(list.value)
};

const edit = (key: string) => {
  editableData[key] = cloneDeep(list.value.filter((item) => key === item.key)[0]);
};
const save = async (key: string) => {
  Object.assign(list.value.filter((item) => key === item.key)[0], editableData[key]);
  const data = list.value.filter((item) => key === item.key)[0];
  delete editableData[key];
  data.expiration = dayjs(data.expiration).format('YYYY-MM-DDTHH:mm:ssZ')
  await SimCardUpdate(data);
  await pageList();
};
const cancel = (key: string) => {
  delete editableData[key];
};

const confirm = async (id: string) => {
  SimCardDelete(id).then(async ({ data }) => {
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
        SimCardCreate({ ...form }).then(async ({ data }) => {
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

onMounted(async () => {
  await pageList();
});

</script>
<style lang="less">

</style>