import { createSSRApp } from "vue";
import App from "./App.vue";
import 'vant/lib/index.css';
import { Button } from 'vant';
import { Tabbar, TabbarItem } from 'vant';

export function createApp() {
  const app = createSSRApp(App);
  app.use(Button);
  app.use(Tabbar);
  app.use(TabbarItem);
  return {
    app,
  };
}
