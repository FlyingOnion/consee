<script setup lang="ts">
import { ref } from "vue";

interface Props {
  open?: boolean;
  title: string;
  rtitle?: string;
}

const props = defineProps<Props>();

const openStatus = ref(props.open);
</script>

<template>
  <div relative>
    <slot name="trigger">
      <p
        m-0
        py-2
        text-bluegray-6
        font-700
        cursor-pointer
        rounded
        flex
        items-center
        gap-2
        @click="openStatus = !openStatus"
      >
        <span v-show="openStatus" i-tabler-chevron-down></span>
        <span v-show="!openStatus" i-tabler-chevron-right></span>
        <span flex-grow>{{ props.title }}</span>
        <span v-if="props.rtitle" ml-auto>{{ props.rtitle }}</span>
      </p>
    </slot>
    <div v-show="openStatus"><slot></slot></div>
  </div>
</template>
