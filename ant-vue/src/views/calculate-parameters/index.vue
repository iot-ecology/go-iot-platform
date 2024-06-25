<template>
  <div class="comp-preview">
    <a-card :bordered="true">
      <a-form layout="inline">
        <a-form-item :label="$t('message.calculationRules')">
          <CalculateSelect v-model="form.calc_rule_id"></CalculateSelect>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="pageList()">{{ $t('message.search') }}</a-button>
        </a-form-item>
      </a-form>
      <a-button style="margin: 10px 0" type="primary" @click="modalVisible = true">{{ $t('message.addition') }}</a-button>

      <a-table :columns="columns" :data-source="list" bordered :pagination="pagination" @change="handleTableChange">
        <template #bodyCell="{ column, text, record }">
          <template v-if="['name', 'reduce', 'signal_name', 'mqtt_client_id', 'mqtt_client_name'].includes(column.dataIndex)">
            <div>
              <a-input
                v-if="editableData[record.key] && !['reduce', 'signal_name', 'mqtt_client_id'].includes(column.dataIndex)"
                v-model:value="editableData[record.key][column.dataIndex]"
                style="margin: -5px 0"
              />
              <a-select v-else-if="editableData[record.key] && column.dataIndex == 'reduce'" v-model:value="editableData[record.key][column.dataIndex]" style="width: 200px">
                <a-select-option value="mean">{{ $t('message.mean') }}</a-select-option>
                <a-select-option value="sum">{{ $t('message.sum') }}</a-select-option>
                <a-select-option value="min">{{ $t('message.min') }}</a-select-option>
                <a-select-option value="max">{{ $t('message.max') }}</a-select-option>
                <a-select-option value="原始">{{ $t('message.original') }}</a-select-option>
              </a-select>
              <MqttSelect v-else-if="editableData[record.key] && column.dataIndex == 'mqtt_client_id'" v-model="editableData[record.key][column.dataIndex]"></MqttSelect>
              <SignalSelect
                v-else-if="editableData[record.key] && column.dataIndex == 'signal_name'"
                v-model="editableData[record.key][column.dataIndex]"
                :mqtt_client_id="editableData[record.key]['mqtt_client_id']"
                name="name"
                :number="true"
              ></SignalSelect>
              <template v-else>
                <div v-if="column.dataIndex == 'mqtt_client_id'">{{ record.mqtt_client_name }}</div>
                <div v-else-if="column.dataIndex == 'reduce'">{{ reduces[record.reduce] }}</div>
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
                <a-popconfirm :title="$t('message.sureDelete')" :okText="$t('message.yes')" :cancelText="$t('message.no')" @confirm="confirm(record.ID)">
                  <a-button type="primary" size="small" danger style="margin-left: 10px;">{{$t('message.delete')}}</a-button>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>

      <a-modal v-model:open="modalVisible" :destroy-on-close="true" :title="$t('message.addition')" @ok="onAddData()" @cancel="clear()">
        <a-form ref="formRef" :label-col="{ style: { width: '110px' } }" :rules="rules" :model="form">
          <a-form-item :label="$t('message.name')" name="name">
            <a-input v-model:value="form.name" style="width: 350px" />
          </a-form-item>
          <a-form-item :label="$t('message.clientID')" name="mqtt_client_id">
            <MqttSelect v-model="form.mqtt_client_id" style="width: 350px" :show="true"></MqttSelect>
          </a-form-item>
          <a-form-item :label="$t('message.signalName')" name="signal_id">
            <SignalSelect v-model="form.signal_id" style="width: 350px" :mqtt_client_id="form.mqtt_client_id" name="ID" :show="true" :number="true" @custom-event="handleCustomEvent"></SignalSelect>
          </a-form-item>
          <a-form-item :label="$t('message.aggregationMethod')" name="reduce">
            <a-select v-model:value="form.reduce" style="width: 350px">
              <a-select-option value="mean">{{ $t('message.mean') }}</a-select-option>
              <a-select-option value="sum">{{ $t('message.sum') }}</a-select-option>
              <a-select-option value="min">{{ $t('message.min') }}</a-select-option>
              <a-select-option value="max">{{ $t('message.max') }}</a-select-option>
              <a-select-option value="原始">{{ $t('message.original') }}</a-select-option>
              <a-select-option value="first">{{ $t('message.first') }}</a-select-option>
              <a-select-option value="last">{{ $t('message.last') }}</a-select-option>
            </a-select>
          </a-form-item>
        </a-form>
      </a-modal>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import type { UnwrapRef } from "vue";
import { reactive, ref, watch } from "vue";
import { message } from "ant-design-vue";
import type { Rule } from "ant-design-vue/es/form";
import { cloneDeep } from "lodash-es";

import { CalcParamCreate, CalcParamDelete, CalcParamPage, CalcParamUpdate } from "@/api";
import { CalculateSelect, MqttSelect, SignalSelect } from "@/components/index.ts";
import {useI18n} from "vue-i18n";

