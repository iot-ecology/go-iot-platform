<template>
  <a-modal :okText="$t('message.confirm')" :cancelText="$t('message.cancel')" :open="isShow" :destroyOnClose="true" :title="$t('message.bindDevice')" @ok="onOk()" @cancel="onCancel()">
    <a-form ref="formRef" :label-col="{ style: { width: '140px' } }" :labelWrap="true" :rules="rules" :model="form">
      <a-form-item :label="$t('message.deviceName')" name="device_id">
        <a-select
          v-model:value="form.device_id"
          mode="tags"
          show-search
          style="width: 200px"
          :options="options"
          :filter-option="filterOption"
          @focus="handleFocus"
          @blur="handleBlur"
          @change="handleChange"
        ></a-select>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">
import { DeviceGroupsQueryBindDevice, DEVICEINFO_PAGE } from "@/api";
import { useQueryPage } from "@/hooks/useQueryPage";
import { useRoles } from "@/hooks/useRoles";
import { SelectProps } from "ant-design-vue";
import { onMounted, reactive, ref, toRef, watch } from "vue";

const props = defineProps({
  isShow: Boolean,
  group_id: Number,
});

const group_id = toRef(() => props.group_id);

const isShow = toRef(() => props.isShow);

const emit = defineEmits(["closeModal", "onSubmit"]);

const formRef = ref<HTMLFormElement | null>(null);
const form = reactive({
  device_id: [],
});

watch(isShow, async (val) => {
  if (val && group_id) {
    try {
      let result = await DeviceGroupsQueryBindDevice({ group_id: group_id.value });
      const { data } = result;
      if (data.code === 20000) {
        if (data.data.length) {
          form.device_id = data.data.map((item: { device_info_id: number }) => {
            return item.device_info_id;
          });
        }
      }
    } catch (error) {
      console.log(error);
    }
  }
});

const { rules } = useRoles({
  device_id: [{ required: true, message: "message.deviceName", trigger: "change" }],
});

const options = ref<SelectProps["options"]>([]);

onMounted(() => {
  getList();
});

const getList = async () => {
  const list = ref([]);
  const pagination = reactive({
    total: 0,
    current: 1,
    pageSize: 9999999,
    showSizeChanger: true, // 显示每页显示条目数选择器
  });
  const { pageList } = useQueryPage({ path: DEVICEINFO_PAGE, type: "get" }, pagination, list);
  await pageList();
  options.value = list.value.map((item) => {
    return {
      value: item.ID,
      label: item.sn,
    };
  });
};

const handleChange = (value: string) => {
  console.log(`selected ${value}`);
};
const handleBlur = () => {
  console.log("blur");
};
const handleFocus = () => {
  console.log("focus");
};
const filterOption = (input: string, option: any) => {
  return option.value.toLowerCase().indexOf(input.toLowerCase()) >= 0;
};

const onOk = () => {
  emit("onSubmit", form);
};

const onCancel = () => {
  emit("closeModal");
};

watch(isShow, () => {
  formRef.value?.resetFields();
});
</script>

<style lang="scss" scoped></style>
