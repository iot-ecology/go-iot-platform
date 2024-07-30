import axios from "axios";
const url: string = import.meta.env.VITE_APP_API_URL;
export async function MqttPage(params: any) {
  return await axios.get(`${url}/mqtt/page`, { params });
}

export async function MqttUpdate(data: any) {
  return await axios.post(`${url}/mqtt/update`, data);
}

export async function MqttCreate(data: any) {
  return await axios.post(`${url}/mqtt/create`, data);
}

export async function MqttSetScript(data: any) {
  return await axios.post(`${url}/mqtt/set-script`, data);
}

export async function MqttCheckScript(data: any) {
  return await axios.post(`${url}/mqtt/check-script`, data);
}
export async function MqttDelete(id: any) {
  // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
  return await axios.post(`${url}/mqtt/delete/${id}`);
}
export async function MqttStart(params: any) {
  return await axios.get(`${url}/mqtt/start`, { params });
}
export async function MqttStop(params: any) {
  return await axios.get(`${url}/mqtt/stop`, { params });
}
export async function MqttSend(data: any) {
  return await axios.post(`${url}/mqtt/send?id=${data.client_id}`, data);
}

export async function MqttNodeUsingStatus(params: any) {
  return await axios.get(`${url}/mqtt/node-using-status`, { params });
}

// 信号
export async function SignalCreate(data: any) {
  return await axios.post(`${url}/signal/create`, data);
}

export async function SignalPage(params: any) {
  return await axios.get(`${url}/signal/page`, { params });
}

export async function SignalUpdate(data: any) {
  return await axios.post(`${url}/signal/update`, data);
}

export async function SignalDelete(id: string) {
  return await axios.post(`${url}/signal/delete/${id}`);
}

// 信息号报警配置
export async function SignalWaringConfigPage(params: any) {
  return await axios.get(`${url}/signal-waring-config/page`, { params });
}

export async function SignalWaringConfigCreate(data: any) {
  return await axios.post(`${url}/signal-waring-config/create`, data);
}

export async function SignalWaringConfigUpdate(data: any) {
  return await axios.post(`${url}/signal-waring-config/update`, data);
}

export async function SignalWaringConfigDelete(id: any) {
  // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
  return await axios.post(`${url}/signal-waring-config/delete/${id}`);
}

export async function SignalWaringConfigQueryRow(data: any) {
  return await axios.post(`${url}/signal-waring-config/query-row`, data);
}
// data
export async function QueryInfluxdb(data: any) {
  return await axios.post(`${url}/query/influxdb`, data);
}

export async function QueryStrInfluxdb(data: any) {
  return await axios.post(`${url}/query/str-influxdb`, data);
}

export async function DashboardPage(params: any) {
  return await axios.get(`${url}/dashboard/page`, { params });
}

export async function DashboardCreate(data: any) {
  return await axios.post(`${url}/dashboard/create`, data);
}

export async function DashboardDelete(id: string) {
  return await axios.post(`${url}/dashboard/delete/${id}`);
}

export async function DashboardId(id: any) {
  return await axios.get(`${url}/dashboard/${id}`);
}

export async function DashboardUpdate(data: any) {
  return await axios.post(`${url}/dashboard/update`, data);
}
// 计算规则
export async function CalcRulePage(params: any) {
  return await axios.get(`${url}/calc-rule/page`, { params });
}

export async function CalcRuleCreate(data: any) {
  return await axios.post(`${url}/calc-rule/create`, data);
}

export async function CalcRuleUpdate(data: any) {
  return await axios.post(`${url}/calc-rule/update`, data);
}

export async function CalcParamStart(id: string) {
  return await axios.post(`${url}/calc-rule/start/${id}`);
}

export async function CalcParamStop(id: string) {
  return await axios.post(`${url}/calc-rule/stop/${id}`);
}

export async function CalcParamMock(data: any) {
  return await axios.post(`${url}/calc-rule/mock`, data);
}
export async function CalcParamRd(params: any) {
  return await axios.get(`${url}/calc-rule/rd`, { params });
}

// 计算参数
export async function CalcParamPage(params: any) {
  return await axios.get(`${url}/calc-param/page`, { params });
}

export async function CalcParamCreate(data: any) {
  return await axios.post(`${url}/calc-param/create`, data);
}

