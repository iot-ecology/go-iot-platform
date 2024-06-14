<template>
  <div class="app-layout">
    <!-- 菜单 -->
    <div class="layout-menu">
      <a-menu v-model:openKeys="menuState.openKeys" v-model:selectedKeys="menuState.selectedKeys" theme="dark" mode="inline" @click="handleMenuClick">
        <template v-for="menu in menuState.menus" :key="menu.path">
          <!-- 隐藏的路由 -->
          <template v-if="menu.meta?.hidden" />
          <!-- 单路由 -->
          <template v-else-if="menu.children?.length">
            <a-menu-item :key="menu.path">
              {{ menu.meta?.title }}
            </a-menu-item>
          </template>
          <!-- 有子路由 -->
          <!--          <template v-else>-->
          <!--            <a-sub-menu :key="menu.path" :title="menu.meta.title">-->
          <!--              <a-menu-item v-for="subMenu in menu.children" :key="subMenu.path">-->
          <!--                {{ subMenu.meta?.title }}-->
          <!--              </a-menu-item>-->
          <!--            </a-sub-menu>-->
          <!--          </template>-->
        </template>
      </a-menu>
    </div>
    <!-- 视图 -->
      <div class="layout-view">
        <div class="header">
          <a-button v-if="language==='en'" @click="onChoiceLanguage('zhCHS')">中文</a-button>
          <a-button v-else @click="onChoiceLanguage('en')">English</a-button>
        </div>
        <div class="layout-view__container">
          <router-view />
        </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, watch,ref } from "vue";
import { useRouter } from "vue-router";
import type { MenuClickEventHandler } from "ant-design-vue/lib/menu/src/interface";
import { useI18n } from 'vue-i18n';
import { useRouterNameStore } from "@/stores/routerPath.ts";

const router = useRouter();
const routes = router.options.routes;
const routerStore = useRouterNameStore();
const language = ref('zhCHS')
const { locale } = useI18n();
// 菜单状态
interface MenuState {
  menus: any;
  openKeys: string[];
  selectedKeys: string[];
}
const menuState = reactive<MenuState>({
  menus: routes,
  openKeys: [routerStore.routerPath],
  selectedKeys: [routerStore.routerPath],
});
const count = computed(() =>routerStore.routerPath);
const onChoiceLanguage = async (lang: string)=>{
  language.value = lang
  locale.value = lang;
}
// 使用 watch 监听 count 的变化
watch(count, (newCount) => {
  if(newCount) {
    menuState.openKeys = [newCount]
    menuState.selectedKeys = [newCount]
  }
});

// 菜单点击回调
const handleMenuClick: MenuClickEventHandler = (menuInfo) => {
  const key = menuInfo.keyPath?.[0];
  routerStore.setRouterName(key)
  const menuKey = menuInfo.key as string;
  const isUrl = menuKey.match(/^https?/);
  if (typeof key === "string") {
    menuState.openKeys = [key];
  }
  if (isUrl) {
    window.open(menuKey);
  } else {
    router.push(menuKey);
  }
};
</script>

<style lang="less" scoped>
.app-layout {
  width: 100%;
  height: 100%;
  display: flex;

  .layout-menu {
    flex: none;
    width: 180px;
    position: relative;

    .ant-menu {
      height: 100%;
      overflow: hidden;
    }
  }

  .layout-view {
    flex: 1;
    box-sizing: border-box;
    min-width: 0;
    overflow: auto;
    .header {
      height: 48px;
      background: white;
      padding: 0 16px;
      display: flex;
      justify-content: right;
      align-items: center;
    }
    .layout-view__container {
      width: 100%;
      min-height: 100%;
      padding: 0 16px 16px;
      border-radius: @containerBorderRadius;
      background-color: @containerBgColor;
      overflow: hidden;
    }
  }
}
</style>
