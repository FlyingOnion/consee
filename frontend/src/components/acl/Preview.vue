<script setup lang="ts">
interface Props {
  rule: string;
}

defineProps<Props>();
</script>

<template>
  <pre v-if="rule">{{ rule }}</pre>
  <p v-else text-bluegray-6>There's still no rule(s) in this token yet.</p>
</template>
