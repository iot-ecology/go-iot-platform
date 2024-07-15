<template>
  <div class="clearfix">
    <a-upload
        :action="uploadUrl"
        list-type="picture-card"
        accept=".jpg,.png,.jpeg"
        :file-list="fileList"
        @preview="handlePreview"
        @change="handleChange"
        @success="handleSuccess"
        @error="handleError"
    >
      <div v-if="fileList.length < 1">
        <plus-outlined />
        <div class="ant-upload-text">Upload</div>
      </div>
    </a-upload>
    <a-modal width="50%" :visible="previewVisible" :footer="null" @cancel="handleCancel">
      <img alt="example" style="width: 100%;margin-top: 30px" :src="previewImage" />
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import {ref, watch} from 'vue'
import { PlusOutlined } from '@ant-design/icons-vue';
import { message } from "ant-design-vue";

const props = defineProps({
  imgStr: {
    type: String,
    default: "",
  },
});
interface FileItem {
  uid?: string;
  name?: string;
  status?: string;
  response?: string;
  percent?: number;
  url?: string;
  preview?: string;
  originFileObj?: any;
}
const emits = defineEmits(["update:modelValue"]);
const uploadUrl = `${import.meta.env.VITE_APP_API_URL}/file/update`;
const fileList = ref<FileItem[]>([]);
const previewVisible = ref<boolean>(false);
const previewImage = ref<string | undefined>('');

watch(()=>props.imgStr,(newValue)=>{
  if(newValue) {
    fileList.value.push({ url:`${import.meta.env.VITE_APP_API_URL}/file/download?path=${newValue}` })
  }
},{immediate:true})
const handleCancel = () => {
  previewVisible.value = false;
};

const handlePreview = async (file: FileItem) => {
  previewImage.value = file.thumbUrl || file.url;
  previewVisible.value = true;
};

const handleSuccess = (response: any, file: any) => {
  if(response.file_path) {
    message.success(response.message);
    updateFileStatus(file, 'done', response.file_path);
    emits("update:modelValue",response.file_path)
  } else {
    message.error(response.message);
    updateFileStatus(file, 'error');
  }
};

const handleError = (error: any, file: any) => {
  message.error(`${file.name} 上传失败`);
  updateFileStatus(file, 'error');
};

const handleChange = ({ fileList: newFileList }: any) => {
  fileList.value = newFileList;
  if(!fileList.value?.length) {
    emits("update:modelValue",'')
  }
};

const updateFileStatus = (file: any, status: string, url: string = '') => {
  const targetFile = fileList.value.find((item: FileItem) => item.uid === file.uid);
  if (targetFile) {
    targetFile.status = status;
    if (url) targetFile.url = url;
  }
};

</script>

<style lang="less">
.clearfix::after {
  content: "";
  display: block;
  clear: both;
}
</style>