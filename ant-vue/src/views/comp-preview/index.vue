<template>
  <div class="comp-preview">
    <a-spin tip="加载中..." size="large" :spinning="showSpinning">
      <a-card :bordered="true">
        <a-form layout="inline">
          <a-form-item label="客户端ID">
            <MqttSelect v-model="value"></MqttSelect>
          </a-form-item>
          <a-form-item>
            <a-button type="primary" @click="pageList()">搜索</a-button>
          </a-form-item>
        </a-form>
        <a-button style="margin: 10px 0" type="primary" @click="modal1Visible = true">新增</a-button>

        <a-table :columns="columns" :data-source="list" bordered :pagination="paginations" @change="handleTableChange">
          <template #bodyCell="{ column, text, record }">
            <template v-if="['name', 'type', 'alias', 'cache_size', 'unit'].includes(column.dataIndex)">
              <div>
                <a-select v-if="editableData[record.key] && column.dataIndex == 'type'" v-model:value="editableData[record.key][column.dataIndex]" style="margin: -5px 0; width: 300px">
                  <a-select-option value="文本">文本</a-select-option>
                  <a-select-option value="数字">数字</a-select-option>
                </a-select>
                <a-input-number
                  v-else-if="editableData[record.key] && column.dataIndex === 'cache_size'"
                  v-model:value="editableData[record.key][column.dataIndex]"
                  :precision="0"
                  :min="1"
                  style="margin: -5px 0"
                />
                <a-input v-else-if="editableData[record.key] && column.dataIndex !== 'type'" v-model:value="editableData[record.key][column.dataIndex]" style="margin: -5px 0" />
                <template v-else>
                  {{ text }}
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
                  <a style="margin-right: 10px" @click="onView(record.mqtt_client_id, record.ID, record.alias, record.unit, record.type)">查看</a>
                  <a @click="edit(record.key)">编辑</a>
                  <a style="margin-left: 10px" @click="onSignal(record.ID, record.mqtt_client_id)">信号报警配置</a>
                  <a style="margin-left: 10px" @click="onHistoryView(record)">历史数据</a>
                  <a-popconfirm title="确认是否删除?" ok-text="是" cancel-text="否" @confirm="confirm(record.ID)">
                    <a style="margin-left: 10px; color: crimson">删除</a>
                  </a-popconfirm>
                </span>
              </div>
            </template>
          </template>
        </a-table>

        <a-modal v-model:open="modal1Visible" :destroy-on-close="true" title="新增" @ok="setModal1Visible()">
          <a-form ref="formRef" :label-col="{ style: { width: '80px' } }" :rules="rules" :model="form" name="nest-messages">
            <a-form-item label="客户端ID" name="mqtt_client_id">
              <MqttModSelect v-model="form.mqtt_client_id" :client-id="value" style="width: 350px"></MqttModSelect>
            </a-form-item>
            <a-form-item label="名称" name="name">
              <a-input v-model:value="form.name" style="width: 350px" />
            </a-form-item>
            <a-form-item label="类型" name="type">
              <a-select v-model:value="form.type" style="width: 350px">
                <a-select-option value="文本">文本</a-select-option>
                <a-select-option value="数字">数字</a-select-option>
              </a-select>
            </a-form-item>
            <a-form-item label="别名" name="alias">
              <a-input v-model:value="form.alias" style="width: 350px"></a-input>
            </a-form-item>
            <a-form-item label="缓存大小" name="cache_size">
              <a-input-number v-model:value="form.cache_size" :min="1" :precision="0" style="width: 350px"></a-input-number>
            </a-form-item>
            <a-form-item label="单位" name="unit">
              <a-input v-model:value="form.unit" style="width: 350px"></a-input>
            </a-form-item>
          </a-form>
        </a-modal>

        <a-modal v-model:open="modalView" style="width: 60%" title="近三十天数据" @ok="modalView = false">
          <a-form :model="form" name="nest-messages">
            <div v-if="!option.series?.length" style="text-align: center; font-size: 18px; height: 200px">暂无数据</div>
            <YcECharts v-else :option="option" :height="300" />
          </a-form>
        </a-modal>

        <!--历史数据-->
        <a-modal v-model:open="historyView" style="width: 60%" title="历史数据" @ok="historyView = false">
          <a-spin tip="加载中..." size="large" :spinning="showSpinning">
            <a-form :model="form" name="nest-messages">
              <a-range-picker value-format="YYYY-MM-DD HH:mm:ss" format="YYYY-MM-DD HH:mm:ss" style="width: 350px" show-time @change="bptjTimeChange" />
              <div v-if="!option.series?.length" style="text-align: center; font-size: 18px; height: 200px">暂无数据</div>
              <YcECharts v-else :option="option" :height="300" />
            </a-form>
          </a-spin>
        </a-modal>
      </a-card>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import type { UnwrapRef } from "vue";
