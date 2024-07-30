import { getAxiosUrl } from "@/utils/setAxiosConfig";

/**
 * 用于处理查询数据的功能。
 *
 * @param apiInfo - 用于指定 API 请求的配置信息。
 * @param pagination - 分页数据
 * @param list - 数据列表。
 *
 * @returns 一个对象，包含以下一个方法：
 *  - pageList: 处理数据查询操作。
 */

export const useQueryPage = (apiInfo: any, pagination: any, list: any) => {
  const pageList = async () => {
    const { data } = await getAxiosUrl(apiInfo, { page: pagination.current, page_size: pagination.pageSize });
    pagination.total = data.data?.total || 0;
    list.value = data.data.data?.map((item: any, index: number) => ({
      key: index,
      ...item,
    }));
  };

  return {
    pageList,
  };
};
