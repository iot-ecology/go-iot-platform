import axios from "axios";

interface ApiInfo {
  path: string;
  type: string;
}

const url: string = import.meta.env.VITE_APP_API_URL;

/**
 * 发送网络请求的通用函数。
 *
 * @param apiInfo - 请求的配置信息对象，包括 API 路径和请求类型。
 *   - `path`: API 的路径部分。
 *   - `type`: HTTP 请求类型，例如 "get"、"post"、"put" 或 "delete"。
 * @param params - 请求的参数。对于 GET 请求，它们会作为查询参数附加到 URL 后；对于 POST 和 PUT 请求，它们会作为请求体发送；对于 DELETE 请求，它们会被附加到 URL。
 *
 * @returns 返回一个 Promise 对象，表示请求的结果。Promise 的解析值是 axios 发起请求后的响应对象。
 */

export async function getAxiosUrl(apiInfo: ApiInfo, params: any) {
  const { type, path } = apiInfo;
  if (type === "get") {
    return await axios[type](`${url}${path}`, { params });
  } else if (type === "delete") {
    return await axios.post(`${url}${path}/${params}`);
  } else {
    return await axios[type](`${url}${path}`, params);
  }
}
