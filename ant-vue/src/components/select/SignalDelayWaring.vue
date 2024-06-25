<script lang="ts" setup>
import {onMounted, ref, watch} from "vue";
import { useRoute } from "vue-router";
import { message } from "ant-design-vue";

import { SignalDelayWaringPage } from "@/api";
import {useI18n} from "vue-i18n";
defineProps({
  modelValue: {
    type: [String, Number, Object, Boolean],
    default: "",
  },
});
const { t } = useI18n();
const page = ref(1);
const pageSelect = ref(1);
const options = ref<any>([]);
const route = useRoute();
const value = ref<any>(Number(route.query.rule_id) || "");
const valueResult = ref<any>(Number(route.query.rule_id) || "");
const valueSearch = ref<any>("");
const showOpen = ref(false);
const emits = defineEmits(["update:modelValue"]);
const List = async () => {
  const { data } = await SignalDelayWaringPage({ name: "", page: page.value, page_size: 100 });
  const listArr = data.data.data.map((item: any) => ({ value: item.ID, label: item.name }));
  options.value = options.value.concat(listArr);
  value.value = value.value || Number(route.query.signal_delay_waring_id) || options.value[0]?.value;
  if (data.data.total > 0 && Number(route.query.signal_delay_waring_id) && !options.value.map((it: any) => it.value).includes(Number(route.query.signal_delay_waring_id))) {
    page.value++;
    await List();
  } else {
    if (data.data.total > 0 && options.value.length < data.data.total) {
      page.value++;
      options.value.push({
        value: -11,
        label: t('message.loadMore'),
      });
    }
  }
  emits("update:modelValue", value.value);
  valueResult.value = value.value;
};
watch(value, async (newValue) => {
  if (!newValue) {
    options.value = [];
    page.value = 1;
    valueSearch.value = "";
    await List();
  }
});
const select = async (ValueClick: any) => {
  if (ValueClick === -11) {
    value.value = valueResult.value;
    options.value.pop();
    if (!valueSearch.value) {
      await List();
    } else {
      const { data } = await SignalDelayWaringPage({ name: valueSearch.value, page: pageSelect.value, page_size: 100 });
      const listArr = data.data.data.map((item: any) => ({ value: item.ID, label: item.name }));
      options.value = options.value.concat(listArr);
      if (data.data.total > 0 && options.value.length < data.data.total) {
        pageSelect.value++;
        options.value.push({
          value: -11,
          label: t('message.loadMore'),
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

  pageSelect.value = 1;
  const { data } = await SignalDelayWaringPage({ name: val, page: pageSelect.value, page_size: 100 });
  if (data.data.total === 0) {
    message.error(t('message.ThereNoSearch'));
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
      label: t('message.loadMore'),
    });
  }
};

onMounted(async ()=>{
  await List();
})
</script>

<template>
  <a-select
    v-model:value="value"
    :show-search="true"
    :open="showOpen"
    allow-clear
    :placeholder="$t('message.pleaseEnter')"
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
