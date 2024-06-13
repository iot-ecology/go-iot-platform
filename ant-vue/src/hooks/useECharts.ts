import { onMounted, onUnmounted, type Ref, unref } from "vue";
import * as echarts from "echarts";
import { debounce } from "lodash-es";

export const useECharts = (el: Ref<HTMLDivElement> | HTMLDivElement, autoUpdateSize = true) => {
  // 实例对象
  let echartsInstance: echarts.ECharts | null = null;

  // 检查DOM对象是否为空
  const isEmptyEl = () => {
    if (!unref(el)) {
      console.warn("[useECharts]传入的DOM对象为空，请检查DOM对象");
      return true;
    }
    return false;
  };

  // 初始化实例
  const initCharts = () => {
    if (isEmptyEl()) return;
    echartsInstance = echarts.init(unref(el));
  };

  // 更新配置
  const setOptions = (option: any) => {
    if (!echartsInstance) initCharts();
    if (!echartsInstance) return;
    echartsInstance.setOption(option);
  };

  // 获取实例
  const getInstance = () => {
    if (!echartsInstance) initCharts();
    return echartsInstance;
  };

  // 更新大小
  const onResize = debounce(() => {
    echartsInstance?.resize();
  }, 100);

  // 监听元素大小变化
  const watchEl = () => {
    if (isEmptyEl()) return;
    const resizeObserve = new ResizeObserver(onResize);
    resizeObserve.observe(unref(el));
  };

  // 组件挂载完成
  onMounted(() => {
    window.addEventListener("resize", onResize);
    if (autoUpdateSize) watchEl();
  });

  // 页面销毁
  onUnmounted(() => {
    window.removeEventListener("resize", onResize);
  });

  return { setOptions, getInstance };
};
