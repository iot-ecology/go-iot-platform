<template>
  <div>
    <a-card title="" :bordered="true">
      <a-button style="margin: 10px 0" type="primary" @click="onAdd()">{{ $t('message.addition') }}</a-button>
      <a-table :data-source="dataSource" bordered :columns="columns" :pagination="pagination" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'operation'">
            <div class="editable-row-operations">
              <span>
                <a-button type="primary" size="small" @click="onDetails(record.ID)">{{ $t('message.check') }}</a-button>
                <a-popconfirm :title="$t('message.sureDelete')" :okText="$t('message.yes')" :cancelText="$t('message.no')" @confirm="confirm(record.ID)" @cancel="cancel1">
                  <a-button type="primary" size="small" danger style="margin-left: 10px;">{{$t('message.delete')}}</a-button>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>
<script setup lang="ts">
import {onMounted, reactive, ref, watch} from "vue";
import { message } from "ant-design-vue";

import { DashboardDelete, DashboardPage } from "@/api";
import { useRouteJump } from "@/hooks/useRouteJump.ts";
import {useI18n} from "vue-i18n";
const { t,locale } = useI18n();

const jump = useRouteJump();
const columns = ref([
  {
    title: t('message.name'),
    dataIndex: "name",
  },
  {
    title: t('message.operation'),
    dataIndex: "operation",
  },
]);
const dataSource = ref([]);
const pagination = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
  showSizeChanger: true, // 显示每页显示条目数选择器
});
watch(locale, () => {
  columns.value = [
    {
      title: t('message.name'),
      dataIndex: "name",
    },
    {
      title: t('message.operation'),
      dataIndex: "operation",
    },
  ]
});


const onAdd = () => {
  jump.routeJump({ path: "/visualization/add" });
};

const listPage = async () => {
  const { data } = await DashboardPage({ page: pagination.current, page_size: pagination.pageSize });
  pagination.total = data.data?.total || 0;
  dataSource.value = data.data.data;
};
const confirm = async (id: string) => {
  DashboardDelete(id).then(({ data }) => {
    if (data.code === 20000) {
      message.success(data.message);
      listPage();
    } else {
      // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
      message.success(data.message);
    }
  });
};

const onDetails = (id: string) => {
  jump.routeJump({ path: "/visualization/add", query: { id } });
};
const handleTableChange = async (page: any) => {
  pagination.current = page.current;
  pagination.pageSize = page.pageSize;
  await listPage();
};
onMounted(async ()=>{
  await listPage();
})
</script>
<style lang="less" scoped></style>