import { onMounted, reactive, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { type FormInstance, message } from "ant-design-vue";
import { type Rule } from "ant-design-vue/es/form";
import dayjs from "dayjs";
import { cloneDeep } from "lodash-es";

import { QueryInfluxdb, QueryStrInfluxdb, SignalCreate, SignalDelete, SignalPage, SignalUpdate } from "@/api";
import { MqttModSelect, MqttSelect, YcECharts } from "@/components";
import { useRouteJump } from "@/hooks/useRouteJump.ts";
import { useRouterNameStore } from "@/stores/routerPath.ts";

interface DataItem {
  client_id: string;
  host: string;
  port: number;
  username: string;
}
const rules: Record<string, Rule[]> = {
  name: [{ required: true, message: "请输入名称", trigger: "blur" }],
  mqtt_client_id: [{ required: true, message: "请选择客户端ID", trigger: "change" }],
  type: [{ required: true, message: "请选择类型", trigger: "change" }],
  alias: [{ required: true, message: "请输入信号别名", trigger: "blur" }],
  cache_size: [{ required: true, message: "请输入缓存大小", trigger: "blur" }],
  unit: [{ required: true, message: "请输入单位", trigger: "blur" }],
};
const jump = useRouteJump();
const formRef = ref<FormInstance>();
const routerStore = useRouterNameStore();
const route = useRoute();
const value = ref("");
const modal1Visible = ref(false);
const modalView = ref(false);
const form = reactive({ mqtt_client_id: Number(route.query.id) || "", name: "", type: "", alias: "", cache_size: 1, unit: "" });
const columns = [
  {
    title: "ID",
    dataIndex: "ID",
  },
  {
    title: "名称",
    dataIndex: "name",
  },
  {
    title: "类型",
    dataIndex: "type",
  },
  {
    title: "别名",
    dataIndex: "alias",
  },
  {
    title: "缓存大小",
    dataIndex: "cache_size",
  },
  {
    title: "单位",
    dataIndex: "unit",
  },
  {
    title: "操作",
    dataIndex: "operation",
  },
];
const list = ref([]);
const editableData: UnwrapRef<Record<string, DataItem>> = reactive({});
const paginations = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
  showSizeChanger: true, // 显示每页显示条目数选择器
});
const option = ref({});
const historyView = ref(false);
const recordObj = ref({});
const showSpinning = ref(false);
watch(value, async (newValue, oldValue) => {
  await pageList();
});

