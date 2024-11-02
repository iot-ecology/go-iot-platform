<script lang="ts" setup>
import {onMounted, ref, watch} from "vue";
import { message } from "ant-design-vue";

import { ProductPage} from "@/api";
import {useI18n} from "vue-i18n";
const props = defineProps({
  modelValue: {
    type: [String, Number, Object, Boolean],
    default: "",
  },
});
const { t } = useI18n();
const page = ref(1);
const pageSelect = ref(1);
const options = ref<any>([]);
const value = ref<any>(props.modelValue);
const valueResult = ref<any>("");
const valueSearch = ref<any>("");
const showOpen = ref(false);
const unit = ref('');

const emits = defineEmits(["update:modelValue",'dataUnit']);
const List = async () => {
  const { data } = await ProductPage({ name: "", page: page.value, page_size: 100 });
  const listArr = data.data.data.map((item: any) => ({ value: item.ID, label: item.name,unit:item.sku }));
  options.value = options.value.concat(listArr);
  if (Number(props.modelValue) && !options.value.map((it: any) => it.value).includes(Number(props.modelValue))) {
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
  if(!unit.value) {
    unit.value = options.value.filter((it:any)=>it.value===value.value)[0]?.unit
    console.log(unit.value)
    emits("dataUnit", unit.value);
  }

  emits("update:modelValue", value.value);
  // emits("dataUnit", unit.value);
  valueResult.value = value.value;
};

watch(()=>props.modelValue,(newValue)=>{
  if(newValue) {
    value.value = newValue
  }
})
const select = async (ValueClick: any) => {
  if (ValueClick === -11) {
    value.value = valueResult.value;
    options.value.pop();
    if (!valueSearch.value) {
      await List();
    } else {
      const { data } = await ProductPage({ name: valueSearch.value, page: pageSelect.value, page_size: 100 });
      const listArr = data.data.data.map((item: any) => ({ value: item.ID, label: item.name, unit:item.sku }));
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
    showOpen.value = false;
    emits("update:modelValue", value.value);
    valueResult.value = value.value;
  }
};
const handleSearch = async (val: string) => {
  if (valueSearch.value === val) {
    return;
  }
  valueSearch.value = val;
  pageSelect.value = 1;
  const { data } = await ProductPage({ name: val, page: pageSelect.value, page_size: 100 });
  if (data.data.total === 0) {
    message.error(t('message.ThereNoSearch'));
    setTimeout(() => {
      value.value = options.value[0].value;
      showOpen.value = false;
      page.value = 2;
    }, 1000);
    return;
  }
  options.value = [];
  const listArr = data.data.data.map((item: any) => ({ value: item.ID, label: item.name, unit:item.sku }));
  options.value = options.value.concat(listArr);
  if (data.data.total > 0 && options.value.length < data.data.total) {
    pageSelect.value++;
    options.value.push({
      value: -11,
      label: t('message.loadMore'),
    });
  }
};

watch(value, async (newValue) => {
  if (!newValue) {
    options.value = [];
    page.value = 1;
    valueSearch.value = "";
    await List();
  }else {
    value.value = newValue
    unit.value = options.value.filter((it:any)=>it.value===value.value)[0]?.unit
    emits("dataUnit", unit.value);
    emits("update:modelValue", newValue);
  }
});

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
  >
  </a-select>
</template>

<style lang="less" scoped></style>
