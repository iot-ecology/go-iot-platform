<template>
  <div>
    <a-card title="" :bordered="true">
      <a-form layout="inline" :model="formState">
        <a-form-item :label="$t('message.scriptAlarm')">
          <SignalDelayWaring v-model="form.signal_delay_waring_id" style="width: 300px" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="pageList()">{{ $t('message.search') }}</a-button>
        </a-form-item>
      </a-form>
      <a-button style="margin: 10px 0" type="primary" @click="onAdd()">{{ $t('message.addition') }}</a-button>
      <!--      表格-->
      <a-table :columns="columns" :data-source="list" bordered :pagination="paginations" @change="handleTableChange">
        <template #bodyCell="{ column, text, record }">
          <template v-if="['name', 'signal_name', 'signal_id', 'mqtt_client_id', 'mqtt_client_name'].includes(column.dataIndex)">
            <div>
              <a-input
                v-if="editableData[record.key] && !['signal_id', 'mqtt_client_id'].includes(column.dataIndex)"
                v-model:value="editableData[record.key][column.dataIndex]"
                style="margin: -5px 0"
              />
              <MqttSelect v-else-if="editableData[record.key] && column.dataIndex == 'mqtt_client_id'" v-model="editableData[record.key][column.dataIndex]"></MqttSelect>
              <SignalSelect
                v-else-if="editableData[record.key] && column.dataIndex == 'signal_id'"
                v-model="editableData[record.key][column.dataIndex]"
                :mqtt_client_id="editableData[record.key]['mqtt_client_id']"
                name="ID"
                :number="true"
                @custom-event="handleCustomEvent"
              ></SignalSelect>
              <template v-else>
                <div v-if="column.dataIndex == 'mqtt_client_id'">{{ record.mqtt_client_name }}</div>
                <div v-else-if="column.dataIndex == 'signal_id'">{{ record.signal_name }}</div>
                <div v-else>{{ text }}</div>
              </template>
            </div>
          </template>
          <template v-else-if="column.dataIndex === 'operation'">
            <div class="editable-row-operations">
              <span v-if="editableData[record.key]">
                <a-button type="primary" size="small" style="margin-right: 10px" @click="save(record.key)">{{$t('message.save')}}</a-button>
                <a-popconfirm :okText="$t('message.yes')" :cancelText="$t('message.no')" :title="$t('message.sureEdit')" @confirm="cancel(record.key)">
                  <a-button type="primary" size="small">{{$t('message.cancel')}}</a-button>
                </a-popconfirm>
              </span>
              <span v-else>
                <a-button type="primary" size="small" @click="edit(record.key)">{{$t('message.edit')}}</a-button>
                <a-popconfirm :title="$t('message.sureDelete')" :okText="$t('message.yes')" :cancelText="$t('message.no')" @confirm="onDelete(record.ID)">
                  <a-button type="primary" size="small" danger style="margin-left: 10px;">{{$t('message.delete')}}</a-button>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>

      <a-modal v-model:open="modalVisible" :destroy-on-close="true" :title="title" class="custom-modal">
        <a-form ref="formRef" :label-col="{ style: { width: '100px' } }" :rules="rules" :model="form">
          <a-form-item :label="$t('message.name')" name="name">
            <a-input v-model:value="form.name" style="width: 350px" />
          </a-form-item>
          <a-form-item :label="$t('message.clientID')" name="mqtt_client_id">
            <MqttSelect v-model="form.mqtt_client_id" style="width: 350px" :show="true"></MqttSelect>
          </a-form-item>
          <a-form-item :label="$t('message.signalName')" name="signal_id">
            <SignalSelect v-model="form.signal_id" style="width: 350px" :mqtt_client_id="form.mqtt_client_id" name="ID" :show="true" :number="true" @custom-event="handleCustomEvent"></SignalSelect>
          </a-form-item>
        </a-form>
        <template #footer>
          <a-button @click="handleCancel">{{$t('message.cancel')}}</a-button>
          <a-button :disabled="loading" type="primary" @click="onAddData()">{{$t('message.confirm')}}</a-button>
        </template>
      </a-modal>
    </a-card>
  </div>
</template>
<script setup lang="ts">
import {onMounted, reactive, ref, type UnwrapRef, watch} from "vue";
import { type FormInstance, message } from "ant-design-vue";
import { type Rule } from "ant-design-vue/es/form";
import { cloneDeep } from "lodash-es";

import { SignalDelayWaringParamCreate, SignalDelayWaringParamDelete, SignalDelayWaringParamPage, SignalDelayWaringParamUpdate } from "@/api";
import { MqttSelect, SignalDelayWaring, SignalSelect } from "@/components/index.ts";
import {useI18n} from "vue-i18n";

interface DataItem {
  name: string;
  mqtt_client_id: string;
  signal_id: string;
}
const { t,locale } = useI18n();
let rules: Record<string, Rule[]> = {
  name: [
    {
      required: true,
      validator: async (rule, value) => {
        if (value) {
          await Promise.resolve();
        } else {
          await Promise.reject(t('message.pleaseName'));
        }
        if (/^[A-Za-z]/.test(value)) {
          await Promise.resolve();
        } else {
          await Promise.reject(t('message.startWithAnEnglish'));
        }
      },
      trigger: "blur",
    },
  ],
  mqtt_client_id: [{ required: true, message: t('message.pleaseSelectClientID'), trigger: "change" }],
  signal_id: [{ required: true, message: t('message.pleaseSignalName'), trigger: "change" }],
};
const title = ref(t('message.addition'));
const columns = ref([
  {
    title: t('message.uniCode'),
    dataIndex: "ID",
  },
  {
    title: t('message.name'),
    dataIndex: "name",
  },
  {
    title: t('message.clientID'),
    dataIndex: "mqtt_client_id",
    render: ({ record }) => {
      return record.mqtt_client_name;
    },
  },
  {
    title: t('message.signalName'),
    dataIndex: "signal_id",
    render: ({ record }) => {
      return record.signal_name;
    },
  },
  {
    title: t('message.operation'),
    dataIndex: "operation",
  },
]);
const formRef = ref<FormInstance>();
const modalVisible = ref(false);
const loading = ref(false);
const list = ref([]);
const editableData: UnwrapRef<Record<string, DataItem>> = reactive({});
const form = reactive({
  id: "",
  name: "",
  mqtt_client_id: "",
  signal_delay_waring_id: "",
  signal_name: "",
  signal_id: "",
});
const formState = reactive({ name: "" });
const signalName = ref("");

