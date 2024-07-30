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
          <template v-if="[ 'password', 'email'].includes(String(column.dataIndex))">
            <div>
              <a-input v-if="editableData[record.key] && column.dataIndex!=='password'" v-model:value="editableData[record.key][column.dataIndex]" style="margin: -5px 0" />
              <a-input-password v-else-if="editableData[record.key] && column.dataIndex==='password'" v-model:value="editableData[record.key][column.dataIndex]" style="margin: -5px 0" />
              <template v-else>
                <div v-if="column.dataIndex==='password'">*****</div>
                <div v-else>{{ text }}</div>
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
                <a-button type="primary" size="small" @click="showRoleBol = true, userId=record.id, getBindRole()" style="margin-left: 10px;">{{$t('message.assigningRoles')}}</a-button>
                                <a-button type="primary" size="small" @click="showDeptBol = true, userId=record.id,getBindDept()" style="margin-left: 10px;">{{$t('message.assigningDept')}}</a-button>
                <a-popconfirm :title="$t('message.sureDelete')" :okText="$t('message.yes')" :cancelText="$t('message.no')" @confirm="confirm(record.id)">
                  <a-button type="primary" size="small" danger style="margin-left: 10px;">{{$t('message.delete')}}</a-button>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>

      <!--      新增-->
      <a-modal :okText="$t('message.confirm')" :cancelText="$t('message.cancel')" v-model:open="modalVisible" :destroy-on-close="true" :title="$t('message.addition')" @cancel="handleCancel()" @ok="onAddData()">
        <a-form ref="formRef" :label-col="{ style: { width: '120px' } }" :labelWrap="true" :rules="rules" :model="form">
          <a-form-item :label="$t('message.username')" name="username">
            <a-input v-model:value="form.username" style="width: 350px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
          <a-form-item :label="$t('message.password')" name="password">
            <a-input-password v-model:value="form.password" style="width: 350px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
          <a-form-item :label="$t('message.email')" name="email">
            <a-input v-model:value="form.email" style="width: 350px" :placeholder="$t('message.pleaseEnter')" />
          </a-form-item>
        </a-form>
      </a-modal>

      <!--分配角色-->
      <a-modal :okText="$t('message.confirm')" :cancelText="$t('message.cancel')" v-model:open="showRoleBol" :destroy-on-close="true" :title="$t('message.assigningRoles')" @ok="onAddRole()">
        <a-transfer
            :data-source="roleData"
            :titles="[$t('message.unbound'), $t('message.bound')]"
            :target-keys="targetKeys"
            :locale="{ itemUnit: '', itemsUnit: '', notFoundContent: $t('message.noData') }"
            :showSelectAll="false"
            :selected-keys="selectedKeys"
            :render="item => item.title"
            @change="handleChange"
            @selectChange="handleSelectChange"
        />
      </a-modal>

      <!--绑定部门-->
      <a-modal :okText="$t('message.confirm')" :cancelText="$t('message.cancel')" v-model:open="showDeptBol" :destroy-on-close="true" title="分配部门" @ok="onAddDept()">
        <DeptSelect :multiple="true" v-model:value="deptId" style="width: 350px" />
      </a-modal>
    </a-card>
  </div>
</template>
<script lang="ts" setup>
import {onMounted, reactive, ref, UnwrapRef, watch} from 'vue'
import {
  RoleList,
  UserBindDept,
  UserBindRole,
  UserCreate,
  UserDelete,
  UserPage, UserQueryBindDept,
  UserQueryBindRole,
  UserUpdate
} from "@/api";
import {useI18n} from "vue-i18n";
import {Rule} from "ant-design-vue/es/form";
import {message} from "ant-design-vue";
import {cloneDeep} from "lodash-es";
import DeptSelect from "@/components/select/DeptSelect.vue";

interface DataItem {
  name: string;
  description: string;
  can_del: boolean;
}

