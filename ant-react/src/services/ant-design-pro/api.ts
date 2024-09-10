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
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.CommonPage<API.MessageListItem>>('/api/MessageList/page', {
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

export async function deptList() {
  return request<API.CommonResp<API.DeptListItem[]>>('/api/Dept/list', {
    method: 'GET',
  });
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
  const request1 = await request<API.SimListItem>('/api/SimCard/page', {
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
export async function productPage(
  params: {
    /** 当前的页码 */
    current?: number /** 页面的容量 */;
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  const request1 = await request<API.ProductItem>('/api/product/page', {
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
  const request1 = await request<API.DeviceGroupItem>('/api/device_group/page', {
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
  return request<API.RuleListItem>('/api/User/create', {
    method: 'POST',
    data: {
      method: 'post',
      ...(options || {}),
    },
  });
}
export async function addSim(options?: { [key: string]: any }) {
  return request<API.RuleListItem>('/api/SimCard/create', {
    method: 'POST',
    data: {
      method: 'post',
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
