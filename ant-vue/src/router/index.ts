import type { RouteRecordRaw } from "vue-router";
import { createRouter, createWebHistory } from "vue-router";

import single from "./single";

export const routes: RouteRecordRaw[] = [
  {
    path: "/",
    name: "Root",
    redirect: "/icon-preview",
    meta: { hidden: true },
  },
  ...single,
];

const router = createRouter({
  history: createWebHistory(import.meta.env.VITE_BASE_URL),
  routes,
});

export default router;
