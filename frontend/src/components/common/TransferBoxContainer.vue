<script setup lang="ts">
import { ref } from "vue";

interface Props {
  boxType: "Available" | "Selected";
}

const props = defineProps<Props>();

const elements = ref<string[]>([]);
const selected = ref<string[]>([]);

function setElements(e: string[]) {
  elements.value = e;
}

function addElements(es: string[]) {
  elements.value.push(...es);
  elements.value.sort();
}

function removeElements(es: string[]) {
  elements.value = elements.value.filter((item) => !es.includes(item));
}

// select
function selectElement(e: string) {
  if (!selected.value.includes(e)) {
    selected.value.push(e);
  } else {
    selected.value = selected.value.filter((item) => item !== e);
  }
}

function selectAll() {
  selected.value = elements.value;
}
function clearSelection() {
  selected.value = [];
}

defineExpose({
  elements,
  selected,
  setElements,
  addElements,
  removeElements,
  selectElement,
  selectAll,
  clearSelection,
});
</script>

<template>
  <div class="flex-1 border border-gray-300 rounded-lg">
    <div class="bg-gray-50 px-4 py-2 border-b border-gray-300 flex justify-between items-center">
      <span class="text-sm font-medium text-gray-700">
        {{ props.boxType }} Items ({{ elements.length }})
      </span>
      <div class="flex gap-2">
        <button
          @click="selectAll"
          class="text-xs px-2 py-1 text-blue-600 hover:bg-blue-50 rounded"
          title="Select All"
        >
          全选
        </button>
        <button
          @click="clearSelection"
          class="text-xs px-2 py-1 text-gray-600 hover:bg-gray-100 rounded"
          title="Clear Selection"
        >
          清空
        </button>
      </div>
    </div>
    <div class="h-64 overflow-y-auto flex flex-col gap-2 p-2">
      <div
        v-for="element in elements"
        :key="element"
        class="px-3 py-2 rounded cursor-pointer"
        :class="{
          'bg-blue-50 text-blue-600 hover:bg-blue-100 active:bg-blue-200':
            selected.includes(element),
          'hover:bg-gray-50 active:bg-gray-100': !selected.includes(element),
        }"
        @click="selectElement(element)"
      >
        <span class="text-sm text-gray-700">{{ element }}</span>
      </div>
    </div>
  </div>
</template>
