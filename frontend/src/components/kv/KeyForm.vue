<script setup lang="ts">
import { computed, ref } from "vue";
import { toast } from "vue3-toastify";
import { valueTypeOptions } from "../../common/kz";
import { kvCreate } from "../../common/alova";
import emitter from "../../common/mitt";

interface Props {
  prefix?: string;
  close: () => void;
}

const props = defineProps<Props>();

const key = ref("");
const value = ref("");
const isPath = computed(() => key.value.trim().endsWith("/"));

const valueType = ref("plaintext");

function createKeyValue() {
  const data = {
    key: props.prefix + key.value.trim(),
    value: isPath.value ? "" : value.value,
    value_type: valueTypeOptions.includes(valueType.value) ? valueType.value : "plaintext",
  };
  if (data.key === props.prefix) {
    toast.error("Key cannot be empty");
    return;
  }
  kvCreate(data)
    .then(() => {
      toast.success("Key/Value created successfully");
      emitter.emit("kvCreate");
      props.close();
    })
    .catch((e: Error) => {
      toast.error(e);
    });

  // const resp = await alova.Post<Response>("/kv/value", data, {
  //   headers: {
  //     [conseeTokenKey]: localStorage.getItem(conseeTokenKey) || "",
  //   },
  // });
  // if (resp.status !== 201) {
  //   toast.error(resp.headers.get(conseeErrorKey) || "Failed to create key and value");
  //   return;
  // }
  // toast.success("Key/Value created successfully");
  // props.close();
}
</script>

<template>
  <div class="bg-white rounded-lg shadow-lg box-border w-150 max-[600px]:w-100vw flex flex-col">
    <!-- Header -->
    <div class="flex items-center justify-between p-6 pb-0">
      <h3 class="text-xl font-semibold text-gray-900 m-0">New Key / Value</h3>
      <button type="button"
        class="w-8 h-8 flex items-center justify-center rounded-full hover:bg-gray-100 transition-colors duration-200"
        @click="close">
        <span class="w-4 h-4 i-tabler-x text-gray-500"></span>
      </button>
    </div>

    <!-- Form Content -->
    <div class="p-6 pt-4 space-y-6">
      <!-- Key Input -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2" for="key">
          Key or folder
        </label>
        <div
          class="flex gap-2 items-center border border-gray-300 rounded-md shadow-sm focus-within:ring-2 focus-within:ring-blue-500 focus-within:border-blue-500">
          <div v-if="prefix" class="px-3 py-2 text-gray-500 text-sm bg-gray-50 border-r border-gray-300 rounded-l-md">
            {{ prefix }}
          </div>
          <input v-model="key"
            class="flex-grow px-3 py-2 border-0 text-gray-700 placeholder-gray-400 focus:outline-none" id="key"
            type="text" placeholder="Key" />
        </div>
      </div>

      <!-- Value Textarea -->
      <div>
        <label class="block text-sm font-medium mb-2" :class="isPath ? 'text-gray-400' : 'text-gray-700'" for="value">
          Value
        </label>
        <div class="relative">
          <textarea :disabled="isPath" v-model="value" id="value"
            class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-y min-h-24 font-mono"
            :class="isPath ? 'bg-gray-50 text-gray-400 cursor-not-allowed' : 'bg-white text-gray-700'
              " placeholder="Enter value..." />
          <div v-if="isPath"
            class="absolute inset-0 flex items-center justify-center bg-gray-50 bg-opacity-90 rounded-md">
            <p class="text-gray-400 text-sm">No value needed for folders</p>
          </div>
        </div>
      </div>

      <!-- Value Type Select -->
      <div>
        <label class="block text-sm font-medium mb-2" :class="isPath ? 'text-gray-400' : 'text-gray-700'"
          for="value_type">
          Value Type
        </label>
        <select :disabled="isPath" v-model="valueType" id="value_type"
          class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          :class="isPath ? 'bg-gray-50 text-gray-400 cursor-not-allowed' : 'bg-white text-gray-700'">
          <option v-for="option in valueTypeOptions" :key="option" :value="option">
            {{ option }}
          </option>
        </select>
      </div>
    </div>

    <!-- Actions -->
    <div class="px-6 py-4 bg-gray-50 border-t border-gray-200 rounded-b-lg flex justify-end gap-3">
      <button type="button"
        class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors duration-200"
        @click="close">
        Cancel
      </button>
      <button type="button"
        class="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors duration-200"
        @click="createKeyValue">
        Save
      </button>
    </div>
  </div>
</template>
