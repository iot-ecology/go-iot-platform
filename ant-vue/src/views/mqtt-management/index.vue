<template>
  <div class="icon-preview">
    <a-card title="" :bordered="true">
      <a-form layout="inline" :model="formState">
        <a-form-item :label="$t('message.clientID')">
          <a-input v-model:value="formState.client_id" style="width: 300px" :placeholder="$t('message.pleaseEnter')" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="onSearch">{{$t('message.search')}}</a-button>
        </a-form-item>
      </a-form>
      <a-button style="margin: 10px 0" type="primary" @click="modalVisible = true">{{ $t('message.addition') }}</a-button>
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
                <a-button type="primary" size="small" style="margin-right: 10px" @click="save(record.key)">{{$t('message.save')}}</a-button>
                <a-popconfirm :title="$t('message.sureEdit')" :ok-text="$t('message.yes')" :cancel-text="$t('message.no')" @confirm="cancel(record.key)">
                  <a-button type="primary" size="small">{{$t('message.cancel')}}</a-button>
                </a-popconfirm>
              </span>
              <span v-else>
                <a-button type="primary" size="small" v-if="!record.start" style="margin-right: 10px" @click="edit(record.key)">{{$t('message.edit')}}</a-button>
                <a-button type="primary" size="small" v-if="!record.start" @click="onStart(record.ID)">{{$t('message.startUp')}}</a-button>
                <a-popconfirm v-else :title="$t('message.sureStop')" :ok-text="$t('message.yes')" :cancel-text="$t('message.no')" @confirm="confirmStop(record.ID)">
                  <a-button type="primary" size="small">{{$t('message.stop')}}</a-button>
                </a-popconfirm>
                <a-button type="primary" size="small" style="margin-left: 10px" @click="(code = record.script), (scriptId = record.ID), (modalScriptShow = true)">{{$t('message.parsingScripts')}}</a-button>
                <a-button type="primary" size="small" style="margin-left: 10px" @click="onSignal(record.ID)">{{$t('message.signalConfig')}}</a-button>
                <a-button type="primary" size="small" v-if="record.start" style="margin-left: 10px" @click="(modalNewShow = true), (formNews.client_id = record.client_id)">{{$t('message.simulated')}}</a-button>
                <a-popconfirm v-if="!record.start" :title="$t('message.sureDelete')" :ok-text="$t('message.yes')" :cancel-text="$t('message.no')" @confirm="confirm(record.ID)">
                  <a-button type="primary" size="small" danger style="margin-left: 10px;">{{$t('message.delete')}}</a-button>
                </a-popconfirm>
              </span>
            </div>
          </template>
        </template>
      </a-table>
    </a-card>

    <!--新增-->
    <a-modal v-model:open="modalVisible" :okText="$t('message.confirm')" :cancelText="$t('message.cancel')" :destroy-on-close="true" :title="$t('message.addition')" @ok="onAddData()">
      <a-form ref="formRef" :label-col="{ style: { width: '90px' } }" :rules="rules" :model="form">
        <a-form-item :label="$t('message.clientID')" name="client_id">
          <a-input v-model:value="form.client_id" style="width: 350px" />
        </a-form-item>
        <a-form-item :label="$t('message.host')" name="host">
          <a-input-password v-model:value="form.host" style="width: 350px" />
        </a-form-item>
        <a-form-item :label="$t('message.port')" name="port">
          <a-input-number v-model:value="form.port" :precision="0" style="width: 350px" />
        </a-form-item>
        <a-form-item :label="$t('message.account')" name="username">
          <a-input v-model:value="form.username" style="width: 350px" />
        </a-form-item>
        <a-form-item :label="$t('message.password')" name="password">
          <a-input-password v-model:value="form.password" style="width: 350px" />
        </a-form-item>
        <a-form-item :label="$t('message.theme')" name="subtopic">
          <a-input v-model:value="form.subtopic" style="width: 350px" />
        </a-form-item>
      </a-form>
    </a-modal>

    <!--    解析脚本-->
    <a-modal v-model:open="modalScriptShow" title="">
      <a-divider>{{ $t('message.parsingScripts') }}</a-divider>
      <a-tooltip placement="right">
        <template #title>
          <span @click="onCopy">{{ scrValue }}</span>
        </template>
        <p style="cursor: pointer; width: 60px">{{ $t('message.listScript') }}</p>
      </a-tooltip>
      <codemirror v-model="code" placeholder="Code here..." :style="{ height: '150px' }" :autofocus="true" :tab-size="2" :extensions="extensions" />
      <a-divider>{{ $t('message.parameters') }}</a-divider>
      <a-input v-model:value="param" style="margin: 10px 0" placeholder="请输入传递参数"></a-input>
      <a-divider>{{ $t('message.analysis') }}</a-divider>
      <codemirror v-model="resultCode" placeholder="Code here..." :style="{ height: '150px' }" :autofocus="true" :tab-size="2" :extensions="extensions" />
      <template #footer>
        <a-button @click="handleCancel">{{ $t('message.cancel') }}</a-button>
        <a-button :disabled="!(loading && code && param)" type="primary" @click="handleReject">{{ $t('message.validate') }}</a-button>
        <a-button :disabled="loading" type="primary" @click="onConfirmScript()">{{ $t('message.confirm') }}</a-button>
      </template>
    </a-modal>

    <!--    模拟发送消息-->
    <a-modal :okText="$t('message.confirm')" :cancelText="$t('message.cancel')" v-model:open="modalNewShow" :destroy-on-close="true" :title="$t('message.sendMessage')" @ok="setModalNew()">
      <a-form ref="formRefNews" :label-col="{ style: { width: '120px' } }" :rules="rules" :model="formNews">
        <a-form-item :label="$t('message.messageContent')" name="payload">
          <a-input v-model:value="formNews.payload" />
        </a-form-item>
        <a-form-item :label="$t('message.serviceQuality')" name="qos">
          <a-select v-model:value="formNews.qos">
            <a-select-option :value="0">{{ $t('message.once') }}</a-select-option>
            <a-select-option :value="1">{{ $t('message.leastOnce') }}</a-select-option>
            <a-select-option :value="2">{{ $t('message.notRepeat') }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item :label="$t('message.theme')" name="topic">
          <a-input v-model:value="formNews.topic" />
        </a-form-item>
        <a-form-item :label="$t('message.keep')" name="retained">
          <a-radio-group v-model:value="formNews.retained">
            <a-radio :value="true">{{ $t('message.yes') }}</a-radio>
            <a-radio :value="false">{{ $t('message.no') }}</a-radio>
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
import { useI18n } from 'vue-i18n';

import { MqttCheckScript, MqttCreate, MqttDelete, MqttPage, MqttSend, MqttSetScript, MqttStart, MqttStop, MqttUpdate } from "@/api";
import { useRouteJump } from "@/hooks/useRouteJump.ts";
import { useRouterNameStore } from "@/stores/routerPath.ts";
const { toClipboard } = useClipboard();
const { t,locale } = useI18n();
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
    message.success(t('message.copySuccess'));
  } catch (e) {
    console.error(e);
  }
};
const jump = useRouteJump();
let rules: Record<string, Rule[]> = {
  client_id: [{ required: true, message: t('message.pleaseClientID'), trigger: "blur" }],
  host: [
    {
      required: true,
      validator: async (rule, value) => {
        if (value) {
          await Promise.resolve();
        } else {
          await Promise.reject(t('message.pleaseHostIP'));
        }
        const ipPattern = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
        if (ipPattern.test(value)) {
          await Promise.resolve();
        } else {
          await Promise.reject(t('message.pleaseCorrectHostIP'));
        }
      },
      trigger: "blur",
    },
  ],
  port: [{ required: true, message: t('message.pleasePort'), trigger: "blur" }],
  username: [{ required: true, message: t('message.pleaseUsername'), trigger: "blur" }],
  password: [{ required: true, message: t('message.pleasePassword'), trigger: "blur" }],
  subtopic: [{ required: true, message: t('message.pleaseTopic'), trigger: "blur" }],
  payload: [{ required: true, message: t('message.pleaseMessage'), trigger: "blur" }],
  qos: [{ required: true, message: t('message.pleaseService'), trigger: "change" }],
  topic: [{ required: true, message: t('message.pleaseTopic'), trigger: "blur" }],
  retained: [{ required: true, message: t('message.pleaseChooseMessage'), trigger: "change" }],
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
const columns = ref([
  {
    title: t('message.uniCode'),
    dataIndex: "ID",
  },
  {
    title: t('message.clientID'),
    dataIndex: "client_id",
  },
  {
    title: t('message.host'),
    dataIndex: "host",
  },
  {
    title: t('message.port'),
    dataIndex: "port",
  },
  {
    title: t('message.account'),
    dataIndex: "username",
  },
  {
    title: t('message.password'),
    dataIndex: "password",
  },
  {
    title: t('message.theme'),
    dataIndex: "subtopic",
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
]);
const paginations = reactive({
  total: 0,
  current: 1,
  pageSize: 10,
  showSizeChanger: true, // 显示每页显示条目数选择器
});
const list = ref([]);
const scrValue =
  'function main(nc) {\n    var dataRows = [\n        { "Name": "Temperature", "Value": "23" },\n        { "Name": "Humidity", "Value": "30" },\n        { "Name": "A", "Value": nc },\n    ];\n    var result = {\n        "Time":  Math.floor(Date.now() / 1000),\n        "DataRows": dataRows,\n        "DeviceUid": "",\n        "Nc": nc \n    };\n    return [result];\n}';

const code = ref("");
const param = ref("");
const loading = ref(true);
const resultCode = ref("");
watch([code, param], () => {
  if (!loading.value) {
    loading.value = true;
  }
});
watch(locale, () => {
  columns.value = [
    {
      title: t('message.uniCode'),
      dataIndex: "ID",
    },
    {
      title: t('message.clientID'),
      dataIndex: "client_id",
    },
    {
      title: t('message.host'),
      dataIndex: "host",
    },
    {
      title: t('message.port'),
      dataIndex: "port",
    },
    {
      title: t('message.account'),
      dataIndex: "username",
    },
    {
      title: t('message.password'),
      dataIndex: "password",
    },
    {
      title: t('message.theme'),
      dataIndex: "subtopic",
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
    client_id: [{ required: true, message: t('message.pleaseClientID'), trigger: "blur" }],
    host: [
      {
        required: true,
        validator: async (rule, value) => {
          if (value) {
            await Promise.resolve();
          } else {
            await Promise.reject(t('message.pleaseHostIP'));
          }
          const ipPattern = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
          if (ipPattern.test(value)) {
            await Promise.resolve();
          } else {
            await Promise.reject(t('message.pleaseCorrectHostIP'));
          }
        },
        trigger: "blur",
      },
    ],
    port: [{ required: true, message: t('message.pleasePort'), trigger: "blur" }],
    username: [{ required: true, message: t('message.pleaseUsername'), trigger: "blur" }],
    password: [{ required: true, message: t('message.pleasePassword'), trigger: "blur" }],
    subtopic: [{ required: true, message: t('message.pleaseTopic'), trigger: "blur" }],
    payload: [{ required: true, message: t('message.pleaseMessage'), trigger: "blur" }],
    qos: [{ required: true, message: t('message.pleaseService'), trigger: "change" }],
    topic: [{ required: true, message: t('message.pleaseTopic'), trigger: "blur" }],
    retained: [{ required: true, message: t('message.pleaseChooseMessage'), trigger: "change" }],
  }
});
const onCopy = async () => {
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
        "Nc": nc
    };
    return [result];
}`;
  await copyText(text);
};
const edit = (key: string) => {
  editableData[key] = cloneDeep(list.value.filter((item) => key === item.key)[0]);
};
const cancel = (key: string) => {
  // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
  delete editableData[key];
};

const onSignal = (id: string) => {
  routerStore.setRouterName("/signal-configuration");
  jump.routeJump({ path: "/signal-configuration", query: { mqtt_client_id: id } });
};
const save = async (key: string) => {
  Object.assign(list.value.filter((item) => key === item.key)[0], editableData[key]);
  const data = list.value.filter((item) => key === item.key)[0];
  // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
  delete editableData[key];
  const ipPattern = /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
  if (!ipPattern.test(data.host)) {
    message.error(t('message.pleaseCorrectHostIP'));
    return;
  }
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
  MqttDelete(id).then(async ({ data }) => {
    if (data.code === 20000) {
      message.success(data.message);
      await pageList();
    } else {
      // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
      message.success(data.message);
    }
  });
};

const confirmStop = async (id: string) => {
  MqttStop({ id }).then(async ({ data }) => {
    if (data.code === 20000) {
      message.success(data.message);
      await pageList();
    } else {
      message.success(data.message);
    }
  });
};

onMounted(async () => {
  await pageList();
});

const onAddData = () => {
  formRef.value
    .validate()
    .then(() => {
      MqttCreate(form.value).then(async ({ data }) => {
        if (data.code === 20000) {
          message.success(data.message);
          modalVisible.value = false;
          formRef.value?.resetFields();
          await pageList();
        } else {
          // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
          message.error(`${t('message.operationFailed')}:${data.data}`);
        }
      }).catch(e=>{
        console.error(e);
      });
    })
    .catch(e => {
      console.error(e);
    });
};

const onConfirmScript = () => {
  MqttSetScript({ id: scriptId.value, script: code.value }).then(async ({ data }) => {
    if (data.code === 20000) {
      modalScriptShow.value = false;
      message.success(data.message);
      resultCode.value = "";
      param.value = "";
      await pageList();
    } else {
      message.error(data.message || t('message.operationFailed'));
    }
  }).catch(e=>{
    console.error(e)
  });
};
const handleCancel = () => {
  modalScriptShow.value = false;
};
const onStart = (id: string) => {
  MqttStart({ id }).then(async ({ data }) => {
    if (data.code === 20000) {
      message.success(data.message);
      await pageList();
    } else {
      message.success(data.message);
    }
  }).catch(e=>{
    console.error(e)
  });
};
const handleReject = () => {
  if (!code.value) {
    message.error(t('message.pleaseScript'));
    return;
  }
  if (!param.value) {
    message.error(t('message.pleaseParameters'));
    return;
  }
  MqttCheckScript({ param: param.value, script: code.value }).then(({ data }) => {
    if (data.code === 20000) {
      resultCode.value = JSON.stringify(data.data);
      loading.value = false;
      message.success(data.message);
    } else {
      message.error(data.message || t('message.operationFailed'));
    }
  }).catch(e=>{
    console.error(e)
  });
};

const setModalNew = () => {
  formRefNews.value
    .validate()
    .then(() => {
      MqttSend(formNews.value).then(async ({ data }) => {
        if (data.code === 20000) {
          message.success(data.message);
          modalNewShow.value = false;
          formRefNews.value?.resetFields();
          await pageList();
        } else {
          // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
          message.error(`${t('message.operationFailed')}:${data.data}`);
        }
      }).catch(e=> {
        console.error(e)
      });
    })
    .catch(e => {
      console.error(e)
    });
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
