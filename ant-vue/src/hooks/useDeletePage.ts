import { getAxiosUrl } from "@/utils/setAxiosConfig";
import { message } from "ant-design-vue";

/**
 * 用于处理删除数据的功能。
 *
 * @param apiInfo - 用于指定 API 请求的配置信息。
 * @param callback - 请求成功后执行的回调函数，可以用于刷新数据或执行其他操作。
 *
 * @returns 一个对象，包含以下一个方法：
 *  - deleteById: 处理数据删除操作。
 */

export const useDeletePage = (apiInfo: any, callback: () => any) => {
  const deleteById = async (id: string) => {
    try {
      const { data } = await getAxiosUrl(apiInfo, id);
      if (data.code === 20000) {
        message.success(data.message);
        await callback();
      } else {
        message.success(data.message);
      }
    } catch (error) {
      console.log(error);
    }
  };
  return {
    deleteById,
  };
};
