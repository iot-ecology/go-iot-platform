import { defineStore } from "pinia";

export const useRouterNameStore = defineStore({
  id: "routerPath",
  state: () => {
    return {
      routerParent:"sub1",
      routerPath: "/mqtt-management",
    };
  },
  getters: {},
  actions: {
    getRouterName() {
      return this.routerPath;
    },
    setRouterName(item: string | any) {
      this.routerPath = item;
    },
    setRouterParent(item: string | any) {
      this.routerParent = item;
    },
  },
  persist: [
    {
      key: "project_template_router",
      paths: ["routerParent","routerPath"],
      storage: localStorage,
    },
  ],
});
