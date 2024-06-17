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
      <a-table :data-source="dataSource" :columns="columns" bordered :pagination="paginations" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'operation'">
            <div class="editable-row-operations">
              <span>
                <a v-if="!record.start" style="margin-left: 10px" @click="confirm(record)">{{$t('message.edit')}}</a>
                <a style="margin-left: 10px" @click="onWaringHistory(record)">{{ $t('message.alarmHistory') }}</a>
                <a style="margin-left: 10px" @click="onGo(record.ID)">{{ $t('message.parameterConfiguration') }}</a>
                <a style="margin-left: 10px" @click="onScript(record.ID)">{{ $t('message.debuggingScripts') }}</a>
                <a-popconfirm :title="$t('message.sureDelete')" :okText="$t('message.yes')" :cancelText="$t('message.no')" @confirm="onDelete(record.ID)">
                  <a style="margin-left: 10px; color: crimson">{{$t('message.delete')}}</a>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>

      <a-modal v-model:open="modalVisible" :destroy-on-close="true" :title="title" class="custom-modal">
        <a-form ref="formRef" :label-col="{ style: { width: '100px' } }" :rules="rules" :model="form">
          <a-form-item :label="$t('message.name')" name="name">
            <a-input v-model:value="form.name" style="width: 300px" />
          </a-form-item>
          <a-form-item :label="$t('message.script')" name="script">
            <a-tooltip placement="right">
              <template #title>
                <span @click="onCopy">{{ scr }}</span>
              </template>
              <div style="cursor: pointer; width: 60px; line-height: 32px">{{ $t('message.listScript') }}</div>
            </a-tooltip>
            <codemirror v-model="form.script" placeholder="Code here..." :style="{ width: '300px', height: '150px' }" :tab-size="2" :extensions="extensions" />
          </a-form-item>
        </a-form>
        <template #footer>
          <a-button @click="handleCancel">{{$t('message.cancel')}}</a-button>
          <a-button :disabled="loading" type="primary" @click="onAddUpdateData()">{{$t('message.confirm')}}</a-button>
        </template>
      </a-modal>

      <!--  调试脚本    -->
      <a-modal v-model:open="modalScript" :destroy-on-close="true" class="custom-modal">
        <a-form ref="formRef" :label-col="{ style: { width: '100px' } }">
          <a-form-item :label="$t('message.simulatedParameters')" name="script">
            <codemirror v-model="mockScript" placeholder="Code here..." :style="{ width: '340px', height: '300px' }" :tab-size="2" :extensions="extensions" />
          </a-form-item>
        </a-form>
        <template #footer>
          <a-button @click="modalScript = false">{{$t('message.cancel')}}</a-button>
          <a-button :disabled="loading" type="primary" @click="onConfirm()">{{$t('message.implement')}}</a-button>
        </template>
      </a-modal>

      <!--  时间选择    -->
      <a-modal v-model:open="modalDate" :title="$t('message.timeframe')" :destroy-on-close="true" class="custom-modal">
        <a-spin :tip="$t('message.loading')" size="large" :spinning="showSpinning">
          <a-form ref="formRefTime" :rules="rules" :model="formObj">
            <a-form-item :label="$t('message.timeframe')" name="date">
              <a-range-picker v-model:value="formObj.date" show-time @change="bptjTimeChange" />
            </a-form-item>
          </a-form>
        </a-spin>
        <template #footer>
          <a-button v-if="!showSpinning" @click="modalDate = false">{{$t('message.cancel')}}</a-button>
          <a-button :loading="showSpinning" type="primary" @click="getQueryRow()">{{$t('message.confirm')}}</a-button>
        </template>
      </a-modal>

      <!--      报警历史-->
      <a-modal v-model:open="modalResult" style="width: 600px" :footer="null" :destroy-on-close="true" :title="$t('message.alarmHistory')">
        <a-tabs>
          <a-tab-pane key="1" :tab="$t('message.table')">
            <a-table bordered :pagination="false" :data-source="dataResult" :columns="columnsResult">
              <template #bodyCell="{ column, text, record }">
                <template v-if="column.dataIndex === 'param' || column.dataIndex === 'script'">
                  <div class="editable-row-operations">
                    <span>
                      <a @click="onView(record, column.dataIndex)">{{ $t('message.check') }}</a>
                    </span>
                  </div>
                </template>
              </template>
            </a-table>
          </a-tab-pane>
          <a-tab-pane key="2" :tab="$t('message.lineChart')">
            <div>
              <div v-if="!option.series[0].data?.length" style="height: 223px; text-align: center; font-size: 18px">{{ $t('message.noData') }}</div>
              <YcECharts v-else :option="option" :height="300" />
            </div>
          </a-tab-pane>
        </a-tabs>
      </a-modal>

      <!--      脚本-->
      <a-modal v-model:open="modalDetail" :footer="null" :title="$t('message.details')">
        <codemirror v-model="scriptDetail" placeholder="Code here..." :style="{ width: '300px', height: '150px' }" :tab-size="2" :extensions="extensions" />
      </a-modal>
    </a-card>
  </div>
