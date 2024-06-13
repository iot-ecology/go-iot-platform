import type { AxiosRequestConfig } from "axios";

// 请求不需要取消的
const notCancel: string[] = [];

// 请求需要取消 但根据url+method取消 默认是url+method+参数
const specialCancel: string[] = [];

// 是否需要取消请求
export const isNeedCancel = (url: string): boolean => {
  return !notCancel.includes(url);
};

// 获取请求体
export const getRequestBody = (config: AxiosRequestConfig): string => {
  const partBody = (config.url as string) + "&" + (config.method as string);
  const allBody = partBody + "&" + JSON.stringify(config.params || config.data);
  const specialFlag = specialCancel.includes(config.url ?? "");
  const requestBody = specialFlag ? partBody : allBody;
  return requestBody;
};

// 取消请求
export const cancelRequest = (config: AxiosRequestConfig, requestArr: any[]): void => {
  const requestBody = getRequestBody(config);
  for (const request of requestArr) {
    if (request.requestBody === requestBody) {
      request.cancelFn();
      const index = requestArr.findIndex((item) => item.requestBody === requestBody);
      requestArr.splice(index, 1);
    }
  }
};
