import "@/styles/main.scss";

import { createApp } from "vue";
import PrimeVue from "primevue/config";
import Aura from "@primeuix/themes/aura";
import { createPinia } from "pinia";
import ToastService from "primevue/toastservice";

import App from "./App.vue";
import { router } from "./router";
import setupAxiosInterceptors from "./plugins/axios";

const pinia = createPinia();

setupAxiosInterceptors();

const app = createApp(App);
app.use(pinia);
app.use(router);
app.use(PrimeVue, {
  theme: {
    preset: Aura,
  },
});
app.use(ToastService);

app.mount("#app");
