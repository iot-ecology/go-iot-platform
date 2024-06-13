<template>
  <div class="box">
    <a-button style="margin-bottom: 10px" type="primary" @click="onAdd">新增</a-button>
    <a-button style="margin-bottom: 10px; margin-left: 10px" type="primary" @click="onSave">保存</a-button>
    <VueDraggable v-model="list1" :animation="150" ghost-class="ghost" group="people" handle=".drag-handle">
      <div v-for="(item, index) in list1" :key="item.id" class="cursor-move">
        <a-collapse :bordered="false" style="background: rgb(255, 255, 255)" :default-active-key="['1']" collapsible="icon">
          <template #expandIcon="{ isActive }">
            <caret-right-outlined :rotate="isActive ? 90 : 0" />
          </template>
          <a-collapse-panel key="1" :style="customStyle" :collapsible="icon">
            <template #header>
              <div class="drag-handle" style="display: flex; justify-content: space-between; border-bottom: 1px solid; padding: 5px">
                <div>
                  <a-input v-model:value="item.name" style="width: 300px"></a-input>
                  <a-button style="margin: 0 12px" type="primary" @click="onSet(index)">配置</a-button>
                  <a-select v-model:value="item.show.type" style="width: 100px" placeholder="请选择" @change="onChange(item.show.type, index)">
                    <a-select-option value="line">折线图</a-select-option>
                    <a-select-option value="bar">柱状图</a-select-option>
                  </a-select>
                  <a-button style="margin: 0 12px" type="primary" @click="onCopy(item)">复制</a-button>
                </div>
                <a-popconfirm title="确认是否删除?" ok-text="是" cancel-text="否" @confirm="onDelete(index)">
                  <close-circle-two-tone style="font-size: 24px; float: right" />
                </a-popconfirm>
              </div>
            </template>
            <div class="chart-container">
              <YcECharts :option="item.chart" :height="300" />
            </div>
          </a-collapse-panel>
        </a-collapse>
      </div>
    </VueDraggable>

    <a-modal v-model:open="modalVisible" title="新增" @ok="setModal1Visible()">
      <a-form :model="form" name="nest-messages">
        <a-form-item label="客户端ID">
          <a-select v-model:value="form.client_id" show-search placeholder="请输入" style="width: 300px" :default-active-first-option="true" :options="options"></a-select>
        </a-form-item>
        <a-form-item label="信号配置名称">
          <a-select v-model:value="form.fields" show-search mode="tags" placeholder="请输入" style="width: 300px" :default-active-first-option="false" :options="options1"></a-select>
        </a-form-item>
        <a-form-item label="时间">
          <a-tabs v-model:activeKey="activeKey">
            <a-tab-pane key="1" tab="动态时间">
              <div style="display: flex; align-items: center">
                最近
                <a-input-number v-model:value="dateTime" style="margin: 0 5px"></a-input-number>
                <a-select v-model:value="dateUnit" style="width: 100px">
                  <a-select-option value="year">年</a-select-option>
                  <a-select-option value="month">月</a-select-option>
                  <a-select-option value="day">日</a-select-option>
                  <a-select-option value="week">周</a-select-option>
                  <a-select-option value="hour">时</a-select-option>
                </a-select>
              </div>
            </a-tab-pane>
            <a-tab-pane key="2" tab="静态时间">
              <a-range-picker v-model:value="time" show-time @change="bptjTimeChange" />
            </a-tab-pane>
          </a-tabs>
        </a-form-item>
        <a-form-item label="时间间隔（秒）">
          <a-input-number v-model:value="form.every"></a-input-number>
        </a-form-item>
        <a-form-item label="脚本">
          <codemirror v-model="form.script" :disabled="form.start" placeholder="Code here..." :style="{ height: '150px' }" :autofocus="true" :tab-size="2" :extensions="extensions" />
        </a-form-item>
        <a-form-item label="统计方式">
          <a-select v-model:value="form.function" style="margin: -5px 0; width: 300px">
            <a-select-option value="mean">平均值</a-select-option>
            <a-select-option value="sum">求和</a-select-option>
            <a-select-option value="min">最小值</a-select-option>
            <a-select-option value="max">最大值</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="是否创建空值">
          <a-radio-group v-model:value="form.create_empty">
            <a-radio :value="true">是</a-radio>
            <a-radio :value="false">否</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal v-model:open="modal1Visible1" title="新增" @ok="setModal1Visible1()">
      <a-form name="nest-messages">
        <a-form-item label="名称">
          <a-input v-model:value="createName"></a-input>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>
<script setup lang="ts">
import { reactive, ref } from "vue";
import { Codemirror } from "vue-codemirror";
import { VueDraggable } from "vue-draggable-plus";
import { useRoute, useRouter } from "vue-router";
import type { Dayjs } from "dayjs";
import dayjs from "dayjs";
import { CaretRightOutlined, CloseCircleTwoTone } from "@ant-design/icons-vue";
import { javascript } from "@codemirror/lang-javascript";
import { EditorView } from "@codemirror/view";

