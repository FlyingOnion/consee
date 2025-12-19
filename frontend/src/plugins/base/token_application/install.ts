import type { App } from "vue";
import Empty from "../../../components/common/Empty.vue";
import Hello from "./Hello.vue";
import ValueArea2 from "../../../components/kv/ValueArea2.vue";

function install<Options = any[]>(app: App, options?: Options) {
  // Register global components
  // app.component("Notification", NotificationBase);

  // Premium components
  app.component("ValueArea2", ValueArea2);
  app.component("Notification", Empty);
  app.component("NotificationEntry2", Empty);
  app.component("Hello", Hello);
  app.component("TokenApplicationEntry", Empty);
}

export { install };
