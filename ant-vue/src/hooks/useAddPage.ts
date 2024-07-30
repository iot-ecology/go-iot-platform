import { getAxiosUrl } from "@/utils/setAxiosConfig";
import { message } from "ant-design-vue";
import { useI18n } from "vue-i18n";

/**
 * 用于处理添加数据的功能。
 *
 * @param apiInfo - 用于指定 API 请求的配置信息。
 * @param formRef - 表单的引用。
 * @param form - 表单数据对象。
 * @param modalVisible - 控制模态框显示状态的引用，设置为 false 会关闭模态框。
 * @param callback - 请求成功后执行的回调函数，可以用于刷新数据或执行其他操作。
 *
 * @returns 一个对象，包含以下两个方法：
 *  - onAddData: 处理数据添加操作。
 *  - onCancelData: 处理取消操作，重置表单数据。
 */

export const useAddPage = (apiInfo: any, formRef: any, form: any, modalVisible: any, callback: () => any) => {
  const { t } = useI18n();
  const onAddData = () => {
    (formRef.value as HTMLFormElement)
      .validate()
      .then(async () => {
        try {
          const { data } = await getAxiosUrl(apiInfo, { ...form });
          if (data.code === 20000) {
            message.success(data.message);
            modalVisible.value = false;
            formRef.value?.resetFields();
            await callback();
          } else {
            message.error(`${t("message.operationFailed")}:${data.data}`);
          }
        } catch (error) {
          console.log(error);
        }
      })
      .catch((e) => {
        console.error(e);
      });
  };

  const onCancelData = () => {
    formRef.value?.resetFields();
  };

  return {
    onAddData,
    onCancelData,
  };
};
