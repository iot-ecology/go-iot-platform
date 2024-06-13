import type { DirectiveBinding } from "vue";

import { hasPermission } from "@/utils/permission";

export default {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    const id = binding.value as string;
    const flag = hasPermission(id);
    if (!flag) {
      el.parentElement?.removeChild(el);
    }
  },
};
