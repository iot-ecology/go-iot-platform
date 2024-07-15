<script lang="ts" setup>
import {onMounted, ref, watch} from "vue";
import { DeptSubs} from "@/api";

const props = defineProps({
  modelValue: {
    type: [String, Number, Object, Boolean],
    default: "",
  },
  childValue:{
    type: [String, Number, Object, Boolean],
    default: "",
  }
});
const value = ref<any>();
const emits = defineEmits(["update:modelValue",'set']);
const treeData = ref([])

watch(value, async (newValue) => {
  emits("set", newValue)
});
const List = async () => {
  const { data } = await DeptSubs({ id: ""});
  treeData.value = data.data?.map((item: any) => ({id: item.ID,value: item.ID,pId:'',title: item.name, children:[]}))
};

const onLoadData = async(treeNode: any) => {
  const { value } = treeNode
  const { data } = await DeptSubs({ id: value});
  const listArr = data.data?.map((item: any) => ({id: item.ID,value: item.ID,pId:'',title: item.name,children:[]}))
  treeData.value.forEach((it: any) => {
    if(it.value === value) {
      it.children = listArr
    }
  })
}

onMounted(async ()=>{
  await List();
})
</script>

<template>
  <a-tree-select
      v-model:value="value"
      tree-data-simple-mode
      style="width: 100%"
      :dropdown-style="{ maxHeight: '400px', overflow: 'auto' }"
      :tree-data="treeData"
      :placeholder="$t('message.pleaseEnter')"
      :load-data="onLoadData"
  />
</template>

<style lang="less" scoped></style>
