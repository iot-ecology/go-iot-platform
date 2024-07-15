<template>
  <template v-for="tag in state.tags" :key="tag">
    <a-tooltip v-if="tag.length > 10" :title="tag">
      <a-tag closable @close="handleClose(tag)">
        {{ `${tag.slice(0, 10)}...` }}
      </a-tag>
    </a-tooltip>
    <a-tag v-else closable @close="handleClose(tag)">
      {{ tag }}
    </a-tag>
  </template>
  <a-input
      v-if="state.inputVisible"
      ref="inputRef"
      type="text"
      size="small"
      :style="{ width: '78px' }"
      v-model:value="state.inputValue"
      @blur="handleInputConfirm"
      @keyup.enter="handleInputConfirm"
  />
  <a-tag v-else @click="showInput" style="background: #fff; border-style: dashed">
    <plus-outlined />
    New Tag
  </a-tag>
</template>
<script lang="ts" setup>
import {nextTick, reactive, ref, watch} from 'vue';
import {PlusOutlined} from '@ant-design/icons-vue';

const props = defineProps({
  tags: {
    type: String,
    default: "",
  },
});
const emits = defineEmits(["update:modelValue"]);
const inputRef = ref();
const state = reactive({
  tags: [],
  inputVisible: false,
  inputValue: '',
});
watch(()=>props.tags,(newValue)=> {
  if(newValue) {
    state.tags = newValue.split(',')
  }
},{immediate:true});
const handleClose = (removedTag: string) => {
  state.tags = state.tags.filter(tag => tag !== removedTag);
  emits("update:modelValue", state.tags.join(','));
};

const showInput = () => {
  state.inputVisible = true;
  nextTick(() => {
    inputRef.value.focus();
  });
};
const handleInputConfirm = () => {
  const inputValue = state.inputValue;
  let tags = state.tags;
  if (inputValue && tags.indexOf(inputValue) === -1) {
    tags = [...tags, inputValue];
  }
  Object.assign(state, {
    tags,
    inputVisible: false,
    inputValue: '',
  });
  emits("update:modelValue", state.tags.join(','));
};
</script>


