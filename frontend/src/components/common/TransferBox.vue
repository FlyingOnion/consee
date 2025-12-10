<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import TransferBoxContainer from "./TransferBoxContainer.vue";

type Container = InstanceType<typeof TransferBoxContainer>;

const leftContainer = ref<Container | null>(null);
const rightContainer = ref<Container | null>(null);

function moveSelected(from: Container, to: Container) {
  to.addElements(from.selected);
  from.removeElements(from.selected);
  from.clearSelection();
}

defineExpose({
  leftContainer,
  rightContainer,
  selected: computed(() => rightContainer.value?.elements || []),
});

const screenWidth = ref(window.innerWidth);
const isMobile = ref(screenWidth.value <= 768);
function updateScreenSize() {
  screenWidth.value = window.innerWidth;
  isMobile.value = screenWidth.value <= 768;
}
onMounted(() => {
  window.addEventListener("resize", updateScreenSize);
});
watch(
  () => isMobile.value,
  (newVal) => {
    isRightPanelVisible.value = !newVal;
  }
);

const isRightPanelVisible = ref(!isMobile.value);
function setRightPanelVisible(visible: boolean) {
  isRightPanelVisible.value = visible;
}
</script>

<template>
  <div class="flex items-center gap-4 relative overflow-hidden">
    <!-- Left Panel -->
    <div class="flex-grow"><TransferBoxContainer box-type="Available" ref="leftContainer" /></div>

    <!-- Transfer Controls -->
    <div v-if="!isMobile" class="flex flex-col gap-2">
      <button
        @click="moveSelected(leftContainer!, rightContainer!)"
        :disabled="!leftContainer?.selected.length"
        class="flex p-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed"
        title="Move to Right"
      >
        <i class="w-4 h-4 i-tabler-chevron-right" />
      </button>
      <button
        @click="moveSelected(rightContainer!, leftContainer!)"
        :disabled="!rightContainer?.selected.length"
        class="flex p-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed"
        title="Move to Left"
      >
        <i class="w-4 h-4 i-tabler-chevron-left" />
      </button>
    </div>

    <transition name="selected">
      <div
        v-show="isRightPanelVisible"
        :class="isMobile ? 'bg-white absolute inset-0 z-40' : 'flex-grow'"
      >
        <TransferBoxContainer box-type="Selected" ref="rightContainer" />
      </div>
    </transition>
  </div>
  <div v-if="isMobile && isRightPanelVisible" flex flex-col gap-2>
    <button
      class="flex p-2 justify-center items-center h-40px hover:bg-gray-50 hover:border-gray-300 rounded cursor-pointer text-bluegray-5"
      @click="setRightPanelVisible(false)"
    >
      Go To Available
    </button>
    <button
      @click="moveSelected(rightContainer!, leftContainer!)"
      :disabled="!rightContainer?.selected.length"
      class="flex p-2 justify-center items-center h-40px bg-blue-600 text-white rounded cursor-pointer hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed"
    >
      <span>Move to Available</span>
    </button>
  </div>
  <div v-if="isMobile && !isRightPanelVisible" flex flex-col gap-2>
    <button
      class="flex p-2 justify-center items-center h-40px hover:bg-gray-50 hover:border-gray-300 rounded cursor-pointer text-bluegray-5"
      @click="setRightPanelVisible(true)"
    >
      Go To Selected
    </button>
    <button
      @click="moveSelected(leftContainer!, rightContainer!)"
      :disabled="!leftContainer?.selected.length"
      class="flex p-2 justify-center items-center h-40px bg-blue-600 text-white rounded cursor-pointer hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed"
    >
      <span>Move to Selected</span>
    </button>
  </div>
</template>

<style scoped>
.selected-enter-active,
.selected-leave-active {
  transition: all 0.3s ease;
}

.selected-enter-from,
.selected-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>
