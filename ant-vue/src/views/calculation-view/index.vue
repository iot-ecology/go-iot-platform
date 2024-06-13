<template>
  <div>
    <a-card title="" :bordered="true">
      <a-button style="margin: 10px 0" type="primary" @click="onAdd()">新增</a-button>
      <a-table :data-source="dataSource" :columns="columns" :pagination="paginations" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'operation'">
            <div class="editable-row-operations">
              <span>
                <a @click="onDetails(record.ID)">查看</a>
                <a-popconfirm title="确认是否删除?" ok-text="是" cancel-text="否" @confirm="confirm(record.ID)" @cancel="cancel1">
                  <a style="margin-left: 10px">删除</a>
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
import { reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { message } from "ant-design-vue";

import { DashboardDelete, DashboardPage } from "@/api";
const columns = ref([
  {
    title: "名称",
    dataIndex: "name",
  },
  {
    title: "操作",
    dataIndex: "operation",
  },
]);
const dataSource = ref([]);
const router = useRouter();
const paginations = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
  showSizeChanger: true, // 显示每页显示条目数选择器
});
const onAdd = () => {
  router.push({ path: "/calculation-view/add" });
};

const listPage = async () => {
  const { data } = await DashboardPage({ page: paginations.current, page_size: paginations.pageSize });
  paginations.total = data.data?.total || 0;
  dataSource.value = data.data.data;
};
listPage();
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
  router.push({ path: "/calculation-view/add", query: { id } });
};
const handleTableChange = async (pagination: any) => {
  paginations.current = pagination.current;
  paginations.pageSize = pagination.pageSize;
  await listPage();
};
</script>
<style lang="less" scoped></style>
