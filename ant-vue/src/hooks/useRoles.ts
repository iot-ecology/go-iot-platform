import { ref, watch } from "vue";
import { useI18n } from "vue-i18n";

/**
 * 自定义组合式函数，用于处理表单验证规则的国际化。
 *
 * @param obj - 一个对象，其中每个键对应一个验证规则数组。每个验证规则对象必须包含 `message` 属性。
 *   - `message`: 验证失败时显示的提示消息，需要经过国际化处理。
 *
 * @returns 一个对象，包含以下属性：
 *   - `rules`: 一个响应式引用，表示处理后的验证规则。每条规则的 `message` 属性经过国际化处理，并且当语言发生变化时，验证规则会自动更新。
 */

export const useRoles = (obj: any) => {
  const { t, locale } = useI18n();
  function setRules() {
    const newRules = {};
    for (const [key, value] of Object.entries(obj)) {
      newRules[key] = value.map((rule: { message: any }) => ({
        ...rule,
        message: t(rule.message), // 使用 t 函数进行翻译
      }));
    }
    return newRules;
  }
  const rules = ref(setRules());
  watch(locale, (newValue) => {
    rules.value = setRules();
  });

  return { rules };
};
