import NProgress from "nprogress";
import "nprogress/nprogress.css";

import router from "./router";
// import {getMetaTitle} from "@/utils/i18n.ts";
import i18n  from "@/i18n";

NProgress.configure({ showSpinner: false });
router.beforeEach((to, from, next) => {
  // 开启进度条
  NProgress.start();
  if (to.name && !from.name && import.meta.env.PROD) {
    // 弹出输入框
    const userInput = window.prompt(i18n.global.t('message.pleaseEnter'));
    if (!userInput) {
      next(false);
    } else {
      if (userInput === import.meta.env.VITE_LOGIN) {
        next();
      } else {
        next(false);
      }
    }
  } else {
    // 替换标题
    if (to.meta.title) {
      document.title = `${String(i18n.global.t(to.meta.title))}`;
    } else {
      document.title = i18n.global.t('message.mqttManagement');
    }
    next();
  }
});

router.afterEach(() => {
  // 关闭进度条
  NProgress.done();
});
