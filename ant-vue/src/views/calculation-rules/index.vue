<template>
  <div>
    <a-card title="" :bordered="true">
      <a-form layout="inline" :model="formState">
        <a-form-item :label="$t('message.name')">
          <a-input v-model:value="formState.name" style="width: 300px" :placeholder="$t('message.pleaseEnter')" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="listPage">{{ $t('message.search') }}</a-button>
        </a-form-item>
      </a-form>
      <a-button style="margin: 10px 0" type="primary" @click="onAdd()">{{ $t('message.addition') }}</a-button>
      <a-table :data-source="dataSource" :columns="columns" bordered :pagination="pagination" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'operation'">
            <div class="editable-row-operations">
              <span>
                <a-button type="primary" size="small" v-if="!record.start" style="margin-left: 10px" @click="confirm(record)">{{$t('message.edit')}}</a-button>
                <a-button type="primary" size="small" v-if="!record.start" style="margin-left: 10px" @click="onGo(record.ID)">{{$t('message.parameterConfiguration')}}</a-button>
                <a-button type="primary" size="small" v-if="!record.start" style="margin-left: 10px" @click="onMock(record.ID)">{{$t('message.simulateExecution')}}</a-button>
                <a-button type="primary" size="small" v-if="!record.start" style="margin-left: 10px" @click="onStart(record.ID, record.mock_value)">{{ $t('message.startUp') }}</a-button>

                <a-popconfirm v-else :title="$t('message.sureStop')" :okText="$t('message.yes')" :cancelText="$t('message.no')" @confirm="confirmStop(record.ID)">
                  <a-button type="primary" size="small" style="margin-left: 10px">{{ $t('message.stop') }}</a-button>
                </a-popconfirm>
                <a-button type="primary" size="small" style="margin-left: 10px" @click="onResult(record.ID)">{{ $t('message.resultViewing') }}</a-button>
              </span>
            </div>
          </template>
        </template>
      </a-table>

      <a-modal v-model:open="modalVisible" :destroy-on-close="true" :title="title" class="custom-modal">
        <a-form ref="formRef" :label-col="{ style: { width: '110px' } }" :rules="rules" :model="form">
          <a-form-item :label="$t('message.name')" name="name">
            <a-input v-model:value="form.name" style="width: 300px" :disabled="form.start" />
          </a-form-item>
          <a-form-item :label="$t('message.forwardTime')" name="offset">
            <a-input-number v-model:value="form.offset" :precision="0" style="width: 300px" :min="11" :disabled="form.start" />
          </a-form-item>
          <a-form-item :label="$t('message.executionCycle')" name="cron">
            <a-input v-model:value="form.cron" style="width: 300px" :disabled="form.start" />
          </a-form-item>
          <a-form-item :label="$t('message.script')" name="script">
            <a-tooltip placement="right">
              <template #title>
                <span @click="onCopy">{{ scr }}</span>
              </template>
              <div style="cursor: pointer; width: 60px; line-height: 32px">{{ $t('message.listScript') }}</div>
            </a-tooltip>
            <codemirror v-model="form.script" :disabled="form.start" placeholder="Code here..." :style="{ width: '300px', height: '150px' }" :tab-size="2" :extensions="extensions" />
          </a-form-item>
        </a-form>
        <template #footer>
          <a-button @click="handleCancel">{{$t('message.cancel')}}</a-button>
          <a-button :disabled="loading" type="primary" @click="onAddData()">{{$t('message.confirm')}}</a-button>
        </template>
      </a-modal>

      <a-modal v-model:open="modalTime" :title="$t('message.timeframe')" class="custom-modal">
        <a-form ref="formRefTime" :rules="rules" :model="formTime">
          <a-form-item :label="$t('message.timeframe')" name="date">
            <a-range-picker :placeholder="[$t('message.startTime'), $t('message.endTime')]" v-model:value="formTime.date" show-time @change="onTimeChange" />
          </a-form-item>
        </a-form>
        <template #footer>
          <a-button @click="modalTime = false">{{$t('message.cancel')}}</a-button>
          <a-button :disabled="loading" type="primary" @click="setMockData()">{{$t('message.confirm')}}</a-button>
        </template>
      </a-modal>

      <a-modal v-model:open="modalDate" :title="$t('message.timeframe')" class="custom-modal">
        <a-spin :tip="$t('message.loading')" size="large" :spinning="showSpinning">
          <a-form ref="formRefDate" :rules="rules" :model="formDate">
            <a-form-item :label="$t('message.timeframe')" name="date">
              <a-range-picker :placeholder="[$t('message.startTime'), $t('message.endTime')]" v-model:value="formDate.date" show-time @change="onResultTime" />
            </a-form-item>
          </a-form>
        </a-spin>
        <template #footer>
          <a-button v-if="!showSpinning" @click="modalDate = false">{{$t('message.cancel')}}</a-button>
          <a-button :loading="showSpinning" type="primary" @click="setTableData()">{{$t('message.confirm')}}</a-button>
        </template>
      </a-modal>

      <a-modal v-model:open="modalMock" :title="$t('message.confirmExpectations')" class="custom-modal">
        <a-form :model="formTime">
          <a-form-item :label="$t('message.data')">
            <div>{{ mockValue }}</div>
          </a-form-item>
        </a-form>
        <template #footer>
          <a-button @click="modalMock = false">{{$t('message.cancel')}}</a-button>
          <a-button :disabled="loading" type="primary" @click="onConfirmStart()">{{$t('message.confirm')}}</a-button>
        </template>
      </a-modal>

      <a-modal v-model:open="modalTable" :footer="null" style="width: 600px; height: 650px; overflow: scroll" title="" class="custom-modal">
        <div style="margin-top: 30px; border: 1px solid #f0f0f0">
          <div style="display: flex; border-bottom: 1px solid #f0f0f0; text-align: center; font-size: 16px">
            <div style="width: 27%; border-right: 1px solid #f0f0f0; padding: 8px 0">{{ $t('message.executionTime') }}</div>
            <div style="width: 27%; border-right: 1px solid #f0f0f0; padding: 8px 0">{{ $t('message.queryStartTime') }}</div>
            <div style="width: 27%; border-right: 1px solid #f0f0f0; padding: 8px 0">{{ $t('message.queryEndTime') }}</div>
            <div style="width: 19%; padding: 8px 0">{{ $t('message.result') }}</div>
          </div>
          <div v-if="!dataResult?.length" style="text-align: center; font-size: 18px; height: 300px">{{ $t('message.noData') }}</div>
          <RecycleScroller v-else v-slot="{ item }" style="height: 480px" class="scroller" :items="dataResult" :item-size="40" key-field="id">
            <div style="display: flex; text-align: center; border-bottom: 1px solid #f0f0f0">
              <div style="width: 27%; border-right: 1px solid #f0f0f0; padding: 8px 0">{{ item.ex_time }}</div>
              <div style="width: 27%; border-right: 1px solid #f0f0f0; padding: 8px 0">{{ item.start_time }}</div>
              <div style="width: 27%; border-right: 1px solid #f0f0f0; padding: 8px 0">{{ item.end_time }}</div>
              <div style="width: 19%; padding: 8px 0; cursor: pointer; color: rgb(24, 144, 255)" @click="onView(item)">{{ $t('message.check') }}</div>
            </div>
          </RecycleScroller>
        </div>
      </a-modal>

      <a-modal v-model:open="modalResult" :footer="null" title="" class="custom-modal">
        <div>{{ textResult }}</div>
      </a-modal>
    </a-card>
  </div>