export async function CalcParamDelete(id: string) {
  return await axios.post(`${url}/calc-param/delete/${id}`);
}

export async function CalcParamUpdate(data: any) {
  return await axios.post(`${url}/calc-param/update`, data);
}

// 脚本报警
export async function SignalDelayWaringPage(params: any) {
  return await axios.get(`${url}/signal-delay-waring/page`, { params });
}
export async function SignalDelayWaringCreate(data: any) {
  return await axios.post(`${url}/signal-delay-waring/create`, data);
}
export async function SignalDelayWaringDelete(id: string) {
  return await axios.post(`${url}/signal-delay-waring/delete/${id}`);
}
export async function SignalDelayWaringUpdate(data: any) {
  return await axios.post(`${url}/signal-delay-waring/update`, data);
}

export async function SignalDelayWaringGenParam(id: number) {
  return await axios.post(`${url}/signal-delay-waring/GenParam/${id}`);
}

export async function SignalDelayWaringMock(id: number) {
  return await axios.post(`${url}/signal-delay-waring/Mock/${id}`);
}

export async function SignalDelayWaringQueryRow(data: any) {
  return await axios.post(`${url}/signal-delay-waring/query-row`, data);
}

// 脚本报警参数
export async function SignalDelayWaringParamPage(params: any) {
  return await axios.get(`${url}/signal-delay-waring-param/page`, { params });
}
export async function SignalDelayWaringParamCreate(data: any) {
  return await axios.post(`${url}/signal-delay-waring-param/create`, data);
}
export async function SignalDelayWaringParamDelete(id: string) {
  return await axios.post(`${url}/signal-delay-waring-param/delete/${id}`);
}
export async function SignalDelayWaringParamUpdate(data: any) {
  return await axios.post(`${url}/signal-delay-waring-param/update`, data);
}
//部门
export async function DeptPage(params: any) {
  return await axios.get(`${url}/Dept/page`, { params });
}
export async function DeptCreate(data: any) {
  return await axios.post(`${url}/Dept/create`, data);
}
export async function DeptDelete(id: any) {
  return await axios.post(`${url}/Dept/delete/${id}`);
}
export async function DeptUpdate(data: any) {
  return await axios.post(`${url}/Dept/update`, data);
}
export async function DeptSubs(params: any) {
  return await axios.get(`${url}/Dept/subs`, { params });
}
//设备详情
export const DEVICEINFO_PAGE = "/DeviceInfo/page";
export const DEVICEINFO_CREATE = "/DeviceInfo/create";
export const DEVICEINFO_DELETE = "/DeviceInfo/delete";
export async function DeviceInfoPage(params: any) {
  return await axios.get(`${url}${DEVICEINFO_PAGE}`, { params });
}
export async function DeviceInfoCreate(data: any) {
  return await axios.post(`${url}${DEVICEINFO_CREATE}`, data);
}
export async function DeviceInfoDelete(id: any) {
  return await axios.post(`${url}${DEVICEINFO_DELETE}/${id}`);
}
export async function DeviceInfoBindMqtt(data: any) {
  return await axios.post(`${url}/DeviceInfo/BindMqtt`, data);
}
export async function DeviceInfoQueryBindMqtt(params: any) {
  return await axios.get(`${url}/DeviceInfo/QueryBindMqtt`, { params });
}
export async function DeviceInfoUpdate(data: any) {
  return await axios.post(`${url}/DeviceInfo/update`, data);
}

// 设备组
export const DEVICE_GROUP_PAGE = "/device_group/page";
export const DEVICE_GROUP_CREATE = "/device_group/create";
export const DEVICE_GROUP_UPDATE = "/device_group/update";
export const DEVICE_GROUP_DELETE = "/device_group/delete";
export const DEVICE_GROUP_BIND_DEVICE = "/device_group/bind_device";
export const DEVICE_GROUP_QUERY_BIND_DEVICE = "/device_group/query_bind_device";
export async function MessageListPage(params: any) {
  return await axios.get(`${url}${DEVICE_GROUP_PAGE}`, { params });
}
export async function DeviceGroupsCreate(data: any) {
  return await axios.post(`${url}${DEVICE_GROUP_CREATE}`, data);
}
export async function DeviceGroupsUpdate(data: any) {
  return await axios.post(`${url}${DEVICE_GROUP_UPDATE}`, data);
}
export async function DeviceGroupsDelete(id: any) {
  return await axios.post(`${url}${DEVICE_GROUP_DELETE}/${id}`);
}
export async function DeviceGroupsBindDevice(data: any) {
  return await axios.post(`${url}${DEVICE_GROUP_BIND_DEVICE}`, data);
}
export async function DeviceGroupsQueryBindDevice(params: any) {
  return await axios.get(`${url}${DEVICE_GROUP_QUERY_BIND_DEVICE}`, { params });
}

