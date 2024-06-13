<template>
  <div>
    <section ref="myEcharts" class="box">
      <div v-for="(item, index) in list" :key="index" @click="onView(item)">
        <YcECharts :option="onSet(item)" :height="300" />
      </div>
    </section>

    <!--      详情-->
    <a-modal v-model:open="showResult" style="width: 700px" :footer="null" :destroy-on-close="true" title="">
      <a-divider>节点信息</a-divider>
      <a-descriptions title="" bordered>
        <a-descriptions-item label="名称">{{ form.name }}</a-descriptions-item>
        <a-descriptions-item label="当前容量">{{ form.size }}</a-descriptions-item>
        <a-descriptions-item label="最大容量">{{ form.max_size }}</a-descriptions-item>
      </a-descriptions>
      <a-divider>客户端信息</a-divider>
      <a-table bordered :pagination="false" :data-source="dataResult" :columns="columnsResult">
        <template #bodyCell="{ column, text, record }">
          <template v-if="column.dataIndex === 'password' || column.dataIndex === 'broker'">
            <a-input-password :value="text" :bordered="false" />
          </template>
        </template>
      </a-table>
    </a-modal>
  </div>
</template>
<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";

import { MqttNodeUsingStatus } from "@/api";
import { YcECharts } from "@/components";

const myEcharts = ref(null);
const list = ref("");
const showResult = ref(false);
const dataResult = ref([]);
const columnsResult = ref([
  {
    title: "客户端ID",
    dataIndex: "client_id",
  },
  {
    title: "主机",
    dataIndex: "broker",
  },
  {
    title: "端口",
    dataIndex: "port",
  },
  {
    title: "账号",
    dataIndex: "username",
  },
  {
    title: "密码",
    dataIndex: "password",
  },
  {
    title: "订阅的主题",
    dataIndex: "sub_topic",
  },
]);
const form = reactive({ name: "", size: "", max_size: "" });
const getNode = async () => {
  const { data } = await MqttNodeUsingStatus();
  list.value = JSON.parse(data.data).data;
};
getNode();
const onSet = (item: any) => {
  const option = {
    series: [
      {
        type: "gauge",
        startAngle: 180,
        endAngle: 0,
        center: ["50%", "75%"],
        radius: "90%",
        min: 0,
        max: item.max_size,
        splitNumber: 8,
        axisLine: {
          lineStyle: {
            width: 6,
            color: [
              [0.8, "#7CFFB2"],
              [1, "#FF6E76"],
            ],
          },
        },
        pointer: {
          animation: false,
          icon: "path://M12.8,0.7l12,40.1H0.7L12.8,0.7z",
          length: "12%",
          width: 20,
          offsetCenter: [0, "-60%"],
          itemStyle: {
            color: "auto",
          },
        },
        axisTick: {
          length: 12,
          lineStyle: {
            color: "auto",
            width: 2,
          },
        },
        splitLine: {
          length: 20,
          lineStyle: {
            color: "auto",
            width: 5,
          },
        },
        axisLabel: {
          color: "#464646",
          show: false,
          fontSize: 12,
          distance: -60,
          rotate: "tangential",
        },
        title: {
          offsetCenter: [0, "-10%"],
          fontSize: 20,
        },
        detail: {
          fontSize: 30,
          offsetCenter: [0, "-35%"],
          valueAnimation: false,
          color: "inherit",
        },
        data: [
          {
            value: item.size,
            name: item.name,
            stateAnimation: false,
          },
        ],
      },
    ],
  };
  return option;
};
const onView = (record: any) => {
  showResult.value = true;
  dataResult.value = record?.client_infos || [];
  form.name = record.name;
  form.size = record.size;
  form.max_size = record.max_size;
};
onMounted(() => {});
</script>
<style lang="less" scoped>
.box {
  display: flex;
  flex-flow: wrap;
  > div {
    width: 33.3%;
  }
}
</style>