interface DataItem {
  name: string;
  mqtt_client_name: string;
  signal_name: string;
  reduce: string;
}
const { t,locale } = useI18n();
let rules: Record<string, Rule[]> = {
  name: [
    {
      required: true,
      validator: async (_, value) => {
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
  mqtt_client_id: [{ required: true, message: t('message.pleaseClientID'), trigger: "change" }],
  signal_id: [{ required: true, message: t('message.pleaseSignalName'), trigger: "change" }],
  reduce: [{ required: true, message: t('message.pleaseAggregationMethod'), trigger: "change" }],
};
const formRef = ref<HTMLFormElement | null>(null);
const modalVisible = ref(false);
const form = reactive({ calc_rule_id: "", mqtt_client_id: "", signal_id: "", signal_name: "", name: "", reduce: "" });
let columns = [
  {
    title: t('message.name'),
    dataIndex: "name",
  },
  {
    title: t('message.clientID'),
    dataIndex: "mqtt_client_id",
    render: ({ record }:any) => {
      return record.mqtt_client_name;
    },
  },
  {
    title: t('message.signalName'),
    dataIndex: "signal_name",
  },
  {
    title: t('message.aggregationMethod'),
    dataIndex: "reduce",
  },
  {
    title: t('message.operation'),
    dataIndex: "operation",
  },
];
const list = ref([]);
const reduces = {
  mean: t('message.mean'),
  sum: t('message.sum'),
  min: t('message.max'),
  max: t('message.min'),
  原始: t('message.original'),
  first: t('message.first'),
  last: t('message.last'),
};
const editableData: UnwrapRef<Record<string, DataItem>> = reactive({});
const pagination = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
  showSizeChanger: true, // 显示每页显示条目数选择器
});
watch(
  () => form.calc_rule_id,
  async () => {
    await pageList();
  },
);
watch(
  () => form.mqtt_client_id,
  () => {
    (formRef.value as HTMLFormElement).clearValidate("mqtt_client_id");
  },
);
watch(
  () => form.signal_id,
  () => {
    (formRef.value as HTMLFormElement).clearValidate("signal_id");
  },
);
watch(locale, () => {
  columns = [
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
      dataIndex: "signal_name",
    },
    {
      title: t('message.aggregationMethod'),
      dataIndex: "reduce",
    },
    {
      title: t('message.operation'),
      dataIndex: "operation",
    }]
  rules = {
    name: [
      {
        required: true,
        validator: async (_, value) => {
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
    mqtt_client_id: [{ required: true, message: t('message.pleaseClientID'), trigger: "change" }],
    signal_id: [{ required: true, message: t('message.pleaseSignalName'), trigger: "change" }],
    reduce: [{ required: true, message: t('message.pleaseAggregationMethod'), trigger: "change" }],
  }
});

const onAddData = () => {
  (formRef.value as HTMLFormElement)
    .validate()
    .then(() => {
      CalcParamCreate({ ...form }).then(async ({ data }) => {
        if (data.code === 20000) {
          message.success(data.message);
          modalVisible.value = false;
          formRef.value?.resetFields();
          await pageList();
        } else {
          message.error(`${t('message.operationFailed')}:${data.data}`);
        }
      });
    }).catch((e:any)=>{
        console.error(e)
      });
};

const clear = () => {
  formRef.value?.resetFields();
};
const edit = (key: string) => {
  editableData[key] = cloneDeep(list.value.filter((item) => key === item.key)[0]);
};
const save = async (key: string) => {
  Object.assign(list.value.filter((item) => key === item.key)[0], editableData[key]);
  const data = list.value.filter((item) => key === item.key)[0];
  delete editableData[key];
  if (!data.mqtt_client_id || !data.signal_name) {
    message.error(t('message.clientSignal'));
    await pageList();
    return;
  }
  await CalcParamUpdate(data);
  await pageList();
};
const cancel = (key: string) => {
  delete editableData[key];
};
const confirm = async (id: string) => {
  CalcParamDelete(id).then(async ({ data }) => {
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
  const { data } = await CalcParamPage({ rule_id: form.calc_rule_id, page: pagination.current, page_size: pagination.pageSize });
  pagination.total = data.data?.total || 0;
  list.value = data.data.data?.map((item: any, index: number) => ({
    key: index,
    ID: item.ID,
    name: item.name,
    reduce: item.reduce,
    mqtt_client_id: item.mqtt_client_id,
    mqtt_client_name: item.mqtt_client_name,
    signal_name: item.signal_name,
  }));
};

const handleCustomEvent = (payload: any) => {
  if (payload.value !== -11) {
    form.signal_name = payload.name;
  }
};

const handleTableChange = async (page: any) => {
  pagination.current = page.current;
  pagination.pageSize = page.pageSize;
  await pageList();
};
</script>

<style lang="less" scoped>
.comp-preview {
  width: 100%;
  height: 100%;
}
</style>