</template>
<script setup lang="ts">
import {h, onMounted, reactive, ref,watch} from "vue";
import useClipboard from "vue-clipboard3";
import { Codemirror } from "vue-codemirror";
import { type FormInstance, message } from "ant-design-vue";
import { type Rule } from "ant-design-vue/es/form";
import dayjs from "dayjs";
import { javascript } from "@codemirror/lang-javascript";
import { EditorView } from "@codemirror/view";

import { SignalDelayWaringCreate, SignalDelayWaringDelete, SignalDelayWaringGenParam, SignalDelayWaringMock, SignalDelayWaringPage, SignalDelayWaringQueryRow, SignalDelayWaringUpdate } from "@/api";
import { YcECharts } from "@/components";
import { useRouteJump } from "@/hooks/useRouteJump.ts";
import { useRouterNameStore } from "@/stores/routerPath.ts";
import {useI18n} from "vue-i18n";
const { toClipboard } = useClipboard();

const { t,locale } = useI18n();
const copyText = async (text: any) => {
  try {
    await toClipboard(text);
    message.success(t('message.copySuccess'));
  } catch (e) {
    console.error(e);
  }
};

const jump = useRouteJump();
let rules: Record<string, Rule[]> = {
  name: [{ required: true, message: t('message.copySuccess'), trigger: "blur" }],
  script: [{ required: true, message: t('message.pleaseScript'), trigger: "blur" }],
  date: [{ required: true, message: t('message.pleaseTime'), trigger: "change" }],
};
const showSpinning = ref(false);
const routerStore = useRouterNameStore();
const title = ref(t('message.addition'));
const columns = ref([
  {
    title: t('message.name'),
    dataIndex: "name",
  },
  {
    title: t('message.script'),
    dataIndex: "script",
  },
  {
    title: t('message.operation'),
    dataIndex: "operation",
  },
]);
const formRef = ref<FormInstance>();
const formRefTime = ref<FormInstance>();
const dataSource = ref([]);
const modalVisible = ref(false);
const modalScript = ref(false);
const loading = ref(false);
const scr = "function main(map){return false}";
const mockScript = ref("");
const mockId = ref();
const modalDate = ref(false);
const modalDetail = ref(false);
const scriptDetail = ref("");
const option = ref({});

const formObj = reactive({ ID: "", up_time_end: "", up_time_start: "", date: "" });
const form = reactive({
  id: "",
  name: "",
  script: "",
});
const formState = reactive({ name: "" });
const modalResult = ref(false);
const dataResult = ref([]);
const columnsResult = ref([
  {
    title: t('message.reportingTime'),
    dataIndex: "up_time",
  },
  {
    title: t('message.complyRules'),
    dataIndex: "value",
    customRender: ({ text }) => {
      return h("span", text ? t('message.yes') : t('message.no'));
    },
  },
  {
    title: t('message.processingTime'),
    dataIndex: "insert_time",
  },
  {
    title: t('message.parameter'),
    dataIndex: "param",
  },
  {
    title: t('message.script'),
    dataIndex: "script",
  },
]);

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
      title: t('message.name'),
      dataIndex: "name",
    },
    {
      title: t('message.script'),
      dataIndex: "script",
    },
    {
      title: t('message.operation'),
      dataIndex: "operation",
    }]
  rules = {
    name: [{ required: true, message: t('message.copySuccess'), trigger: "blur" }],
    script: [{ required: true, message: t('message.pleaseScript'), trigger: "blur" }],
    date: [{ required: true, message: t('message.pleaseTime'), trigger: "change" }]
  },
  columnsResult.value = [
    {
      title: t('message.reportingTime'),
      dataIndex: "up_time",
    },
    {
      title: t('message.complyRules'),
      dataIndex: "value",
      customRender: ({ text }) => {
        return h("span", text ? t('message.yes') : t('message.no'));
      },
    },
    {
      title: t('message.processingTime'),
      dataIndex: "insert_time",
    },
    {
      title: t('message.parameter'),
      dataIndex: "param",
    },
    {
      title: t('message.script'),
      dataIndex: "script",
    },
  ]
});