</template>
<script setup lang="ts">
import {h, onMounted, reactive, ref,watch} from "vue";
import useClipboard from "vue-clipboard3";
import { Codemirror } from "vue-codemirror";
import { message } from "ant-design-vue";
import { type Rule } from "ant-design-vue/es/form";
import cronParse from "cron-parser";
import dayjs from "dayjs";
import { javascript } from "@codemirror/lang-javascript";
import { EditorView } from "@codemirror/view";

import { CalcParamMock, CalcParamRd, CalcParamStart, CalcParamStop, CalcRuleCreate, CalcRulePage, CalcRuleUpdate } from "@/api";
import { useRouteJump } from "@/hooks/useRouteJump.ts";
import { useRouterNameStore } from "@/stores/routerPath.ts";
import {useI18n} from "vue-i18n";

const { t,locale } = useI18n();
const { toClipboard } = useClipboard();
const copyText = async (text: any) => {
  try {
    await toClipboard(text);
    message.success(t('message.copySuccess'));
  } catch (e) {
    console.error(e);
  }
};
let rules: Record<string, Rule[]> = {
  name: [{ required: true, message: t('message.pleaseName'), trigger: "blur" }],
  offset: [{ required: true, message: t('message.pleaseForwardTime'), trigger: "blur" }],
  cron: [{ required: true, message: t('message.pleaseExecutionCycle'), trigger: "blur" }],
  script: [{ required: true, message: t('message.pleaseScript'), trigger: "blur" }],
  date: [{ required: true, message: t('message.pleaseTime'), trigger: "change" }],
};
const showSpinning = ref(false);
const jump = useRouteJump();
const routerStore = useRouterNameStore();
const title = ref(t('message.addition'));
const columns = ref([
  {
    title: t('message.uniCode'),
    dataIndex: "ID",
  },
  {
    title:  t('message.name'),
    dataIndex: "name",
  },
  {
    title: t('message.forwardTime'),
    dataIndex: "offset",
  },
  {
    title: t('message.executionCycle'),
    dataIndex: "cron",
  },
  {
    title: t('message.script'),
    dataIndex: "script",
  },
  {
    title: t('message.start'),
    dataIndex: "start",
    customRender: ({ text }:any) => {
      return h("span", text ? t('message.yes') : t('message.no'));
    },
  },
  {
    title: t('message.operation'),
    dataIndex: "operation",
  },
]);
const formRef = ref<HTMLFormElement | null>(null);
const formRefTime = ref<HTMLFormElement | null>(null);
const formRefDate = ref<HTMLFormElement | null>(null);
const dataSource = ref([]);
const modalVisible = ref(false);
const modalTime = ref(false);
const modalMock = ref(false);
const modalDate = ref(false);
const modalTable = ref(false);
const modalResult = ref(false);
const loading = ref(false);
const textResult = ref("");
const dataResult = ref([]);
const scr = "function main(map){return map}";
// 0/10 * * * * ?
const form = reactive({
  id: "",
  name: "",
  cron: "",
  script: "",
  offset: null,
  start: false,
});
const formTime = reactive({ start_time: "", date: "", id: "", end_time: "" });
const formDate = reactive({ start_time: "", date: "", id: "", end_time: "" });
const formState = reactive({ name: "" });
const mockValue = ref("");
const startId = ref("");
const myTheme = EditorView.theme(
  {
    // 输入的字体颜色
    "&": {
      color: "#0052D9",
      backgroundColor: "#FFFFFF",
    },
    ".cm-content": {
      caretColor: "#0052D9",
    },
    // 激活背景色
    ".cm-activeLine": {
      backgroundColor: "#FAFAFA",
    },
    // 激活序列的背景色
    ".cm-activeLineGutter": {
      backgroundColor: "#FAFAFA",
    },
    // 光标的颜色
    "&.cm-focused .cm-cursor": {
      borderLeftColor: "#0052D9",
    },
    // 选中的状态
    "&.cm-focused .cm-selectionBackground, ::selection": {
      backgroundColor: "#0052D9",
      color: "#FFFFFF",
    },
    // 左侧侧边栏的颜色
    ".cm-gutters": {
      backgroundColor: "#FFFFFF",
      color: "#ddd", // 侧边栏文字颜色
      border: "none",
    },
  },
  { dark: true },
);
const extensions = [javascript(), myTheme, EditorView.lineWrapping];


