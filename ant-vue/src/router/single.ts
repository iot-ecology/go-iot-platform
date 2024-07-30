import Layout from "@/layout/index.vue";
export default [
  {
    path: "",
    key: "sub1",
    component: Layout,
    meta: { title: "message.deviceManagement" },
    redirect: "/mqtt-management/index",
    children: [
      {
        path: "/deviceInfos/index",
        name: "deviceInfos",
        meta: { title: "message.deviceInfos" },
        component: async () => await import("@/views/deviceInfos/index.vue"),
      },
      {
        path: "/deviceGroup/index",
        name: "deviceGroup",
        meta: { title: "message.deviceGroup" },
        component: async () => await import("@/views/deviceGroup/index.vue"),
      },
    ],
  },
  {
    path: "",
    component: Layout,
    key: "sub2",
    meta: { title: "message.personnelManagement" },
    children: [
      {
        path: "/user/index",
        name: "user",
        meta: { title: "message.user" },
        component: async () => await import("@/views/user/index.vue"),
      },
      {
        path: "/dept/index",
        name: "dept",
        meta: { title: "message.dept" },
        component: async () => await import("@/views/dept/index.vue"),
      },
      {
        path: "/role/index",
        name: "role",
        meta: { title: "message.role" },
        component: async () => await import("@/views/role/index.vue"),
      },
      {
        path: "/messageList/index",
        name: "messageList",
        meta: { title: "message.messageList" },
        component: async () => await import("@/views/message-list/index.vue"),
      },
    ],
  },
  {
    path: "",
    component: Layout,
    key: "sub3",
    meta: { title: "message.lifeCycle" },
    children: [
      {
        path: "/product/index",
        name: "product",
        meta: { title: "产品" },
        component: async () => await import("@/views/product/index.vue"),
      },
      // {
      //   path: "/messageList/index",
      //   name: "messageList",
      //   meta: { title: "安装记录" },
      //   component: async () => await import("@/views/message-list/index.vue"),
      // },
      {
        path: "/shipmentRecords/index",
        name: "shipmentRecords",
        meta: { title: "message.shipmentRecords" },
        component: async () => await import("@/views/shipment-records/index.vue"),
      },
      {
        path: "/productionPlans/index",
        name: "productionPlans",
        meta: { title: "message.productionPlans" },
        component: async () => await import("@/views/production-plans/index.vue"),
      },
      {
        path: "/repairRecords/index",
        name: "repairRecords",
        meta: { title: "message.repairRecords" },
        component: async () => await import("@/views/repair-records/index.vue"),
      },
      {
        path: "/sim-card/index",
        name: "sim",
        meta: { title: "SIM卡" },
        component: async () => await import("@/views/sim-card/index.vue"),
      },
      // {
      //   path: "/messageList/index",
      //   name: "messageList",
      //   meta: { title: "SIM卡历史" },
      //   component: async () => await import("@/views/message-list/index.vue"),
      // },
    ],
  },
  {
    path: "",
    component: Layout,
    key: "sub4",
    meta: { title: "message.mqttProtocolManagement" },
    children: [
      {
        path: "/mqtt-management/index",
        name: "IconPreview",
        meta: { title: "message.mqttManagement" },
        component: async () => await import("@/views/mqtt-management/index.vue"),
      },
      {
        path: "/signal-configuration/index",
        name: "CompPreview",
        meta: { title: "message.signalConfig" },
        component: async () => await import("@/views/signal-configuration/index.vue"),
      },
      {
        path: "/signal/index",
        name: "Signal",
        meta: { title: "message.SignalAlarmConfig" },
        component: async () => await import("@/views/signal/index.vue"),
      },
      {
        path: "/visualization",
        component: async () => await import("@/views/visualization/index.vue"),
        meta: { title: "message.visualization" },
        redirect: "/visualization/list",
        children: [
          {
            path: "/visualization/list",
            name: "List",
            meta: { title: "列表" },
            component: async () => await import("@/views/visualization/list.vue"),
          },
          {
            path: "/visualization/add",
            name: "Draggable",
            meta: { title: "message.addition", hidden: true },
            component: async () => await import("@/views/visualization/add.vue"),
          },
        ],
      },
      {
        path: "/calculation-rules/index",
        name: "calculation-rules",
        meta: { title: "message.calculationRules" },
        component: async () => await import("@/views/calculation-rules/index.vue"),
      },
      {
        path: "/calculate-parameters/index",
        name: "calculate-parameters",
        meta: { title: "message.calculateParameters" },
        component: async () => await import("@/views/calculate-parameters/index.vue"),
      },
      {
        path: "/script-alarm/index",
        name: "script-alarm",
        meta: { title: "message.scriptAlarm" },
        component: async () => await import("@/views/script-alarm/index.vue"),
      },
      {
        path: "/script-alarm-parameters/index",
        name: "script-alarm-parameters",
        meta: { title: "message.scriptAlarmParameters" },
        component: async () => await import("@/views/script-alarm-parameters/index.vue"),
      },
      {
        path: "/node-details/index",
        name: "node-details",
        meta: { title: "message.nodeDetails" },
        component: async () => await import("@/views/node-details/index.vue"),
      },
    ],
  },
  {
    path: "",
    component: Layout,
    key: "sub5",
    meta: { title: "message.protocolDocking" },
    children: [
      {
        path: "/http/index",
        name: "http",
        meta: { title: "message.http" },
        component: async () => await import("@/views/http/index.vue"),
      },
      {
        path: "/tcp/index",
        name: "tcp",
        meta: { title: "message.tcp" },
        component: async () => await import("@/views/tcp/index.vue"),
      },
      {
        path: "/coap/index",
        name: "coap",
        meta: { title: "message.coap" },
        component: async () => await import("@/views/coap/index.vue"),
      },
      {
        path: "/websocket/index",
        name: "websocket",
        meta: { title: "message.websocket" },
        component: async () => await import("@/views/websocket/index.vue"),
      },
    ],
  },
  {
    path: "",
    component: Layout,
    key: "sub6",
    meta: { title: "message.noticeManagement" },
    children: [
      {
        path: "/feishuRobot/index",
        name: "feishuRobot",
        meta: { title: "message.feishuRobot" },
        component: async () => await import("@/views/feishuRobot/index.vue"),
      },
      {
        path: "/dingdingRobot/index",
        name: "dingdingRobot",
        meta: { title: "message.dingdingRobot" },
        component: async () => await import("@/views/dingdingRobot/index.vue"),
      },
    ],
  },
  {
    path: "",
    component: Layout,
    key: "sub7",
    meta: { title: "message.forwardManagement" },
    children: [
      {
        path: "/cassandra/index",
        name: "cassandra",
        meta: { title: "message.cassandra" },
        component: async () => await import("@/views/cassandra/index.vue"),
      },
      {
        path: "/clickhouse/index",
        name: "clickhouse",
        meta: { title: "message.clickhouse" },
        component: async () => await import("@/views/clickhouse/index.vue"),
      },
      {
        path: "/influxdb2/index",
        name: "influxdb2",
        meta: { title: "message.influxdb2" },
        component: async () => await import("@/views/influxdb2/index.vue"),
      },
      {
        path: "/mongo/index",
        name: "mongo",
        meta: { title: "message.mongo" },
        component: async () => await import("@/views/mongo/index.vue"),
      },
      {
        path: "/mysql/index",
        name: "mysql",
        meta: { title: "message.mysql" },
        component: async () => await import("@/views/mysql/index.vue"),
      },
    ],
  },
  // {
  //   path: "/mqtt-management",
  //   meta: { title: 'message.mqttManagement' },
  //   component: Layout,
  //   redirect: "/mqtt-management/index",
  //   children: [
  //     {
  //       path: "/mqtt-management/index",
  //       name: "IconPreview",
  //       meta: { title: "message.mqttManagement" },
  //       component: async () => await import("@/views/mqtt-management/list.vue"),
  //     },
  //   ],
  // },
  // {
  //   path: "/signal-configuration",
  //   component: Layout,
  //   meta: { title: "message.signalConfig" },
  //   redirect: "/signal-configuration/index",
  //   children: [
  //     {
  //       path: "/signal-configuration/index",
  //       name: "CompPreview",
  //       meta: { title: "message.signalConfig" },
  //       component: async () => await import("@/views/signal-configuration/list.vue"),
  //     },
  //   ],
  // },
  // {
  //   path: "/signal",
  //   component: Layout,
  //   meta: { title: "message.SignalAlarmConfig" },
  //   redirect: "/signal/index",
  //   children: [
  //     {
  //       path: "/signal/index",
  //       name: "Signal",
  //       meta: { title: "message.SignalAlarmConfig" },
  //       component: async () => await import("@/views/signal/list.vue"),
  //     },
  //   ],
  // },
  // {
  //   path: "/visualization",
  //   component: Layout,
  //   meta: { title: "message.visualization" },
  //   redirect: "/visualization/index",
  //   children: [
  //     {
  //       path: "/visualization/index",
  //       name: "List",
  //       meta: { title: "列表" },
  //       component: async () => await import("@/views/visualization/list.vue"),
  //     },
  //     {
  //       path: "/visualization/add",
  //       name: "Draggable",
  //       meta: { title: "新增", hidden: true },
  //       component: async () => await import("@/views/visualization/add.vue"),
  //     },
  //   ],
  // },
  // {
  //   path: "/calculation-rules",
  //   component: Layout,
  //   meta: { title: "message.calculationRules" },
  //   redirect: "/calculation-rules/index",
  //   children: [
  //     {
  //       path: "/calculation-rules/index",
  //       name: "calculation-rules",
  //       meta: { title: "message.calculationRules" },
  //       component: async () => await import("@/views/calculation-rules/list.vue"),
  //     },
  //   ],
  // },

  // {
  //   path: "/calculate-parameters",
  //   component: Layout,
  //   meta: { title: "message.calculateParameters" },
  //   redirect: "/calculate-parameters/index",
  //   children: [
  //     {
  //       path: "/calculate-parameters/index",
  //       name: "calculate-parameters",
  //       meta: { title: "message.calculateParameters" },
  //       component: async () => await import("@/views/calculate-parameters/list.vue"),
  //     },
  //   ],
  // },
  // {
  //   path: "/script-alarm",
  //   component: Layout,
  //   meta: { title: "message.scriptAlarm" },
  //   redirect: "/script-alarm/index",
  //   children: [
  //     {
  //       path: "/script-alarm/index",
  //       name: "script-alarm",
  //       meta: { title: "message.scriptAlarm" },
  //       component: async () => await import("@/views/script-alarm/list.vue"),
  //     },
  //   ],
  // },
  // {
  //   path: "/script-alarm-parameters",
  //   component: Layout,
  //   meta: { title: "message.scriptAlarmParameters" },
  //   redirect: "/script-alarm-parameters/index",
  //   children: [
  //     {
  //       path: "/script-alarm-parameters/index",
  //       name: "script-alarm-parameters",
  //       meta: { title: "message.scriptAlarmParameters" },
  //       component: async () => await import("@/views/script-alarm-parameters/list.vue"),
  //     },
  //   ],
  // },
  // {
  //   path: "/node-details",
  //   component: Layout,
  //   meta: { title: "message.nodeDetails" },
  //   redirect: "/node-details/index",
  //   children: [
  //     {
  //       path: "/node-details/index",
  //       name: "node-details",
  //       meta: { title: "message.nodeDetails" },
  //       component: async () => await import("@/views/node-details/list.vue"),
  //     },
  //   ],
  // },
  // {
  //   path: "/dept",
  //   component: Layout,
  //   meta: { title: "message.dept" },
  //   redirect: "/dept/index",
  //   children: [
  //     {
  //       path: "/dept/index",
  //       name: "dept",
  //       meta: { title: "message.dept" },
  //       component: async () => await import("@/views/dept/list.vue"),
  //     },
  //   ],
  // },
  // {
  //   path: "/deviceInfos",
  //   component: Layout,
  //   meta: { title: "message.deviceInfos" },
  //   redirect: "/deviceInfos/index",
  //   children: [
  //     {
  //       path: "/deviceInfos/index",
  //       name: "deviceInfos",
  //       meta: { title: "message.deviceInfos" },
  //       component: async () => await import("@/views/deviceInfos/list.vue"),
  //     },
  //   ],
  // },
  // {
  //   path: "/messageList",
  //   component: Layout,
  //   meta: { title: "message.messageList" },
  //   redirect: "/messageList/index",
  //   children: [
  //     {
  //       path: "/messageList/index",
  //       name: "messageList",
  //       meta: { title: "message.messageList" },
  //       component: async () => await import("@/views/message-list/list.vue"),
  //     },
  //   ],
  // },
  // {
  //   path: "/productionPlans",
  //   component: Layout,
  //   meta: { title: "message.productionPlans" },
  //   redirect: "/productionPlans/index",
  //   children: [
  //     {
  //       path: "/productionPlans/index",
  //       name: "productionPlans",
  //       meta: { title: "message.productionPlans" },
  //       component: async () => await import("@/views/production-plans/list.vue"),
  //     },
  //   ],
  // },
  // {
  //   path: "/repairRecords",
  //   component: Layout,
  //   meta: { title: "message.repairRecords" },
  //   redirect: "/repairRecords/index",
  //   children: [
  //     {
  //       path: "/repairRecords/index",
  //       name: "repairRecords",
  //       meta: { title: "message.repairRecords" },
  //       component: async () => await import("@/views/repair-records/list.vue"),
  //     },
  //   ],
  // },
  // {
  //   path: "/role",
  //   component: Layout,
  //   meta: { title: "message.role" },
  //   redirect: "/role/index",
  //   children: [
  //     {
  //       path: "/role/index",
  //       name: "role",
  //       meta: { title: "message.role" },
  //       component: async () => await import("@/views/role/list.vue"),
  //     },
  //   ],
  // },
  // {
  //   path: "/shipmentRecords",
  //   component: Layout,
  //   meta: { title: "message.shipmentRecords" },
  //   redirect: "/shipmentRecords/index",
  //   children: [
  //     {
  //       path: "/shipmentRecords/index",
  //       name: "shipmentRecords",
  //       meta: { title: "message.shipmentRecords" },
  //       component: async () => await import("@/views/shipment-records/list.vue"),
  //     },
  //   ],
  // }
];