const setModal1Visible = () => {
  formRef.value
    .validate()
    .then(() => {
      SignalCreate({ ...form }).then(({ data }) => {
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
    .catch((error) => {});
};

const save = async (key: string) => {
  Object.assign(list.value.filter((item) => key === item.key)[0], editableData[key]);
  const data = list.value.filter((item) => key === item.key)[0];
  // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
  delete editableData[key];
  // eslint-disable-next-line no-debugger
  await SignalUpdate(data);
  await pageList();
};
const cancel = (key: string) => {
  // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
  delete editableData[key];
};
const edit = (key: string) => {
  editableData[key] = cloneDeep(list.value.filter((item) => key === item.key)[0]);
};
const confirm = async (id: string) => {
  SignalDelete(id).then(({ data }) => {
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
  const { data } = await SignalPage({ mqtt_client_id: value.value, page: paginations.current, page_size: paginations.pageSize });
  paginations.total = data.data?.total || 0;
  list.value = data.data.data?.map((item: any, index: number) => ({
    key: index,
    ID: item.ID,
    mqtt_client_id: item.mqtt_client_id,
    mqtt_client_name: item.mqtt_client_name,
    name: item.name,
    type: item.type,
    alias: item.alias,
    cache_size: item.cache_size,
    unit: item.unit,
  }));
};
onMounted(async () => {});

const onSignal = (id: string, mqtt_client_id: string) => {
  routerStore.setRouterName("/signal");
  jump.routeJump({ path: "/signal", query: { id, mqtt_client_id } });
};

const handleTableChange = async (pagination: any) => {
  paginations.current = pagination.current;
  paginations.pageSize = pagination.pageSize;
  await pageList();
};

const onView = (id: number, rowId: number, alias: string, unit: string, type: string) => {
  const time = dayjs();
  const data = type === "数字" ? QueryInfluxdb : QueryStrInfluxdb;
  showSpinning.value = true;
  data({
    measurement: String(id),
    fields: [String(rowId), "storage_time", "push_time"],
    start_time: time.subtract(30, "day").unix(),
    end_time: time.unix(),
    aggregation: {
      every: 1,
      function: type === "数字" ? "mean" : "first",
      create_empty: false,
    },
  })
    .then(({ data }) => {
      if (data.code === 20000) {
        modalView.value = true;
        const series: any = [];
        const list: any = data.data;
        let xAxis: any = list?.push_time?.map((it) => dayjs(it._time).format("YYYY-MM-DD HH:mm:ss"));
        delete list.push_time;
        delete list.storage_time;
        Object.keys(list).forEach((item: string) => {
          series.push({
            name: alias,
            type: "line",
            barWidth: "20",
            data: list[item].map((it) => it._value),
          });
        });
        xAxis = xAxis?.length ? xAxis.slice(0, series[0].data?.length) : [];
        option.value = {
          tooltip: {
            trigger: "axis",
            formatter: (params) => {
              return `${params[0].name}: ${params[0].value} ${unit}`;
            },
          },
          legend: {
            data: [alias],
          },
          grid: {
            left: "3%",
            right: "4%",
            bottom: "8%",
            containLabel: true,
          },
          toolbox: {
            show: true,
            feature: {
              mark: { show: false },
              dataView: { show: false, readOnly: false },
              magicType: { show: true, type: ["line", "bar"] },
              restore: { show: false },
              saveAsImage: { show: false },
            },
          },
          xAxis: {
            show: false,
            type: "category",
            boundaryGap: true,
            data: [...xAxis],
          },
          yAxis: {
            type: "value",
          },
          series: [...series],
          dataZoom: [
            {
              type: "inside",
            },
            {
              type: "slider",
              top: 260,
            },
          ],
        };
      } else {
        message.error(data.message || "操作失败");
      }
    })
    .finally(() => {
      showSpinning.value = false;
    });
};

const onHistoryView = (record: any) => {
  recordObj.value = record;
  historyView.value = true;
};

const bptjTimeChange = (date: any) => {
  if (!date) {
    option.value = {};
    return;
  }
  const start_time = dayjs(date[0]).unix();
  const end_time = dayjs(date[1]).unix();
  showSpinning.value = true;
  const data = recordObj.value.type === "数字" ? QueryInfluxdb : QueryStrInfluxdb;
  data({
    measurement: String(recordObj.value.mqtt_client_id),
    fields: [String(recordObj.value.ID), "storage_time", "push_time"],
    start_time,
    end_time,
    aggregation: {
      every: 1,
      function: recordObj.value.type === "数字" ? "mean" : "first",
      create_empty: false,
    },
  })
    .then(({ data }) => {
      if (data.code === 20000) {
        historyView.value = true;
        const series: any = [];
        const list: any = data.data;
        let xAxis: any = list?.push_time?.map((it) => dayjs(it._time).format("YYYY-MM-DD HH:mm:ss")) || [];
        delete list.push_time;
        delete list.storage_time;

        Object.keys(list).forEach((item: string) => {
          series.push({
            name: recordObj.value.alias,
            type: "line",
            barWidth: "20",
            data: list[item].map((it) => it._value),
          });
        });
        xAxis = xAxis?.length ? xAxis.slice(0, series[0].data?.length) : [];
        option.value = {
          tooltip: {
            trigger: "axis",
            formatter: (params) => {
              return `${params[0].name}: ${params[0].value} ${recordObj.value.unit}`;
            },
          },
          legend: {
            data: [recordObj.value.alias],
          },
          grid: {
            left: "3%",
            right: "4%",
            bottom: "8%",
            containLabel: true,
          },
          toolbox: {
            show: true,
            feature: {
              mark: { show: false },
              dataView: { show: false, readOnly: false },
              magicType: { show: true, type: ["line", "bar"] },
              restore: { show: false },
              saveAsImage: { show: false },
            },
          },
          xAxis: {
            show: false,
            type: "category",
            boundaryGap: true,
            data: [...xAxis],
          },
          yAxis: {
            type: "value",
          },
          series: [...series],
          dataZoom: [
            {
              type: "inside",
            },
            {
              type: "slider",
              top: 260,
            },
          ],
        };
      } else {
        message.error(data.message || "操作失败");
      }
    })
    .finally(() => {
      showSpinning.value = false;
    });
};
</script>

<style lang="less" scoped>
.comp-preview {
  width: 100%;
  height: 100%;
}
</style>
