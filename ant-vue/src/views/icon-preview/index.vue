<template>
  <div class="icon-preview">
    <a-card title="" :bordered="true">
      <a-form layout="inline" :model="formState">
        <a-form-item label="客户端ID">
          <a-input v-model:value="formState.client_id" style="width: 300px" placeholder="请输入" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="onSearch">搜索</a-button>
        </a-form-item>
      </a-form>
      <a-button style="margin: 10px 0" type="primary" @click="modalVisible = true">新增</a-button>
      <!--      表格-->
      <a-table :columns="columns" :data-source="list" bordered :pagination="paginations" @change="handleTableChange">
        <template #bodyCell="{ column, text, record }">
          <template v-if="['host', 'port', 'username', 'password', 'subtopic'].includes(column.dataIndex)">
            <div>
              <a-input-number v-if="editableData[record.key] && column.dataIndex === 'port'" v-model:value="editableData[record.key][column.dataIndex]" :precision="0" style="margin: -5px 0" />
              <a-input-password v-else-if="editableData[record.key] && column.dataIndex === 'host'" v-model:value="editableData[record.key][column.dataIndex]" :precision="0" style="margin: -5px 0" />
              <a-input-password
                v-else-if="editableData[record.key] && column.dataIndex === 'password'"
                v-model:value="editableData[record.key][column.dataIndex]"
                :precision="0"
                style="margin: -5px 0"
              />
              <a-input v-else-if="editableData[record.key] && column.dataIndex !== 'port'" v-model:value="editableData[record.key][column.dataIndex]" style="margin: -5px 0" />
              <template v-else>
                <a-input-password v-if="column.dataIndex === 'password' || column.dataIndex === 'host'" :value="text" :bordered="false" />
                <div v-else>{{ text }}</div>
              </template>
            </div>
          </template>
          <template v-else-if="column.dataIndex === 'operation'">
            <div class="editable-row-operations">
              <span v-if="editableData[record.key]">
                <a-typography-link style="margin-right: 10px" @click="save(record.key)">保存</a-typography-link>
                <a-popconfirm title="确定取消编辑吗?" @confirm="cancel(record.key)">
                  <a>取消</a>
                </a-popconfirm>
              </span>
              <span v-else>
                <a v-if="!record.start" style="margin-right: 10px" @click="edit(record.key)">编辑</a>
                <a v-if="!record.start" @click="onStart(record.ID)">启动</a>
                <a-popconfirm v-else title="确认是否停止?" ok-text="是" cancel-text="否" @confirm="confirmStop(record.ID)">
                  <a>停止</a>
                </a-popconfirm>

                <a style="margin-left: 10px" @click="(code = record.script), (scriptId = record.ID), (modalScriptShow = true)">解析脚本</a>
                <a style="margin-left: 10px" @click="onSignal(record.ID)">信号配置</a>
                <a v-if="record.start" style="margin-left: 10px" @click="(modalNewShow = true), (formNews.client_id = record.client_id)">模拟发送</a>
                <a-popconfirm v-if="!record.start" title="确认是否删除?" ok-text="是" cancel-text="否" @confirm="confirm(record.ID)">
                  <a style="margin-left: 10px; color: crimson">删除</a>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>
    </a-card>

    <!--新增-->
    <a-modal v-model:open="modalVisible" :destroy-on-close="true" title="新增" @ok="setModal1Visible()">
      <a-form ref="formRef" :label-col="{ style: { width: '80px' } }" :rules="rules" :model="form" name="nest-messages">
        <a-form-item label="客户端ID" name="client_id">
          <a-input v-model:value="form.client_id" style="width: 350px" />
        </a-form-item>
        <a-form-item label="主机" name="host">
          <a-input-password v-model:value="form.host" style="width: 350px" />
        </a-form-item>
        <a-form-item label="端口" name="port">
          <a-input-number v-model:value="form.port" :precision="0" style="width: 350px" />
        </a-form-item>
        <a-form-item label="账号" name="username">
          <a-input v-model:value="form.username" style="width: 350px" />
        </a-form-item>
        <a-form-item label="密码" name="password">
          <a-input-password v-model:value="form.password" style="width: 350px" />
        </a-form-item>
        <a-form-item label="订阅的主题" name="subtopic">
          <a-input v-model:value="form.subtopic" style="width: 350px" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!--    解析脚本-->
    <a-modal v-model:open="modalScriptShow" title="">
      <a-divider>解析脚本</a-divider>
      <a-tooltip placement="right">
        <template #title>
          <span @click="onCopy">{{ scrValue }}</span>
        </template>
        <p style="cursor: pointer; width: 60px">示列脚本</p>
      </a-tooltip>
      <codemirror v-model="code" placeholder="Code here..." :style="{ height: '150px' }" :autofocus="true" :tab-size="2" :extensions="extensions" />
      <a-divider>传递参数</a-divider>
      <a-input v-model:value="param" style="margin: 10px 0" placeholder="请输入传递参数"></a-input>
      <a-divider>解析结果</a-divider>
      <codemirror v-model="resultCode" placeholder="Code here..." :style="{ height: '150px' }" :autofocus="true" :tab-size="2" :extensions="extensions" />
      <template #footer>
        <a-button @click="handleCancel">取消</a-button>
        <a-button :disabled="!(loading && code && param)" type="primary" @click="handleReject">验证</a-button>
        <a-button :disabled="loading" type="primary" @click="setModalScript()">确定</a-button>
      </template>
    </a-modal>

    <!--    模拟发送消息-->
    <a-modal v-model:open="modalNewShow" :destroy-on-close="true" title="发送消息" @ok="setModalNew()">
      <a-form ref="formRefNews" :label-col="{ style: { width: '80px' } }" :rules="rules" :model="formNews" name="nest-messages">
        <a-form-item label="消息内容" name="payload">
          <a-input v-model:value="formNews.payload" />
        </a-form-item>
        <a-form-item label="服务质量" name="qos">
          <a-select v-model:value="formNews.qos">
            <a-select-option :value="0">最多发送一次，可能会丢失</a-select-option>
            <a-select-option :value="1">至少发送一次，可能会重复</a-select-option>
            <a-select-option :value="2">发送一次，确保消息不会重复</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="主题" name="topic">
          <a-input v-model:value="formNews.topic" />
        </a-form-item>
        <a-form-item label="保留消息" name="retained">
          <a-radio-group v-model:value="formNews.retained">
            <a-radio :value="true">是</a-radio>
            <a-radio :value="false">否</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import type { UnwrapRef } from "vue";
