// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 获取当前的用户 GET /api/currentUser */
export async function currentUser(options?: { [key: string]: any }) {
  return request<{
    data: API.CurrentUser;
  }>('/api/currentUser', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 退出登录接口 POST /api/login/outLogin */
export async function outLogin(options?: { [key: string]: any }) {
  return request<Record<string, any>>('/api/login/outLogin', {
    method: 'POST',
    ...(options || {}),
  });
}

/** 登录接口 POST /api/login/account */
export async function login(body: API.LoginParams, options?: { [key: string]: any }) {
  return request<API.LoginResult>('/api/login/account', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 此处后端没有提供注释 GET /api/notices */
export async function getNotices(options?: { [key: string]: any }) {
  return request<API.NoticeIconList>('/api/notices', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 获取规则列表 GET /api/rule */
export async function rule(
  params: {
    // query
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  return request<API.RuleList>('/api/rule', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

export async function deptPage(
  params: {
    // query
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.DeptListItem>>('/api/Dept/page', {
    method: 'GET',
    params: {
      page: params.current,
      page_size: params.pageSize,
      ...params,
    },
    ...(options || {}),
  });
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function devicePage(
  params: {
    // query
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
    sn?: string;
    protocol?: string;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.DeviceInfoItem>>('/api/DeviceInfo/page', {
    method: 'GET',
    params: {
      page: params.current,
      page_size: params.pageSize,
      ...params,
    },
    ...(options || {}),
  });
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function rolePage(
  params: {
    // query
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.RoleListItem>>('/api/Role/page', {
    method: 'GET',
    params: {
      page: params.current,
      page_size: params.pageSize,
      ...params,
    },
    ...(options || {}),
  });
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function messagePage(
  params: {
    // query
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
    message_type_id?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.MessageListItem>>('/api/MessageList/page', {
    method: 'GET',
    params: {
      page: params.current,
      page_size: params.pageSize,
      message_type_id: params.message_type_id,
    },
    ...(options || {}),
  });
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function deptList() {
  return request<API.CommonResp<API.DeptListItem[]>>('/api/Dept/list', {
    method: 'GET',
  });
}

export async function mqttList() {
  return request<API.CommonResp<API.MqttListItem[]>>('/api/mqtt/list', {
    method: 'GET',
  });
}

export async function CassandraTransmitList() {
  return request<API.CommonResp<API.MqttListItem[]>>('/api/CassandraTransmit/list', {
    method: 'GET',
  });
}

export async function ClickhouseTransmitList() {
  return request<API.CommonResp<API.MqttListItem[]>>('/api/ClickhouseTransmit/list', {
    method: 'GET',
  });
}

export async function InfluxdbTransmitList() {
  return request<API.CommonResp<API.MqttListItem[]>>('/api/InfluxdbTransmit/list', {
    method: 'GET',
  });
}

export async function MongoTransmitList() {
  return request<API.CommonResp<API.MqttListItem[]>>('/api/MongoTransmit/list', {
    method: 'GET',
  });
}

export async function MySQLTransmitList() {
  return request<API.CommonResp<API.MqttListItem[]>>('/api/MySQLTransmit/list', {
    method: 'GET',
  });
}

export async function deviceList() {
  return request<API.CommonResp<API.DeviceInfoListItem[]>>('/api/DeviceInfo/list', {
    method: 'GET',
  });
}

export async function scriptWaringList() {
  return request<API.CommonResp<API.MqttListItem[]>>('/api/signal-delay-waring/list', {
    method: 'GET',
  });
}

export async function calcRuleList() {
  return request<API.CommonResp<API.CalcRuleListItem[]>>('/api/calc-rule/list', {
    method: 'GET',
  });
}

export async function signalList(data: any) {
  const request1 = await request<API.CommonResp<string>>('/api/signal/list', {
    method: 'GET',
    params: data,
  });
  return request1;
}

export async function mqttById(id: any) {
  return request<API.CommonResp<API.MqttListItem>>('/api/mqtt/byId/' + id, {
    method: 'GET',
  });
}

export async function signalById(id: any) {
  return request<API.CommonResp<API.MqttListItem>>('/api/signal/byId/' + id, {
    method: 'GET',
  });
}

export async function genScriptParam(data: string | number) {
  const request1 = await request<API.CommonResp<string>>(
    '/api/signal-delay-waring/GenParam/' + data,
    {
      method: 'POST',
    },
  );
  return request1;
}

export async function mockScriptWaringRun(data: string | number) {
  const request1 = await request<API.CommonResp<string>>('/api/signal-delay-waring/Mock/' + data, {
    method: 'POST',
  });
  return request1;
}

export async function scriptWaringById(id: any) {
  return request<API.CommonResp<API.MqttListItem>>('/api/signal-delay-waring/byId/' + id, {
    method: 'GET',
  });
}

export async function calcRuleMockRun(data: any) {
  const request1 = await request<API.CommonResp<string>>('/api/calc-rule/mock', {
    method: 'POST',
    data: data,
  });
  return request1;
}

export async function calcRuleStart(data: string | number) {
  const request1 = await request<API.CommonResp<string>>('/api/calc-rule/start/' + data, {
    method: 'POST',
  });
  return request1;
}

export async function calcRuleStop(data: string | number) {
  const request1 = await request<API.CommonResp<string>>('/api/calc-rule/stop/' + data, {
    method: 'POST',
  });
  return request1;
}

export async function influxdbQuery(data: any) {
  const request1 = await request<API.CommonResp<string>>('/api/query/influxdb', {
    method: 'POST',
    data: data,
  });
  return request1;
}

export async function productList() {
  return request<API.CommonResp<API.ProductItem[]>>('/api/product/list', {
    method: 'GET',
  });
}

export async function FindByShipmentProductDetail(id: string) {
  return request<API.CommonResp<API.ProductItem[]>>(
    '/api/ShipmentRecord/FindByShipmentProductDetail/' + id,
    {
      method: 'GET',
    },
  );
}

export async function userPage(
  params: {
    // query
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.UserList>('/api/User/page', {
    method: 'GET',
    params: {
      page: params.current,
      page_size: params.pageSize,
      ...params,
    },
    ...(options || {}),
  });
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function calcRuleResultQuery(
  params: {
    rule_id?: number;
    start_time?: number;
    end_time?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonResp<string>>('/api/calc-rule/rd', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
  return request1;
}

export async function simCardPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
    access_number?: string;
    iccid?: string;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.SimListItem>>('/api/SimCard/page', {
    method: 'GET',
    params: {
      page: params.current,
      page_size: params.pageSize,
      AccessNumber: params.access_number ? params.access_number : '',

      iccid: params.iccid ? params.iccid : '',
    },
    ...(options || {}),
  });
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function CassandraTransmitPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.CassandraTransmitListItem>>(
    '/api/CassandraTransmit/page',
    {
      method: 'GET',
      params: {
        page: params.current,
        page_size: params.pageSize,

        ...params,
      },
      ...(options || {}),
    },
  );
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function CassandraTransmitBindPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.CassandraTransmitListItem>>(
    '/api/CassandraTransmitBind/page',
    {
      method: 'GET',
      params: {
        page: params.current,
        page_size: params.pageSize,

        ...params,
      },
      ...(options || {}),
    },
  );
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function ClickhouseTransmitBindPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.CassandraTransmitListItem>>(
    '/api/ClickhouseTransmitBind/page',
    {
      method: 'GET',
      params: {
        page: params.current,
        page_size: params.pageSize,

        ...params,
      },
      ...(options || {}),
    },
  );
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function InfluxdbTransmitBindPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.CassandraTransmitListItem>>(
    '/api/InfluxdbTransmitBind/page',
    {
      method: 'GET',
      params: {
        page: params.current,
        page_size: params.pageSize,

        ...params,
      },
      ...(options || {}),
    },
  );
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function MongoTransmitBindPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.CassandraTransmitListItem>>(
    '/api/MongoTransmitBind/page',
    {
      method: 'GET',
      params: {
        page: params.current,
        page_size: params.pageSize,

        ...params,
      },
      ...(options || {}),
    },
  );
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function MySQLTransmitBindPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.CassandraTransmitListItem>>(
    '/api/MySQLTransmitBind/page',
    {
      method: 'GET',
      params: {
        page: params.current,
        page_size: params.pageSize,

        ...params,
      },
      ...(options || {}),
    },
  );
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function ClickhouseTransmitPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.CassandraTransmitListItem>>(
    '/api/ClickhouseTransmit/page',
    {
      method: 'GET',
      params: {
        page: params.current,
        page_size: params.pageSize,

        ...params,
      },
      ...(options || {}),
    },
  );
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function MongoTransmitPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.CassandraTransmitListItem>>(
    '/api/MongoTransmit/page',
    {
      method: 'GET',
      params: {
        page: params.current,
        page_size: params.pageSize,

        ...params,
      },
      ...(options || {}),
    },
  );
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function InfluxdbTransmitUpdateFormPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.CassandraTransmitListItem>>(
    '/api/InfluxdbTransmit/page',
    {
      method: 'GET',
      params: {
        page: params.current,
        page_size: params.pageSize,

        ...params,
      },
      ...(options || {}),
    },
  );
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function MySQLTransmitPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.CassandraTransmitListItem>>(
    '/api/MySQLTransmit/page',
    {
      method: 'GET',
      params: {
        page: params.current,
        page_size: params.pageSize,

        ...params,
      },
      ...(options || {}),
    },
  );
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function httpHandlerPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.HttpHandlerListItem>>('/api/HttpHandler/page', {
    method: 'GET',
    params: {
      page: params.current,
      page_size: params.pageSize,
      ...params,
    },
    ...(options || {}),
  });
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function WebsocketHandlerPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.HttpHandlerListItem>>(
    '/api/WebsocketHandler/page',
    {
      method: 'GET',
      params: {
        page: params.current,
        page_size: params.pageSize,
        ...params,
      },
      ...(options || {}),
    },
  );
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function TcpHandlerPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.HttpHandlerListItem>>('/api/TcpHandler/page', {
    method: 'GET',
    params: {
      page: params.current,
      page_size: params.pageSize,
      ...params,
    },
    ...(options || {}),
  });
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function CoapHandlerPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.HttpHandlerListItem>>('/api/CoapHandler/page', {
    method: 'GET',
    params: {
      page: params.current,
      page_size: params.pageSize,
      ...params,
    },
    ...(options || {}),
  });
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function scriptWaringPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
    name?: string;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.SimListItem>>('/api/signal-delay-waring/page', {
    method: 'GET',
    params: {
      page: params.current,
      page_size: params.pageSize,
      ...params,
    },
    ...(options || {}),
  });
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function calcParamPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.CalcParamListItem>>('/api/calc-param/page', {
    method: 'GET',
    params: {
      page: params.current,
      page_size: params.pageSize,
    },
    ...(options || {}),
  });
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function calcRulePage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.SimListItem>>('/api/calc-rule/page', {
    method: 'GET',
    params: {
      page: params.current,
      page_size: params.pageSize,
      ...params,
    },
    ...(options || {}),
  });
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function scritpParamPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.SimListItem>>(
    '/api/signal-delay-waring-param/page',
    {
      method: 'GET',
      params: {
        page: params.current,
        page_size: params.pageSize,
        ...params,
      },
      ...(options || {}),
    },
  );
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function signalWaringPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.SignalWaringItem>>(
    '/api/signal-waring-config/page',
    {
      method: 'GET',
      params: {
        page: params.current,
        page_size: params.pageSize,
        ...params,
      },
      ...(options || {}),
    },
  );
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function feishuPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
    access_token?: string;
    content?: string;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.FeiShuListItem>>('/api/FeiShuId/page', {
    method: 'GET',
    params: {
      page: params.current,
      page_size: params.pageSize,
      ...params,
    },
    ...(options || {}),
  });
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function dingDingPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
    access_token?: string;

    content?: string;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.DingDingListItem>>('/api/DingDing/page', {
    method: 'GET',
    params: {
      page: params.current,
      page_size: params.pageSize,
      ...params,
    },
    ...(options || {}),
  });
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function signalPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.SimListItem>>('/api/signal/page', {
    method: 'GET',
    params: {
      page: params.current,
      page_size: params.pageSize,
      ...params,
    },
    ...(options || {}),
  });
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function mqttPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.SimListItem>>('/api/mqtt/page', {
    method: 'GET',
    params: {
      page: params.current,
      page_size: params.pageSize,
      ...params,
    },
    ...(options || {}),
  });
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function mqttScriptCheck(param: string, script: string) {
  const request1 = await request<API.CommonResp<string>>('/api/mqtt/check-script', {
    method: 'POST',
    data: {
      param: param,
      script: script,
    },
  });
  return request1;
}

export async function setMqttScriptCheck(id: number, script: string) {
  const request1 = await request<API.CommonResp<string>>('/api/mqtt/set-script', {
    method: 'POST',
    data: {
      id: id,
      script: script,
    },
  });
  return request1;
}

export async function sendMqttClientMessage(id: string, data: any) {
  const request1 = await request<API.CommonResp<string>>('/api/mqtt/send', {
    method: 'POST',
    params: {
      id: id,
    },
    data: data,
  });
  return request1;
}

export async function startMqttClient(id: number) {
  const request1 = await request<API.CommonResp<string>>('/api/mqtt/start', {
    method: 'GET',
    params: {
      id: id,
    },
  });
  return request1;
}

export async function stopMqttClient(id: number) {
  const request1 = await request<API.CommonResp<string>>('/api/mqtt/stop', {
    method: 'GET',
    params: {
      id: id,
    },
  });
  return request1;
}

export async function shipmentPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
    customer_name?: string;
    status?: string;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.ShipmentRecordListItem>>(
    '/api/ShipmentRecord/page',
    {
      method: 'GET',
      params: {
        page: params.current,
        page_size: params.pageSize,
        customer_name: params.customer_name,
        status: params.status,
      },
      ...(options || {}),
    },
  );
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function productPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.ProductItem>>('/api/product/page', {
    method: 'GET',
    params: {
      page: params.current,
      page_size: params.pageSize,
      ...params,
    },
    ...(options || {}),
  });
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

export async function deviceGroupPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.DeviceGroupItem>>('/api/device_group/page', {
    method: 'GET',
    params: {
      page: params.current,
      page_size: params.pageSize,
      ...params,
    },
    ...(options || {}),
  });
  return {
    data: request1.data?.data,
    success: request1.code === 20000,
    total: request1.data?.total,
  };
}

/** 更新规则 PUT /api/rule */
export async function updateRule(options?: { [key: string]: any }) {
  return request<API.RuleListItem>('/api/rule', {
    method: 'POST',
    data: {
      method: 'update',
      ...(options || {}),
    },
  });
}

/** 新建规则 POST /api/rule */
export async function addRule(options?: { [key: string]: any }) {
  return request<API.RuleListItem>('/api/rule', {
    method: 'POST',
    data: {
      method: 'post',
      ...(options || {}),
    },
  });
}

export async function addUser(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/User/create', {
    method: 'POST',
    data: {
      method: 'post',
      ...(options || {}),
    },
  });
}

export async function addSim(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/SimCard/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addCassandraTransmitPage(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/CassandraTransmit/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addCassandraTransmitBind(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/CassandraTransmitBind/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addClickhouseTransmitBind(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/ClickhouseTransmitBind/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addInfluxdbTransmitBind(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/InfluxdbTransmitBind/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addMongoTransmitBind(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/MongoTransmitBind/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addMySQLTransmitBind(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/MySQLTransmitBind/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addClickhouseTransmit(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/ClickhouseTransmit/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addMongoTransmit(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/MongoTransmit/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addInfluxdbTransmit(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/InfluxdbTransmit/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addMySQLTransmit(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/MySQLTransmit/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addHttpHandler(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/HttpHandler/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addWebsocketHandler(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/WebsocketHandler/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addTcpHandler(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/TcpHandler/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addCoapHandler(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/CoapHandler/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addScriptWaring(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/signal-delay-waring/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addCalcParam(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/calc-param/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addCalcRule(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/calc-rule/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addScriptWaringParam(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/signal-delay-waring-param/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addSignalWaring(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/signal-waring-config/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addFeishu(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/FeiShuId/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addDingDing(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/DingDing/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addSignal(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/signal/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addMQTT(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/mqtt/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addShipmentRecord(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/ShipmentRecord/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addDevice(options?: { [key: string]: any }) {
  return request<API.CommonResp<string>>('/api/DeviceInfo/create', {
    method: 'POST',
    data: {
      ...(options || {}),
    },
  });
}

export async function addProduct(options?: { [key: string]: any }) {
  return request<API.ProductItem>('/api/product/create', {
    method: 'POST',
    data: {
      method: 'post',
      ...(options || {}),
    },
  });
}

export async function addDeviceGroup(options?: { [key: string]: any }) {
  return request<API.RuleListItem>('/api/device_group/create', {
    method: 'POST',
    data: {
      method: 'post',
      ...(options || {}),
    },
  });
}

export async function addDept(options?: { [key: string]: any }) {
  return request<API.CommonResp<any>>('/api/Dept/create', {
    method: 'POST',
    data: {
      method: 'post',
      ...(options || {}),
    },
  });
}

export async function addRole(options?: { [key: string]: any }) {
  return request<API.CommonResp<any>>('/api/Role/create', {
    method: 'POST',
    data: {
      method: 'post',
      ...(options || {}),
    },
  });
}

export async function deleteUser(id: any) {
  return request<API.CommonResp<string>>('/api/User/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteSimCard(id: any) {
  return request<API.CommonResp<string>>('/api/SimCard/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteCassandraTransmit(id: any) {
  return request<API.CommonResp<string>>('/api/CassandraTransmit/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteCassandraTransmitBind(id: any) {
  return request<API.CommonResp<string>>('/api/CassandraTransmitBind/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteClickhouseTransmitBind(id: any) {
  return request<API.CommonResp<string>>('/api/ClickhouseTransmitBind/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteInfluxdbTransmitBind(id: any) {
  return request<API.CommonResp<string>>('/api/InfluxdbTransmitBind/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteMongoTransmitBind(id: any) {
  return request<API.CommonResp<string>>('/api/MongoTransmitBind/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteMySQLTransmitBind(id: any) {
  return request<API.CommonResp<string>>('/api/MySQLTransmitBind/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteClickhouseTransmit(id: any) {
  return request<API.CommonResp<string>>('/api/ClickhouseTransmit/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteMongoTransmit(id: any) {
  return request<API.CommonResp<string>>('/api/MongoTransmit/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteInfluxdbTransmit(id: any) {
  return request<API.CommonResp<string>>('/api/InfluxdbTransmit/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteMySQLTransmit(id: any) {
  return request<API.CommonResp<string>>('/api/MySQLTransmit/delete/' + id, {
    method: 'POST',
  });
}

export async function delteHttpHandler(id: any) {
  return request<API.CommonResp<string>>('/api/HttpHandler/delete/' + id, {
    method: 'POST',
  });
}

export async function deltedWebsocketHandler(id: any) {
  return request<API.CommonResp<string>>('/api/WebsocketHandler/delete/' + id, {
    method: 'POST',
  });
}

export async function deltedTcpHandler(id: any) {
  return request<API.CommonResp<string>>('/api/TcpHandler/delete/' + id, {
    method: 'POST',
  });
}

export async function delteCoapHandler(id: any) {
  return request<API.CommonResp<string>>('/api/CoapHandler/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteScriptWaring(id: any) {
  return request<API.CommonResp<string>>('/api/signal-delay-waring/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteCalcParam(id: any) {
  return request<API.CommonResp<string>>('/api/calc-param/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteCalcRule(id: any) {
  return request<API.CommonResp<string>>('/api/calc-rule/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteScriptWaringParam(id: any) {
  return request<API.CommonResp<string>>('/api/signal-delay-waring-param/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteSignalWaring(id: any) {
  return request<API.CommonResp<string>>('/api/signal-waring-config/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteFeishu(id: any) {
  return request<API.CommonResp<string>>('/api/FeiShuId/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteDingDing(id: any) {
  return request<API.CommonResp<string>>('/api/DingDing/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteSignal(id: any) {
  return request<API.CommonResp<string>>('/api/signal/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteMQTT(id: any) {
  return request<API.CommonResp<string>>('/api/mqtt/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteDeviceInfo(id: any) {
  return request<API.CommonResp<string>>('/api/DeviceInfo/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteShipmentRecord(id: any) {
  return request<API.CommonResp<string>>('/api/ShipmentRecord/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteProduct(id: any) {
  return request<API.CommonResp<string>>('/api/product/delete/' + id, {
    method: 'POST',
  });
}

export async function deleteDeviceGroup(id: any) {
  return request<API.CommonResp<string>>('/api/device_group/delete/' + id, {
    method: 'POST',
  });
}

export async function updateUser(dt: any) {
  return request<API.CommonResp<string>>('/api/User/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateSimCard(dt: any) {
  return request<API.CommonResp<string>>('/api/SimCard/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateCassandraTransmit(dt: any) {
  return request<API.CommonResp<string>>('/api/CassandraTransmit/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateCassandraTransmitBind(dt: any) {
  return request<API.CommonResp<string>>('/api/CassandraTransmitBind/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateClickhouseTransmitBind(dt: any) {
  return request<API.CommonResp<string>>('/api/ClickhouseTransmitBind/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateInfluxdbTransmitBind(dt: any) {
  return request<API.CommonResp<string>>('/api/InfluxdbTransmitBind/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateMongoTransmitBind(dt: any) {
  return request<API.CommonResp<string>>('/api/MongoTransmitBind/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateMySQLTransmitBind(dt: any) {
  return request<API.CommonResp<string>>('/api/MySQLTransmitBind/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateClickhouseTransmit(dt: any) {
  return request<API.CommonResp<string>>('/api/ClickhouseTransmit/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateMongoTransmitPage(dt: any) {
  return request<API.CommonResp<string>>('/api/MongoTransmit/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateInfluxdbTransmit(dt: any) {
  return request<API.CommonResp<string>>('/api/InfluxdbTransmit/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateMySQLTransmit(dt: any) {
  return request<API.CommonResp<string>>('/api/MySQLTransmit/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateHttpHandler(dt: any) {
  return request<API.CommonResp<string>>('/api/HttpHandler/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateWsHandler(dt: any) {
  return request<API.CommonResp<string>>('/api/WebsocketHandler/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateTcpHandler(dt: any) {
  return request<API.CommonResp<string>>('/api/TcpHandler/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateCoapHandler(dt: any) {
  return request<API.CommonResp<string>>('/api/CoapHandler/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateScriptWaring(dt: any) {
  return request<API.CommonResp<string>>('/api/signal-delay-waring/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateCalcParam(dt: any) {
  return request<API.CommonResp<string>>('/api/calc-param/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateSignalWaring(dt: any) {
  return request<API.CommonResp<string>>('/api/signal-waring-config/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateCalcRule(dt: any) {
  return request<API.CommonResp<string>>('/api/calc-rule/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateScriptWaringParam(dt: any) {
  return request<API.CommonResp<string>>('/api/signal-delay-waring-param/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateFeishu(dt: any) {
  return request<API.CommonResp<string>>('/api/FeiShuId/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateDingDing(dt: any) {
  return request<API.CommonResp<string>>('/api/DingDing/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateSignal(dt: any) {
  return request<API.CommonResp<string>>('/api/signal/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateMQTT(dt: any) {
  return request<API.CommonResp<string>>('/api/mqtt/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateDeviceInfo(dt: any) {
  return request<API.CommonResp<string>>('/api/DeviceInfo/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateShipmentRecord(dt: any) {
  return request<API.CommonResp<string>>('/api/ShipmentRecord/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateProduct(dt: any) {
  return request<API.CommonResp<string>>('/api/product/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateDeviceGroup(dt: any) {
  return request<API.CommonResp<string>>('/api/device_group/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateRole(dt: any) {
  return request<API.CommonResp<string>>('/api/Role/update', {
    method: 'POST',
    data: dt,
  });
}

export async function updateDept(dt: any) {
  return request<API.CommonResp<string>>('/api/Dept/update', {
    method: 'POST',
    data: dt,
  });
}

/** 删除规则 DELETE /api/rule */
export async function removeRule(options?: { [key: string]: any }) {
  return request<Record<string, any>>('/api/rule', {
    method: 'POST',
    data: {
      method: 'delete',
      ...(options || {}),
    },
  });
}

export async function waringHistory(d: any) {
  return request<API.CommonResp<string>>('/api/signal-waring-config/query-row', {
    method: 'POST',
    data: d,
  });
}

export async function scriptWaringHistory(d: any) {
  return request<API.CommonResp<string>>('/api/signal-delay-waring/query-row', {
    method: 'POST',
    data: d,
  });
}

enum MessageType {
  // 计划开始通知
  StartNotification = 1, // 计划临期通知
  DueSoonNotification = 2, // 计划到期通知
  DueNotification = 3, // 生产开始通知
  ProductionStartNotification = 4, // 生产完成通知
  ProductionCompleteNotification = 5, // 维修通知
  MaintenanceNotification = 6, // 维修开始通知
  MaintenanceStartNotification = 7, // 维修结束通知
  MaintenanceEndNotification = 8, // SIM卡超时通知
  SimCardExpireTime = 9, // 设备掉线通知
  DeviceOffMessage = 10,
}

export function cc() {
  const maap = [];
  maap.push({ label: getMessageTypeDescription(1), value: 1 });
  maap.push({ label: getMessageTypeDescription(2), value: 2 });
  maap.push({ label: getMessageTypeDescription(3), value: 3 });
  maap.push({ label: getMessageTypeDescription(4), value: 4 });
  maap.push({ label: getMessageTypeDescription(5), value: 5 });
  maap.push({ label: getMessageTypeDescription(6), value: 6 });
  maap.push({ label: getMessageTypeDescription(7), value: 7 });
  maap.push({ label: getMessageTypeDescription(8), value: 8 });
  maap.push({ label: getMessageTypeDescription(9), value: 9 });
  maap.push({ label: getMessageTypeDescription(10), value: 10 });
  return maap;
}

// 获取枚举值对应的中文描述
export function getMessageTypeDescription(type?: number): string {
  if (type == MessageType.StartNotification) {
    return '计划开始通知';
  } else if (type == MessageType.DueSoonNotification) {
    return '计划临期通知';
  } else if (type == MessageType.DueNotification) {
    return '计划到期通知';
  } else if (type == MessageType.ProductionStartNotification) {
    return '生产开始通知';
  } else if (type == MessageType.ProductionCompleteNotification) {
    return '生产完成通知';
  } else if (type == MessageType.MaintenanceNotification) {
    return '维修通知';
  } else if (type == MessageType.MaintenanceStartNotification) {
    return '维修开始通知';
  } else if (type == MessageType.MaintenanceEndNotification) {
    return '维修结束通知';
  } else if (type == MessageType.SimCardExpireTime) {
    return 'SIM卡超时通知';
  } else if (type == MessageType.DeviceOffMessage) {
    return '设备掉线通知';
  } else {
    return '未知类型';
  }
}
