<script lang="ts" setup>
import { ref, watch } from "vue";
import { useRoute } from "vue-router";
import { message } from "ant-design-vue";

import { CalcRulePage } from "@/api";
defineProps({
  modelValue: {
    type: [String, Number, Object, Boolean],
    default: "",
  },
});
const page = ref(1);
const pageSelect = ref(1);
const options = ref([]);
const route = useRoute();
const value = ref(Number(route.query.rule_id) || "");
const valueResult = ref(Number(route.query.rule_id) || "");
const valueSearch = ref("");
const showOpen = ref(false);
const emits = defineEmits(["update:modelValue"]);
const List = async () => {
  const { data } = await CalcRulePage({ name: "", page: page.value, page_size: 100 });
  const listArr = data.data.data.map((item: any) => ({ value: item.ID, label: item.name }));
  options.value = options.value.concat(listArr);
  value.value = value.value || Number(route.query.rule_id) || options.value[0]?.value;
  if (data.data.total > 0 && Number(route.query.rule_id) && !options.value.map((it) => it.value).includes(Number(route.query.rule_id))) {
    page.value++;
    await List();
  } else {
    if (data.data.total > 0 && options.value.length < data.data.total) {
      page.value++;
      options.value.push({
        value: -11,
        label: "加载更多",
      });
    }
  }
  emits("update:modelValue", value.value);
  valueResult.value = value.value;
};
List();
watch(value, (newValue, oldValue) => {
  if (!newValue) {
    options.value = [];
    page.value = 1;
    valueSearch.value = "";
    List();
  }
});
const select = async (ValueClick: any) => {
  if (ValueClick === -11) {
    value.value = valueResult.value;
    options.value.pop();
    if (!valueSearch.value) {
      await List();
    } else {
      const { data } = await CalcRulePage({ name: valueSearch.value, page: pageSelect.value, page_size: 100 });
      const listArr = data.data.data.map((item: any) => ({ value: item.ID, label: item.name }));
      options.value = options.value.concat(listArr);
      if (data.data.total > 0 && options.value.length < data.data.total) {
        pageSelect.value++;
        options.value.push({
          value: -11,
          label: "加载更多",
        });
      }
    }
  } else {
    emits("update:modelValue", value.value);
    valueResult.value = value.value;
    showOpen.value = false;
  }
};
const handleSearch = async (val: string) => {
  if (valueSearch.value === val) {
    return;
  }
  valueSearch.value = val;

  // options.value = [];
  pageSelect.value = 1;
  const { data } = await CalcRulePage({ name: val, page: pageSelect.value, page_size: 100 });
  if (data.data.total === 0) {
    message.error("当前搜索没有相关数据");
    setTimeout(() => {
      value.value = options.value[0].value;
      valueSearch.value = "";
      showOpen.value = false;
      page.value = 2;
    }, 1000);
    return;
  }
  options.value = [];
  const listArr = data.data.data.map((item: any) => ({ value: item.ID, label: item.name }));
  options.value = options.value.concat(listArr);
  if (data.data.total > 0 && options.value.length < data.data.total) {
    pageSelect.value++;
    options.value.push({
      value: -11,
      label: "加载更多",
    });
  }
};
</script>

<template>
  <a-select
    v-model:value="value"
    :show-search="true"
    :open="showOpen"
    allow-clear
    placeholder="请输入"
    style="width: 300px"
    :default-active-first-option="false"
    :show-arrow="false"
    :filter-option="false"
    :not-found-content="null"
    :options="options"
    @search="handleSearch"
    @select="select"
    @click="showOpen = true"
    @blur="showOpen = false"
  ></a-select>
</template>

<style lang="less" scoped></style>
