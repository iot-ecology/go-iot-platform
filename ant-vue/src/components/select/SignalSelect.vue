<script lang="ts" setup>
import {onMounted, ref, watch} from "vue";
import { useRoute } from "vue-router";

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
  number: {
    type: Boolean,
    required: false,
    default: () => false,
  },
});

const { t } = useI18n();
const mqttClientId = ref<number | string>(props.mqtt_client_id);
const page = ref(1);
const page1 = ref(1);
const options = ref<any>([]);
const value = ref<any>(props.modelValue);
const valueResult = ref<any>(props.modelValue);
const valueSearch = ref<any>("");
const showOpen = ref(false);
const route = useRoute();
const emits = defineEmits(["update:modelValue", "custom-event"]);
watch(
  () => props.mqtt_client_id,
  async (newValue, oldValue) => {
    options.value = [];
    value.value = "";
    valueResult.value = "";
    page.value = 1;
    page1.value = 1;
    mqttClientId.value = newValue;
    if (newValue && newValue !== oldValue) {
      await List();
    }
  },
);
if (props.show) {
  value.value = "";
  valueSearch.value = "";
}
const List = async () => {
  const params = {
    mqtt_client_id: mqttClientId.value,
    page: page.value,
    page_size: 100,
  };
  if (props.number) {
    params.type = "数字";
  }
  const { data } = await SignalPage(params);
  const list1 = data.data.data.map((item: any) => ({ value: item[props.name], label: item.name + "（" + item.alias + "）", name: item.name, alias: item.alias, unit: item.unit }));
  options.value = options.value.concat(list1);
  if (!props.show) {
    value.value = value.value || Number(route.query.id) || options.value[0]?.value;
  }
  if (mqttClientId.value != route.query.mqtt_client_id && route.query.mqtt_client_id) {
    value.value = options.value[0]?.value;
  }
  if (
    data.data?.total > 0 &&
    ((value.value && !options.value.map((it: any) => it.value).includes(value.value)) || (Number(route.query.id) && !options.value.map((it: any) => it.value).includes(Number(route.query.id))))
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
  emits("update:modelValue", value.value);
  emits(
    "custom-event",
    options.value.find((it:any) => it.value === value.value),
  );
  valueResult.value = value.value;
};


const select = async (ValueClick: any) => {
  if (ValueClick === -11) {
    value.value = valueResult.value;
    options.value.pop();
    if (!valueSearch.value) {
      await List();
    } else {
      const params = {
        mqtt_client_id: mqttClientId.value,
        page: page.value,
        page_size: 100,
      };
      if (props.number) {
        params.type = "数字";
      }
      const { data } = await SignalPage(params);
      const list1 = data.data.data.map((item: any) => ({ value: item[props.name], label: item.name + "（" + item.alias + "）", name: item.name, alias: item.alias, unit: item.unit }));
      options.value = options.value.concat(list1);
      if (data.data.total > 0 && options.value.length < data.data.total) {
        page1.value++;
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
const onChange = (val:any, option:any) => {
  emits("custom-event", option);
  if (val !== -11) {
    emits("update:modelValue", value);
  }
};

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
    :show-search="true"
    :open="showOpen"
    style="width: 300px"
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
