<template>
  <div class="comp-preview">
    <a-card :bordered="true">
      <a-form layout="inline">
        <a-form-item label="客户端ID">
          <MqttSelect v-model="form.mqtt_client_id"></MqttSelect>
        </a-form-item>
        <a-form-item label="信号名称">
          <SignalSelect v-model="form.signal_id" :mqtt_client_id="form.mqtt_client_id"></SignalSelect>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="pageList()">搜索</a-button>
        </a-form-item>
      </a-form>
      <a-button style="margin: 10px 0" type="primary" @click="modal1Visible = true">新增</a-button>

      <a-table :columns="columns" :data-source="list" bordered :pagination="paginations" @change="handleTableChange">
        <template #bodyCell="{ column, text, record }">
          <template v-if="['min', 'max', 'in_or_out'].includes(column.dataIndex)">
            <div>
              <a-switch
                v-if="editableData[record.key] && column.dataIndex === 'in_or_out'"
                v-model:checked="editableData[record.key][column.dataIndex]"
                checked-children="内报警"
                un-checked-children="外报警"
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
                <div v-if="column.dataIndex !== 'in_or_out'">{{ text }}</div>
                <div v-else>{{ text ? "内报警" : "外报警" }}</div>
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
                <a style="margin-left: 10px" @click="onWaringHistory(record)">报警历史</a>
                <a-popconfirm title="确认是否删除?" ok-text="是" cancel-text="否" @confirm="confirm(record.ID)">
                  <a style="margin-left: 10px; color: crimson">删除</a>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>

      <a-modal v-model:open="modalTime" title="时间范围" class="custom-modal">
        <a-spin tip="加载中..." size="large" :spinning="showSpinning">
          <a-form ref="formRefTime" :rules="rules" :model="formObj" name="nest-messages">
            <a-form-item label="时间范围" name="date">
              <a-range-picker v-model:value="formObj.date" show-time @change="bptjTimeChange" />
            </a-form-item>
          </a-form>
        </a-spin>
        <template #footer>
          <a-button v-if="!showSpinning" @click="modalTime = false">取消</a-button>
          <a-button :loading="showSpinning" type="primary" @click="setModalTime()">确定</a-button>
        </template>
      </a-modal>

      <a-modal v-model:open="modal1Visible" :destroy-on-close="true" title="新增" @ok="setModal1Visible()">
        <a-form ref="formRef" :label-col="{ style: { width: '100px' } }" :rules="rules" :model="form" name="nest-messages">
          <a-form-item label="最小值" name="min">
            <a-input-number v-model:value="form.min" style="width: 200px" :min="0" :max="form.max" />
          </a-form-item>
          <a-form-item label="最大值" name="max">
            <a-input-number v-model:value="form.max" style="width: 200px" :min="form.min" />
          </a-form-item>
          <a-form-item label="内报警外报警" name="checked">
            <a-switch v-model:checked="form.checked" checked-children="内报警" un-checked-children="外报警" @change="handleChange" />
          </a-form-item>
        </a-form>
      </a-modal>

      <a-modal v-model:open="modalHistory" :footer="null" :destroy-on-close="true" title="报警历史">
        <a-table bordered :pagination="false" :data-source="dataResult" :columns="columnsResult"> </a-table>
      </a-modal>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import type { UnwrapRef } from "vue";
import { reactive, ref, watch } from "vue";
import type { FormInstance } from "ant-design-vue";
import { message } from "ant-design-vue";
import { type Rule } from "ant-design-vue/es/form";
import dayjs from "dayjs";
import { cloneDeep } from "lodash-es";

import { SignalWaringConfigCreate, SignalWaringConfigDelete, SignalWaringConfigPage, SignalWaringConfigQueryRow, SignalWaringConfigUpdate } from "@/api";
import { MqttSelect, SignalSelect } from "@/components/index.ts";

interface DataItem {
  client_id: string;
  host: string;
  port: number;
  username: string;
}
const formRef = ref<FormInstance>();
const formRefTime = ref<FormInstance>();
const modal1Visible = ref(false);
const modalTime = ref(false);
const modalHistory = ref(false);
const form = reactive({ mqtt_client_id: "", signal_id: "", max: "", min: "", in_or_out: 1, checked: true });
const columns = [
  {
    title: "ID",
    dataIndex: "ID",
  },
  {
    title: "最大值",
    dataIndex: "max",
  },
  {
    title: "最小值",
    dataIndex: "min",
  },
  {
    title: "内报警外报警",
    dataIndex: "in_or_out",
    render: ({ record }) => {
      return record.in_or_out ? "内报警" : "外报警";
    },
  },
  {
    title: "操作",
    dataIndex: "operation",
  },
];
const list = ref([]);
const dataResult = ref([]);
const rules: Record<string, Rule[]> = {
  min: [{ required: true, message: "请输入最小值", trigger: "blur" }],
  max: [{ required: true, message: "请输入最大值", trigger: "blur" }],
  checked: [{ required: true, message: "请选择内报警外报警", trigger: "change" }],
  date: [{ required: true, message: "请选择时间", trigger: "change" }],
};
const showSpinning = ref(false);

