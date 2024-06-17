<template>
  <div class="box">
    <a-button style="margin-bottom: 10px" type="primary" @click="onAdd">{{ $t('message.addition') }}</a-button>
    <a-button style="margin-bottom: 10px; margin-left: 10px" type="primary" @click="onSave">{{$t('message.save')}}</a-button>
    <VueDraggable v-model="listArr" :animation="150" ghost-class="ghost" group="people" handle=".drag-handle">
      <div v-for="(item, index) in listArr" :key="item.id" class="cursor-move">
        <a-spin :tip="$t('message.loading')" size="large" :spinning="item.showSpinning">
          <a-collapse :bordered="false" style="background: rgb(255, 255, 255)" :default-active-key="['1']" collapsible="icon">
            <template #expandIcon="{ isActive }">
              <caret-right-outlined :rotate="isActive ? 90 : 0" />
            </template>
            <a-collapse-panel key="1" :style="customStyle" :collapsible="icon">
              <template #header>
                <div class="drag-handle" style="display: flex; justify-content: space-between; border-bottom: 1px solid; padding: 5px">
                  <div>
                    <a-input v-model:value="item.name" style="width: 300px"></a-input>
                    <a-button style="margin: 0 12px" type="primary" @click="onSet(index)">{{ $t('message.allocation') }}</a-button>
                    <a-select v-model:value="item.show.type" style="width: 100px" :placeholder="$t('message.pleaseSelect')" @change="onChange(item.show.type, index)">
                      <a-select-option value="line">{{ $t('message.lineChart') }}</a-select-option>
                      <a-select-option value="bar">{{ $t('message.barChart') }}</a-select-option>
                    </a-select>
                    <a-button style="margin: 0 12px" type="primary" @click="onCopy(item)">{{ $t('message.copy') }}</a-button>
                  </div>
                  <a-popconfirm :title="$t('message.sureDelete')" :okText="$t('message.yes')" :cancelText="$t('message.no')" @confirm="onDelete(index)">
                    <close-circle-two-tone two-tone-color="crimson" style="font-size: 24px" />
                  </a-popconfirm>
                </div>
              </template>
              <div class="chart-container">
                <div v-if="!item.chart?.series?.length" style="text-align: center; font-size: 18px; height: 200px">{{ $t('message.noData') }}</div>
                <YcECharts v-else :option="item.chart" :height="300" />
              </div>
            </a-collapse-panel>
          </a-collapse>
        </a-spin>
      </div>
    </VueDraggable>

    <a-modal v-model:open="modalVisible" style="width: 35%" :title="$t('message.addition')">
      <a-form ref="formRef" :label-col="{ style: { width: '130px' } }" :model="form" :rules="rules" name="nest-messages">
        <a-form-item :label="$t('message.time')">
          <a-tabs v-model:activeKey="activeKey">
            <a-tab-pane key="1" :tab="$t('message.dynamicTime')">
              <div style="display: flex; align-items: center">
                {{ $t('message.recently') }}
                <a-input-number v-model:value="dateTime" style="margin: 0 5px; width: 150px"></a-input-number>
                <a-select v-model:value="dateUnit" style="width: 165px">
                  <a-select-option value="year">{{ $t('message.year') }}</a-select-option>
                  <a-select-option value="month">{{ $t('message.month') }}</a-select-option>
                  <a-select-option value="day">{{ $t('message.day') }}</a-select-option>
                  <a-select-option value="week">{{ $t('message.week') }}</a-select-option>
                  <a-select-option value="hour">{{ $t('message.hour') }}</a-select-option>
                </a-select>
              </div>
            </a-tab-pane>
            <a-tab-pane key="2" :tab="$t('message.staticTime')">
              <a-range-picker v-model:value="time" value-format="YYYY-MM-DD HH:mm:ss" format="YYYY-MM-DD HH:mm:ss" style="width: 350px" show-time @change="bptjTimeChange" />
            </a-tab-pane>
          </a-tabs>
        </a-form-item>
        <a-form-item :label="$t('message.timeInterval')" name="every">
          <a-input-number v-model:value="form.every" :precision="0" style="width: 350px"></a-input-number>
        </a-form-item>

        <a-form-item :label="$t('message.createValue')" name="create_empty">
          <a-radio-group v-model:value="form.create_empty">
            <a-radio :value="true">{{ $t('message.yes') }}</a-radio>
            <a-radio :value="false">{{ $t('message.no') }}</a-radio>
          </a-radio-group>
        </a-form-item>

        <a-divider>{{ $t('message.signalConfig') }}</a-divider>
        <a-button type="primary" @click="onAddSignal">{{ $t('message.addition') }}</a-button>

        <div style="border: 1px solid #d9d9d9; margin-top: 12px">
          <div style="display: flex; line-height: 32px; border-bottom: 1px solid #d9d9d9">
            <div style="width: 33.3%; text-align: center; border-right: 1px solid #d9d9d9">{{ $t('message.clientID') }}</div>
            <div style="width: 33.3%; text-align: center; border-right: 1px solid #d9d9d9">{{ $t('message.signalName') }}</div>
            <div style="width: 33.3%; text-align: center">{{ $t('message.statisticalMethods') }}</div>
          </div>
          <div v-for="(item, index) in form.list" :key="index" style="display: flex; justify-content: space-between">
            <div style="width: 33.3%; text-align: center; padding: 4px 0; border-right: 1px solid #d9d9d9">
              <MqttSelect v-model="item.client_id" style="width: 160px"></MqttSelect>
            </div>
            <div style="width: 33.3%; text-align: center; padding: 4px 0; border-right: 1px solid #d9d9d9">
              <SignalModeSelect v-model="item.fields" style="width: 160px" :mqtt_client_id="item.client_id" name="ID" @custom-event="handleCustomEvent"></SignalModeSelect>
            </div>
            <div style="width: 33.3%; text-align: center; padding: 4px 0">
              <a-select v-model:value="item.function" style="width: 160px">
                <a-select-option value="mean">{{ $t('message.mean') }}</a-select-option>
                <a-select-option value="sum">{{ $t('message.sum') }}</a-select-option>
                <a-select-option value="min">{{ $t('message.min') }}</a-select-option>
                <a-select-option value="max">{{ $t('message.max') }}</a-select-option>
                <a-select-option value="first">{{ $t('message.first') }}</a-select-option>
                <a-select-option value="last">{{ $t('message.last') }}</a-select-option>
              </a-select>
            </div>
          </div>
        </div>
      </a-form>
      <template #footer>
        <a-button @click="modalVisible = false">{{$t('message.cancel')}}</a-button>
        <a-button type="primary" @click="onConfirm()">{{$t('message.confirm')}}</a-button>
      </template>
    </a-modal>

    <a-modal v-model:open="modalSave" :title="$t('message.save')" @ok="onSaveInformation()">
      <a-form name="nest-messages">
        <a-form-item :label="$t('message.name')">
          <a-input v-model:value="createName"></a-input>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>
