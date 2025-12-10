<script setup lang="ts">
import { computed, ref, watch } from "vue";

import MonacoEditor from "@guolao/vue-monaco-editor";
import {
  b64Encode,
  type KeyValue,
  valueTypeOptions,
  debounce,
} from "../../common/kz";
import { useRouter } from "vue-router";
import {
  kvUpdate,
  kvUpdateValueType,
  kvDelete,
  kvDeleteHints,
  kvGetValueType,
  kvGetValue,
} from "../../common/alova";
import { toast } from "vue3-toastify";
import FullScreenModal from "../common/FullScreenModal.vue";
import DeleteConfirm from "../common/DeleteConfirm.vue";
import emitter from "../../common/mitt";

interface Props {
  k: string;
}

const props = defineProps<Props>();

const key = ref(props.k);
const b64key = computed(() => b64Encode(key.value));

function resetVersion() {}

const valueType = ref("plaintext");

function setValue() {
  kvGetValue(b64key.value)
    .then((kv: KeyValue) => {
      code.value = kv.value;
      originalCode.value = kv.value;
      router.replace(`/kv/${b64key.value}`);
    })
    .catch((e: Error) => {
      toast.error(e);
    });
  kvGetValueType(b64key.value)
    .then((vt) => {
      valueType.value = vt || "plaintext";
    })
    .catch((e: Error) => {
      toast.error(e);
    });
}
setValue();

watch(
  () => props.k,
  (k) => {
    key.value = k;
    debouncedSetValue();
  },
  { immediate: false }
);

// 使用 debounce 包装 setValue 函数，延迟 50ms
const debouncedSetValue = debounce(setValue, 50);

const router = useRouter();

const code = ref("");
const originalCode = ref("");

// Save 功能
function saveKeyValue() {
  kvUpdate(b64key.value, code.value)
    .then(() => {
      toast.success("Key/Value saved successfully");
      originalCode.value = code.value;
    })
    .catch((e: Error) => {
      toast.error(e);
    });
}

function updateValueType() {
  kvUpdateValueType(b64key.value, valueType.value)
    // .then(() => {
    //   invalidateCache(alova.snapshots.match("kvGetValueType"));
    // })
    .catch((e) => {
      toast.error(e);
    });
}

// Delete 功能
function deleteKeyValue() {
  kvDelete(b64key.value)
    .then(() => {
      toast.success("Key/Value deleted successfully");
      emitter.emit("kvDelete");
    })
    .catch((e: Error) => {
      toast.error(e);
    });
}

defineExpose({
  resetVersion,
});
</script>

<template>
  <div class="p-4 bg-gray-50 rounded-lg">
    <!-- Value Type -->
    <div class="space-y-2">
      <label class="block text-sm font-medium text-gray-700" for="value-type">Value Type</label>
      <select id="value-type" v-model="valueType"
        class="w-full px-2 py-1 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white"
        @change="updateValueType">
        <option v-for="o in valueTypeOptions" :key="o" :value="o">{{ o }}</option>
      </select>
    </div>
  </div>
  <div flex-grow>
    <MonacoEditor v-model:value="code" :options="{
      language: valueType || 'plaintext',
      minimap: { enabled: false },
      contextmenu: false,
      automaticLayout: true,
    }" />
  </div>

  <!-- Save 和 Delete 按钮 -->
  <div class="md:px-6 md:py-4">
    <div class="flex flex-col md:flex-row md:justify-end gap-3">
      <!-- Save 按钮 -->
      <button type="button"
        class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
        :disabled="code === originalCode" @click="saveKeyValue">
        Save
      </button>

      <FullScreenModal>
        <template #trigger="{ open }">
          <!-- Delete 按钮 -->
          <button type="button"
            class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 transition-colors duration-200"
            @click="open">
            Delete
          </button>
        </template>
        <template #default="{ close }">
          <DeleteConfirm :close :deleted-elements="[key]" :hints="kvDeleteHints" @delete="deleteKeyValue" />
        </template>
      </FullScreenModal>
    </div>
  </div>
</template>