import { h, onMounted, reactive, ref, watch } from "vue";
import useClipboard from "vue-clipboard3";
import { Codemirror } from "vue-codemirror";
import { type FormInstance, message } from "ant-design-vue";
import { type Rule } from "ant-design-vue/es/form";
import { cloneDeep } from "lodash-es";
import { javascript } from "@codemirror/lang-javascript";
import { EditorView } from "@codemirror/view";

import { MqttCheckScript, MqttCreate, MqttDelete, MqttPage, MqttSend, MqttSetScript, MqttStart, MqttStop, MqttUpdate, SignalWaringConfigCreate } from "@/api";
import { useRouteJump } from "@/hooks/useRouteJump.ts";
import { useRouterNameStore } from "@/stores/routerPath.ts";
const { toClipboard } = useClipboard();
interface DataItem {
  client_id: string;
  host: string;
  port: number;
  username: string;
  password: string;
  subtopic: string;
}
const copyText = async (text: string) => {
  try {
    await toClipboard(text);
    message.success("复制成功");
  } catch (e) {
    console.error(e);
  }
};
const jump = useRouteJump();
const rules: Record<string, Rule[]> = {
  client_id: [{ required: true, message: "请输入客户端ID", trigger: "blur" }],
  host: [
    {
      required: true,
      validator: async (rule, value) => {
        if (value) {
          await Promise.resolve();
        } else {
          await Promise.reject("请输入主机IP");
        }
        const ipPattern = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
        if (ipPattern.test(value)) {
          await Promise.resolve();
        } else {
          await Promise.reject("请输入正确的主机IP");
        }
      },
      trigger: "blur",
    },
  ],
  port: [{ required: true, message: "请输入端口", trigger: "blur" }],
  username: [{ required: true, message: "请输入账号", trigger: "blur" }],
  password: [{ required: true, message: "请输入密码", trigger: "blur" }],
  subtopic: [{ required: true, message: "请输入订阅的主题", trigger: "blur" }],
  payload: [{ required: true, message: "请输入消息内容", trigger: "blur" }],
  qos: [{ required: true, message: "请选择服务质量", trigger: "change" }],
  topic: [{ required: true, message: "请输入主题", trigger: "blur" }],
  retained: [{ required: true, message: "请选择是否保留消息", trigger: "change" }],
};
const routerStore = useRouterNameStore();
const scriptId = ref("");
const formRef = ref<FormInstance>();
const formRefNews = ref<FormInstance>();
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
const modalVisible = ref(false);
const modalScriptShow = ref(false);
const modalNewShow = ref(false);
const formState = reactive({ client_id: "" });
const form = ref({ client_id: "", host: "", port: "", username: "", password: "", subtopic: "" });
const formNews = ref({ client_id: "", payload: "", qos: "", retained: false, topic: "" });
const editableData: UnwrapRef<Record<string, DataItem>> = reactive({});
const columns = [
  {
    title: "唯一码",
    dataIndex: "ID",
  },
  {
    title: "客户端ID",
    dataIndex: "client_id",
  },
  {
    title: "主机",
    dataIndex: "host",
  },
  {
    title: "端口",
    dataIndex: "port",
  },
  {
    title: "账号",
    dataIndex: "username",
  },
  {
    title: "密码",
    dataIndex: "password",
  },
  {
    title: "订阅的主题",
    dataIndex: "subtopic",
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
];
const paginations = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
  showSizeChanger: true, // 显示每页显示条目数选择器
});
const list = ref([]);
const scrValue =
  'function main(nc) {\n    var dataRows = [\n        { "Name": "Temperature", "Value": "23" },\n        { "Name": "Humidity", "Value": "30" },\n        { "Name": "A", "Value": nc },\n    ];\n    var result = {\n        "Time":  Math.floor(Date.now() / 1000),\n        "DataRows": dataRows,\n        "DeviceUid": "",\n        "Nc": nc // 确保结果对象中包含nc参数\n    };\n    return [result];\n}';

