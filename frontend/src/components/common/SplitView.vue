<template>
  <div class="overflow-y-auto flex-grow relative p-2 flex transition-all duration-300 ease-in-out">
    <slot v-if="loading" name="loading">
      <p text-bluegray-5>Connecting to consee backend</p>
    </slot>
    <slot v-else-if="error" name="error">
      <p text-bluegray-5>{{ error }}</p>
    </slot>
    <template v-else>
      <!-- 左侧列表区域 -->
      <transition name="list">
        <div
          v-if="isListVisible"
          :class="isMobile ? 'bg-white absolute inset-2 z-40' : 'w-320px overflow-y-auto'"
        >
          <div flex items-center p-2 gap-2 h-52px>
            <button
              v-if="isMobile"
              @click="setListVisible(false)"
              class="flex p-2 rounded hover:bg-gray-100 transition-colors"
            >
              <i w-4 h-4 i-tabler-chevron-left />
            </button>
            <h2 class="text-lg font-semibold flex-grow">{{ title }}</h2>
            <slot name="controls"></slot>
          </div>
          <div class="flex-1 overflow-y-auto">
            <slot name="list"></slot>
          </div>
        </div>
      </transition>

      <!-- 右侧内容区域 -->
      <div class="relative flex-1 flex flex-col">
        <div flex items-center h-52px flex-shrink-0>
          <button
            @click="setListVisible(!isListVisible)"
            class="flex p-2 rounded cursor-pointer hover:bg-gray-100 transition-colors"
          >
            <i
              w-4
              h-4
              :class="isListVisible ? 'i-tabler-chevron-left' : 'i-tabler-chevron-right'"
              :title="`show / hide ${title.toLocaleLowerCase()}`"
            />
          </button>
          <template v-if="current">
            <h2 class="text-lg font-semibold flex-grow">{{ current }}</h2>
            <slot name="controls-right"></slot>
          </template>
          <p v-else h-32px text-bluegray-5 py-1>{{ defaultHint }}</p>
        </div>
        <div class="flex-grow flex flex-col overflow-y-auto">
          <slot v-if="current"></slot>
          <slot v-else name="empty">
            <p pl-3 text-gray-500>Choose one to see details.</p>
            <p v-for="h in hints" :key="h" pl-3 text-gray-500>{{ h }}</p>
          </slot>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";

interface Props {
  loading: boolean;
  error?: Error;
  title: string;
  current?: string;
  hints?: string[];
}
const props = defineProps<Props>();

const screenWidth = ref(window.innerWidth);
const isMobile = computed(() => screenWidth.value <= 960);
const isListVisibleStorageKey = "G-Consee-SplitView:Is-List-Visible";
const isListVisible = ref(!isMobile.value && localStorage.getItem(isListVisibleStorageKey) === "1");

const defaultHint = computed(
  () => `${isListVisible.value ? "Hide" : "Show"} ${props.title.toLocaleLowerCase()}`
);

function setListVisible(visible: boolean) {
  isListVisible.value = visible;
  if (visible) {
    localStorage.setItem(isListVisibleStorageKey, "1");
  } else {
    localStorage.removeItem(isListVisibleStorageKey);
  }
}

function updateScreenSize() {
  screenWidth.value = window.innerWidth;
}

onMounted(() => {
  window.addEventListener("resize", updateScreenSize);
});
</script>

<style scoped>
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(-100%);
}
</style>
