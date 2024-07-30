<template>
  <div>
    <a-card :bordered="true">
      <a-button style="margin: 10px 0" type="primary" @click="modalVisible = true">{{ $t("message.addition") }}</a-button>
      <a-table :columns="columns" :data-source="list" bordered :pagination="pagination" @change="handleTableChange">
        <template #headerCell="{ column }">
          <template v-if="column.dataIndex === 'content'">
            <span>
              {{ $t("message.dingdingContent") }}
            </span>
            <a-tooltip>
              <template #title><span v-html="contentToolTip"></span></template>
              <QuestionCircleOutlined :style="{ marginLeft: '5px' }" />
            </a-tooltip>
          </template>
        </template>
        <template #bodyCell="{ column, text, record }">
          <template v-if="['name', 'access_token', 'content', 'secret'].includes(String(column.dataIndex))">
            <div>
              <template v-if="editableData[record.key] && column.dataIndex === 'access_token'">
                <a-input v-model:value="editableData[record.key][column.dataIndex]" style="width: 200px" />
              </template>
              <template v-else-if="editableData[record.key] && column.dataIndex === 'content'">
                <a-textarea v-model:value="editableData[record.key][column.dataIndex]" style="width: 200px" />
              </template>
              <template v-else-if="editableData[record.key] && column.dataIndex === 'name'">
                <a-input v-model:value="editableData[record.key][column.dataIndex]" style="width: 200px" />
              </template>
              <template v-else-if="editableData[record.key] && column.dataIndex === 'secret'">
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
          <a-form-item :label="$t('message.dingdingAccessToken')" name="access_token">
            <a-input v-model:value="form.access_token" style="width: 250px" />
          </a-form-item>
          <a-form-item :label="$t('message.dingdingContent')" name="content">
            <template #tooltip>
              <a-tooltip>
                <template #title><span v-html="contentToolTip"></span></template>
                <QuestionCircleOutlined :style="{ marginLeft: '5px' }" />
              </a-tooltip>
            </template>
            <a-textarea v-model:value="form.content" style="width: 250px" />
          </a-form-item>
          <a-form-item :label="$t('message.dingdingName')" name="name">
            <a-input v-model:value="form.name" style="width: 250px" />
          </a-form-item>
          <a-form-item :label="$t('message.dingdingSecret')" name="secret">
            <a-input v-model:value="form.secret" style="width: 250px" />
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
import { DingDing_CREATE, DingDing_DELETE, DingDing_PAGE, DingDing_UPDATE } from "@/api";
import { QuestionCircleOutlined } from "@ant-design/icons-vue";

const contentToolTip = `
  设备名称: {{device_name}}<br>
  信号名称: {{signal_name}}<br>
  小值: {{min}}<br>
  大值: {{max}}<br>
  当前值: {{signal_value}}<br>
  单位: {{unit}}<br>
  范围内，范围外: {{in_or_out}}
`;

const { t } = useI18n();
const { columns } = useColumns([
  {
    title: "message.id",
    dataIndex: "ID",
  },
  {
    title: "message.dingdingName",
    dataIndex: "name",
  },
  {
    title: "message.dingdingAccessToken",
    dataIndex: "access_token",
  },
  {
    title: "message.dingdingContent",
    dataIndex: "content",
  },
  {
    title: "message.dingdingSecret",
    dataIndex: "secret",
  },
  {
    title: "message.operation",
    dataIndex: "operation",
  },
]);

const { rules } = useRoles({
  access_token: [{ required: true, message: "message.dingdingAccessToken", trigger: "blur" }],
  content: [{ required: true, message: "message.dingdingContent", trigger: "blur" }],
  name: [{ required: true, message: "message.dingdingName", trigger: "blur" }],
  secret: [{ required: true, message: "message.dingdingSecret", trigger: "blur" }],
});

const { pagination } = usePagination();

const formRef = ref<HTMLFormElement | null>(null);
const modalVisible = ref(false);
const form = reactive({
  name: "",
  access_token: "",
  content: "",
  secret: "",
});
const list = ref([]);

const { editableData, edit, cancel } = useEditableData(list);

const { pageList } = useQueryPage({ path: DingDing_PAGE, type: "get" }, pagination, list);

const { onAddData, onCancelData } = useAddPage({ path: DingDing_CREATE, type: "post" }, formRef, form, modalVisible, pageList);

const { save } = useUpdatePage({ path: DingDing_UPDATE, type: "post" }, list, editableData, pageList);

const { deleteById } = useDeletePage({ path: DingDing_DELETE, type: "delete" }, pageList);

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