<script setup lang="ts">
import { reactive, ref ,watch} from "vue";
import { VueDraggable } from "vue-draggable-plus";
import { useRoute, useRouter } from "vue-router";
import { type FormInstance, message } from "ant-design-vue";
import { type Rule } from "ant-design-vue/es/form";
import type { Dayjs } from "dayjs";
import dayjs from "dayjs";
import { cloneDeep } from "lodash-es";
import { CaretRightOutlined, CloseCircleTwoTone } from "@ant-design/icons-vue";

import { DashboardCreate, DashboardId, DashboardUpdate, QueryInfluxdb, SignalPage } from "@/api";
import { YcECharts } from "@/components";
import { MqttSelect, SignalModeSelect } from "@/components/index.ts";
import {useI18n} from "vue-i18n";
const { t,locale } = useI18n();
interface Item {
  id?: number;
  name: string;
  show: any;
  param: any;
  chart: any;
  showSpinning: boolean;
}

type RangeValue = [Dayjs, Dayjs];
const formRef = ref<FormInstance>();
const listArr = ref<Item[]>([
  {
    name: "测试",
    id: 1,
    show: { type: "line" },
    param: {
      measurement: "",
      fields: [],
      sub: null,
      dateUnit: null,
      start_time: null,
      end_time: null,
      aggregation: { every: 1, function: "mean", create_empty: false },
      list: [{ client_id: "", fields: "", function: "mean" }],
    },
    showSpinning: false,
    chart: {},
  },
]);
const modalVisible = ref(false);
const modalSave = ref(false);
const time = ref<RangeValue>();
const form = reactive({ client_id: "", fields: [], start_time: "", end_time: "", every: 1, function: "mean", create_empty: false, list: [{ client_id: "", fields: "", function: "mean" }] });
const indexNumber = ref(0);
const activeKey = ref("1");
const dateTime = ref(null);
const dateUnit = ref("day");
const createName = ref("");
const id = ref("");
const router = useRouter();
const route = useRoute();
const optionList = ref([]);
const customStyle = "background: #f7f7f7;border-radius: 4px;margin-bottom: 24px;border: 0;overflow: hidden";
let rules: Record<string, Rule[]> = {
  client_id: [{ required: true, message: t('message.pleaseSelectClientID'), trigger: "change" }],
  fields: [{ required: true, message: t('message.pleaseSignalName'), trigger: "change" }],
  every: [{ required: true, message: t('message.pleaseTimeInterval'), trigger: "blur" }],
  function: [{ required: true, message: t('message.pleaseStatisticalMethod'), trigger: "change" }],
  create_empty: [{ required: true, message: t('message.pleaseChooseValue'), trigger: "change" }],
};
const pageSize = 2000
const reduces = {
  mean: t('message.mean'),
  sum: t('message.sum'),
  min: t('message.max'),
  max: t('message.min'),
  原始: t('message.original'),
  first: t('message.first'),
  last: t('message.last'),
};
watch(locale, () => {
  rules = {
    client_id: [{ required: true, message: t('message.pleaseSelectClientID'), trigger: "change" }],
    fields: [{ required: true, message: t('message.pleaseSignalName'), trigger: "change" }],
    every: [{ required: true, message: t('message.pleaseTimeInterval'), trigger: "blur" }],
    function: [{ required: true, message: t('message.pleaseStatisticalMethod'), trigger: "change" }],
    create_empty: [{ required: true, message: t('message.pleaseChooseValue'), trigger: "change" }],
  }
});

