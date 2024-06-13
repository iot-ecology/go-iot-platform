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
                <a v-if="!record.start" style="margin-left: 10px" @click="onGo(record.ID)">参数配置</a>
                <a v-if="!record.start" style="margin-left: 10px" @click="onMock(record.ID)">模拟执行</a>
                <a v-if="!record.start" style="margin-left: 10px" @click="onStart(record.ID, record.mock_value)">启动</a>

                <a-popconfirm v-else title="确认是否停止?" ok-text="是" cancel-text="否" @confirm="confirmStop(record.ID)">
                  <a style="margin-left: 10px">停止</a>
                </a-popconfirm>
                <a style="margin-left: 10px" @click="onResult(record.ID)">结果查看</a>
              </span>
            </div>
          </template>
        </template>
      </a-table>

      <a-modal v-model:open="modalVisible" :destroy-on-close="true" :title="title" class="custom-modal">
        <a-form ref="formRef" :label-col="{ style: { width: '100px' } }" :rules="rules" :model="form" name="nest-messages">
          <a-form-item label="名称" name="name">
            <a-input v-model:value="form.name" style="width: 300px" :disabled="form.start" />
          </a-form-item>
          <a-form-item label="前移时间（s）" name="offset">
            <a-input-number v-model:value="form.offset" :precision="0" style="width: 300px" :min="11" :disabled="form.start" />
          </a-form-item>
          <a-form-item label="执行周期" name="cron">
            <a-input v-model:value="form.cron" style="width: 300px" :disabled="form.start" />
          </a-form-item>
          <a-form-item label="脚本" name="script">
            <a-tooltip placement="right">
              <template #title>
                <span @click="onCopy">{{ scr }}</span>
              </template>
              <div style="cursor: pointer; width: 60px; line-height: 32px">示列脚本</div>
            </a-tooltip>
            <codemirror v-model="form.script" :disabled="form.start" placeholder="Code here..." :style="{ width: '300px', height: '150px' }" :tab-size="2" :extensions="extensions" />
          </a-form-item>
        </a-form>
        <template #footer>
          <a-button @click="handleCancel">取消</a-button>
          <!--          <a-button :disabled="!loading" type="primary" @click="handleReject">验证</a-button>-->
          <a-button :disabled="loading" type="primary" @click="setModal1Visible()">确定</a-button>
        </template>
      </a-modal>

      <a-modal v-model:open="modalTime" title="时间范围" class="custom-modal">
        <a-form ref="formRefTime" :rules="rules" :model="formTime" name="nest-messages">
          <a-form-item label="时间范围" name="date">
            <a-range-picker v-model:value="formTime.date" show-time @change="bptjTimeChange" />
          </a-form-item>
        </a-form>
        <template #footer>
          <a-button @click="modalTime = false">取消</a-button>
          <a-button :disabled="loading" type="primary" @click="setModalMock()">确定</a-button>
        </template>
      </a-modal>

      <a-modal v-model:open="modalDate" title="时间范围" class="custom-modal">
        <a-spin tip="加载中..." size="large" :spinning="showSpinning">
          <a-form ref="formRefDate" :rules="rules" :model="formDate" name="nest-messages">
            <a-form-item label="时间范围" name="date">
              <a-range-picker v-model:value="formDate.date" show-time @change="bptjTime" />
            </a-form-item>
          </a-form>
        </a-spin>
        <template #footer>
          <a-button v-if="!showSpinning" @click="modalDate = false">取消</a-button>
          <a-button :loading="showSpinning" type="primary" @click="setModalData()">确定</a-button>
        </template>
      </a-modal>

      <a-modal v-model:open="modalMock" title="确认模拟结果是否符合预期" class="custom-modal">
        <a-form :model="formTime" name="nest-messages">
          <a-form-item label="数据">
            <div>{{ mockValue }}</div>
          </a-form-item>
        </a-form>
        <template #footer>
          <a-button @click="modalMock = false">取消</a-button>
          <a-button :disabled="loading" type="primary" @click="onConfimStart()">确定</a-button>
        </template>
      </a-modal>

      <a-modal v-model:open="modalTable" :footer="null" style="width: 600px; height: 650px; overflow: scroll" title="" class="custom-modal">
        <div style="margin-top: 30px; border: 1px solid #f0f0f0">
          <div style="display: flex; border-bottom: 1px solid #f0f0f0; text-align: center; font-size: 16px">
            <div style="width: 27%; border-right: 1px solid #f0f0f0; padding: 8px 0">执行时间</div>
            <div style="width: 27%; border-right: 1px solid #f0f0f0; padding: 8px 0">查询开始时间</div>
            <div style="width: 27%; border-right: 1px solid #f0f0f0; padding: 8px 0">查询结束时间</div>
            <div style="width: 19%; padding: 8px 0">结果</div>
          </div>
          <div v-if="!dataResult?.length" style="text-align: center; font-size: 18px; height: 300px">暂无数据</div>
          <RecycleScroller v-else v-slot="{ item }" style="height: 480px" class="scroller" :items="dataResult" :item-size="40" key-field="start_time">
            <div style="display: flex; text-align: center; border-bottom: 1px solid #f0f0f0">
              <div style="width: 27%; border-right: 1px solid #f0f0f0; padding: 8px 0">{{ item.ex_time }}</div>
              <div style="width: 27%; border-right: 1px solid #f0f0f0; padding: 8px 0">{{ item.start_time }}</div>
              <div style="width: 27%; border-right: 1px solid #f0f0f0; padding: 8px 0">{{ item.end_time }}</div>
              <div style="width: 19%; padding: 8px 0; cursor: pointer; color: rgb(24, 144, 255)" @click="onView(item)">查看</div>
            </div>
          </RecycleScroller>
        </div>
      </a-modal>

      <a-modal v-model:open="showResult" :footer="null" title="" class="custom-modal">
        <div>{{ textResult }}</div>
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
import cronParse from "cron-parser";
import dayjs from "dayjs";
import { javascript } from "@codemirror/lang-javascript";
import { EditorView } from "@codemirror/view";

