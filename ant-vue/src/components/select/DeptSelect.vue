<script lang="ts" setup>
import {onMounted, ref, watch} from "vue";
import {DeptSubs} from "@/api";

const props = defineProps({
  modelValue: {
    type: [String, Number, Object, Boolean],
    default: "",
  },
  multiple: {
    type: Boolean,
    default: false,
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
  treeData.value = data.data?.map((item: any) => ({id: item.ID,value: item.ID,title: item.name, children:[]}))
  treeData.value.forEach(item=>{
    getTree(item,item.id)
  })
};

// const onLoadData = async(treeNode: any) => {
//   const { value } = treeNode
//   const { data } = await DeptSubs({ id: value});
//   const listArr = data.data?.map((item: any) => ({id: item.ID,value: item.ID,pId:item.parent_id,title: item.name,children:[]}))
//   treeData.value.forEach((it: any) => {
//     if(it.value === value) {
//       it.children = listArr
//     }
//   })
// }

async function getTree(item:any,id:string) {
  if(!item) return
  const { data } = await DeptSubs({ id: id});
  item.children = data.data?.map((item: any) => ({id: item.ID, value: item.ID, pId: item.parent_id, title: item.name, children: []}))
  if(item.id) {
    item.children.forEach((text:any)=> {
      getTree(text,text.id)
    })
  }
}

onMounted(async ()=>{
  await List();
})
</script>

<template>
  <a-tree-select
      v-model:value="value"
      style="width: 100%"
      allow-clear
      treeDefaultExpandAll
      showCheckedStrategy
      :multiple="props.multiple"
      :dropdown-style="{ maxHeight: '400px', overflow: 'auto' }"
      :tree-data="treeData"
      :placeholder="$t('message.pleaseEnter')"
  />
</template>

<style lang="less" scoped></style>