const paginations = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
  showSizeChanger: true, // 显示每页显示条目数选择器
});
const formObj = reactive({ ID: "", up_time_start: "", up_time_end: "", date: "" });
const editableData: UnwrapRef<Record<string, DataItem>> = reactive({});
const columnsResult = ref([
  {
    title: "上报时间",
    dataIndex: "up_time",
  },
  {
    title: "数据值",
    dataIndex: "value",
  },
  {
    title: "处理时间",
    dataIndex: "insert_time",
  },
]);

watch([() => form.mqtt_client_id, () => form.signal_id], async ([newParam1, newParam2], [oldParam1, oldParam2]) => {
  if (newParam1 && newParam2) {
    await pageList();
  }
});

const setModal1Visible = () => {
  formRef.value
    .validate()
    .then(() => {
      if (!form.signal_id || !form.mqtt_client_id) {
        message.error("客户端ID和信号名称必选");
        return;
      }
      SignalWaringConfigCreate({ ...form }).then(({ data }) => {
        if (data.code === 20000) {
          message.success(data.message);
          modal1Visible.value = false;
          formRef.value?.resetFields();
          pageList();
        } else {
          // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
          message.error(`操作失败:${data.data}`);
        }
      });
    })
    .catch(() => {});
};

const setModalTime = () => {
  formRefTime.value
    .validate()
    .then(() => {
      showSpinning.value = true;
      SignalWaringConfigQueryRow(formObj)
        .then(({ data }) => {
          dataResult.value = data.data?.map(({ insert_time, up_time, value }) => ({
            insert_time: dayjs(insert_time * 1000).format("YYYY-MM-DD HH:mm:ss"),
            up_time: dayjs(up_time * 1000).format("YYYY-MM-DD HH:mm:ss"),
            value,
          }));
        })
        .finally(() => {
          modalHistory.value = true;
          modalTime.value = false;
          showSpinning.value = false;
          formObj.date = "";
        });
    })
    .catch(() => {});
};
const edit = (key: string) => {
  editableData[key] = cloneDeep(list.value.filter((item) => key === item.key)[0]);
};
const save = async (key: string) => {
  Object.assign(list.value.filter((item) => key === item.key)[0], editableData[key]);
  const data = list.value.filter((item) => key === item.key)[0];
  // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
  delete editableData[key];
  data.in_or_out = data.in_or_out ? 1 : 0;
  // eslint-disable-next-line no-debugger
  await SignalWaringConfigUpdate(data);
  await pageList();
};
const cancel = (key: string) => {
  // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
  delete editableData[key];
};
const confirm = async (id: string) => {
  SignalWaringConfigDelete(id).then(({ data }) => {
    if (data.code === 20000) {
      message.success(data.message);
      pageList();
    } else {
      // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
      message.success(data.message);
    }
  });
};

const pageList = async () => {
  const { data } = await SignalWaringConfigPage({ mqtt_client_id: form.mqtt_client_id, signal_id: form.signal_id, page: paginations.current, page_size: paginations.pageSize });
  paginations.total = data.data?.total || 0;
  list.value = data.data.data?.map((item: any, index: number) => ({
    key: index,
    ID: item.ID,
    mqtt_client_id: item.mqtt_client_id,
    signal: item.signal,
    signal_id: item.signal_id,
    in_or_out: item.in_or_out === 1,
    max: item.max,
    min: item.min,
  }));
};
const handleChange = () => {
  form.in_or_out = form.checked ? 1 : 0;
};

const handleTableChange = async (pagination: any) => {
  paginations.current = pagination.current;
  paginations.pageSize = pagination.pageSize;
  await pageList();
};

const onWaringHistory = (record: any) => {
  formObj.ID = record.ID;
  modalTime.value = true;
};
const bptjTimeChange = (date: any, dataString: any) => {
  formObj.up_time_start = dayjs(dataString[0]).unix();
  formObj.up_time_end = dayjs(dataString[1]).unix();
};
</script>

<style lang="less" scoped>
.comp-preview {
  width: 100%;
  height: 100%;
}
</style>