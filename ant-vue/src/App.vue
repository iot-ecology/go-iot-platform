<template>
  <a-config-provider :locale="zhCN" :theme="antDesignVueTheme" :auto-insert-space-in-button="false">
    <router-view />
  </a-config-provider>
</template>

<script setup lang="ts">
import zhCN from "ant-design-vue/es/locale/zh_CN";
import dayjs from "dayjs";
import "dayjs/locale/zh-cn";

import { useTheme } from "@/hooks/useTheme";

const { antDesignVueTheme } = useTheme();

dayjs.locale("zh-cn");
</script>

<style lang="less">
* {
  box-sizing: border-box;
}
html,
body {
  width: 100%;
  height: 100%;
  background-color: @layoutBgColor;
}
#app {
  width: 100%;
  height: 100%;
  min-width: 1280px;
  color: @colorText;
  font-size: @fontSize;
  font-family: @fontFamily;
}
#nprogress .bar {
  background-color: @colorPrimary !important;
}
</style>