watch(
  () => form.signal_delay_waring_id,
  async () => {
    await pageList();
  },
);
watch(
  () => form.mqtt_client_id,
  () => {
    formRef.value.clearValidate("mqtt_client_id");
  },
);
watch(
  () => form.signal_id,
  () => {
    formRef.value.clearValidate("signal_id");
  },
);
watch(locale, () => {
  columns.value = [
    {
      title: t('message.uniCode'),
      dataIndex: "ID",
    },
    {
      title: t('message.name'),
      dataIndex: "name",
    },
    {
      title: t('message.clientID'),
      dataIndex: "mqtt_client_id",
      render: ({ record }) => {
        return record.mqtt_client_name;
      },
    },
    {
      title: t('message.signalName'),
      dataIndex: "signal_id",
      render: ({ record }) => {
        return record.signal_name;
      },
    },
    {
      title: t('message.operation'),
      dataIndex: "operation",
    },
  ]
  rules = {
    name: [
      {
        required: true,
        validator: async (rule, value) => {
          if (value) {
            await Promise.resolve();
          } else {
            await Promise.reject(t('message.pleaseName'));
          }
          if (/^[A-Za-z]/.test(value)) {
            await Promise.resolve();
          } else {
            await Promise.reject(t('message.startWithAnEnglish'));
          }
        },
        trigger: "blur",
      },
    ],
    mqtt_client_id: [{ required: true, message: t('message.pleaseSelectClientID'), trigger: "change" }],
    signal_id: [{ required: true, message: t('message.pleaseSignalName'), trigger: "change" }]
  }
});

const onAdd = () => {
  modalVisible.value = true;
  title.value = t('message.addition');
};
const paginations = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
  showSizeChanger: true, // 显示每页显示条目数选择器
});
const pageList = async () => {
  const { data } = await SignalDelayWaringParamPage({ signal_delay_waring_id: form.signal_delay_waring_id, page: paginations.current, page_size: paginations.pageSize });
  paginations.total = data.data?.total || 0;
  list.value = data.data.data?.map((item: any, index: number) => ({
    key: index,
    ID: item.ID,
    name: item.name,
    mqtt_client_id: item.mqtt_client_id,
    mqtt_client_name: item.mqtt_client_name,
    signal_name: item.signal_name,
    signal_id: item.signal_id,
  }));
};
const edit = (key: string) => {
  editableData[key] = cloneDeep(list.value.filter((item) => key === item.key)[0]);
};

const cancel = (key: string) => {
  // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
  delete editableData[key];
};

const handleTableChange = async (pagination: any) => {
  paginations.current = pagination.current;
  paginations.pageSize = pagination.pageSize;
  await pageList();
};
const onAddData = () => {
  formRef.value
    .validate()
    .then(() => {
      if (title.value === t('message.addition')) {
        const data = { ...form };
        delete data.id;
        SignalDelayWaringParamCreate(data).then(async ({ data }) => {
          if (data.code === 20000) {
            modalVisible.value = false;
            message.success(t('message.newSuccessfullyAdded'));
            formRef.value?.resetFields();
            await pageList();
          } else {
            message.error(data.message);
          }
        }).catch(e=>{
          console.error(e)
        });
      } else {
        SignalDelayWaringParamUpdate(form).then(async ({ data }) => {
          if (data.code === 20000) {
            modalVisible.value = false;
            message.success(t('message.editSuccessful'));
            await pageList();
          } else {
            message.error(data.message);
          }
        }).catch(e=>{
          console.error(e)
        });
      }
    })
    .catch(e => {
      console.error(e)
    });
};

const handleCancel = () => {
  modalVisible.value = false;
};

const save = async (key: string) => {
  const englishLetterRegex = /^[A-Za-z]$/;
  Object.assign(list.value.filter((item) => key === item.key)[0], editableData[key]);
  const data = list.value.filter((item) => key === item.key)[0];
  // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
  delete editableData[key];
  // eslint-disable-next-line no-debugger
  if (!data.mqtt_client_id || !data.signal_name) {
    message.error(t('message.clientSignal'));
    return;
  }
  if (!englishLetterRegex.test(data.name.charAt(0))) {
    message.error(t('message.startWithAnEnglish'));
    return;
  }

  data.signal_name = signalName.value;
  await SignalDelayWaringParamUpdate(data);
  await pageList();
};

// 删除
const onDelete = async (id: string) => {
  SignalDelayWaringParamDelete(id).then(async ({ data }) => {
    if (data.code === 20000) {
      message.success(data.message);
      await pageList();
    } else {
      // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
      message.success(data.message);
    }
  }).catch(e=>{
    console.error(e)
  });
};

const handleCustomEvent = (payload: any) => {
  if (payload?.value !== -11) {
    form.signal_name = payload.name;
    signalName.value = payload.name;
  }
};
onMounted(async ()=>{
  await pageList();
})
</script>
<style lang="less" scoped></style>
