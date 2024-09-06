/**
 * @name umi 的路由配置
 * @description 只支持 path,component,routes,redirect,wrappers,name,icon 的配置
 * @param path  path 只支持两种占位符配置，第一种是动态参数 :id 的形式，第二种是 * 通配符，通配符只能出现路由字符串的最后。
 * @param component 配置 location 和 path 匹配后用于渲染的 React 组件路径。可以是绝对路径，也可以是相对路径，如果是相对路径，会从 src/pages 开始找起。
 * @param routes 配置子路由，通常在需要为多个路径增加 layout 组件时使用。
 * @param redirect 配置路由跳转
 * @param wrappers 配置路由组件的包装组件，通过包装组件可以为当前的路由组件组合进更多的功能。 比如，可以用于路由级别的权限校验
 * @param name 配置路由的标题，默认读取国际化文件 menu.ts 中 menu.xxxx 的值，如配置 name 为 login，则读取 menu.ts 中 menu.login 的取值作为标题
 * @param icon 配置路由的图标，取值参考 https://admin/components/icon-cn， 注意去除风格后缀和大小写，如想要配置图标为 <StepBackwardOutlined /> 则取值应为 stepBackward 或 StepBackward，如想要配置图标为 <UserOutlined /> 则取值应为 user 或者 User
 * @doc https://umijs.org/docs/guides/routes
 */
export default [
  {
    path: '/user',
    layout: false,
    routes: [
      {
        name: 'login',
        path: '/user/login',
        component: './User/Login',
      },
    ],
  },
  {
    path: '/welcome',
    name: 'welcome',
    icon: 'smile',
    component: './Welcome',
  },
  {
    path: '/device',
    name: 'device',
    icon: 'crown',
    access: 'canAdmin',
    routes: [
      {
        path: '/device',
        redirect: '/device/device-info',
      },
      {
        path: '/device/device-info',
        name: 'device-info',
        component: './Device/DeviceInfo/DeviceInfo',
      }, {
        path: '/device/device-group',
        name: 'device-group',
        component: './Device/DeviceGroup/DeviceGroup',
      },
    ],
  },
  {
    path: '/user',
    name: 'user',
    icon: 'TeamOutlined',
    access: 'canAdmin',
    routes: [
      {
        path: '/user',
        redirect: '/user/user-list',
      },
      {
        path: '/user/user-list',
        name: 'user-list',
        component: './User/UserList/UserList',
      }, {
        path: '/user/dept-list',
        name: 'dept-list',
        component: './User/DeptList/DeptList',
      },  {
        path: '/user/role-list',
        name: 'role-list',
        component: './User/RoleList/RoleList',
      }, {
        path: '/user/message-list',
        name: 'message-list',
        component: './User/MessageList/MessageList',
      },
    ],
  },


  {
    path: '/lifecycle',
    name: 'lifecycle',
    icon: 'UnorderedListOutlined',
    access: 'canAdmin',
    routes: [
      {
        path: '/lifecycle',
        redirect: '/lifecycle/product-list',
      },
      {
        path: '/lifecycle/product-list',
        name: 'product-list',
        component: './Lifecycle/ProductList/ProductList',
      }, {
        path: '/lifecycle/send-list',
        name: 'send-list',
        component: './Lifecycle/SendList/SendList',
      },{
        path: '/lifecycle/produce-list',
        name: 'produce-list',
        component: './Lifecycle/ProduceList/ProduceList',
      },{
        path: '/lifecycle/operation-list',
        name: 'operation-list',
        component: './Lifecycle/OperationList/OperationList',
      },{
        path: '/lifecycle/sim-list',
        name: 'sim-list',
        component: './Lifecycle/SimList/SimList',
      },
    ],
  },
  {
    path: '/protocol',
    name: 'protocol',
    icon: 'ReadOutlined',
    access: 'canAdmin',
    routes: [
      {
        path: '/protocol/mqtt',
        name: 'mqtt',
        routes: [
          {
            path: '/protocol/mqtt/clients',
            name: 'clients',
            component: './Protocol/mqtt/clients',
          }, {
            path: '/protocol/mqtt/nodes',
            name: 'nodes',
            component: './Protocol/mqtt/nodes',
          },
        ]
      },
      {
        path: '/protocol/http',
        name: 'http',
        routes: [
          {
            path: '/protocol/http/clients',
            name: 'clients',
            component: './Protocol/http/clients',
          },
        ]
      },{
        path: '/protocol/tcp',
        name: 'tcp',
        routes: [
          {
            path: '/protocol/tcp/clients',
            name: 'clients',
            component: './Protocol/tcp/clients',
          },
        ]
      },{
        path: '/protocol/coap',
        name: 'coap',
        routes: [
          {
            path: '/protocol/coap/clients',
            name: 'clients',
            component: './Protocol/coap/clients',
          },
        ]
      },{
        path: '/protocol/ws',
        name: 'ws',
        routes: [
          {
            path: '/protocol/ws/clients',
            name: 'clients',
            component: './Protocol/ws/clients',
          },
        ]
      },
    ]
  },

  {
    path: '/data',
    name: 'data',
    icon: 'TableOutlined',
    access: 'canAdmin',
    routes: [
      {
        path: '/data/signal',
        name: 'signal',
        component: './Data/Signal',
      },  {
        path: '/data/signal-waring',
        name: 'signal-waring',
        component: './Data/SignalWaring',
      },  {
        path: '/data/signal-vis',
        name: 'signal-vis',
        component: './Data/SignalVis',
      },  {
        path: '/data/calc',
        name: 'calc',
        component: './Data/Calc',
      },  {
        path: '/data/calc-param',
        name: 'calc-param',
        component: './Data/CalcParam',
      },  {
        path: '/data/script-waring',
        name: 'script-waring',
        component: './Data/ScriptWaring',
      },  {
        path: '/data/script-param',
        name: 'script-param',
        component: './Data/ScriptParam',
      },
    ]
  },


  {
    path: '/notice',
    name: 'notice',
    icon: "NotificationOutlined",
    access: 'canAdmin',
    routes: [
      {
        path: '/notice/feishu',
        name: 'feishu',
        component: './notice/feishu',
      },{
        path: '/notice/dingding',
        name: 'dingding',
        component: './notice/dingding',
      },
    ]
  },
  {
    path: '/forward',
    name: 'forward',
    icon: 'DeliveredProcedureOutlined',
    access: 'canAdmin',
    routes: [
      {
        path: '/forward/cassandra',
        name: 'cassandra',
        component: './forward/cassandra',
      },{
        path: '/forward/clickhouse',
        name: 'clickhouse',
        component: './forward/clickhouse',
      },{
        path: '/forward/influxdb2',
        name: 'influxdb2',
        component: './forward/influxdb2',
      },{
        path: '/forward/mongo',
        name: 'mongo',
        component: './forward/mongo',
      },{
        path: '/forward/mysql',
        name: 'mysql',
        component: './forward/mysql',
      },
    ]
  },
  {
    name: 'list.table-list',
    icon: 'table',
    path: '/list',
    component: './TableList',
  },
  {
    path: '/',
    redirect: '/welcome',
  },
  {
    path: '*',
    layout: false,
    component: './404',
  },
];