const deptId = ref([]);
const showRoleBol = ref(false);
const showDeptBol = ref(false)
const roleData = ref([]);
const selectedKeys = ref([]);
const targetKeys = ref([]);
const userId = ref('');
const { t, locale } = useI18n();
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
const form = reactive({username: '', password: '', email: ''});
const list = ref([]);
const columns = ref([
  {
    title: t('message.uniCode'),
    dataIndex: "id",
  },
  {
    title: t('message.username'),
    dataIndex: "username",
  },
  {
    title: t('message.password'),
    dataIndex: "password",
  },
  {
    title: t('message.email'),
    dataIndex: "email",
  },
  {
    title: t('message.operation'),
    dataIndex: "operation",
  }
]);
let rules: Record<string, Rule[]> = {
  username: [{ required: true, message: t('message.pleaseName'), trigger: "blur" }],
  password: [{ required: true, message: t('message.pleaseEnterDescription'), trigger: "blur" }],
  email: [{ required: false, message: t('message.pleaseEnterDescription'), trigger: "blur" },
    { type: 'email', message: t('message.pleaseEnterEmailAddress'),trigger: "blur" }
  ],
};

watch(locale, () => {
  columns.value = [
    {
      title: t('message.uniCode'),
      dataIndex: "id",
    },
    {
      title: t('message.username'),
      dataIndex: "username",
    },
    {
      title: t('message.password'),
      dataIndex: "password",
    },
    {
      title: t('message.email'),
      dataIndex: "email",
    },
    {
      title: t('message.operation'),
      dataIndex: "operation",
    }
  ]
  rules = {
    username: [{ required: true, message: t('message.pleaseName'), trigger: "blur" }],
    password: [{ required: true, message: t('message.pleaseEnterDescription'), trigger: "blur" }],
    email: [{ required: false, message: t('message.pleaseEnterDescription'), trigger: "blur" },
      { type: 'email', message: t('message.pleaseEnterEmailAddress'),trigger: "blur" }
    ],
  }
});

const pageList = async () => {
  const { data } = await UserPage({name: nameStr.value, page: pagination.current, page_size: pagination.pageSize });
  pagination.total = data.data?.total || 0;
  list.value = data.data.data?.map((item: any, index: number) => ({
    key: index,
    id: item.ID,
    username: item.username,
    password: item.password,
    email: item.email
  }));
};

const edit = (key: string) => {
  editableData[key] = cloneDeep(list.value.filter((item) => key === item.key)[0]);
};
const save = async (key: string) => {
  Object.assign(list.value.filter((item) => key === item.key)[0], editableData[key]);
  const data = list.value.filter((item) => key === item.key)[0];
  delete editableData[key];
  await UserUpdate(data);
  await pageList();
};
const cancel = (key: string) => {
  delete editableData[key];
};

const confirm = async (id: string) => {
  UserDelete(id).then(async ({ data }) => {
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
        UserCreate({ ...form }).then(async ({ data }) => {
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


const onAddDept = async() => {
console.log(deptId.value);
  UserBindDept({ user_id: userId.value,dept_id: deptId.value, }).then(async ({ data }) => {
    if (data.code === 20000) {
      message.success(data.message);
      showDeptBol.value = false;
      await pageList();
    } else {
      message.error(`${t('message.operationFailed')}:${data.data}`);
    }
  }).catch(e=>{
    console.error(e)
  });
}
const handleCancel = ()=>{
  formRef.value?.resetFields();
}

const onAddRole = async()=> {
  UserBindRole({ role_id:targetKeys.value, user_id: userId.value }).then(async ({ data }) => {
    if (data.code === 20000) {
      message.success(data.message);
      showRoleBol.value = false;
      await pageList();
    } else {
      message.error(`${t('message.operationFailed')}:${data.data}`);
    }
  }).catch(e=>{
    console.error(e)
  });
}

const handleTableChange = async (page: any) => {
  pagination.current = page.current;
  pagination.pageSize = page.pageSize;
  await pageList();
};

const handleChange = (nextTargetKeys: string[], direction: string, moveKeys: string[]) => {
  targetKeys.value = nextTargetKeys;
};
const handleSelectChange = (sourceSelectedKeys: string[], targetSelectedKeys: string[]) => {
  selectedKeys.value = [...sourceSelectedKeys, ...targetSelectedKeys];
};

const getRoleList = async()=> {
  const { data } = await RoleList();
  roleData.value = data.data?.map((item: any)=>({key:item.ID,title:item.name}))
}


const getBindRole = async()=>{
  targetKeys.value = []
  const { data } = await UserQueryBindRole(userId.value);
  targetKeys.value = data.data?.map(it=>it.role_id)
}

const getBindDept = async()=> {
  deptId.value = []
  const { data } = await UserQueryBindDept(userId.value);
  deptId.value = data.data?.map(it=>it.dept_id)
}

onMounted(async () => {
  await pageList();
  await getRoleList();
});

</script>
<style lang="less">

</style>