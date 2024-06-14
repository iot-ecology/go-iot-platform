<template>
  <div class="comp-preview">
    <a-card :bordered="true">
      <a-form layout="inline">
        <a-form-item :label="$t('message.clientID')">
          <MqttSelect v-model="form.mqtt_client_id"></MqttSelect>
        </a-form-item>
        <a-form-item :label="$t('message.signalName')">
          <SignalSelect v-model="form.signal_id" :mqtt_client_id="form.mqtt_client_id"></SignalSelect>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="pageList()">{{ $t('message.search') }}</a-button>
        </a-form-item>
      </a-form>
      <a-button style="margin: 10px 0" type="primary" @click="modalVisible = true">{{ $t('message.addition') }}</a-button>

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
                <a-typography-link style="margin-right: 10px" @click="save(record.key)">{{$t('message.save')}}</a-typography-link>
                <a-popconfirm :title="$t('message.sureEdit')" @confirm="cancel(record.key)">
                  <a>{{$t('message.cancel')}}</a>
                </a-popconfirm>
              </span>
              <span v-else>
                <a @click="edit(record.key)">{{$t('message.edit')}}</a>
                <a style="margin-left: 10px" @click="onWaringHistory(record)">{{ $t('message.alarmHistory') }}</a>
                <a-popconfirm :title="$t('message.sureDelete')" :okText="$t('message.yes')" :cancelText="$t('message.no')" @confirm="confirm(record.ID)">
                  <a style="margin-left: 10px; color: crimson">{{$t('message.delete')}}</a>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>

      <a-modal v-model:open="modalTime" :title="$t('message.timeframe')" class="custom-modal">
        <a-spin :tip="$t('message.loading')" size="large" :spinning="showSpinning">
          <a-form ref="formRefTime" :rules="rules" :model="formObj">
            <a-form-item :label="$t('message.timeframe')" name="date">
              <a-range-picker v-model:value="formObj.date" show-time @change="bptjTimeChange" />
            </a-form-item>
          </a-form>
        </a-spin>
        <template #footer>
          <a-button v-if="!showSpinning" @click="modalTime = false">{{$t('message.cancel')}}</a-button>
          <a-button :loading="showSpinning" type="primary" @click="getHistoryData()">{{$t('message.confirm')}}</a-button>
        </template>
      </a-modal>

      <a-modal :okText="$t('message.confirm')" :cancelText="$t('message.cancel')" v-model:open="modalVisible" :destroy-on-close="true" :title="$t('message.addition')" @ok="onAddData()">
        <a-form ref="formRef" :label-col="{ style: { width: '120px' } }" :labelWrap="true" :rules="rules" :model="form">
          <a-form-item :label="$t('message.min')" name="min">
            <a-input-number v-model:value="form.min" style="width: 200px" :min="0" :max="form.max" />
          </a-form-item>
          <a-form-item :label="$t('message.max')" name="max">
            <a-input-number v-model:value="form.max" style="width: 200px" :min="form.min" />
          </a-form-item>
          <a-form-item :label="$t('message.internalExternalAlarm')" name="checked">
            <a-switch v-model:checked="form.checked" :checked-children="$t('message.internalAlarm')" :un-checked-children="$t('message.externalAlarm')" @change="handleChange" />
          </a-form-item>
        </a-form>
      </a-modal>

      <a-modal :okText="$t('message.confirm')" :cancelText="$t('message.cancel')" v-model:open="modalHistory" :footer="null" :destroy-on-close="true" :title="$t('message.alarmHistory')">
        <a-table bordered :pagination="false" :data-source="dataResult" :columns="columnsResult"> </a-table>
      </a-modal>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import type { UnwrapRef } from "vue";
import { reactive, ref, watch} from "vue";
import type { FormInstance } from "ant-design-vue";
import { message } from "ant-design-vue";
import { type Rule } from "ant-design-vue/es/form";
import dayjs from "dayjs";
import { cloneDeep } from "lodash-es";