// 飞书机器人
export const FEISHU_ID_PAGE = "/FeiShuId/page";
export const FEISHU_ID_CREATE = "/FeiShuId/create";
export const FEISHU_ID_UPDATE = "/FeiShuId/update";
export const FEISHU_ID_DELETE = "/FeiShuId/delete";
export async function FeiShuPage(params: any) {
  return await axios.get(`${url}${FEISHU_ID_PAGE}`, { params });
}
export async function FeiShuCreate(data: any) {
  return await axios.post(`${url}${FEISHU_ID_CREATE}`, data);
}
export async function FeiShuUpdate(data: any) {
  return await axios.post(`${url}${FEISHU_ID_UPDATE}`, data);
}
export async function FeiShuDelete(id: any) {
  return await axios.post(`${url}${FEISHU_ID_DELETE}/${id}`);
}

// 钉钉机器人
export const DingDing_PAGE = "/DingDing/page";
export const DingDing_CREATE = "/DingDing/create";
export const DingDing_UPDATE = "/DingDing/update";
export const DingDing_DELETE = "/DingDing/delete";
export async function DingDingPage(params: any) {
  return await axios.get(`${url}${DingDing_PAGE}`, { params });
}
export async function DingDingCreate(data: any) {
  return await axios.post(`${url}${DingDing_CREATE}`, data);
}
export async function DingDingUpdate(data: any) {
  return await axios.post(`${url}${DingDing_UPDATE}`, data);
}
export async function DingDingDelete(id: any) {
  return await axios.post(`${url}${DingDing_DELETE}/${id}`);
}

// CassandraTransmit
export const CASSANDRA_PAGE = "/CassandraTransmit/page";
export const CASSANDRA_CREATE = "/CassandraTransmit/create";
export const CASSANDRA_UPDATE = "/CassandraTransmit/update";
export const CASSANDRA_DELETE = "/CassandraTransmit/delete";
export async function CassandraTransmitPage(params: any) {
  return await axios.get(`${url}${CASSANDRA_PAGE}`, { params });
}
export async function CassandraTransmitCreate(data: any) {
  return await axios.post(`${url}${CASSANDRA_CREATE}`, data);
}
export async function CassandraTransmitUpdate(data: any) {
  return await axios.post(`${url}${CASSANDRA_UPDATE}`, data);
}
export async function CassandraTransmitDelete(id: any) {
  return await axios.post(`${url}${CASSANDRA_DELETE}/${id}`);
}

// ClickhouseTransmit
export const CLICK_HOUSE_PAGE = "/ClickhouseTransmit/page";
export const CLICK_HOUSE_CREATE = "/ClickhouseTransmit/create";
export const CLICK_HOUSE_UPDATE = "/ClickhouseTransmit/update";
export const CLICK_HOUSE_DELETE = "/ClickhouseTransmit/delete";
export async function ClickhouseTransmitPage(params: any) {
  return await axios.get(`${url}${CLICK_HOUSE_PAGE}`, { params });
}
export async function ClickhouseTransmitCreate(data: any) {
  return await axios.post(`${url}${CLICK_HOUSE_CREATE}`, data);
}
export async function ClickhouseTransmitUpdate(data: any) {
  return await axios.post(`${url}${CLICK_HOUSE_UPDATE}`, data);
}
export async function ClickhouseTransmitDelete(id: any) {
  return await axios.post(`${url}${CLICK_HOUSE_DELETE}/${id}`);
}

