import { createApp } from "vue";
import App from "./App.vue";
import { router } from "./router";
import "./style.css";

import "virtual:uno.css";
import Vue3Toastify, { type ToastContainerOptions } from "vue3-toastify";
import "vue3-toastify/dist/index.css";
import { install as VueMonacoEditorPlugin } from "@guolao/vue-monaco-editor";
import i18n from "./i18n";

const app = createApp(App);

app.use(Vue3Toastify, {
  clearOnUrlChange: false,
} as ToastContainerOptions);
app.use(VueMonacoEditorPlugin, {
  paths: {
    // CDN 配置
    vs: "https://cdn.jsdelivr.net/npm/monaco-editor@0.52.2/min/vs",
  },
});

// Basic plugins
import { install as basicPlugins } from "./plugins/base/token_application/install";
app.use(basicPlugins);

// Premium plugins
// import { install as premiumPlugins } from "./plugins/premium/token_application/install";
// app.use(premiumPlugins, { router });

app.use(router);
app.use(i18n);
app.mount("#app");
