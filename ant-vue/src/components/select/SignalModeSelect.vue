<script lang="ts" setup>
import {onMounted, ref, watch} from "vue";

import { SignalPage } from "@/api";
import {useI18n} from "vue-i18n";
const props = defineProps({
  modelValue: {
    type: [String, Number, Object, Boolean],
    default: "",
  },
  mqtt_client_id: {
    type: Number,
    required: true,
    default: () => null,
  },
  name: {
    type: String,
    required: false,
    default: () => "ID",
  },
  show: {
    type: Boolean,
    required: false,
    default: () => false,
  },
});
const { t } = useI18n();
const mqttClientId = ref<number | string>(props.mqtt_client_id);
const page = ref(1);
const pageSelect = ref(1);
const options = ref<any>([]);
const value = ref<any>(String(props.modelValue) || "");
const valueResult = ref<any>(String(props.modelValue) || "");
const valueSearch = ref<any>("");
const showOpen = ref(false);
const emits = defineEmits(["update:modelValue", "custom-event"]);
watch(
  () => props.mqtt_client_id,
  async (newValue, oldValue) => {
    options.value = [];
    value.value = "";
    valueResult.value = "";
    page.value = 1;
    pageSelect.value = 1;
    mqttClientId.value = newValue;
    if (newValue && newValue !== oldValue) {
      await List();
    }
  },
);
const List = async () => {
  const { data } = await SignalPage({ mqtt_client_id: mqttClientId.value, page: page.value, page_size: 100, type: "数字" });
  const listArr = data.data.data.map((item: any) => ({ value: String(item[props.name]), label: item.name + "（" + item.alias + "）", alias: item.alias, unit: item.unit }));
  options.value = options.value.concat(listArr);
  if (!props.show) {
    if (!value.value) {
      value.value = options.value[0]?.value ? options.value[0]?.value : "";
    }
  }
  if (
    value.value &&
    options.value?.length &&
    !containsAllElements(
      value.value,
      options.value.map((it: any) => it.value),
    )
  ) {
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
  emits(
    "custom-event",
    options.value.filter((it: any) => value.value == it.value),
  );
  emits("update:modelValue", value.value);
  valueResult.value = value.value;
};


const select = async (valueClick: any) => {
  if (valueClick === -11) {
    value.value = valueResult.value;
    options.value.pop();
    if (!valueSearch.value) {
      await List();
    } else {
      const { data } = await SignalPage({ mqtt_client_id: valueSearch.value, page: pageSelect.value, page_size: 100, type: "数字" });
      const listArr = data.data.data.map((item: any) => ({ value: String(item[props.name]), label: item.name + "（" + item.alias + "）", alias: item.alias, unit: item.unit }));
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
  }
};

const onChange = (value:any, option:any) => {
  emits("custom-event", option);
  emits("update:modelValue", value);
};

function containsAllElements<T>(value: any, array: T[]): boolean {
  return array.includes(value);
}

onMounted(async ()=>{
  if (mqttClientId.value) {
    await List();
  }
})
</script>

<template>
  <a-select
    v-model:value="value"
    :placeholder="$t('message.pleaseEnter')"
    style="width: 300px"
    :open="showOpen"
    :default-active-first-option="false"
    :show-arrow="false"
    :filter-option="false"
    :not-found-content="null"
    :options="options"
    @select="select"
    @click="showOpen = true"
    @blur="showOpen = false"
    @change="onChange"
  ></a-select>
</template>

<style lang="less" scoped></style>
