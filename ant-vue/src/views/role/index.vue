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
          <template v-if="['name', 'description', 'can_del'].includes(String(column.dataIndex))">
            <div>
              <a-input v-if="editableData[record.key] && column.dataIndex==='name'" v-model:value="editableData[record.key][column.dataIndex]" style="margin: -5px 0" />
              <a-textarea v-else-if="editableData[record.key] && column.dataIndex==='description'" v-model:value="editableData[record.key][column.dataIndex]" style="margin: -5px 0" />
              <a-radio-group v-else-if="editableData[record.key] && column.dataIndex==='can_del'" name="radioGroup" v-model:value="editableData[record.key][column.dataIndex]">
                <a-radio :value="true">是</a-radio>
                <a-radio :value="false">否</a-radio>
              </a-radio-group>
              <template v-else>
                <div v-if="column.dataIndex!=='can_del'">{{ text }}</div>
                <div v-else>{{ text?'是':'否' }}</div>
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
                <a-popconfirm v-if="record.can_del" :title="$t('message.sureDelete')" :okText="$t('message.yes')" :cancelText="$t('message.no')" @confirm="confirm(record.id)">
                  <a-button type="primary" size="small" danger style="margin-left: 10px;">{{$t('message.delete')}}</a-button>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>
      <!--新增-->
      <a-modal :okText="$t('message.confirm')" :cancelText="$t('message.cancel')" v-model:open="modalVisible" :destroy-on-close="true" :title="$t('message.addition')" @cancel="handleCancel()" @ok="onAddData()">
        <a-form ref="formRef" :label-col="{ style: { width: '120px' } }" :labelWrap="true" :rules="rules" :model="form">
          <a-form-item :label="$t('message.name')" name="name">
            <a-input v-model:value="form.name" style="width: 350px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
          <a-form-item :label="$t('message.description')" name="description">
            <a-textarea v-model:value="form.description" style="width: 350px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
        </a-form>
      </a-modal>
    </a-card>
  </div>
</template>
<script lang="ts" setup>
import {onMounted, reactive, ref, UnwrapRef, watch} from 'vue'
import { RoleCreate, RoleDelete, RolePage, RoleUpdate} from "@/api";
import {useI18n} from "vue-i18n";
import {Rule} from "ant-design-vue/es/form";
import {message} from "ant-design-vue";
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
const form = reactive({can_del: true, description: '', name: ''});
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
    title: t('message.description'),
    dataIndex: "description",
  },
  {
    title: t('message.operation'),
    dataIndex: "operation",
  }
]);
let rules: Record<string, Rule[]> = {
  name: [{ required: true, message: t('message.pleaseName'), trigger: "blur" }],
  description: [{ required: true, message: t('message.pleaseEnterDescription'), trigger: "blur" }],
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
      title: t('message.description'),
      dataIndex: "description",
    },
    {
      title: t('message.operation'),
      dataIndex: "operation",
    }
  ]
  rules = {
    name: [{ required: true, message: t('message.pleaseName'), trigger: "blur" }],
    description: [{ required: true, message: t('message.pleaseEnterDescription'), trigger: "blur" }],
  }
});

const pageList = async () => {
  const { data } = await RolePage({name: nameStr.value, page: pagination.current, page_size: pagination.pageSize });
  pagination.total = data.data?.total || 0;
  list.value = data.data.data?.map((item: any, index: number) => ({
    key: index,
    id: item.ID,
    name: item.name,
    description: item.description,
    can_del: item.can_del
  }));
};

const edit = (key: string) => {
  editableData[key] = cloneDeep(list.value.filter((item) => key === item.key)[0]);
};
const save = async (key: string) => {
  Object.assign(list.value.filter((item) => key === item.key)[0], editableData[key]);
  const data = list.value.filter((item) => key === item.key)[0];
  delete editableData[key];
  await RoleUpdate(data);
  await pageList();
};
const cancel = (key: string) => {
  delete editableData[key];
};

const confirm = async (id: string) => {
  RoleDelete(id).then(async ({ data }) => {
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

const handleCancel = ()=>{
  formRef.value?.resetFields();
}


const onAddData = async() => {
  (formRef.value as HTMLFormElement)
      .validate()
      .then(() => {
        RoleCreate({ ...form }).then(async ({ data }) => {
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