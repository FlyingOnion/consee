<script setup lang="ts">
import { computed, ref } from "vue";
import TransferBox from "./common/TransferBox.vue";
import { kvList, wExport } from "../common/alova";
import { toast } from "vue3-toastify";

const transferBoxRef = ref<InstanceType<typeof TransferBox>>();
const aclEnabled = ref(false);
const exportFormat = ref("zip");
const loading = ref(false);

const formatOptions = [
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
  if (!selectedKeys.length && !aclEnabled.value) {
    toast.error("Please put at least one key into 'Selected Items' or enable ACL export");
    return;
  }
  wExport({
    keys: selectedKeys,
    include_acl: aclEnabled.value,
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
    <h1 class="text-2xl font-semibold text-gray-900">Export Resources</h1>

    <div class="bg-white rounded-lg shadow">
      <div class="p-6 border-b border-gray-200">
        <h2 class="text-lg font-medium text-gray-900 mb-4">Key Values</h2>
        <div v-if="loading" class="text-sm text-gray-500 py-4">Loading keys...</div>
        <TransferBox v-else ref="transferBoxRef" />
      </div>
    </div>

    <!-- ACL Resources Option -->
    <div class="bg-white rounded-lg shadow">
      <div class="p-6 border-b border-gray-200">
        <h2 class="text-lg font-medium text-gray-900">ACL Resources</h2>
        <input v-model="aclEnabled" type="checkbox" class="text-blue-600 focus:ring-blue-500" />
      </div>
    </div>

    <!-- Export Format -->
    <div class="bg-white rounded-lg shadow mb-6">
      <div class="p-6">
        <h2 class="text-lg font-medium text-gray-900 mb-4">Export Format</h2>
        <select v-model="exportFormat"
          class="w-full md:w-64 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500">
          <option v-for="option in formatOptions" :key="option.value" :value="option.value"
            :disabled="!option.available">
            {{ option.value }}{{ !option.available ? " (WIP)" : "" }}
          </option>
        </select>
      </div>
    </div>

    <!-- Action Buttons -->
    <div class="flex justify-end gap-4">
      <button
        @click="handleExport"
        class="px-6 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
      >
        Confirm Export
      </button>
    </div>
  </div>
</template>
