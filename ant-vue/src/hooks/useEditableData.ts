import { cloneDeep } from "lodash-es";
import { reactive, UnwrapRef } from "vue";

/**
 * 自定义组合式函数，用于管理可编辑数据项。
 *
 * @param list - 包含数据项的响应式引用列表，每个数据项应有一个唯一标识符（如 `key`）。
 *
 * @returns 一个对象，包含以下属性：
 *   - `editableData`: 一个响应式对象，用于存储当前正在编辑的数据项。每个数据项以 `key` 为属性存储。
 *   - `edit`: 一个函数，用于启动对指定 `key` 对应的数据项的编辑。将数据项的深拷贝存储在 `editableData` 中。
 *   - `cancel`: 一个函数，用于取消对指定 `key` 对应的数据项的编辑。将 `editableData` 中的相关条目删除。
 */

export const useEditableData = (list: any) => {
  const editableData: UnwrapRef<Record<string, any>> = reactive({});
  const edit = (key: string) => {
    editableData[key] = cloneDeep(list.value.filter((item) => key === item.key)[0]);
  };
  const cancel = (key: string) => {
    delete editableData[key];
  };
  return {
    editableData,
    edit,
    cancel,
  };
};
