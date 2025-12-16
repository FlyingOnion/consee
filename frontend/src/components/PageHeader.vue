<script setup lang="ts">
import { ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import { RouterLink, useRoute, useRouter } from "vue-router";
import emitter from "../common/mitt";
import { conseeTokenKey } from "../common/const";
import type { AuthResult } from "../common/alova";

const { t, locale } = useI18n();
const mobileMenuOpen = ref(false);
const router = useRouter();
interface Link {
  name: string;
  path: string;
  show?: () => boolean;
}

const links: Link[] = [
  { name: "header.home", path: "/home", show: alwaysShow },
  { name: "header.kv", path: "/kv", show: isAuthenticated },
  { name: "header.tokens", path: "/acl/tokens", show: isAuthenticated },
  { name: "header.policies", path: "/acl/policies", show: isAuthenticated },
  { name: "header.roles", path: "/acl/roles", show: isAuthenticated },
  { name: "header.import", path: "/import", show: isAdmin },
  { name: "header.export", path: "/export", show: isAdmin },
];

// App.vue 设置了检查认证的逻辑，这里只需要设置一个标记，等待 App.vue 发送 login 事件信号
const authenticated = ref(false);
const admin = ref(false);

function alwaysShow() {
  return true;
}

function isAuthenticated() {
  return authenticated.value;
}

function isAdmin() {
  return admin.value;
}

emitter.on("login", setAuthValue);
function setAuthValue(authResult: AuthResult) {
  authenticated.value = true;
  admin.value = authResult.admin === 1;
}

function logout() {
  localStorage.removeItem(conseeTokenKey);
  authenticated.value = false;
  admin.value = false;
  mobileMenuOpen.value = false;
  emitter.emit("openNotificationsChange", 0);
  emitter.emit("logout");
  router.push(`/home`);
}

const route = useRoute();
const mobileShowLogo = ref(route.path !== "/home");
watch(
  () => route.path,
  (newPath) => {
    mobileShowLogo.value = newPath !== "/home";
  },
  { immediate: true }
);

function toggleLanguage() {
  const newLocale = locale.value === "en" ? "zh" : "en";
  locale.value = newLocale;
  localStorage.setItem("locale", newLocale);
}
</script>

<template>
  <div class="hidden md:flex items-center p-2 gap-2">
    <img src="/consee.svg" w-20 alt="Consee" />
    <div v-for="l in links" :key="l.path" flex>
      <RouterLink v-if="l.show?.()" :to="l.path" px-3 py-2 hover:bg-gray-2 active:bg-gray-3 rounded text-blue-4>
        <p flex items-center h-4>{{ t(l.name) }}</p>
      </RouterLink>
    </div>

    <div ml-auto flex items-center gap-2>
      <Notification v-if="admin" />
      <button @click="toggleLanguage" class="flex p-2 rounded-md text-gray-700 hover:bg-gray-100 cursor-pointer"
        :title="locale === 'en' ? '切换到中文' : 'Switch to English'">
        <i w-4 h-4 i-tabler-language />
        <span class="text-xs">{{ locale === "en" ? "中" : "EN" }}</span>
      </button>
      <button v-if="authenticated" @click="logout"
        class="flex p-2 rounded-md text-gray-700 hover:bg-gray-100 cursor-pointer" :title="t('header.logout')">
        <i w-4 h-4 i-tabler-logout />
      </button>
    </div>
  </div>

  <div class="md:hidden flex flex-col p-2">
    <div class="relative w-full flex justify-between items-center">
      <button v-if="authenticated" @click="mobileMenuOpen = !mobileMenuOpen"
        class="flex p-2 rounded-md text-gray-700 hover:bg-gray-100 z-10">
        <i w-4 h-4 i-tabler-baseline-density-medium />
      </button>
      <button @click="toggleLanguage" class="flex p-2 rounded-md text-gray-700 hover:bg-gray-100 cursor-pointer z-10"
        :class="{ 'ml-auto': !authenticated }" :title="locale === 'en' ? '切换到中文' : 'Switch to English'">
        <i w-4 h-4 i-tabler-language />
        <span class="text-xs">{{ locale === "en" ? "中" : "EN" }}</span>
      </button>
      <div v-if="mobileShowLogo" absolute left-0 right-0 h-32px flex justify-center items-center>
        <img src="/consee.svg" w-20 alt="Consee" />
      </div>
    </div>
    <div class="relative">
      <Transition name="mobile-menu">
        <div v-show="mobileMenuOpen" class="absolute flex flex-col bg-white opacity-90 overflow-hidden z-60">
          <hr />
          <div v-for="l in links" :key="l.path">
            <RouterLink v-if="l.show?.()" :to="l.path" flex px-4 py-2 hover:bg-gray-2 active:bg-gray-3 rounded
              text-blue-4 :title="t(l.name)">
              {{ t(l.name) }}
            </RouterLink>
          </div>
          <div v-if="admin">
            <NotificationEntry2 />
          </div>
          <div>
            <a flex items-center gap-2 px-4 py-2 text-blue-4 hover:bg-gray-2 cursor-pointer rounded @click="logout"
              :title="t('header.logout')">
              {{ t("header.logout") }}
              <i w-4 h-4 text-bluegray-6 i-tabler-logout />
            </a>
          </div>
          <!-- <template v-for="l in links">
          <a v-if="route.path.startsWith(l.path)"
            class="no-underline block rounded px-4 py-2 text-blue-4 hover:bg-gray-2">{{ l.name
            }}</a>
          <RouterLink v-else class="no-underline block rounded px-4 py-2 text-blue-4 hover:bg-gray-2" :to="l.path">{{
            l.name }}
          </RouterLink>
        </template> -->
          <hr />
        </div>
      </Transition>
    </div>
  </div>
</template>

<style scoped>
.mobile-menu-enter-from,
.mobile-menu-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.mobile-menu-enter-active,
.mobile-menu-leave-active {
  transition: all 0.3s;
}
</style>
