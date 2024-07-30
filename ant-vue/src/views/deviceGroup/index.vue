<template>
  <div>
    <a-card :bordered="true">
      <a-button style="margin: 10px 0" type="primary" @click="modalVisible = true">{{ $t("message.addition") }}</a-button>
      <a-table :columns="columns" :data-source="list" bordered :pagination="pagination" @change="handleTableChange">
        <template #bodyCell="{ column, text, record }">
          <template v-if="['name'].includes(String(column.dataIndex))">
            <div>
              <template v-if="editableData[record.key] && column.dataIndex === 'name'">
                <a-input v-model:value="editableData[record.key][column.dataIndex]" style="width: 200px" />
              </template>
              <template v-else>
                {{ text }}
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
                <a-button size="small" @click="handleBindDevice(record.ID)">{{ $t("message.bindDevice") }}</a-button>
                <a-button style="margin-left: 10px" type="primary" size="small" @click="edit(record.key)">{{ $t("message.edit") }}</a-button>
                <a-popconfirm :title="$t('message.sureDelete')" :okText="$t('message.yes')" :cancelText="$t('message.no')" @confirm="deleteById(record.ID)">
                  <a-button type="primary" size="small" danger style="margin-left: 10px">{{ $t("message.delete") }}</a-button>
                </a-popconfirm>
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
          <a-form-item :label="$t('message.deviceGroupsName')" name="name">
            <a-input v-model:value="form.name" style="width: 200px" />
          </a-form-item>
        </a-form>
      </a-modal>
      <!-- 绑定设备 -->
      <bind-device :isShow="showBindDevice" @closeModal="handleCloseModal" @onSubmit="handlerSubmit" :group_id="group_id" />
    </a-card>
  </div>
</template>
<script setup lang="ts">
import { reactive, ref, watch } from "vue";
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
import { DEVICE_GROUP_CREATE, DEVICE_GROUP_DELETE, DEVICE_GROUP_PAGE, DEVICE_GROUP_UPDATE, DeviceGroupsBindDevice } from "@/api";
import bindDevice from "./modal/bindDevice.vue";
import { message } from "ant-design-vue";

const { t } = useI18n();
const { columns } = useColumns([
  {
    title: "message.id",
    dataIndex: "ID",
  },
  {
    title: "message.deviceGroupsName",
    dataIndex: "name",
  },
  {
    title: "message.operation",
    dataIndex: "operation",
  },
]);

const { rules } = useRoles({
  name: [{ required: true, message: "message.deviceGroupsName", trigger: "blur" }],
});

const { pagination } = usePagination();

const formRef = ref<HTMLFormElement | null>(null);
const modalVisible = ref(false);
const form = reactive({
  name: "",
});
const list = ref([]);

const { editableData, edit, cancel } = useEditableData(list);

const { pageList } = useQueryPage({ path: DEVICE_GROUP_PAGE, type: "get" }, pagination, list);

const { onAddData, onCancelData } = useAddPage({ path: DEVICE_GROUP_CREATE, type: "post" }, formRef, form, modalVisible, pageList);

const { save } = useUpdatePage({ path: DEVICE_GROUP_UPDATE, type: "post" }, list, editableData, pageList);

const { deleteById } = useDeletePage({ path: DEVICE_GROUP_DELETE, type: "delete" }, pageList);

const handleTableChange = async (page: any) => {
  pagination.current = page.current;
  pagination.pageSize = page.pageSize;
  await pageList();
};

const showBindDevice = ref(false);
const group_id = ref("");

const handleBindDevice = (id: string) => {
  showBindDevice.value = true;
  group_id.value = id;
};

const handleCloseModal = () => {
  showBindDevice.value = false;
};

// 确认绑定
const handlerSubmit = async (info: any) => {
  const params = {
    ...info,
    group_id: group_id.value,
  };
  try {
    let result = await DeviceGroupsBindDevice(params);
    const { data } = result;
    if (data.code === 20000) {
      message.success(data.message);
      modalVisible.value = false;
      formRef.value?.resetFields();
      showBindDevice.value = false;
    } else {
      message.error(`${t("message.operationFailed")}:${data.data}`);
    }
  } catch (error) {
    console.log(error);
  }
};

onMounted(async () => {
  await pageList();
});
</script>
<style></style>
