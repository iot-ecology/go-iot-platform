<template>
  <div>
    <a-card :bordered="true">
      <a-button style="margin: 10px 0" type="primary" @click="modalVisible = true">{{ $t("message.addition") }}</a-button>
      <a-table :columns="columns" :data-source="list" bordered :pagination="pagination" @change="handleTableChange">
        <template #bodyCell="{ column, text, record }">
          <template v-if="['name', 'password', 'host'].includes(String(column.dataIndex))">
            <div>
              <template v-if="editableData[record.key] && column.dataIndex === 'name'">
                <a-input v-model:value="editableData[record.key][column.dataIndex]" style="width: 200px" />
              </template>
              <template v-else>
                <a-input-password v-if="column.dataIndex === 'password' || column.dataIndex === 'host'" :value="text" :bordered="false" />
                <data v-else>{{ text }}</data>
              </template>
            </div>
          </template>
          <template v-else-if="column.dataIndex === 'operation'">
            <div class="editable-row-operations">
              <span v-if="editableData[record.key]">
                <a-button type="primary" size="small" style="margin-right: 10px" @click="save(record.key)">{{ $t("message.save") }}</a-button>
                <a-popconfirm :title="$t('message.sureEdit')" @confirm="cancel(record.key)">
                  <a-button type="primary" size="small">{{ $t("message.cancel") }}</a-button>
                </a-popconfirm>
              </span>
              <span v-else>
                <a-button style="margin-left: 10px" type="primary" size="small" @click="edit(record.key)">{{ $t("message.edit") }}</a-button>
                <!-- <a-popconfirm :title="$t('message.sureDelete')" :okText="$t('message.yes')" :cancelText="$t('message.no')" @confirm="deleteById(record.ID)">
                  <a-button type="primary" size="small" danger style="margin-left: 10px">{{ $t("message.delete") }}</a-button>
                </a-popconfirm> -->
              </span>
            </div>
          </template>
        </template>
      </a-table>
      <a-modal
        :okText="$t('message.confirm')"
        :cancelText="$t('message.cancel')"
        v-model:open="modalVisible"
        :destroyOnClose="true"
        :title="$t('message.addition')"
        @ok="onAddData()"
        @cancel="onCancelData()"
      >
        <a-form ref="formRef" :label-col="{ style: { width: '140px' } }" :labelWrap="true" :rules="rules" :model="form">
          <a-form-item :label="$t('message.name')" name="name">
            <a-input v-model:value="form.name" style="width: 250px" />
          </a-form-item>
          <a-form-item :label="$t('message.host')" name="host">
            <a-input-password v-model:value="form.host" style="width: 250px" />
          </a-form-item>
          <a-form-item :label="$t('message.port')" name="port">
            <a-input-number v-model:value="form.port" style="width: 250px" />
          </a-form-item>
          <a-form-item :label="$t('message.username')" name="username">
            <a-input v-model:value="form.username" style="width: 250px" />
          </a-form-item>
          <a-form-item :label="$t('message.password')" name="password">
            <a-input-password v-model:value="form.password" style="width: 250px" />
          </a-form-item>
        </a-form>
      </a-modal>
    </a-card>
  </div>
</template>
<script setup lang="ts">
import { reactive, ref } from "vue";
import { onMounted } from "vue";
import { useI18n } from "vue-i18n";
import { useColumns } from "@/hooks/useColumns";
import { usePagination } from "@/hooks/usePagination";
import { useEditableData } from "@/hooks/useEditableData";
import { useRoles } from "@/hooks/useRoles";
import { useQueryPage } from "@/hooks/useQueryPage";
import { useAddPage } from "@/hooks/useAddPage";
import { useUpdatePage } from "@/hooks/useUpdatePage";
import { useDeletePage } from "@/hooks/useDeletePage";
import { CLICK_HOUSE_CREATE, CLICK_HOUSE_DELETE, CLICK_HOUSE_PAGE, CLICK_HOUSE_UPDATE } from "@/api";

const { t } = useI18n();
const { columns } = useColumns([
  {
    title: "message.id",
    dataIndex: "ID",
  },
  {
    title: "message.name",
    dataIndex: "name",
  },
  {
    title: "message.host",
    dataIndex: "host",
  },
  {
    title: "message.port",
    dataIndex: "port",
  },
  {
    title: "message.username",
    dataIndex: "username",
  },
  {
    title: "message.password",
    dataIndex: "password",
  },
  {
    title: "message.operation",
    dataIndex: "operation",
  },
]);

const { rules } = useRoles({
  name: [{ required: true, message: "message.name", trigger: "blur" }],
  host: [{ required: true, message: "message.host", trigger: "blur" }],
  port: [{ required: true, message: "message.port", trigger: "blur" }],
  username: [{ required: true, message: "message.username", trigger: "blur" }],
  password: [{ required: true, message: "message.password", trigger: "blur" }],
});

const { pagination } = usePagination();

const formRef = ref<HTMLFormElement | null>(null);
const modalVisible = ref(false);
const form = reactive({
  name: "",
  host: "",
  port: 0,
  username: "",
  password: "",
});
const list = ref([]);

const { editableData, edit, cancel } = useEditableData(list);

const { pageList } = useQueryPage({ path: CLICK_HOUSE_PAGE, type: "get" }, pagination, list);

const { onAddData, onCancelData } = useAddPage({ path: CLICK_HOUSE_CREATE, type: "post" }, formRef, form, modalVisible, pageList);

const { save } = useUpdatePage({ path: CLICK_HOUSE_UPDATE, type: "post" }, list, editableData, pageList);

// const { deleteById } = useDeletePage({ path: CLICK_HOUSE_DELETE, type: "delete" }, pageList);

const handleTableChange = async (page: any) => {
  pagination.current = page.current;
  pagination.pageSize = page.pageSize;
  await pageList();
};

onMounted(async () => {
  await pageList();
});
</script>
<style></style>
