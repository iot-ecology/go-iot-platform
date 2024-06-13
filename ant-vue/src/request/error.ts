import { message } from "ant-design-vue";
export const errorHandler = (code: number, msg: string) => {
  switch (code) {
    case 401:
      message.error("没有访问权限");
      break;
    case 403:
      message.error("资源不可用");
      break;
    case 404:
      message.error("找不到对应资源");
      break;
    case 500:
      message.error("服务器异常");
      break;
    case 502:
      message.error("服务器异常");
      break;
    default:
      if (typeof msg === "string") {
        message.error(msg);
      } else {
        message.error(JSON.stringify(msg));
      }
      break;
  }
};
