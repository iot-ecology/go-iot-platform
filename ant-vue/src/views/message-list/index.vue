<template>
  <div class="message-preview">
    <a-card :bordered="true">
      <a-table :columns="columns" :data-source="list" bordered :pagination="pagination" @change="handleTableChange">
      </a-table>
    </a-card>
  </div>
</template>
<script setup lang="ts">

import { reactive, ref, watch} from "vue";
import { MessageListPage} from "@/api";
import {onMounted} from "vue";
import {useI18n} from "vue-i18n";

const { t,locale } = useI18n();
const list = ref([]);

const pagination = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
  showSizeChanger: true, // 显示每页显示条目数选择器
});
let columns = ref([
  {
    title: t('message.uniCode'),
    dataIndex: "ID",
  },
  {
    title: t('message.messageContent'),
    dataIndex: "content",
  },
  {
    title: t('message.messageTypeId'),
    dataIndex: "message_type_id",
  },
  {
    title: t('message.refId'),
    dataIndex: "ref_id",
  },
]);
watch(locale, () => {
  columns.value = [
    {
      title: t('message.uniCode'),
      dataIndex: "ID",
    },
    {
      title: t('message.messageContent'),
      dataIndex: "content",
    },
    {
      title: t('message.messageTypeId'),
      dataIndex: "message_type_id",
    },
    {
      title: t('message.refId'),
      dataIndex: "ref_id",
    },
  ]
});

const pageList = async () => {
  const { data } = await MessageListPage({ page: pagination.current, page_size: pagination.pageSize });
  pagination.total = data.data?.total || 0;
  list.value = data.data.data?.map((item: any, index: number) => ({
    key: index,
    ID: item.ID,
    content: item.content,
    en_content: item.en_content,
    message_type_id: item.message_type_id,
    ref_id: item.ref_id,
  }));
};


onMounted(async()=>{
  await pageList()
})


</script>
<style>


</style>