watch(locale, () => {
  columns.value = [
    {
      title: t('message.uniCode'),
      dataIndex: "ID",
    },
    {
      title:  t('message.name'),
      dataIndex: "name",
    },
    {
      title: t('message.forwardTime'),
      dataIndex: "offset",
    },
    {
      title: t('message.executionCycle'),
      dataIndex: "cron",
    },
    {
      title: t('message.script'),
      dataIndex: "script",
    },
    {
      title: t('message.start'),
      dataIndex: "start",
      customRender: ({ text }) => {
        return h("span", text ? t('message.yes') : t('message.no'));
      },
    },
    {
      title: t('message.operation'),
      dataIndex: "operation",
    },
  ]
  rules = {
    name: [{ required: true, message: t('message.pleaseName'), trigger: "blur" }],
    offset: [{ required: true, message: t('message.pleaseForwardTime'), trigger: "blur" }],
    cron: [{ required: true, message: t('message.pleaseExecutionCycle'), trigger: "blur" }],
    script: [{ required: true, message: t('message.pleaseScript'), trigger: "blur" }],
    date: [{ required: true, message: t('message.pleaseTime'), trigger: "change" }],
  }
});
const onAdd = () => {
  modalVisible.value = true;
  title.value = t('message.addition');
};
const pagination = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
  showSizeChanger: true, // 显示每页显示条目数选择器
});
const listPage = async () => {
  const { data } = await CalcRulePage({ name: formState.name, page: pagination.current, page_size: pagination.pageSize });
  dataSource.value = data.data.data;
  pagination.total = data.data?.total || 0;
};
const onCopy = async () => {
  await copyText(scr);
};
const confirm = async (record: any) => {
  modalVisible.value = true;
  title.value = t('message.edit');
  form.id = record.ID;
  form.name = record.name;
  form.cron = record.cron;
  form.script = record.script;
  form.offset = record.offset;
};

