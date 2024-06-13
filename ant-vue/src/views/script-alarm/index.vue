<template>
  <div>
    <a-card title="" :bordered="true">
      <a-form layout="inline" :model="formState">
        <a-form-item label="名称">
          <a-input v-model:value="formState.name" style="width: 300px" placeholder="请输入" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="listPage">搜索</a-button>
        </a-form-item>
      </a-form>
      <a-button style="margin: 10px 0" type="primary" @click="onAdd()">新增</a-button>
      <a-table :data-source="dataSource" :columns="columns" bordered :pagination="paginations" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'operation'">
            <div class="editable-row-operations">
              <span>
                <a v-if="!record.start" style="margin-left: 10px" @click="confirm(record)">编辑</a>
                <a style="margin-left: 10px" @click="onWaringHistory(record)">报警历史</a>
                <a style="margin-left: 10px" @click="onGo(record.ID)">参数配置</a>
                <a style="margin-left: 10px" @click="onScript(record.ID)">调试脚本</a>
                <a-popconfirm title="确认是否删除?" ok-text="是" cancel-text="否" @confirm="onDelete(record.ID)">
                  <a style="margin-left: 10px; color: crimson">删除</a>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>

      <a-modal v-model:open="modalVisible" :destroy-on-close="true" :title="title" class="custom-modal">
        <a-form ref="formRef" :label-col="{ style: { width: '100px' } }" :rules="rules" :model="form" name="nest-messages">
          <a-form-item label="名称" name="name">
            <a-input v-model:value="form.name" style="width: 300px" />
          </a-form-item>
          <a-form-item label="脚本" name="script">
            <a-tooltip placement="right">
              <template #title>
                <span @click="onCopy">{{ scr }}</span>
              </template>
              <div style="cursor: pointer; width: 60px; line-height: 32px">示列脚本</div>
            </a-tooltip>
            <codemirror v-model="form.script" placeholder="Code here..." :style="{ width: '300px', height: '150px' }" :tab-size="2" :extensions="extensions" />
          </a-form-item>
        </a-form>
        <template #footer>
          <a-button @click="handleCancel">取消</a-button>
          <a-button :disabled="loading" type="primary" @click="setModal1Visible()">确定</a-button>
        </template>
      </a-modal>

      <!--  调试脚本    -->
      <a-modal v-model:open="modalScript" :destroy-on-close="true" class="custom-modal">
        <a-form ref="formRef" :label-col="{ style: { width: '100px' } }" name="nest-messages">
          <a-form-item label="模拟参数" name="script">
            <codemirror v-model="mockScript" placeholder="Code here..." :style="{ width: '340px', height: '300px' }" :tab-size="2" :extensions="extensions" />
          </a-form-item>
        </a-form>
        <template #footer>
          <a-button @click="modalScript = false">取消</a-button>
          <a-button :disabled="loading" type="primary" @click="onConfirm()">执行</a-button>
        </template>
      </a-modal>

      <!--  时间选择    -->
      <a-modal v-model:open="showDate" title="时间范围" :destroy-on-close="true" class="custom-modal">
        <a-spin tip="加载中..." size="large" :spinning="showSpinning">
          <a-form ref="formRefTime" :rules="rules" :model="formObj" name="nest-messages">
            <a-form-item label="时间范围" name="date">
              <a-range-picker v-model:value="formObj.date" show-time @change="bptjTimeChange" />
            </a-form-item>
          </a-form>
        </a-spin>
        <template #footer>
          <a-button v-if="!showSpinning" @click="showDate = false">取消</a-button>
          <a-button :loading="showSpinning" type="primary" @click="getQueryRow()">确定</a-button>
        </template>
      </a-modal>

      <!--      报警历史-->
      <a-modal v-model:open="showResult" style="width: 600px" :footer="null" :destroy-on-close="true" title="报警历史">
        <a-tabs>
          <a-tab-pane key="1" tab="表格">
            <a-table bordered :pagination="false" :data-source="dataResult" :columns="columnsResult">
              <template #bodyCell="{ column, text, record }">
                <template v-if="column.dataIndex === 'param' || column.dataIndex === 'script'">
                  <div class="editable-row-operations">
                    <span>
                      <a @click="onView(record, column.dataIndex)">查看</a>
                    </span>
                  </div>
                </template>
              </template>
            </a-table>
          </a-tab-pane>
          <a-tab-pane key="2" tab="折线图">
            <div>
              <div v-if="!option.series[0].data?.length" style="height: 223px; text-align: center; font-size: 18px">暂无数据</div>
              <YcECharts v-else :option="option" :height="300" />
            </div>
          </a-tab-pane>
        </a-tabs>
      </a-modal>

      <!--      脚本-->
      <a-modal v-model:open="showDetail" :footer="null" title="详情">
        <codemirror v-model="scriptDetail" placeholder="Code here..." :style="{ width: '300px', height: '150px' }" :tab-size="2" :extensions="extensions" />
      </a-modal>
    </a-card>
  </div>
