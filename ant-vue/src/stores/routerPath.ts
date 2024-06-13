import { defineStore } from "pinia";

export const useRouterNameStore = defineStore({
  id: "routerPath",
  state: () => {
    return {
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
  },
  persist: [
    {
      key: "project_template_router",
      paths: ["routerPath"],
      storage: localStorage,
    },
  ],
});