if (route.query.id) {
  DashboardId(route.query.id).then(({ data }) => {
    if (data.code === 20000) {
      createName.value = data.data.name;
      id.value = data.data.ID;
      if (JSON.parse(data.data.config) && JSON.parse(data.data.config)?.length) {
        listArr.value = JSON.parse(data.data.config);
        listArr.value.forEach((item, index) => {
          getData(index);
        });
      }
    }
  }).catch(e=>{
    console.error(e)
  });
}
const onCopy = (item: any) => {
  const obj = cloneDeep(item);
  const { name, show, param, chart } = obj;
  listArr.value.push({ name: name + `（${t('message.copy')}）`, show, param, chart, showSpinning: false });
};

function onAdd() {
  listArr.value.push({
    name: "测试",
    id: listArr.value.length + 1,
    show: { type: "line" },
    param: {
      measurement: "",
      fields: [],
      sub: null,
      dateUnit: null,
      start_time: null,
      end_time: null,
      aggregation: { every: 1, function: "mean", create_empty: false },
      list: [{ client_id: "", fields: "", function: "mean" }],
    },
    chart: {},
    showSpinning: false,
  });
}

const onSet = (index: number) => {
  indexNumber.value = index;
  if (route.query.id) {
    form.client_id = Number(listArr.value[index].param.measurement) || "";
    form.fields = listArr.value[index].param.fields || [];
    form.every = listArr.value[index].param.aggregation.every || 1;
    form.function = listArr.value[index].param.aggregation.function || "mean";
    form.create_empty = listArr.value[index].param.aggregation.create_empty || false;
    form.list = listArr.value[index].param.list;
    activeKey.value = listArr.value[index].param.dateUnit && listArr.value[index].param.sub ? "1" : "2";
    if (listArr.value[index].param?.start_time && listArr.value[index].param?.end_time) {
      activeKey.value = "2";
    } else {
      activeKey.value = "1";
    }
    dateUnit.value = listArr.value[index].param?.dateUnit || "day";
    dateTime.value = listArr.value[index].param?.sub || "";
    form.start_time = listArr.value[index].param?.start_time || "";
    form.end_time = listArr.value[index].param?.end_time || "";
    if (activeKey.value === "1") {
      time.value = [];
    } else {
      if (listArr.value[index].param.start_time && listArr.value[index].param.end_time) {
        time.value = [dayjs(listArr.value[index].param.start_time * 1000).format("YYYY-MM-DD HH:mm:ss"), dayjs(listArr.value[index].param.end_time * 1000).format("YYYY-MM-DD HH:mm:ss")];
      } else {
        time.value = [];
      }
    }
  }
  modalVisible.value = true;
};
const onSave = () => {
  modalSave.value = true;
};