</template>
<script setup lang="ts">
import { h, reactive, ref } from "vue";
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
const { toClipboard } = useClipboard();
const copyText = async (text: any) => {
  try {
    await toClipboard(text);
    message.success("复制成功");
  } catch (e) {
    console.error(e);
  }
};

const jump = useRouteJump();
const rules: Record<string, Rule[]> = {
  name: [{ required: true, message: "请输入名称", trigger: "blur" }],
  script: [{ required: true, message: "请输入脚本", trigger: "blur" }],
  date: [{ required: true, message: "请选择时间", trigger: "change" }],
};
const showSpinning = ref(false);
const routerStore = useRouterNameStore();
const title = ref("新增");
const columns = ref([
  {
    title: "名称",
    dataIndex: "name",
  },
  {
    title: "脚本",
    dataIndex: "script",
  },
  {
    title: "操作",
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
const showDate = ref(false);
const showDetail = ref(false);
const scriptDetail = ref("");
const option = ref({});

const formObj = reactive({ ID: "", up_time_end: "", up_time_start: "", date: "" });
const form = reactive({
  id: "",
  name: "",
  script: "",
});
const formState = reactive({ name: "" });
const showResult = ref(false);
const dataResult = ref([]);
const columnsResult = ref([
  {
    title: "上报时间",
    dataIndex: "up_time",
  },
  {
    title: "是否符合规则",
    dataIndex: "value",
    customRender: ({ text }) => {
      return h("span", text ? "是" : "否");
    },
  },
  {
    title: "处理时间",
    dataIndex: "insert_time",
  },
  {
    title: "参数",
    dataIndex: "param",
  },
  {
    title: "脚本",
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
const onAdd = () => {
  modalVisible.value = true;
  title.value = "新增";
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
listPage();
const onCopy = () => {
  copyText(scr);
};
const confirm = async (record: any) => {
  modalVisible.value = true;
  title.value = "编辑";
  form.id = record.ID;
  form.name = record.name;
  form.script = record.script;
};

const handleTableChange = async (pagination: any) => {
  paginations.current = pagination.current;
  paginations.pageSize = pagination.pageSize;
  await listPage();
};
const setModal1Visible = () => {
  formRef.value
    .validate()
    .then(() => {
      if (title.value === "新增") {
        const data = { ...form };
        delete data.id;
        SignalDelayWaringCreate(data).then(({ data }) => {
          if (data.code === 20000) {
            modalVisible.value = false;
            message.success("新增成功");
            formRef.value?.resetFields();
            listPage();
          } else {
            message.error(data.message);
          }
        });
      } else {
        SignalDelayWaringUpdate(form).then(({ data }) => {
          if (data.code === 20000) {
            modalVisible.value = false;
            message.success("编辑成功");
            listPage();
          } else {
            message.error(data.message);
          }
        });
      }
    })
    .catch(() => {});
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
    .catch(() => {})
    .finally(() => {
      modalScript.value = true;
    });
};
// 执行
const onConfirm = () => {
  SignalDelayWaringMock(mockId.value)
    .then(({ data }) => {
      message.success("执行结果:" + data.data);
      modalScript.value = false;
    })
    .catch(() => {})
    .finally(() => {});
};

// 删除
const onDelete = async (id: string) => {
  SignalDelayWaringDelete(id).then(({ data }) => {
    if (data.code === 20000) {
      message.success(data.message);
      listPage();
    } else {
      // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
      message.success(data.message);
    }
  });
};

// 报警历史
const onWaringHistory = (record: any) => {
  formObj.ID = record.ID;
  showDate.value = true;
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
        .finally(() => {
          showResult.value = true;
          showDate.value = false;
          showSpinning.value = false;
          formObj.date = "";
        });
    })
    .catch(() => {});
};
const bptjTimeChange = (date: any, dataString: any) => {
  formObj.up_time_start = dayjs(dataString[0]).unix();
  formObj.up_time_end = dayjs(dataString[1]).unix();
};

// 查看报警历史内容
const onView = (record: any, str: string) => {
  showDetail.value = true;
  scriptDetail.value = str === "param" ? JSON.stringify(record[str]) : record[str];
};
</script>
<style lang="less" scoped></style>
