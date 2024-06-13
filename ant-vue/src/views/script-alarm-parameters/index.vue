<template>
  <div>
    <a-card title="" :bordered="true">
      <a-form layout="inline" :model="formState">
        <a-form-item label="脚本报警">
          <SignalDelayWaring v-model="form.signal_delay_waring_id" style="width: 300px" />
          <!--          <a-input v-model:value="formState.name" style="width: 300px" placeholder="请输入" />-->
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="pageList()">搜索</a-button>
        </a-form-item>
      </a-form>
      <a-button style="margin: 10px 0" type="primary" @click="onAdd()">新增</a-button>
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
                <a-typography-link style="margin-right: 10px" @click="save(record.key)">保存</a-typography-link>
                <a-popconfirm title="确定取消编辑吗?" @confirm="cancel(record.key)">
                  <a>取消</a>
                </a-popconfirm>
              </span>
              <span v-else>
                <a @click="edit(record.key)">编辑</a>
                <!--                <a style="margin-left: 10px" @click="onSignal(record.ID, record.mqtt_client_id)">信号报警配置</a>-->
                <a-popconfirm title="确认是否删除?" ok-text="是" cancel-text="否" @confirm="onDelete(record.ID)">
                  <a style="margin-left: 10px; color: crimson">删除</a>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>

      <a-modal v-model:open="modal1Visible" :destroy-on-close="true" :title="title" class="custom-modal">
        <a-form ref="formRef" :label-col="{ style: { width: '100px' } }" :rules="rules" :model="form" name="nest-messages">
          <a-form-item label="名称" name="name">
            <a-input v-model:value="form.name" style="width: 350px" />
          </a-form-item>
          <a-form-item label="客户端ID" name="mqtt_client_id">
            <MqttSelect v-model="form.mqtt_client_id" style="width: 350px" :show="true"></MqttSelect>
          </a-form-item>
          <a-form-item label="信号名称" name="signal_id">
            <SignalSelect v-model="form.signal_id" style="width: 350px" :mqtt_client_id="form.mqtt_client_id" name="ID" :show="true" :number="true" @custom-event="handleCustomEvent"></SignalSelect>
          </a-form-item>
        </a-form>
        <template #footer>
          <a-button @click="handleCancel">取消</a-button>
          <a-button :disabled="loading" type="primary" @click="setModal1Visible()">确定</a-button>
        </template>
      </a-modal>
    </a-card>
  </div>
</template>
<script setup lang="ts">
import { reactive, ref, type UnwrapRef, watch } from "vue";
import { type FormInstance, message } from "ant-design-vue";
import { type Rule } from "ant-design-vue/es/form";
import { cloneDeep } from "lodash-es";

import { SignalDelayWaringParamCreate, SignalDelayWaringParamDelete, SignalDelayWaringParamPage, SignalDelayWaringParamUpdate } from "@/api";
import { MqttSelect, SignalDelayWaring, SignalSelect } from "@/components/index.ts";

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
};
const title = ref("新增");
const columns = ref([
  {
    title: "ID",
    dataIndex: "ID",
  },
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
    dataIndex: "signal_id",
    render: ({ record }) => {
      return record.signal_name;
    },
  },
  {
    title: "操作",
    dataIndex: "operation",
  },
]);
const formRef = ref<FormInstance>();
const modal1Visible = ref(false);
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
  (newValue, oldValue) => {
    console.log(newValue);
    pageList();
  },
);
watch(
  () => form.mqtt_client_id,
  (newValue, oldValue) => {
    formRef.value.clearValidate("mqtt_client_id");
  },
);
watch(
  () => form.signal_id,
  (newValue, oldValue) => {
    formRef.value.clearValidate("signal_id");
  },
);

const onAdd = () => {
  modal1Visible.value = true;
  title.value = "新增";
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
pageList();
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
const setModal1Visible = () => {
  formRef.value
    .validate()
    .then(() => {
      if (title.value === "新增") {
        const data = { ...form };
        delete data.id;
        SignalDelayWaringParamCreate(data).then(({ data }) => {
          if (data.code === 20000) {
            modal1Visible.value = false;
            message.success("新增成功");
            formRef.value?.resetFields();
            pageList();
          } else {
            message.error(data.message);
          }
        });
      } else {
        SignalDelayWaringParamUpdate(form).then(({ data }) => {
          if (data.code === 20000) {
            modal1Visible.value = false;
            message.success("编辑成功");
            pageList();
          } else {
            message.error(data.message);
          }
        });
      }
    })
    .catch(() => {});
};

const handleCancel = () => {
  modal1Visible.value = false;
};

const save = async (key: string) => {
  const englishLetterRegex = /^[A-Za-z]$/;
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
  if (!englishLetterRegex.test(data.name.charAt(0))) {
    message.error("名称必须以英文字母开头");
    return;
  }

  data.signal_name = signalName.value;
  await SignalDelayWaringParamUpdate(data);
  await pageList();
};

// 删除
const onDelete = async (id: string) => {
  SignalDelayWaringParamDelete(id).then(({ data }) => {
    if (data.code === 20000) {
      message.success(data.message);
      pageList();
    } else {
      // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
      message.success(data.message);
    }
  });
};

const handleCustomEvent = (payload: any) => {
  if (payload?.value !== -11) {
    form.signal_name = payload.name;
    signalName.value = payload.name;
  }
};
</script>
<style lang="less" scoped></style>
