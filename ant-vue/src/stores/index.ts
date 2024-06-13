import { createPinia } from "pinia";
import piniaPluginPersist from "pinia-plugin-persistedstate";
const store = createPinia();
store.use(piniaPluginPersist);
export default store;
