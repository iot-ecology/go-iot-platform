import { getAxiosUrl } from "@/utils/setAxiosConfig";
import { message } from "ant-design-vue";
import { useI18n } from "vue-i18n";

/**
 * 用于处理更新数据的功能。
 *
 * @param apiInfo - 用于指定 API 请求的配置信息。
 * @param list - 数据列表
 * @param editableData - 编辑的信息
 * @param callback - 请求成功后执行的回调函数，可以用于刷新数据或执行其他操作。
 *
 * @returns 一个对象，包含以下一个方法：
 *  - save: 处理数据更新操作。
 */

export const useUpdatePage = (apiInfo: any, list: any, editableData: any, callback: () => any) => {
  const { t } = useI18n();
  const save = async (key: string) => {
    Object.assign(list.value.filter((item: { key: string }) => key === item.key)[0], editableData[key]);
    const params = list.value.filter((item: { key: string }) => key === item.key)[0];
    // eslint-disable-next-line @typescript-eslint/no-dynamic-delete
    delete editableData[key];
    // await DeviceInfoUpdate(data);
    try {
      const { data } = await getAxiosUrl(apiInfo, params);
      if (data.code === 20000) {
        message.success(data.message);
      } else {
        message.error(`${t("message.operationFailed")}:${data.data}`);
      }
    } catch (error) {
      console.log(error);
    }
    await callback();
  };

  return {
    save,
  };
};
