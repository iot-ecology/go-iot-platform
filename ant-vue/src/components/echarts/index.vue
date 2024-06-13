<script lang="ts" setup>
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from "vue";
import * as echarts from "echarts";
import { type EChartOption } from "echarts";

import echartsConfig from "./walden.project.json";

const props = defineProps({
  option: {
    type: Object as () => EChartOption | any,
    required: true,
  },
  width: {
    type: Number,
  },
  height: {
    type: Number,
  },
});

let chart: echarts.ECharts;

const chartRef = ref<HTMLDivElement | undefined>();
console.log(props.option.series[0].type);
if (props.option.series[0].type === "gauge") {
  nextTick(() => {
    chart.setOption(props.option || {}, true);
  });
}
watch(
  () => props.option,
  async (value) => {
    if (value.series[0].type === "gauge") {
      return;
    }
    await nextTick(() => {
      // 重置为数值类、（false）的时候开启加载动画
      // if (!Object.keys(value || {}).length) {
      //   chart.showLoading("default", {
      //     text: "加载中...",
      //     color: "rgb(22, 93, 255)",
      //     textColor: "#000",
      //     maskColor: "rgba(255, 255, 255, 0.7)",
      //   });
      // } else {
      //   chart.hideLoading();
      // }
      chart.clear();
      chart.setOption(value || {}, true);
    });
  },
  { immediate: true },
);

function resize() {
  chart?.resize();
}

window.addEventListener("resize", resize);

onMounted(() => {
  if (chartRef.value) {
    chart = echarts.init(chartRef.value, echartsConfig.theme, {
      width: props.width,
      height: props.height,
      locale: "ZH",
    });
  }
});

onBeforeUnmount(() => {
  window.removeEventListener("resize", resize);
});
</script>

<template>
  <div ref="chartRef" class="t-chart" />
</template>

<style scoped>
.t-chart {
  width: 100%;
  height: 100%;
}
</style>
