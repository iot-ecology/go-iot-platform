import type { DirectiveBinding } from "vue";
import { message } from "ant-design-vue";
import { useClipboard } from "@vueuse/core";

interface CopyEl extends HTMLElement {
  copyStr: string;
  copyHandler: () => void;
}

export default {
  mounted(el: CopyEl, binding: DirectiveBinding) {
    el.copyStr = binding.value;
    el.copyHandler = () => {
      const { copy } = useClipboard({ legacy: true });
      copy(el.copyStr)
        .then(() => {
          message.success("复制成功");
        })
        .catch((error) => {
          message.error("复制失败");
          console.log("复制失败原因", error);
        });
    };
    el.addEventListener("click", el.copyHandler);
  },
  updated(el: CopyEl, binding: DirectiveBinding) {
    el.copyStr = binding.value;
  },
  unmounted(el: CopyEl) {
    el.removeEventListener("click", el.copyHandler);
  },
};
