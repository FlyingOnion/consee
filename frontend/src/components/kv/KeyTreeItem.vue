<script setup lang="ts">
import { computed, inject, ref } from "vue";
import { b64Encode, type TreeItem } from "../../common/kz";
import { kvDelete, kvDeleteHints } from "../../common/alova";
import FullScreenModal from "../common/FullScreenModal.vue";
import KeyForm from "./KeyForm.vue";
import DeleteConfirm from "../common/DeleteConfirm.vue";
import { toast } from "vue3-toastify";
import emitter from "../../common/mitt";

interface Props {
  allKeys: string[];
  item: TreeItem;
}

const props = defineProps<Props>();
const isRoot = computed(() => props.item.path === "/");

const isKeyOpen: (key: string) => boolean = inject("isKeyOpen")!;
const addOpenKey: (key: string) => void = inject("addOpenKey")!;
const removeOpenKey: (key: string) => void = inject("removeOpenKey")!;
const onItemClick: (key: string) => void = inject("onItemClick")!;

const openStatus = ref(!props.item.isLeaf && isKeyOpen(props.item.path));

function toggleOpenStatus() {
  openStatus.value ? removeOpenKey(props.item.path) : addOpenKey(props.item.path);
  openStatus.value = !openStatus.value;
}

function onClick() {
  props.item.isLeaf ? onItemClick(props.item.path) : toggleOpenStatus();
}

const paddingStyle = {
  paddingLeft: `${props.item.depth < 1 ? 0.25 : props.item.depth - 0.75}rem`,
  paddingRight: "0.25rem",
  paddingTop: "0.25rem",
  paddingBottom: "0.25rem",
};

function selfAndChildren(): string[] {
  return props.item.path.endsWith("/")
    ? props.allKeys.filter((item) => item.startsWith(props.item.path))
    : [props.item.path];
}

function deleteKeyValue() {
  kvDelete(b64Encode(props.item.path))
    .then(() => {
      toast.success("Key/Value deleted successfully");
      emitter.emit("kvDelete");
    })
    .catch((e: Error) => {
      toast.error(e);
    });
}
</script>

<script lang="ts">
export default {
  name: "KeyTreeItem",
  props: {
    item: {
      type: Object,
      required: true,
    },
  },
};
</script>

<template>
  <div flex items-center gap-2>
    <p v-if="isRoot" flex-grow m-0 rounded p-1>/</p>
    <p
      v-else
      flex-grow
      flex
      items-center
      gap-1
      m-0
      :style="paddingStyle"
      rounded
      cursor-pointer
      hover:bg-gray-2
      active:bg-gray-3
      :title="props.item.path"
      :onclick="onClick"
    >
      <span v-if="props.item.isLeaf" i-tabler-file-description></span>
      <span v-else-if="openStatus" i-tabler-folder-open></span>
      <span v-else i-tabler-folder></span>
      {{ props.item.key }}
    </p>
    <FullScreenModal v-if="isRoot || !item.isLeaf">
      <template #trigger="{ open }">
        <span i-tabler-plus cursor-pointer @click="open" />
      </template>
      <template #default="{ close }">
        <KeyForm :close="close" :prefix="isRoot ? '' : props.item.path" />
      </template>
    </FullScreenModal>
    <span v-if="isRoot || item.isLeaf" w-16px h-16px />
    <FullScreenModal v-if="!isRoot || item.isLeaf">
      <template #trigger="{ open }">
        <span i-tabler-trash cursor-pointer @click="open" />
      </template>
      <template #default="{ close }">
        <DeleteConfirm
          :close="close"
          :deleted-elements="selfAndChildren()"
          :hints="kvDeleteHints"
          @delete="deleteKeyValue"
        />
      </template>
    </FullScreenModal>
  </div>
  <div v-show="openStatus" v-for="item in props.item.children">
    <KeyTreeItem :allKeys :item="item" />
  </div>
</template>
