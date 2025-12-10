<script setup lang="ts">
import { ref } from "vue";
import { useI18n } from "vue-i18n";
import emitter from "../common/mitt";
import { uuidRegexp } from "../common/kz";
import { toast } from "vue3-toastify";
import { conseeTokenKey } from "../common/const";
import { authenticate } from "../common/alova";

const { t } = useI18n();

const currentToken = ref(localStorage.getItem(conseeTokenKey) || "");

const token = ref("");

function login() {
  token.value = token.value.trim();
  if (token.value.length === 0 || token.value === currentToken.value) {
    token.value = "";
    return;
  }
  if (!uuidRegexp.test(token.value)) {
    token.value = "";
    toast.error(t("home.tokenError"));
    return;
  }
  // 使用 alovaCall 的 withToken 来调用，如果 token 不对，会自动清除 token
  localStorage.setItem(conseeTokenKey, token.value);
  authenticate().then((data) => {
    console.log("nOpenMsg:", data.n)
    currentToken.value = token.value;
    emitter.emit("login", data);
    emitter.emit("openNotificationsChange", data.n || 0);
  }).catch((e: Error) => {
    localStorage.removeItem(conseeTokenKey);
    toast.error(e);
  }).finally(() => {
    token.value = "";
  });

}

emitter.on("logout", clearToken);
function clearToken() {
  currentToken.value = "";
}
</script>

<template>
  <div class="flex-grow">
    <div class="w-full max-w-md mx-auto flex flex-col gap-6">
      <!-- Logo Section -->
      <div class="text-center">
        <img class="w-32 h-32 mx-auto" src="/consee.svg" alt="Consee" />
        <h1 class="text-3xl font-bold text-gray-900 mb-2">{{ t("home.title") }}</h1>
        <p class="text-gray-600">{{ t("home.subtitle") }}</p>
      </div>

      <!-- Token Display Section -->
      <div v-if="currentToken" class="bg-green-50 border border-green-200 rounded-lg p-4">
        <div class="flex items-center mb-2">
          <i class="w-5 h-5 i-tabler-circle-check text-green-500 mr-2"></i>
          <span class="text-green-800 font-medium">{{ t("home.loggedIn.title") }}</span>
        </div>
        <p class="text-green-700 text-sm mb-3">
          {{ t("home.loggedIn.description") }}
        </p>
        <div class="bg-white border border-green-300 rounded p-3">
          <code class="text-green-900 font-mono text-sm break-all">{{ currentToken }}</code>
        </div>
        <p class="text-green-600 text-xs mt-2">{{ t("home.loggedIn.switchPrompt") }}</p>
      </div>

      <!-- Login Form Section -->
      <div class="bg-white border border-gray-200 rounded-lg flex flex-col p-4 gap-2">
        <div v-if="!currentToken">
          <h2 class="text-lg font-semibold text-gray-900 mb-2">{{ t("home.accessRequired.title") }}</h2>
          <p class="text-gray-600 text-sm">{{ t("home.accessRequired.description") }}</p>
        </div>

        <label for="token-input" class="block text-sm font-medium text-gray-700">{{ t("home.tokenLabel") }}</label>
        <input id="token-input" type="text" v-model="token" :placeholder="t('home.tokenPlaceholder')"
          class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors duration-200"
          :class="{ 'border-red-300': token && token.length === 0 }" />

        <button type="button" @click="login" :disabled="!token.trim()"
          class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200">
          <span v-if="currentToken">{{ t("home.switchUser") }}</span>
          <span v-else>{{ t("home.login") }}</span>
        </button>

        <!-- Token Application Entry -->
        <TokenApplicationEntry v-if="!currentToken" />
      </div>
    </div>
  </div>
</template>