const confirmStop = (id: string) => {
  CalcParamStop(id)
    .then(({ data }) => {
      if (data.code === 20000) {
        message.success(data.message);
      } else {
        message.error(data.message);
      }
    }).catch(e=>{
        console.error(e)
      })
    .finally(async () => {
      await listPage();
    });
};
const handleTableChange = async (page: any) => {
  pagination.current = page.current;
  pagination.pageSize = page.pageSize;
  await listPage();
};
const onAddData = () =>
    {
      (formRef.value as HTMLFormElement)
          .validate()
          .then(() => {
            try {
              cronParse.parseExpression(form.cron);
              if (form.cron.split(" ").length !== 6) {
                message.error(t('message.PleaseCorrectExecutionCycle'));
                return;
              }
              if (title.value === t('message.addition')) {
                const data = { ...form };
                delete data.id;
                delete data.start;
                CalcRuleCreate(data).then(async({ data }) => {
                  if (data.code === 20000) {
                    modalVisible.value = false;
                    message.success(t('message.newSuccessfullyAdded'));
                    formRef.value?.resetFields();
                    await listPage();
                  } else {
                    message.error(data.message);
                  }
                }).catch(e=>{
                  console.error(e)
                });
              } else {
                CalcRuleUpdate(form).then(async({ data }) => {
                  if (data.code === 20000) {
                    modalVisible.value = false;
                    message.success(t('message.editSuccessful'));
                    await listPage();
                  } else {
                    message.error(data.message);
                  }
                }).catch(e=>{
                  console.error(e)
                });
              }
            } catch (error) {
              message.error(t('message.pleaseCorrectCycle'));
            }
          }).catch((e: any)=>{
        console.error(e)
      })
    }
;

const onGo = (id: string) => {
  routerStore.setRouterName("/calculate-parameters");
  jump.routeJump({ path: "/calculate-parameters", query: { rule_id: id } });
};
const handleCancel = () => {
  modalVisible.value = false;
};

// 启动
const onStart = (id: string, value: string) => {
  startId.value = id;
  mockValue.value = value;
  if (!mockValue.value) {
    message.warning(t('message.pleaseClickButtonFirst'));
    return;
  }
  modalMock.value = true;
};

const onResult = (id: string) => {
  modalDate.value = true;
  formDate.id = id;
};

const onConfirmStart = () => {
  CalcParamStart(startId.value)
    .then(({ data }) => {
      if (data.code === 20000) {
        message.success(data.message);
      } else {
        message.error(data.message);
      }
    }).catch(e=>{
        console.error(e)
      })
    .finally(async () => {
      modalMock.value = false;
      await listPage();
    });
};

const onMock = (id: string) => {
  formTime.id = id;
  modalTime.value = true;
};

const setMockData = () => {
  (formRefTime.value as HTMLFormElement)
    .validate()
    .then(() => {
      CalcParamMock({ id: formTime.id, start_time: formTime.start_time, end_time: formTime.end_time })
        .then(async ({ data }) => {
          if (data.code === 20000) {
            message.success(JSON.stringify(data.data));
            formRefTime.value?.resetFields();
            await listPage();
          } else {
            message.error(data.message);
          }
        })
          .catch(e=>{
            console.error(e)
          })
        .finally(() => {
          modalTime.value = false;
        });
    })
};

const setTableData = () => {
  (formRefDate.value as HTMLFormElement)
    .validate()
    .then(() => {
      showSpinning.value = true;
      CalcParamRd({ rule_id: formDate.id, start_time: formDate.start_time, end_time: formDate.end_time })
        .then(async ({ data }) => {
          if (data.code === 20000) {
            dataResult.value =
              data.data?.map(({ start_time, end_time, ex_time, result }: any,index: number) => ({
                id: index+1,
                start_time: dayjs(start_time * 1000).format("YYYY-MM-DD HH:mm:ss"),
                end_time: dayjs(end_time * 1000).format("YYYY-MM-DD HH:mm:ss"),
                ex_time: dayjs(ex_time * 1000).format("YYYY-MM-DD HH:mm:ss"),
                result,
              })) || [];
            formRefDate.value?.resetFields();
            await listPage();
          } else {
            message.error(data.message);
          }
        }).catch(e=>{
            console.error(e)
          })
        .finally(() => {
          modalDate.value = false;
          modalTable.value = true;
          showSpinning.value = false;
        });
    })
};
const onTimeChange = (dataString: any) => {
  formTime.start_time = dayjs(dataString[0]).unix();
  formTime.end_time = dayjs(dataString[1]).unix();
};

const onResultTime = (dataString: any) => {
  formDate.start_time = dayjs(dataString[0]).unix();
  formDate.end_time = dayjs(dataString[1]).unix();
};

const onView = (record: any) => {
  textResult.value = JSON.stringify(record.result);
  modalResult.value = true;
};

onMounted(async ()=>{
  await listPage();
})
</script>
<style lang="less" scoped></style>
