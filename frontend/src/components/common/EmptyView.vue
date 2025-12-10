<template>
  <div class="h-full flex items-center justify-center">
    <div class="text-center text-gray-500">
      <i class="w-16 h-16 i-tabler-key-off mx-auto mb-4 text-gray-300" />
      <slot name="title">
        <p class="text-lg font-medium">No {{ resource || "item" }} selected</p>
      </slot>
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  resource?: string;
}

defineProps<Props>();
</script>