const onAdd = () => {
  modalVisible.value = true;
  title.value = t('message.addition');
};
const paginations = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
  showSizeChanger: true, // 显示每页显示条目数选择器
});
const listPage = async () => {
  const { data } = await SignalDelayWaringPage({ name: formState.name, page: paginations.current, page_size: paginations.pageSize });
  paginations.total = data.data?.total || 0;
  dataSource.value = data.data.data;
};
const onCopy = () => {
  copyText(scr);
};
const confirm = async (record: any) => {
  modalVisible.value = true;
  title.value = t('message.edit');
  form.id = record.ID;
  form.name = record.name;
  form.script = record.script;
};

const handleTableChange = async (pagination: any) => {
  paginations.current = pagination.current;
  paginations.pageSize = pagination.pageSize;
  await listPage();
};
const onAddUpdateData = () => {
  formRef.value
    .validate()
    .then(() => {
      if (title.value === t('message.addition')) {
        const data = { ...form };
        delete data.id;
        SignalDelayWaringCreate(data).then(async ({ data }) => {
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
        SignalDelayWaringUpdate(form).then(async ({ data }) => {
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
    })
    .catch(e => {
      console.error(e)
    });
};

const onGo = (id: string) => {
  routerStore.setRouterName("/script-alarm-parameters");
  jump.routeJump({ path: "/script-alarm-parameters", query: { signal_delay_waring_id: id } });
};
const handleCancel = () => {
  modalVisible.value = false;
};

// 调试脚本
const onScript = (id: number) => {
  mockId.value = id;
  SignalDelayWaringGenParam(id)
    .then(({ data }) => {
      mockScript.value = JSON.stringify(data.data);
    })
    .catch(e => {
      console.error(e)
    })
    .finally(() => {
      modalScript.value = true;
    });
};
// 执行
const onConfirm = () => {
  SignalDelayWaringMock(mockId.value)
    .then(({ data }) => {
      message.success(`${t('message.resultsOfEnforcement')}:` + data.data);
      modalScript.value = false;
    })
    .catch(e => {
      console.error(e)
    })
    .finally(() => {});
};

// 删除
const onDelete = async (id: string) => {
  SignalDelayWaringDelete(id).then(async ({ data }) => {
    if (data.code === 20000) {
      message.success(data.message);
      await listPage();
    } else {
      // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
      message.success(data.message);
    }
  }).catch(e=>{
    console.error(e)
  });
};

// 报警历史
const onWaringHistory = (record: any) => {
  formObj.ID = record.ID;
  modalDate.value = true;
};

const getQueryRow = () => {
  formRefTime.value
    .validate()
    .then(() => {
      showSpinning.value = true;
      SignalDelayWaringQueryRow(formObj)
        .then(({ data }) => {
          dataResult.value = data.data?.map(({ insert_time, up_time, value, param, script }) => ({
            insert_time: dayjs(insert_time * 1000).format("YYYY-MM-DD HH:mm:ss"),
            up_time: dayjs(up_time * 1000).format("YYYY-MM-DD HH:mm:ss"),
            value,
            param,
            script,
          }));
          option.value = {
            tooltip: {
              trigger: "axis",
              formatter: (params) => {
                return `${params[0].name}: ${params[0].value === 1 ? "是" : "否"}`;
              },
            },
            xAxis: {
              type: "category",
              show: true,
              data: dataResult.value?.map((it) => it.insert_time),
            },
            yAxis: {
              show: false,
              type: "value",
            },
            series: [
              {
                data: dataResult.value?.map((it) => (it.value ? 1 : 0)),
                type: "line",
              },
            ],
            dataZoom: [
              {
                type: "inside",
              },
              {
                type: "slider",
                top: 250,
              },
            ],
          };
        })
        .catch(e=>{
            console.error(e)
          })
        .finally(() => {
          modalResult.value = true;
          modalDate.value = false;
          showSpinning.value = false;
          formObj.date = "";
        });
    })
    .catch(e => {
      console.error(e)
    });
};
const bptjTimeChange = (date: any, dataString: any) => {
  formObj.up_time_start = dayjs(dataString[0]).unix();
  formObj.up_time_end = dayjs(dataString[1]).unix();
};

// 查看报警历史内容
const onView = (record: any, str: string) => {
  modalDetail.value = true;
  scriptDetail.value = str === "param" ? JSON.stringify(record[str]) : record[str];
};

onMounted(async ()=>{
  await listPage();
})

</script>
<style lang="less" scoped></style>
