import { ref } from "vue";
import { theme } from "ant-design-vue";

import { useThemeStore } from "@/stores/theme";
import { antDesignVueDark, cssVariableDark } from "@/theme/dark";
import { antDesignVueLight, cssVariableLight } from "@/theme/light";

export const useTheme = () => {
  const themeStore = useThemeStore();
  const { defaultAlgorithm, darkAlgorithm } = theme;
  const antDesignVueTheme = ref<Record<string, any>>({});
  const cssVariable = ref<Record<string, string>>({});

  // 切换主题
  const switchTheme = () => {
    antDesignVueTheme.value = themeStore.themeMode === "dark" ? antDesignVueDark : antDesignVueLight;
    antDesignVueTheme.value = {
      ...antDesignVueTheme.value,
      algorithm: themeStore.themeMode === "dark" ? darkAlgorithm : defaultAlgorithm,
    };
    cssVariable.value = themeStore.themeMode === "dark" ? cssVariableDark : cssVariableLight;
    addCssToRoot();
  };

  // 添加css属性到:root
  const addCssToRoot = () => {
    let cssStr = "";
    for (const key in cssVariable.value) {
      if (Object.hasOwnProperty.call(cssVariable.value, key)) {
        const value = cssVariable.value[key];
        cssStr += `${key}: ${value};`;
      }
    }
    let style = document.querySelector("style#theme-variable") as HTMLStyleElement;
    if (!style) {
      style = document.createElement("style");
      style.id = "theme-variable";
      document.getElementsByTagName("head")[0].append(style);
    }
    style.innerText = `:root{${cssStr}}`;
  };

  switchTheme();
  themeStore.$subscribe(switchTheme);

  return {
    antDesignVueTheme,
    cssVariable,
  };
};
