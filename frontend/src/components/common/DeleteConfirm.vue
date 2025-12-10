<template>
  <div
    class="relative bg-white p-4 rounded-lg box-border w-150 max-[600px]:w-full flex flex-col gap-2"
  >
    <span i-tabler-x cursor-pointer absolute right-4 top-4 @click="close" />
    <h3 text-center text-base text-gray-9>Delete Confirm</h3>

    <p text-gray-5 text-sm>Are you sure you want to delete? This action cannot be undone.</p>

    <p v-for="h in props.hints" text-gray-5 text-sm>{{ h }}</p>
    
    <p v-for="item in props.deletedElements" text-sm text-gray-6 font-mono>- {{ item }}</p>
    
    <div bg-gray-50 px-6 py-3 flex flex-row-reverse gap-2>
      <button
        type="button"
        rounded-md
        bg-white
        px-3
        py-2
        text-sm
        font-semibold
        text-gray-900
        ring-1
        shadow-xs
        ring-gray-300
        ring-inset
        hover:bg-gray-50
        w-20
        @click="close"
      >
        Cancel
      </button>
      <button
        type="button"
        rounded-md
        bg-red-5
        px-3
        py-2
        text-sm
        font-semibold
        text-white
        shadow-xs
        hover:bg-red-4
        w-20
        @click="doDelete"
      >
        Delete
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  close: () => void;
  hints?: string[];
  deletedElements?: string[];
  onDelete: () => void;
}

const props = defineProps<Props>();

function doDelete() {
  if (!props.deletedElements) {
    props.close();
    return;
  }
  props.onDelete();
  props.close();
}
</script>
