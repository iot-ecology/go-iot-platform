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
  return await axios.post(`${url}/mqtt/set-script1`, data);
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

export async function DashboardId(id:number | string ) {
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