import { DashboardCreate, DashboardId, DashboardUpdate, MqttPage, QueryInfluxdb, SignalPage } from "@/api";
import { YcECharts } from "@/components";

const myTheme = EditorView.theme(
  {
    // 输入的字体颜色
    "&": {
      color: "#0052D9",
      backgroundColor: "#FFFFFF",
    },
    ".cm-content": {
      caretColor: "#0052D9",
    },
    // 激活背景色
    ".cm-activeLine": {
      backgroundColor: "#FAFAFA",
    },
    // 激活序列的背景色
    ".cm-activeLineGutter": {
      backgroundColor: "#FAFAFA",
    },
    // 光标的颜色
    "&.cm-focused .cm-cursor": {
      borderLeftColor: "#0052D9",
    },
    // 选中的状态
    "&.cm-focused .cm-selectionBackground, ::selection": {
      backgroundColor: "#0052D9",
      color: "#FFFFFF",
    },
    // 左侧侧边栏的颜色
    ".cm-gutters": {
      backgroundColor: "#FFFFFF",
      color: "#ddd", // 侧边栏文字颜色
      border: "none",
    },
  },
  { dark: true },
);
const extensions = [javascript(), myTheme];

interface Item {
  id?: number;
  name: string;
  show: any;
  param: any;
  chart: any;
}
type RangeValue = [Dayjs, Dayjs];
const list1 = ref<Item[]>([]);
const modalVisible = ref(false);
const modal1Visible1 = ref(false);
const time = ref<RangeValue>();
const form = reactive({ client_id: "", fields: [], start_time: "", end_time: "", every: 1, function: "mean", create_empty: false });
const options = ref([]);
const options1 = ref([]);
const indexNumber = ref(0);
const activeKey = ref("1");
const dateTime = ref(null);
const dateUnit = ref("");
const createName = ref("");
const id = ref("");
const router = useRouter();
const route = useRoute();
const customStyle = "background: #f7f7f7;border-radius: 4px;margin-bottom: 24px;border: 0;overflow: hidden";
if (route.query.id) {
  DashboardId(route.query.id).then(({ data }) => {
    if (data.code === 20000) {
      list1.value = JSON.parse(data.data.config);
      createName.value = data.data.name;
      id.value = data.data.ID;
      list1.value.forEach((item, index) => {
        getModal(index);
      });
    }
  });
}
const onCopy = (item: any) => {
  const { name, show, param, chart } = item;
  list1.value.push({ name, show, param, chart });
};

function onAdd() {
  list1.value.push({
    name: "测试",
    id: list1.value.length,
    show: { type: "line" },
    param: { measurement: "", fields: [], sub: null, dateUnit: null, start_time: null, end_time: null, aggregation: { every: 1, function: "mean", create_empty: false } },
    chart: {},
  });
}

const onSet = (index: number) => {
  indexNumber.value = index;
  if (route.query.id) {
    form.client_id = list1.value[index].param.measurement || options.value[0]?.value;
    form.fields = list1.value[index].param.fields || [];
    form.every = list1.value[index].param.aggregation.every || 1;
    form.function = list1.value[index].param.aggregation.function || "mean";
    form.create_empty = list1.value[index].param.aggregation.create_empty || false;
    activeKey.value = list1.value[index].param.dateUnit ? "1" : "2";
    dateUnit.value = list1.value[index].param?.dateUnit || "";
    dateTime.value = list1.value[index].param?.sub || "";
  }
  modalVisible.value = true;
};

const List = async () => {
  const { data } = await MqttPage({ client_id: "", page: 1, page_size: 100 });
  options.value = data.data.data.map((item: any) => ({ value: String(item.ID), label: item.client_id }));
  form.client_id = options.value[0]?.value;
  List1();
};
List();

const List1 = async () => {
  const { data } = await SignalPage({ mqtt_client_id: form.client_id, page: 1, page_size: 100 });
  options1.value = data.data.data.map((item: any) => ({ value: String(item.name), label: item.name }));
};
const onSave = () => {
  modal1Visible1.value = true;
};

