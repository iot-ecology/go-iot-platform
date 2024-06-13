import { defineStore } from "pinia";

interface PermissionStore {
  openValidator: boolean; // 是否开启权限拦截
  btnPermissionList: string[]; // 当前路由下按钮权限集合
}

const isDev = import.meta.env.MODE === "development";

export const usePermissionStore = defineStore({
  id: "permission",
  state: (): PermissionStore => {
    return {
      openValidator: !isDev,
      btnPermissionList: [],
    };
  },
});
