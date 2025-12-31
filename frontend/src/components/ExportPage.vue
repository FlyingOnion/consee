<script setup lang="ts">
import { ref } from "vue";
import TransferBox from "./common/TransferBox.vue";
import { kvList, wExport } from "../common/alova";
import { toast } from "vue3-toastify";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const transferBoxRef = ref<InstanceType<typeof TransferBox>>();
const includeACL = ref(false);
const exportFormat = ref("zip");
const loading = ref(false);

const formatOptions = [
  { value: "json", available: true },
  { value: "zip", available: true },
  { value: "tar.gz", available: false },
  { value: "tar.bz", available: false },
  { value: "7z", available: false },
  { value: "rar", available: false },
];

kvList()
  .then((res) => {
    transferBoxRef.value?.leftContainer?.setElements(res);
  })
  .catch((error) => {
    toast.error(error);
  })
  .finally(() => {
    loading.value = false;
  });

// const exportUrl = computed<string>(() => {
//   const selectedKeys = transferBoxRef.value?.selected || [];
//   return `/export?keys=${selectedKeys.join(",")}&include_acl=${aclEnabled.value ? "1" : "0"
//     }&format=${exportFormat.value}`;
// });

function handleExport() {
  const selectedKeys = transferBoxRef.value?.selected || [];
  if (!selectedKeys.length && !includeACL.value) {
    toast.error("Please put at least one key into 'Selected Items' or enable ACL export");
    return;
  }
  wExport({
    keys: selectedKeys,
    acl: includeACL.value,
    format: exportFormat.value,
  }).then((blob) => {
    const date = new Date();
    const year = date.getFullYear();
    const month = (date.getMonth() + 1).toString().padStart(2, "0");
    const day = date.getDate().toString().padStart(2, "0");
    const hours = date.getHours().toString().padStart(2, "0");
    const minutes = date.getMinutes().toString().padStart(2, "0");
    const seconds = date.getSeconds().toString().padStart(2, "0");

    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.target = "_blank";
    a.download = `consee-export-${year}${month}${day}-${hours}${minutes}${seconds}.${exportFormat.value}`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
  });
}
</script>

<template>
  <div class="flex-grow flex flex-col p-6 gap-3 overflow-y-auto">
    <h1 class="text-2xl font-semibold text-gray-900">{{ t("export.title") }}</h1>

    <div class="bg-white rounded-lg shadow">
      <div class="p-6 border-b border-gray-200">
        <h2 class="text-lg font-medium text-gray-900 mb-4">{{ t("export.kv") }}</h2>
        <div v-if="loading" class="text-sm text-gray-500 py-4">Loading keys...</div>
        <TransferBox v-else ref="transferBoxRef" />
      </div>
    </div>

    <!-- ACL Resources Option -->
    <div class="flex flex-col p-6 gap-4 bg-white rounded-lg shadow">
      <h2 class="text-lg font-medium text-gray-900">{{ t("export.acl") }}</h2>
      <div class="flex items-center gap-4">
        <label for="include-acl" text-sm text-gray-500>{{ t("export.aclLabel") }}</label>
        <input id="include-acl" v-model="includeACL" type="checkbox"
          class="w-5 h-5 text-blue-600 border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 cursor-pointer transition-colors duration-200" />
      </div>
    </div>

    <!-- Export Format -->
    <div class="flex flex-col p-6 gap-4 bg-white rounded-lg shadow">
      <div class="flex flex-col md:flex-row gap-4">
        <h2 class="text-lg font-medium text-gray-900 min-w-20">{{ t("export.format") }}</h2>
        <select v-model="exportFormat"
          class="flex-grow md:max-w-60 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
          <option v-for="option in formatOptions" :key="option.value" :value="option.value"
            :disabled="!option.available">
            {{ option.value }}{{ !option.available ? " (WIP)" : "" }}
          </option>
        </select>
      </div>
      <div v-if="exportFormat === 'json'" class="flex gap-2">
        <i class="i-tabler-info-circle p-2 text-blue-500 h-1lh"></i>
        <p text-sm text-gray-500>{{ t("export.jsonPrompt") }}</p>
      </div>
    </div>

    <!-- Action Buttons -->
    <div class="flex justify-end gap-4">
      <button @click="handleExport"
        class="px-6 py-2 text-sm font-medium text-white bg-blue-600 rounded-md cursor-pointer not-disabled:hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
        :disabled="!includeACL && transferBoxRef?.selected.length === 0">
        {{ t("export.confirm") }}
      </button>
    </div>
  </div>
</template>