import { CalcParamMock, CalcParamRd, CalcParamStart, CalcParamStop, CalcRuleCreate, CalcRulePage, CalcRuleUpdate, SignalCreate } from "@/api";
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
const rules: Record<string, Rule[]> = {
  name: [{ required: true, message: "请输入名称", trigger: "blur" }],
  offset: [{ required: true, message: "请输入前移时间", trigger: "blur" }],
  cron: [{ required: true, message: "请输入执行周期", trigger: "blur" }],
  script: [{ required: true, message: "请输入脚本", trigger: "blur" }],
  date: [{ required: true, message: "请选择时间", trigger: "change" }],
};
const showSpinning = ref(false);
const jump = useRouteJump();
const routerStore = useRouterNameStore();
const title = ref("新增");
const columns = ref([
  {
    title: "ID",
    dataIndex: "ID",
  },
  {
    title: "名称",
    dataIndex: "name",
  },
  {
    title: "前移时间(s)",
    dataIndex: "offset",
  },
  {
    title: "执行周期",
    dataIndex: "cron",
  },
  {
    title: "脚本",
    dataIndex: "script",
  },
  {
    title: "是否启动",
    dataIndex: "start",
    customRender: ({ text }) => {
      return h("span", text ? "是" : "否");
    },
  },
  {
    title: "操作",
    dataIndex: "operation",
  },
]);
const formRef = ref<FormInstance>();
const formRefTime = ref<FormInstance>();
const formRefDate = ref<FormInstance>();
const dataSource = ref([]);
const modalVisible = ref(false);
const modalTime = ref(false);
const modalMock = ref(false);
const modalDate = ref(false);
const modalTable = ref(false);
const showResult = ref(false);
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
  const { data } = await CalcRulePage({ name: formState.name, page: paginations.current, page_size: paginations.pageSize });
  dataSource.value = data.data.data;
  paginations.total = data.data?.total || 0;
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
    })
    .finally(() => {
      listPage();
    });
};
const handleTableChange = async (pagination) => {
  paginations.current = pagination.current;
  paginations.pageSize = pagination.pageSize;
  await listPage();
};
const setModal1Visible = () => {
  formRef.value
    .validate()
    .then(() => {
      try {
        const interval = cronParse.parseExpression(form.cron);
        if (form.cron.split(" ").length !== 6) {
          message.error("请输入正确执行周期");
          return;
        }
        if (title.value === "新增") {
          const data = { ...form };
          delete data.id;
          delete data.start;
          CalcRuleCreate(data).then(({ data }) => {
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
          CalcRuleUpdate(form).then(({ data }) => {
            if (data.code === 20000) {
              modalVisible.value = false;
              message.success("编辑成功");
              listPage();
            } else {
              message.error(data.message);
            }
          });
        }
      } catch (error) {
        message.error("请输入正确周期");
      }
    })
    .catch(() => {});
};

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
    message.warning("请先点击模拟执行按钮");
    return;
  }
  modalMock.value = true;
};