const code = ref("");
const param = ref("");
const loading = ref(true);
const resultCode = ref("");
watch([code, param], () => {
  if (!loading.value) {
    loading.value = true;
  }
});
const onCopy = () => {
  const text = `function main(nc) {
    var dataRows = [
        { "Name": "Temperature", "Value": "23" },
        { "Name": "Humidity", "Value": "30" },
        { "Name": "A", "Value": nc },
    ];
    var result = {
        "Time":  Math.floor(Date.now() / 1000),
        "DataRows": dataRows,
        "DeviceUid": "${scriptId.value}",
        "Nc": nc // 确保结果对象中包含nc参数
    };
    return [result];
}`;
  copyText(text);
};
const edit = (key: string) => {
  editableData[key] = cloneDeep(list.value.filter((item) => key === item.key)[0]);
};
const cancel = (key: string) => {
  // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
  delete editableData[key];
};

const onSignal = (id: string) => {
  routerStore.setRouterName("/comp-preview");
  jump.routeJump({ path: "/comp-preview", query: { mqtt_client_id: id } });
};
const save = async (key: string) => {
  Object.assign(list.value.filter((item) => key === item.key)[0], editableData[key]);
  const data = list.value.filter((item) => key === item.key)[0];
  // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
  delete editableData[key];
  const ipPattern = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
  if (!ipPattern.test(data.host)) {
    message.error("请输入正确的主机IP");
    await pageList();
    return;
  }
  // eslint-disable-next-line no-debugger
  await MqttUpdate(data);
  await pageList();
};
const pageList = async () => {
  const { data } = await MqttPage({ client_id: formState.client_id, page: paginations.current, page_size: paginations.pageSize });
  paginations.total = data.data?.total || 0;
  list.value = data.data.data?.map((item: any, index: number) => ({
    key: index,
    ID: item.ID,
    client_id: item.client_id,
    host: item.host,
    password: item.password,
    port: item.port,
    script: item.script,
    start: item.start,
    subtopic: item.subtopic,
    username: item.username,
  }));
};
const onSearch = async () => {
  await pageList();
};

const confirm = (id: string) => {
  MqttDelete(id).then(({ data }) => {
    if (data.code === 20000) {
      message.success(data.message);
      pageList();
    } else {
      // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
      message.success(data.message);
    }
  });
};

const confirmStop = async (id: string) => {
  MqttStop({ id }).then(({ data }) => {
    if (data.code === 20000) {
      message.success(data.message);
      pageList();
    } else {
      message.success(data.message);
    }
  });
};

onMounted(async () => {
  await pageList();
});

const setModal1Visible = () => {
  formRef.value
    .validate()
    .then(() => {
      MqttCreate(form.value).then(({ data }) => {
        if (data.code === 20000) {
          message.success(data.message);
          modalVisible.value = false;
          formRef.value?.resetFields();
          pageList();
        } else {
          // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
          message.error(`操作失败:${data.data}`);
        }
      });
    })
    .catch(() => {});
};

const setModalScript = () => {
  MqttSetScript({ id: scriptId.value, script: code.value }).then(({ data }) => {
    if (data.code === 20000) {
      modalScriptShow.value = false;
      message.success(data.message);
      resultCode.value = "";
      param.value = "";
      pageList();
    } else {
      message.error(data.message || "操作失败");
    }
  });
};
const handleCancel = () => {
  modalScriptShow.value = false;
};
const onStart = (id: string) => {
  MqttStart({ id }).then(({ data }) => {
    if (data.code === 20000) {
      message.success(data.message);
      pageList();
    } else {
      message.success(data.message);
    }
  });
};
const handleReject = () => {
  if (!code.value) {
    message.error("请输入脚本");
    return;
  }
  if (!param.value) {
    message.error("请输入传递参数");
    return;
  }
  MqttCheckScript({ param: param.value, script: code.value }).then(({ data }) => {
    if (data.code === 20000) {
      resultCode.value = JSON.stringify(data.data);
      loading.value = false;
      message.success(data.message);
    } else {
      message.error(data.message || "操作失败");
    }
  });
};

const setModalNew = () => {
  formRefNews.value
    .validate()
    .then(() => {
      MqttSend(formNews.value).then(({ data }) => {
        if (data.code === 20000) {
          message.success(data.message);
          modalNewShow.value = false;
          formRefNews.value?.resetFields();
          pageList();
        } else {
          // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
          message.error(`操作失败:${data.data}`);
        }
      });
    })
    .catch(() => {});
};

const handleTableChange = async (pagination: any) => {
  paginations.current = pagination.current;
  paginations.pageSize = pagination.pageSize;
  await pageList();
};
</script>

<style lang="less" scoped>
.icon-preview {
  width: 100%;
  height: 100%;
}
</style>
