<template>
  <div class="comp-preview">
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
          <template v-if="['name', 'parent_name'].includes(String(column.dataIndex))">
            <div>
              <a-input v-if="editableData[record.key] && column.dataIndex==='name'" v-model:value="editableData[record.key][column.dataIndex]" style="margin: -5px 0" />
              <DeptSelect v-else-if="editableData[record.key]" :childValue="editableData[record.key]['id']" @set="setRoleFun(record.key, column.dataIndex, $event)" style="margin: -5px 0" />
              <template v-else>
                <div>{{ text }}</div>
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

      <a-modal :okText="$t('message.confirm')" :cancelText="$t('message.cancel')" v-model:open="modalVisible" :destroy-on-close="true" :title="$t('message.addition')" @ok="onAddData()">
        <a-form ref="formRef" :label-col="{ style: { width: '120px' } }" :labelWrap="true" :rules="rules" :model="form">
          <a-form-item :label="$t('message.name')" name="name">
            <a-input v-model:value="form.name" style="width: 350px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
          <a-form-item :label="$t('message.superior')" name="parent_id">
            <DeptSelect v-model:value="form.parent_id" style="width: 350px" />
          </a-form-item>
        </a-form>
      </a-modal>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import type { UnwrapRef } from "vue";
import {onMounted, reactive, ref, watch} from "vue";
import { message } from "ant-design-vue";
import { type Rule } from "ant-design-vue/es/form";
import { cloneDeep } from "lodash-es";
import DeptSelect from  '@/components/select/DeptSelect.vue';

import {
  DeptCreate,
  DeptDelete,
  DeptPage,
  DeptUpdate,
} from "@/api";
import {useI18n} from "vue-i18n";
const { t,locale } = useI18n();
interface DataItem {
  name: string;
  parent_id: number;
}
const formRef = ref<HTMLFormElement | null>(null);
const modalVisible = ref(false);
const form = reactive({ name: "", parent_id: null });
const nameStr = ref('')
let columns = ref([
  {
    title: t('message.uniCode'),
    dataIndex: "id",
  },
  {
    title: t('message.name'),
    dataIndex: "name",
  },
  {
    title: t('message.superior'),
    dataIndex: "parent_name",
  },
  {
    title: t('message.operation'),
    dataIndex: "operation",
  },
]);
const list = ref([]);
let rules: Record<string, Rule[]> = {
  name: [{ required: true, message: t('message.pleaseName'), trigger: "blur" }],
};

const pagination = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
  showSizeChanger: true, // 显示每页显示条目数选择器
});
const editableData: UnwrapRef<Record<string, DataItem>> = reactive({});
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
      title: t('message.superior'),
      dataIndex: "parent_name",
    },
    {
      title: t('message.operation'),
      dataIndex: "operation",
    },
  ]
  rules = {
    name: [{ required: true, message: t('message.pleaseName'), trigger: "blur" }],
  }
});

const setRoleFun = (key: any,index: any,value: any)=>{
  console.log(key,index,value)
  editableData[key]['parent_id'] = value;
}

const onAddData = () => {
  (formRef.value as HTMLFormElement)
      .validate()
      .then(() => {
        DeptCreate({ ...form }).then(async ({ data }) => {
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
  delete editableData[key];
  await DeptUpdate(data);
  await pageList();
};
const cancel = (key: string) => {
  delete editableData[key];
};
const confirm = async (id: string) => {
  DeptDelete(id).then(async ({ data }) => {
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
  const { data } = await DeptPage({name: nameStr.value, page: pagination.current, page_size: pagination.pageSize });
  pagination.total = data.data?.total || 0;
  list.value = data.data.data?.map((item: any, index: number) => ({
    key: index,
    id: item.ID,
    name: item.name,
    parent_id: item.parent_id,
    parent_name: item.parent_name
  }));
};

const handleTableChange = async (page: any) => {
  pagination.current = page.current;
  pagination.pageSize = page.pageSize;
  await pageList();
};
onMounted(async()=>{
  await pageList();
})
</script>

<style lang="less" scoped>
.comp-preview {
  width: 100%;
  height: 100%;
}
</style>
