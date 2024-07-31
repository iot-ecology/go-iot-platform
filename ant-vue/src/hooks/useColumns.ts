import { ref, watch } from "vue";
import { useI18n } from "vue-i18n";

// 定义列的类型
interface Column {
  title: string;
  dataIndex: string;
}

/**
 * 自定义组合式函数，用于处理表格列的国际化显示。
 *
 * @param initialList - 初始的列配置列表，每个列配置对象包含 `title` 和 `dataIndex` 属性。
 *   - `title`: 列的标题，需要经过国际化处理。
 *   - `dataIndex`: 列的数据索引，用于在数据源中查找对应的值。
 *
 * @returns 一个对象，包含以下属性：
 *   - `columns`: 一个响应式引用，表示处理后的列配置列表。列标题经过国际化处理，并且当语言发生变化时，列配置会自动更新。
 *
 * @example
 * const { columns } = useColumns([
 *   { title: 'name', dataIndex: 'name' },
 *   { title: 'age', dataIndex: 'age' }
 * ]);
 *
 *
 */

export const useColumns = (initialList: Column[]) => {
  const { t, locale } = useI18n();
  function setList() {
    let list = initialList.map((item) => {
      return {
        title: t(item.title),
        dataIndex: item.dataIndex,
      };
    });
    return list;
  }
  const columns = ref(setList());
  watch(locale, (newValue) => {
    columns.value = setList();
  });

  return { columns };
};
