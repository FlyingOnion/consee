import { createRouter, createWebHistory, type RouteRecordRaw } from "vue-router";

import HomePage from "./components/HomePage.vue";
import KeyValuePage from "./components/KeyValuePage.vue";
import TokenPage from "./components/TokenPage.vue";
import PolicyPage from "./components/PolicyPage.vue";
import ImportPage from "./components/ImportPage.vue";
import ExportPage from "./components/ExportPage.vue";
import RolePage from "./components/RolePage.vue";

const routes: RouteRecordRaw[] = [
  { path: "/", redirect: "home" },
  { path: "/home", component: HomePage, meta: { title: "Consee | Home" } },
  { path: "/kv", component: KeyValuePage, meta: { title: "Consee | KV" } },
  { path: "/kv/:key", component: KeyValuePage, meta: { title: "Consee | KV" } },
  { path: "/acl/tokens", component: TokenPage, meta: { title: "Consee | ACL Token" } },
  { path: "/acl/token/:id", component: TokenPage, meta: { title: "Consee | ACL Token" } },
  { path: "/acl/policies", component: PolicyPage, meta: { title: "Consee | ACL Policy" } },
  {
    path: "/acl/policy/:b64name",
    component: PolicyPage,
    meta: { title: "Consee | ACL Policy" },
  },
  { path: "/acl/roles", component: RolePage, meta: { title: "Consee | ACL Role" } },
  { path: "/acl/roles/:id", component: RolePage, meta: { title: "Consee | ACL Role" } },
  { path: "/import", component: ImportPage, meta: { title: "Consee | Import" } },
  { path: "/export", component: ExportPage, meta: { title: "Consee | Export" } },
];

const base = import.meta.env.CONSEE_BASE || '';
const router = createRouter({
  history: createWebHistory(base),
  routes,
});

router.beforeEach((to, _, next) => {
  document.title = (to.meta.title as string) || "Consee";
  next();
});

export { router };
