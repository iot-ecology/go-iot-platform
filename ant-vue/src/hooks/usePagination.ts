import { reactive } from "vue";

export const usePagination = () => {
  const pagination = reactive({
    total: 0,
    current: 1,
    pageSize: 10,
    showSizeChanger: true, // 显示每页显示条目数选择器
  });

  return {
    pagination,
  };
};