const onDelete = (index: number) => {
  list1.value.splice(index, 1);
};
const setModal1Visible = () => {
  let start_time = null;
  const end_time = dayjs().unix();
  if (activeKey.value === "1") {
    start_time = dayjs().subtract(dateTime.value, dateUnit.value).unix();
  }
  QueryInfluxdb({
    measurement: form.client_id,
    fields: [...form.fields, "storage_time", "push_time"],
    start_time: activeKey.value === "1" ? start_time : form.start_time,
    end_time: activeKey.value === "1" ? end_time : form.end_time,
    aggregation: {
      every: form.every,
      function: form.function,
      create_empty: form.create_empty,
    },
  }).then(({ data }) => {
    if (data.code === 20000) {
      modalVisible.value = false;
      const series: any = [];
      const list: any = data.data;
      const xAxis: any = list?.push_time.map((it) => dayjs(it._start).format("YYYY-MM-DD HH:mm:ss"));
      delete list.push_time;
      delete list.storage_time;

      Object.keys(list).forEach((item: string) => {
        series.push({
          name: item,
          type: list1.value[indexNumber.value].show.type,
          barWidth: "20",
          data: list[item].map((it) => it._value),
        });
      });
      list1.value[indexNumber.value].chart = {
        tooltip: {
          trigger: "axis",
        },
        legend: {
          data: Object.keys(list),
        },
        grid: {
          left: "3%",
          right: "4%",
          bottom: "8%",
          containLabel: true,
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
      list1.value[indexNumber.value].param = {
        measurement: form.client_id,
        fields: form.fields,
        start_time: activeKey.value === "1" ? start_time : form.start_time,
        end_time: activeKey.value === "1" ? end_time : form.end_time,
        sub: activeKey.value === "1" ? dateTime.value : "",
        dateUnit: activeKey.value === "1" ? dateUnit.value : "",
        aggregation: {
          every: form.every,
          function: form.function,
          create_empty: form.create_empty,
        },
      };
      if (activeKey.value === "1") {
        delete list1.value[indexNumber.value].param.start_time;
        delete list1.value[indexNumber.value].param.end_time;
      }
      if (activeKey.value === "2") {
        delete list1.value[indexNumber.value].param.sub;
        delete list1.value[indexNumber.value].param.dateUnit;
      }
    }
  });
};

const setModal1Visible1 = () => {
  const list = list1.value.map((it) => ({ name: it.name, id: it.id, show: it.show, param: it.param }));
  const data = {
    config: JSON.stringify(list),
    name: createName.value,
  };
  if (route.query.id) {
    data.id = id.value;
    DashboardUpdate(data).then(({ data }) => {
      if (data.code === 20000) {
        router.push({ path: "/draggable/index" });
      }
    });
  } else {
    DashboardCreate(data).then(({ data }) => {
      if (data.code === 20000) {
        router.push({ path: "/draggable/index" });
      }
    });
  }
};

// 查看编辑的时候获取数据
const getModal = (index: number) => {
  let start_time = null;
  const end_time = dayjs().unix();
  if (list1.value[index].param.sub && list1.value[index].param.dateUnit) {
    start_time = dayjs().subtract(list1.value[index].param.sub, list1.value[index].param.dateUnit).unix();
  }
  QueryInfluxdb({
    measurement: list1.value[index].param.measurement,
    fields: [...list1.value[index].param.fields, "storage_time", "push_time"],
    start_time: list1.value[index].param.sub ? start_time : list1.value[index].param.start_time,
    end_time: list1.value[index].param.sub ? end_time : list1.value[index].param.end_time,
    aggregation: {
      every: list1.value[index].param.aggregation.every,
      function: list1.value[index].param.aggregation.function,
      create_empty: list1.value[index].param.aggregation.create_empty,
    },
  }).then(({ data }) => {
    if (data.code === 20000) {
      modalVisible.value = false;
      const series: any = [];
      const list: any = data.data;
      const xAxis: any = list?.push_time?.map((it) => dayjs(it._start).format("YYYY-MM-DD HH:mm:ss"));
      delete list.push_time;
      delete list.storage_time;

      Object.keys(list).forEach((item: string) => {
        series.push({
          name: item,
          type: list1.value[index].show.type,
          barWidth: "20",
          data: list[item].map((it) => it._value),
        });
      });
      list1.value[index].chart = {
        tooltip: {
          trigger: "axis",
        },
        legend: {
          data: Object.keys(list),
        },
        grid: {
          left: "3%",
          right: "4%",
          bottom: "8%",
          containLabel: true,
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
      if (list1.value[index].param?.dateUnit) {
        delete list1.value[index].param.start_time;
        delete list1.value[index].param.end_time;
      }
      if (!list1.value[index].param?.dateUnit) {
        delete list1.value[index].param.sub;
        delete list1.value[index].param.dateUnit;
      }
    }
  });
};

const bptjTimeChange = (date: any) => {
  form.start_time = dayjs(date[0]).unix();
  form.end_time = dayjs(date[1]).unix();
};

const onChange = (text: string, index: number) => {
  const data = list1.value[index].chart;
  if (!Object.keys(data).length) return;
  data.series.forEach((item: any) => {
    item.type = text;
  });
  list1.value[index].chart = { ...data };
};
</script>
<style lang="less" scoped>
:deep(.ant-tabs-tab) {
  padding: 5px 0;
}
.drag-handle {
  cursor: grab;
}
.chart-container {
  position: relative;
  height: 300px;
  width: 100%;
  transition: 0.3s;
}
.box {
  padding: 16px;

  .cursor-move {
    padding: 10px;
    border: 1px solid rgba(0, 0, 0, 0.3);
    margin-bottom: 20px;
  }
}
</style>