import { SignalWaringConfigCreate, SignalWaringConfigDelete, SignalWaringConfigPage, SignalWaringConfigQueryRow, SignalWaringConfigUpdate } from "@/api";
import { MqttSelect, SignalSelect } from "@/components/index.ts";
import {useI18n} from "vue-i18n";
const { t,locale } = useI18n();
interface DataItem {
  max: number;
  min: number;
  in_or_out: boolean;
}
const formRef = ref<FormInstance>();
const formRefTime = ref<FormInstance>();
const modalVisible = ref(false);
const modalTime = ref(false);
const modalHistory = ref(false);
const form = reactive({ mqtt_client_id: "", signal_id: "", max: "", min: "", in_or_out: 1, checked: true });
let columns = [
  {
    title: t('message.uniCode'),
    dataIndex: "ID",
  },
  {
    title: t('message.max'),
    dataIndex: "max",
  },
  {
    title: t('message.min'),
    dataIndex: "min",
  },
  {
    title: t('message.internalExternalAlarm'),
    dataIndex: "in_or_out",
    render: ({ record }) => {
      return record.in_or_out ? t('message.internalAlarm') : t('message.externalAlarm');
    },
  },
  {
    title: t('message.operation'),
    dataIndex: "operation",
  },
];
const list = ref([]);
const dataResult = ref([]);
let rules: Record<string, Rule[]> = {
  min: [{ required: true, message: t('message.pleaseMinimum'), trigger: "blur" }],
  max: [{ required: true, message: t('message.pleaseMaximum'), trigger: "blur" }],
  checked: [{ required: true, message: t('message.pleaseAlarm'), trigger: "change" }],
  date: [{ required: true, message: t('message.pleaseTime'), trigger: "change" }],
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
    title: t('message.reportingTime'),
    dataIndex: "up_time",
  },
  {
    title: t('message.dataValue'),
    dataIndex: "value",
  },
  {
    title: t('message.processingTime'),
    dataIndex: "insert_time",
  },
]);

watch([() => form.mqtt_client_id, () => form.signal_id], async ([newParam1, newParam2], []) => {
  if (newParam1 && newParam2) {
    await pageList();
  }
});
watch(locale, () => {
  columns = [
    {
      title: t('message.uniCode'),
      dataIndex: "ID",
    },
    {
      title: t('message.max'),
      dataIndex: "max",
    },
    {
      title: t('message.min'),
      dataIndex: "min",
    },
    {
      title: t('message.internalExternalAlarm'),
      dataIndex: "in_or_out",
      render: ({ record }) => {
        return record.in_or_out ? t('message.internalAlarm') : t('message.externalAlarm');
      },
    },
    {
      title: t('message.operation'),
      dataIndex: "operation",
    },
  ]
  rules = {
    min: [{ required: true, message: t('message.pleaseMinimum'), trigger: "blur" }],
    max: [{ required: true, message: t('message.pleaseMaximum'), trigger: "blur" }],
    checked: [{ required: true, message: t('message.pleaseAlarm'), trigger: "change" }],
    date: [{ required: true, message: t('message.pleaseTime'), trigger: "change" }],
  }
  columnsResult.value = [
    {
      title: t('message.reportingTime'),
      dataIndex: "up_time",
    },
    {
      title: t('message.dataValue'),
      dataIndex: "value",
    },
    {
      title: t('message.processingTime'),
      dataIndex: "insert_time",
    },
  ]
});

const onAddData = () => {
  formRef.value
    .validate()
    .then(() => {
      if (!form.signal_id || !form.mqtt_client_id) {
        message.error(t('message.clientSignal'));
        return;
      }
      SignalWaringConfigCreate({ ...form }).then(async ({ data }) => {
        if (data.code === 20000) {
          message.success(data.message);
          modalVisible.value = false;
          formRef.value?.resetFields();
          await pageList();
        } else {
          // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
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

const getHistoryData = () => {
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
        .catch(e=>{
            console.error(e)
          })
        .finally(() => {
          modalHistory.value = true;
          modalTime.value = false;
          showSpinning.value = false;
          formObj.date = "";
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
  // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
  delete editableData[key];
  data.in_or_out = data.in_or_out ? 1 : 0;
  await SignalWaringConfigUpdate(data);
  await pageList();
};
const cancel = (key: string) => {
  delete editableData[key];
};
const confirm = async (id: string) => {
  SignalWaringConfigDelete(id).then(async ({ data }) => {
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
