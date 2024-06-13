import NProgress from "nprogress";
import "nprogress/nprogress.css";

import router from "./router";
NProgress.configure({ showSpinner: false });

router.beforeEach((to, from, next) => {
  // 开启进度条
  NProgress.start();
  if (to.name && !from.name && import.meta.env.PROD) {
    // 弹出输入框
    const userInput = window.prompt("请输入一个值:");
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
      document.title = `${String(to.meta.title)}`;
    } else {
      document.title = "客户端管理";
    }
    next();
  }
});

router.afterEach(() => {
  // 关闭进度条
  NProgress.done();
});
