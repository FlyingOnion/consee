<script setup lang="ts">
import { computed } from "vue";
import KeyTreeItem from "./KeyTreeItem.vue";
import { parseTree, type TreeItem } from "../../common/kz";

interface Props {
  data: string[];
}

const props = defineProps<Props>();

const itemList = computed<TreeItem[]>(() => {
  return parseTree(props.data || [], { withRoot: true });
});
</script>

<template>
  <div v-for="item in itemList">
    <KeyTreeItem :allKeys="data" :item="item" />
  </div>
</template>
