<template>
  <div class="comp-preview">
    <a-card :bordered="true">
      <a-form layout="inline">
        <a-form-item label="计算规则">
          <CalculateSelect v-model="form.calc_rule_id"></CalculateSelect>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="pageList()">搜索</a-button>
        </a-form-item>
      </a-form>
      <a-button style="margin: 10px 0" type="primary" @click="modalVisible = true">新增</a-button>

      <a-table :columns="columns" :data-source="list" bordered :pagination="paginations" @change="handleTableChange">
        <template #bodyCell="{ column, text, record }">
          <template v-if="['name', 'reduce', 'signal_name', 'mqtt_client_id', 'mqtt_client_name'].includes(column.dataIndex)">
            <div>
              <a-input
                v-if="editableData[record.key] && !['reduce', 'signal_name', 'mqtt_client_id'].includes(column.dataIndex)"
                v-model:value="editableData[record.key][column.dataIndex]"
                style="margin: -5px 0"
              />
              <a-select v-else-if="editableData[record.key] && column.dataIndex == 'reduce'" v-model:value="editableData[record.key][column.dataIndex]" style="width: 200px">
                <a-select-option value="mean">平均值</a-select-option>
                <a-select-option value="sum">求和</a-select-option>
                <a-select-option value="min">最小值</a-select-option>
                <a-select-option value="max">最大值</a-select-option>
                <a-select-option value="原始">原始</a-select-option>
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
                <a-typography-link style="margin-right: 10px" @click="save(record.key)">保存</a-typography-link>
                <a-popconfirm title="确定取消编辑吗?" @confirm="cancel(record.key)">
                  <a>取消</a>
                </a-popconfirm>
              </span>
              <span v-else>
                <a @click="edit(record.key)">编辑</a>
                <a-popconfirm title="确认是否删除?" ok-text="是" cancel-text="否" @confirm="confirm(record.ID)">
                  <a style="margin-left: 10px; color: crimson">删除</a>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>

      <a-modal v-model:open="modalVisible" :destroy-on-close="true" title="新增" @ok="onAddData()" @cancel="clear()">
        <a-form ref="formRef" :label-col="{ style: { width: '100px' } }" :rules="rules" :model="form">
          <a-form-item label="名称" name="name">
            <a-input v-model:value="form.name" style="width: 350px" />
          </a-form-item>
          <a-form-item label="客户端ID" name="mqtt_client_id">
            <MqttSelect v-model="form.mqtt_client_id" style="width: 350px" :show="true"></MqttSelect>
          </a-form-item>
          <a-form-item label="信号名称" name="signal_id">
            <SignalSelect v-model="form.signal_id" style="width: 350px" :mqtt_client_id="form.mqtt_client_id" name="ID" :show="true" :number="true" @custom-event="handleCustomEvent"></SignalSelect>
          </a-form-item>
          <a-form-item label="聚合方式" name="reduce">
            <a-select v-model:value="form.reduce" style="width: 350px">
              <a-select-option value="mean">平均值</a-select-option>
              <a-select-option value="sum">求和</a-select-option>
              <a-select-option value="min">最小值</a-select-option>
              <a-select-option value="max">最大值</a-select-option>
              <a-select-option value="原始">原始</a-select-option>
              <a-select-option value="first">首条</a-select-option>
              <a-select-option value="last">尾条</a-select-option>
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
import { type FormInstance, message } from "ant-design-vue";
import type { Rule } from "ant-design-vue/es/form";
import { cloneDeep } from "lodash-es";

import { CalcParamCreate, CalcParamDelete, CalcParamPage, CalcParamUpdate } from "@/api";
import { CalculateSelect, MqttSelect, SignalSelect } from "@/components/index.ts";

interface DataItem {
  client_id: string;
  host: string;
  port: number;
  username: string;
}
const rules: Record<string, Rule[]> = {
  name: [
    {
      required: true,
      validator: async (rule, value) => {
        if (value) {
          await Promise.resolve();
        } else {
          await Promise.reject("请输入名称");
        }
        if (/^[A-Za-z]/.test(value)) {
          await Promise.resolve();
        } else {
          await Promise.reject("名称必须以英文字母开头");
        }
      },
      trigger: "blur",
    },
  ],
  mqtt_client_id: [{ required: true, message: "请选择客户端ID", trigger: "change" }],
  signal_id: [{ required: true, message: "请选择信号名称", trigger: "change" }],
  reduce: [{ required: true, message: "请选择聚合方式", trigger: "change" }],
};
const formRef = ref<FormInstance>();
const modalVisible = ref(false);
const form = reactive({ calc_rule_id: "", mqtt_client_id: "", signal_id: "", signal_name: "", name: "", reduce: "" });
const columns = [
  {
    title: "名称",
    dataIndex: "name",
  },
  {
    title: "客户端ID",
    dataIndex: "mqtt_client_id",
    render: ({ record }) => {
      return record.mqtt_client_name;
    },
  },
  {
    title: "信号名称",
    dataIndex: "signal_name",
  },
  {
    title: "聚合方式",
    dataIndex: "reduce",
  },
  {
    title: "操作",
    dataIndex: "operation",
  },
];
const list = ref([]);
const reduces = {
  mean: "平均值",
  sum: "求和",
  min: "最小值",
  max: "最大值",
  原始: "原始",
  first: "首条",
  last: "尾条",
};
const editableData: UnwrapRef<Record<string, DataItem>> = reactive({});
const paginations = reactive({
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
    formRef.value.clearValidate("mqtt_client_id");
  },
);
watch(
  () => form.signal_id,
  () => {
    formRef.value.clearValidate("signal_id");
  },
);

const onAddData = () => {
  formRef.value
    .validate()
    .then(() => {
      CalcParamCreate({ ...form }).then(async ({ data }) => {
        if (data.code === 20000) {
          message.success(data.message);
          modalVisible.value = false;
          formRef.value?.resetFields();
          await pageList();
        } else {
          // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
          message.error(`操作失败:${data.data}`);
        }
      });
    }).catch(e=>{
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
  // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
  delete editableData[key];
  // eslint-disable-next-line no-debugger
  if (!data.mqtt_client_id || !data.signal_name) {
    message.error("客户端ID和信号名称必选");
    await pageList();
    return;
  }
  await CalcParamUpdate(data);
  await pageList();
};
const cancel = (key: string) => {
  // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
  delete editableData[key];
};
const confirm = async (id: string) => {
  CalcParamDelete(id).then(async ({ data }) => {
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


const pageList = async () => {
  const { data } = await CalcParamPage({ rule_id: form.calc_rule_id, page: paginations.current, page_size: paginations.pageSize });
  paginations.total = data.data?.total || 0;
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

const handleTableChange = async (pagination: any) => {
  paginations.current = pagination.current;
  paginations.pageSize = pagination.pageSize;
  await pageList();
};
</script>

<style lang="less" scoped>
.comp-preview {
  width: 100%;
  height: 100%;
}
</style>
