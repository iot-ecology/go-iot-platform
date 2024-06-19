import Layout from "@/layout/index.vue";
export default [
  {
    path: "/mqtt-management",
    component: Layout,
    meta: { title: 'message.mqttManagement' },
    redirect: "/mqtt-management/index",
    children: [
      {
        path: "/mqtt-management/index",
        name: "IconPreview",
        meta: { title: "message.mqttManagement" },
        component: async () => await import("@/views/mqtt-management/index.vue"),
      },
    ],
  },
  {
    path: "/signal-configuration",
    component: Layout,
    meta: { title: "message.signalConfig" },
    redirect: "/signal-configuration/index",
    children: [
      {
        path: "/signal-configuration/index",
        name: "CompPreview",
        meta: { title: "message.signalConfig" },
        component: async () => await import("@/views/signal-configuration/index.vue"),
      },
    ],
  },
  {
    path: "/signal",
    component: Layout,
    meta: { title: "message.SignalAlarmConfig" },
    redirect: "/signal/index",
    children: [
      {
        path: "/signal/index",
        name: "Signal",
        meta: { title: "message.SignalAlarmConfig" },
        component: async () => await import("@/views/signal/index.vue"),
      },
    ],
  },
  {
    path: "/visualization",
    component: Layout,
    meta: { title: "message.visualization" },
    redirect: "/visualization/index",
    children: [
      {
        path: "/visualization/index",
        name: "List",
        meta: { title: "列表" },
        component: async () => await import("@/views/visualization/index.vue"),
      },
      {
        path: "/visualization/add",
        name: "Draggable",
        meta: { title: "新增", hidden: true },
        component: async () => await import("@/views/visualization/add.vue"),
      },
    ],
  },
  {
    path: "/calculation-rules",
    component: Layout,
    meta: { title: "message.calculationRules" },
    redirect: "/calculation-rules/index",
    children: [
      {
        path: "/calculation-rules/index",
        name: "calculation-rules",
        meta: { title: "message.calculationRules" },
        component: async () => await import("@/views/calculation-rules/index.vue"),
      },
    ],
  },
  {
    path: "/calculate-parameters",
    component: Layout,
    meta: { title: "message.calculateParameters" },
    redirect: "/calculate-parameters/index",
    children: [
      {
        path: "/calculate-parameters/index",
        name: "calculate-parameters",
        meta: { title: "message.calculateParameters" },
        component: async () => await import("@/views/calculate-parameters/index.vue"),
      },
    ],
  },
  {
    path: "/script-alarm",
    component: Layout,
    meta: { title: "message.scriptAlarm" },
    redirect: "/script-alarm/index",
    children: [
      {
        path: "/script-alarm/index",
        name: "script-alarm",
        meta: { title: "message.scriptAlarm" },
        component: async () => await import("@/views/script-alarm/index.vue"),
      },
    ],
  },
  {
    path: "/script-alarm-parameters",
    component: Layout,
    meta: { title: "message.scriptAlarmParameters" },
    redirect: "/script-alarm-parameters/index",
    children: [
      {
        path: "/script-alarm-parameters/index",
        name: "script-alarm-parameters",
        meta: { title: "message.scriptAlarmParameters" },
        component: async () => await import("@/views/script-alarm-parameters/index.vue"),
      },
    ],
  },
  {
    path: "/node-details",
    component: Layout,
    meta: { title: "message.nodeDetails" },
    redirect: "/node-details/index",
    children: [
      {
        path: "/node-details/index",
        name: "node-details",
        meta: { title: "message.nodeDetails" },
        component: async () => await import("@/views/node-details/index.vue"),
      },
    ],
  },
];