// InfluxdbTransmit
export const INFLUXDB_PAGE = "/InfluxdbTransmit/page";
export const INFLUXDB_CREATE = "/InfluxdbTransmit/create";
export const INFLUXDB_UPDATE = "/InfluxdbTransmit/update";
export const INFLUXDB_DELETE = "/InfluxdbTransmit/delete";
export async function InfluxdbTransmitPage(params: any) {
  return await axios.get(`${url}${INFLUXDB_PAGE}`, { params });
}
export async function InfluxdbTransmitCreate(data: any) {
  return await axios.post(`${url}${INFLUXDB_CREATE}`, data);
}
export async function InfluxdbTransmitUpdate(data: any) {
  return await axios.post(`${url}${INFLUXDB_UPDATE}`, data);
}
export async function InfluxdbTransmitDelete(id: any) {
  return await axios.post(`${url}${INFLUXDB_DELETE}/${id}`);
}

// MongoTransmit
export const MONGO_PAGE = "/MongoTransmit/page";
export const MONGO_CREATE = "/MongoTransmit/create";
export const MONGO_UPDATE = "/MongoTransmit/update";
export const MONGO_DELETE = "/MongoTransmit/delete";
export async function MongoTransmitPage(params: any) {
  return await axios.get(`${url}${MONGO_PAGE}`, { params });
}
export async function MongoTransmitCreate(data: any) {
  return await axios.post(`${url}${MONGO_CREATE}`, data);
}
export async function MongoTransmitUpdate(data: any) {
  return await axios.post(`${url}${MONGO_UPDATE}`, data);
}
export async function MongoTransmitDelete(id: any) {
  return await axios.post(`${url}${MONGO_DELETE}/${id}`);
}

// MySQLTransmit
export const MYSQL_PAGE = "/MySQLTransmit/page";
export const MYSQL_CREATE = "/MySQLTransmit/create";
export const MYSQL_UPDATE = "/MySQLTransmit/update";
export const MYSQL_DELETE = "/MySQLTransmit/delete";
export async function MySQLTransmitPage(params: any) {
  return await axios.get(`${url}${MYSQL_PAGE}`, { params });
}
export async function MySQLTransmitCreate(data: any) {
  return await axios.post(`${url}${MYSQL_CREATE}`, data);
}
export async function MySQLTransmitUpdate(data: any) {
  return await axios.post(`${url}${MYSQL_UPDATE}`, data);
}
export async function MySQLTransmitDelete(id: any) {
  return await axios.post(`${url}${MYSQL_DELETE}/${id}`);
}

//角色
export async function RolePage(params: any) {
  return await axios.get(`${url}/Role/page`, { params });
}
export async function RoleList() {
  return await axios.get(`${url}/Role/list`);
}
export async function RoleCreate(data: any) {
  return await axios.post(`${url}/Role/create`, data);
}
export async function RoleUpdate(data: any) {
  return await axios.post(`${url}/Role/update`, data);
}
export async function RoleDelete(id: any) {
  return await axios.post(`${url}/Role/delete/${id}`);
}

//用户
export async function UserPage(params: any) {
  return await axios.get(`${url}/User/page`, { params });
}
export async function UserCreate(data: any) {
  return await axios.post(`${url}/User/create`, data);
}
export async function UserUpdate(data: any) {
  return await axios.post(`${url}/User/update`, data);
}
export async function UserDelete(id: any) {
  return await axios.post(`${url}/User/delete/${id}`);
}
export async function UserBindRole(data: any) {
  return await axios.post(`${url}/User/BindRole`, data);
}
export async function UserQueryBindRole(params: any) {
  return await axios.get(`${url}/User/QueryBindRole?user_id=${params}`);
}

//产品
export async function ProductPage(params: any) {
  return await axios.get(`${url}/product/page`, { params });
}
export async function ProductCreate(data: any) {
  return await axios.post(`${url}/product/create`, data);
}
export async function ProductUpdate(data: any) {
  return await axios.post(`${url}/product/update`, data);
}
export async function ProductDelete(id: any) {
  return await axios.post(`${url}/product/delete/${id}`);
}
//sim卡
export async function SimCardPage(params: any) {
  return await axios.get(`${url}/SimCard/page`, { params });
}
export async function SimCardCreate(data: any) {
  return await axios.post(`${url}/SimCard/create`, data);
}
export async function SimCardUpdate(data: any) {
  return await axios.post(`${url}/SimCard/update`, data);
}
export async function SimCardDelete(id: any) {
  return await axios.post(`${url}/SimCard/delete/${id}`);
}
