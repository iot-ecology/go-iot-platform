import { createApp } from "vue";
import { RecycleScroller } from "vue-virtual-scroller";
import "normalize.css";
import "virtual:svg-icons-register";
import "./permission";
import "./styles/index.less";
import "vue-virtual-scroller/dist/vue-virtual-scroller.css";

import App from "./App.vue";
import directiveRegister from "./directives";
import VxeTable from "./plugins/vxe-table";
import router from "./router";
import store from "./stores";
import i18n from './i18n';

const app = createApp(App);

app.use(router);
app.use(store);
app.component("RecycleScroller", RecycleScroller);
app.use(directiveRegister);
app.use(VxeTable);
app.use(i18n);

app.mount("#app");
