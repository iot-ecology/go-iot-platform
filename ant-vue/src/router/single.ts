import Layout from "@/layout/index.vue";
export default [
  {
    path: "/icon-preview",
    component: Layout,
    meta: { title: "MQTT客户端管理" },
    redirect: "/icon-preview/index",
    children: [
      {
        path: "/icon-preview/index",
        name: "IconPreview",
        meta: { title: "MQTT客户端" },
        component: async () => await import("@/views/icon-preview/index.vue"),
      },
    ],
  },
  {
    path: "/comp-preview",
    component: Layout,
    meta: { title: "信号配置" },
    redirect: "/comp-preview/index",
    children: [
      {
        path: "/comp-preview/index",
        name: "CompPreview",
        meta: { title: "信号配置" },
        component: async () => await import("@/views/comp-preview/index.vue"),
      },
    ],
  },
  {
    path: "/signal",
    component: Layout,
    meta: { title: "信号报警配置" },
    redirect: "/signal/index",
    children: [
      {
        path: "/signal/index",
        name: "Signal",
        meta: { title: "信号报警配置" },
        component: async () => await import("@/views/signal/index.vue"),
      },
    ],
  },
  {
    path: "/draggable",
    component: Layout,
    meta: { title: "可视化" },
    redirect: "/draggable/index",
    children: [
      {
        path: "/draggable/index",
        name: "List",
        meta: { title: "列表" },
        component: async () => await import("@/views/draggable/index.vue"),
      },
      {
        path: "/draggable/add",
        name: "Draggable",
        meta: { title: "新增", hidden: true },
        component: async () => await import("@/views/draggable/add.vue"),
      },
    ],
  },
  {
    path: "/calculation-rules",
    component: Layout,
    meta: { title: "计算规则" },
    redirect: "/calculation-rules/index",
    children: [
      {
        path: "/calculation-rules/index",
        name: "calculation-rules",
        meta: { title: "计算规则" },
        component: async () => await import("@/views/calculation-rules/index.vue"),
      },
    ],
  },
  {
    path: "/calculate-parameters",
    component: Layout,
    meta: { title: "计算参数" },
    redirect: "/calculate-parameters/index",
    children: [
      {
        path: "/calculate-parameters/index",
        name: "calculate-parameters",
        meta: { title: "计算参数" },
        component: async () => await import("@/views/calculate-parameters/index.vue"),
      },
    ],
  },
  {
    path: "/script-alarm",
    component: Layout,
    meta: { title: "脚本报警" },
    redirect: "/script-alarm/index",
    children: [
      {
        path: "/script-alarm/index",
        name: "script-alarm",
        meta: { title: "脚本报警" },
        component: async () => await import("@/views/script-alarm/index.vue"),
      },
    ],
  },
  {
    path: "/script-alarm-parameters",
    component: Layout,
    meta: { title: "脚本报警参数" },
    redirect: "/script-alarm-parameters/index",
    children: [
      {
        path: "/script-alarm-parameters/index",
        name: "script-alarm-parameters",
        meta: { title: "脚本报警参数" },
        component: async () => await import("@/views/script-alarm-parameters/index.vue"),
      },
    ],
  },
  {
    path: "/node-details",
    component: Layout,
    meta: { title: "节点详情" },
    redirect: "/node-details/index",
    children: [
      {
        path: "/node-details/index",
        name: "node-details",
        meta: { title: "节点详情" },
        component: async () => await import("@/views/node-details/index.vue"),
      },
    ],
  },
];
