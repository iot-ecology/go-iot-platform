# ☀️ 初衷

每当新开一个项目时，都要浪费很长时间去搞一个项目的基建  
久而久之，我就在想，为什么不搞一个模板出来，把每个项目中需要用到的工具和组件给抽离出来  
然后这个项目就应运而生了，当然目前只是一个雏形，后面会在空闲时间逐渐完善  
如果你也有跟我一样的想法，可以提交PR，咱们一起维护它

# ✨ 基建

- [Vue](https://cn.vuejs.org/)
- [Vite](https://vitejs.cn/)
- [Typescript](https://www.typescriptlang.org/zh/)
- [Ant-Design-Vue](https://next.antdv.com/docs/vue/introduce-cn)
- [Pinia](https://pinia.web3doc.top/introduction.html)
- [Axios](https://www.axios-http.cn/docs/intro)
- [Dayjs](https://dayjs.fenxianglu.cn/category/)
- [Lodash](https://www.lodashjs.com/)

# ⛵ 规范

- 使用 [ESLint](https://zh-hans.eslint.org/) + [Prettier](https://www.prettier.cn/docs/options.html) 的形式去做代码规范，并使用了`standard`标准  
- 使用 [husky](https://github.com/typicode/husky) + [lint-staged](https://github.com/okonet/lint-staged) 做提交时的拦截和自动修复
  - ⚠️`mac`用户注意`.husky`文件夹下的执行权限 可使用`ls -a .husky`查看
- 使用 [commitlint](https://github.com/conventional-changelog/commitlint) 来限制提交信息的规范  