const onDelete = (index: number) => {
  listArr.value.splice(index, 1);
};
const onConfirm = () => {
  let start_time = null;
  const end_time = dayjs().unix();
  if (activeKey.value === "1") {
    start_time = dayjs().subtract(dateTime.value, dateUnit.value).unix();
  }
  if (activeKey.value === "1" && !dateTime.value) {
    message.error(t('message.pleaseTime'));
    return;
  }

  if (activeKey.value === "2" && !(form.start_time && form.end_time)) {
    message.error(t('message.pleaseTime'));
    return;
  }
  const series: any = [];
  const legend: any = [];
  listArr.value[indexNumber.value].showSpinning = true;
  modalVisible.value = false;
  form.list.forEach((item, index) => {
    QueryInfluxdb({
      measurement: String(item.client_id),
      fields: [String(item.fields), "storage_time", "push_time"],
      start_time: activeKey.value === "1" ? start_time : form.start_time,
      end_time: activeKey.value === "1" ? end_time : form.end_time,
      aggregation: {
        every: form.every,
        function: item.function,
        create_empty: form.create_empty,
      },
    })
      .then(async ({ data }) => {
        if (data.code === 20000) {
          const dataArr = await SignalPage({ mqtt_client_id: item.client_id, page: 1, page_size: pageSize });
          const signalList = dataArr.data.data?.data || [];
          listArr.value[indexNumber.value].param = {
            measurement: form.client_id,
            fields: form.fields,
            aggregation: {
              every: form.every,
              function: form.function,
              create_empty: form.create_empty,
            },
            list: form.list,
          };
          if (activeKey.value === "1") {
            listArr.value[indexNumber.value].param.sub = activeKey.value === "1" ? dateTime.value : "";
            listArr.value[indexNumber.value].param.sub = activeKey.value === "1" ? dateUnit.value : "";
          }
          if (activeKey.value === "2") {
            listArr.value[indexNumber.value].param.start_time = activeKey.value === "1" ? start_time : form.start_time;
            listArr.value[indexNumber.value].param.end_time = activeKey.value === "1" ? end_time : form.end_time;
          }
          const list: any = data.data;
          const xAxis: any = list?.push_time?.map((it: any) => dayjs(it._time).format("YYYY-MM-DD HH:mm:ss")) || [];
          delete list.push_time;
          delete list.storage_time;

          listArr.value[indexNumber.value].chart = {
            tooltip: {
              trigger: "axis",
            },
            legend: {
              data: [],
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
              data: [],
            },
            yAxis: {
              type: "value",
            },
            series: [],
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
          const items = signalList.find((it: any) => it.ID === Number(item.fields));
          const mqttName = items?.mqtt_client_name;
          const signalName = items?.name;
          const alias = items?.alias;
          const tipName = `${mqttName}-${signalName}(${alias})-${reduces[item.function]}`
          legend.push(tipName);
          Object.keys(list).forEach((text: string) => {
            series.push({
              name: tipName,
              type: listArr.value[indexNumber.value].show.type,
              barWidth: "20",
              data: list[text].map((it) => it._value),
            });
          });
          if (Object.keys(list).length) {
            listArr.value[indexNumber.value].chart.tooltip = {
              trigger: "axis",
              formatter: (params: any) => {
                let res = `${params[0].name}` + "<br/>";
                params.forEach(function (item: any) {
                  res += `${item.seriesName}: ${item.value} ${optionList.value[0]?.unit}` + "<br/>";
                });
                return res;
              },
            };
            listArr.value[indexNumber.value].chart.legend.data = legend;
            listArr.value[indexNumber.value].chart.xAxis.data = [...xAxis];
            listArr.value[indexNumber.value].chart.series = [...series];
          }
        } else {
          message.error(data.message);
        }
      }).catch(e=>{
      console.error(e)
    })
      .finally(() => {
        if (index === form.list.length - 1) {
          listArr.value[indexNumber.value].showSpinning = false;
        }
      });
  });
};

const onSaveInformation = () => {
  const list = listArr.value.map((it) => ({ name: it.name, id: it.id, show: it.show, param: it.param, showSpinning: it.showSpinning }));
  const data = {
    config: JSON.stringify(list),
    name: createName.value,
  };
  if (route.query.id) {
    data.id = id.value;
    DashboardUpdate(data).then(({ data }) => {
      if (data.code === 20000) {
        router.push({ path: "/visualization/index" });
      }
    }).catch(e=>{
      console.error(e)
    });
  } else {
    DashboardCreate(data).then(({ data }) => {
      if (data.code === 20000) {
        router.push({ path: "/visualization/index" });
      }
    }).catch(e=>{
      console.error(e)
    });
  }
};

// 查看编辑的时候获取数据
const getData = async (index: number) => {
  const { data } = await SignalPage({ mqtt_client_id: listArr.value[index].param.measurement, page: 1, page_size: pageSize });
  const signalList = data.data?.data || [];
  let start_time = null;
  const end_time = dayjs().unix();
  if (listArr.value[index].param.sub && listArr.value[index].param.dateUnit) {
    start_time = dayjs().subtract(listArr.value[index].param.sub, listArr.value[index].param.dateUnit).unix();
  }
  const series: any = [];
  const aliasList: any = [];
  listArr.value[index].showSpinning = true;
  listArr.value[index].param.list.forEach((text: any, number: number) => {
    QueryInfluxdb({
      measurement: String(text.client_id),
      fields: [String(text.fields), "storage_time", "push_time"],
      start_time: listArr.value[index].param.sub ? start_time : listArr.value[index].param.start_time,
      end_time: listArr.value[index].param.sub ? end_time : listArr.value[index].param.end_time,
      aggregation: {
        every: listArr.value[index].param.aggregation.every,
        function: text.function,
        create_empty: listArr.value[index].param.aggregation.create_empty,
      },
    })
      .then(({ data }) => {
        if (data.code === 20000) {
          modalVisible.value = false;
          const list: any = data.data;
          const xAxis: any = list?.push_time?.map((it: any) => dayjs(it._time).format("YYYY-MM-DD HH:mm:ss")) || [];
          delete list.push_time;
          delete list.storage_time;

          Object.keys(list).forEach((item: string) => {
            const items = signalList.find((it: any) => it.ID === Number(item));
            const mqttName = items?.mqtt_client_name;
            const signalName = items?.name;
            const alias = items?.alias;
            const tipName = `${mqttName}-${signalName}(${alias})-${reduces[text.function]}`
            series.push({
              name: tipName,
              type: listArr.value[index].show.type,
              barWidth: "20",
              data: list[item].map((it: any) => it._value),
            });
          });
          listArr.value[index].chart = {
            tooltip: {
              trigger: "axis",
            },
            legend: {
              data: [],
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
              data: [],
            },
            yAxis: {
              type: "value",
            },
            series: [],
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
          if (listArr.value[index].param?.dateUnit) {
            delete listArr.value[index].param.start_time;
            delete listArr.value[index].param.end_time;
          }
          if (!listArr.value[index].param?.dateUnit) {
            delete listArr.value[index].param.sub;
            delete listArr.value[index].param.dateUnit;
          }
          const items = signalList.find((it: any) => it.ID === Number(text.fields));
          const mqttName = items?.mqtt_client_name;
          const signalName = items?.name;
          const alias = items?.alias;
          const unit = items?.unit;
          const tipName = `${mqttName}-${signalName}(${alias})-${reduces[text.function]}`
          aliasList.push(tipName);
          if (Object.keys(list).length) {
            listArr.value[index].chart.tooltip = {
              trigger: "axis",
              formatter: (params: any) => {
                let res = `${params[0].name}` + "<br/>";
                params.forEach(function (item: any) {
                  res += `${item.seriesName}: ${item.value} ${unit}` + "<br/>";
                });
                return res;
              },
            };
            listArr.value[index].chart.legend.data = aliasList;
            listArr.value[index].chart.xAxis.data = [...xAxis];
            listArr.value[index].chart.series = [...series];
          }
        }
      }).catch(e=>{
        console.error(e)
      })
      .finally(() => {
        if (number === listArr.value[index].param.list.length - 1) {
          listArr.value[index].showSpinning = false;
        }
      });
  });
};

const bptjTimeChange = (date: any) => {
  if (!date) {
    form.start_time = "";
    form.end_time = "";
    return;
  }
  form.start_time = dayjs(date[0]).unix();
  form.end_time = dayjs(date[1]).unix();
};

const onChange = (text: string, index: number) => {
  const data = listArr.value[index].chart;
  if (!Object.keys(data).length) return;
  data.series.forEach((item: any) => {
    item.type = text;
  });
  listArr.value[index].chart = { ...data };
};

const handleCustomEvent = (payload: any) => {
  if (payload.value !== -11) {
    optionList.value.push(...payload);
  }
  // 处理从子组件传递来的参数
};

const onAddSignal = () => {
  form.list.push({ client_id: "", fields: "", function: "mean" });
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
