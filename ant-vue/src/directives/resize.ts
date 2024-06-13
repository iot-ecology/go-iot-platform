import type { DirectiveBinding } from "vue";

interface ResizeEl extends HTMLElement {
  _observer: InstanceType<typeof ResizeObserver>;
}

export default {
  mounted(el: ResizeEl, binding: DirectiveBinding) {
    let width = el.clientWidth;
    let height = el.clientHeight;

    el._observer = new ResizeObserver(() => {
      if (el.clientWidth !== width || el.clientHeight !== height) {
        width = el.clientWidth;
        height = el.clientHeight;
        binding.value(width, height);
      }
    });
    el._observer.observe(el);
  },
  unmounted(el: ResizeEl) {
    el._observer?.disconnect();
  },
};
