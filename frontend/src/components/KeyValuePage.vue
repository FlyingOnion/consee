<script setup lang="ts">
import { provide, ref } from "vue";
import KeyTree from "./kv/KeyTree.vue";
import { useRoute, useRouter } from "vue-router";
import { b64Decode, b64Encode, debounce } from "../common/kz";
import SplitView from "./common/SplitView.vue";
import emitter from "../common/mitt";
import { toast } from "vue3-toastify";
import { kvList } from "../common/alova";
import EmptyView from "./common/EmptyView.vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

function loadOpenKeys(): Set<string> {
  const openKeys = localStorage.getItem("kv-open-keys");
  if (openKeys) {
    return new Set(JSON.parse(openKeys) as string[]);
  }
  return new Set();
}
const openKeySet: Set<string> = loadOpenKeys();

/**
 * This function should be debounced and not be called directly.
 * Call debouncedRecordOpenKeys instead.
 */
function recordOpenKeys() {
  localStorage.setItem("kv-open-keys", JSON.stringify(Array.from(openKeySet)));
}

const debouncedRecordOpenKeys = debounce(recordOpenKeys, 3000);

function addOpenKey(key: string) {
  openKeySet.add(key);
  debouncedRecordOpenKeys();
}
function removeOpenKey(key: string) {
  openKeySet.delete(key);
  debouncedRecordOpenKeys();
}
function isKeyOpen(key: string) {
  return openKeySet.has(key);
}

function onItemClick(key: string) {
  router.push(`/kv/${b64Encode(key)}`);
  currentKey.value = key;
  valuearea2.value?.resetVersion?.();
}

provide("addOpenKey", addOpenKey);
provide("removeOpenKey", removeOpenKey);
provide("isKeyOpen", isKeyOpen);
provide("onItemClick", onItemClick);

function initialKey() {
  const key = route.params.key as string;
  return key ? b64Decode(key) || "" : "";
}

const route = useRoute();
const router = useRouter();
const keys = ref<string[]>([]);
const currentKey = ref(initialKey());
/**
 * loading is only used when loading key list, not when loading value.
 */
const loading = ref(true);
const error = ref<Error | undefined>(undefined);

kvList()
  .then((data) => {
    keys.value = data;
    loading.value = false;
  })
  .catch((e: Error) => {
    currentKey.value = ""; // to prevent the value area from showing the wrong key
    loading.value = false;
    error.value = e;
  });

emitter.on("kvCreate", refresh);
emitter.on("kvDelete", refresh);

function refresh() {
  setTimeout(() => {
    kvList()
      .then((data) => {
        keys.value = data;
        if (currentKey.value && !keys.value.includes(currentKey.value)) {
          currentKey.value = "";
          router.replace("/kv");
        }
      })
      .catch((e: Error) => {
        toast.warn(t("keyValue.refreshFailed", { message: e.message }));
      });
  }, 300);
}

const valuearea2 = ref<{ resetVersion: () => void } | null>(null);
</script>

<template>
  <SplitView :loading :error :title="t('keyValue.title')" :current="currentKey">
    <template #list>
      <KeyTree :data="keys" />
    </template>
    <ValueArea2 ref="valuearea2" :k="currentKey"></ValueArea2>
    <template #empty>
      <EmptyView :resource="t('keyValue.resource')">
        <p class="text-sm">{{ t('keyValue.selectPrompt') }}</p>
      </EmptyView>
    </template>
  </SplitView>
</template>