const onResult = (id: string) => {
  modalDate.value = true;
  formDate.id = id;
};

const onConfimStart = () => {
  CalcParamStart(startId.value)
    .then(({ data }) => {
      if (data.code === 20000) {
        message.success(data.message);
      } else {
        message.error(data.message);
      }
    })
    .finally(() => {
      modalMock.value = false;
      listPage();
    });
};

const onMock = (id: string) => {
  formTime.id = id;
  modalTime.value = true;
};

const setModalMock = () => {
  formRefTime.value
    .validate()
    .then(() => {
      CalcParamMock({ id: formTime.id, start_time: formTime.start_time, end_time: formTime.end_time })
        .then(({ data }) => {
          if (data.code === 20000) {
            message.success(JSON.stringify(data.data));
            formRefTime.value?.resetFields();
            listPage();
          } else {
            message.error(data.message);
          }
        })
        .finally(() => {
          modalTime.value = false;
        });
    })
    .catch(() => {});
};

const setModalData = () => {
  formRefDate.value
    .validate()
    .then(() => {
      showSpinning.value = true;
      CalcParamRd({ rule_id: formDate.id, start_time: formDate.start_time, end_time: formDate.end_time })
        .then(({ data }) => {
          if (data.code === 20000) {
            // message.success(JSON.stringify(data.data));
            dataResult.value =
              data.data?.map(({ start_time, end_time, ex_time, result }) => ({
                start_time: dayjs(start_time * 1000).format("YYYY-MM-DD HH:mm:ss"),
                end_time: dayjs(end_time * 1000).format("YYYY-MM-DD HH:mm:ss"),
                ex_time: dayjs(ex_time * 1000).format("YYYY-MM-DD HH:mm:ss"),
                result,
              })) || [];
            formRefDate.value?.resetFields();
            listPage();
          } else {
            message.error(data.message);
          }
        })
        .finally(() => {
          modalDate.value = false;
          modalTable.value = true;
          showSpinning.value = false;
        });
    })
    .catch(() => {});
};
const bptjTimeChange = (date: any, dataString: any) => {
  formTime.start_time = dayjs(dataString[0]).unix();
  formTime.end_time = dayjs(dataString[1]).unix();
};

const bptjTime = (date: any, dataString: any) => {
  formDate.start_time = dayjs(dataString[0]).unix();
  formDate.end_time = dayjs(dataString[1]).unix();
};

const onView = (record: any) => {
  textResult.value = JSON.stringify(record.result);
  showResult.value = true;
};
</script>
<style lang="less" scoped></style>
