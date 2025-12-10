<template>
  <div flex flex-wrap gap-2>
    <p m-0 w-full>Admin notification is available in premium version.</p>
    <a px-2 py-1 rounded bg-gray-50 hover:bg-gray-2 @click="ok">OK</a>
    <a px-2 py-1 rounded bg-gray-50 hover:bg-gray-2 @click="hide">do not show again</a>
  </div>
</template>

<script setup lang="ts">
interface Props {
  ok: () => void;
  hide: () => void;
}

defineProps<Props>();
</script>
