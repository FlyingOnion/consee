<script setup lang="ts">
import { ref } from "vue";
import PolicySelectAll from "./PolicySelectAll.vue";
import { aclRoleCreate } from "../../common/alova";
import { toast } from "vue3-toastify";
import emitter from "../../common/mitt";

interface Props {
  close: () => void;
}
const props = defineProps<Props>();

const name = ref("");
const description = ref("");
const policySelect = ref<InstanceType<typeof PolicySelectAll> | null>(null);

async function createRole() {
  if (!name.value) {
    toast.error("Role name cannot be empty");
    return;
  }
  
  const selectedPolicies = policySelect.value?.selected || [];
  
  try {
    await aclRoleCreate({
      name: name.value,
      description: description.value,
      policies: selectedPolicies,
    });
    toast.success("Role created");
    props.close();
    emitter.emit("roleCreate");
  } catch (e: any) {
    toast.error(e.message || "Failed to create role");
  }
}
</script>

<template>
  <div
    class="relative bg-white p-4 rounded-lg box-border w-150 max-[600px]:w-full flex flex-col gap-2"
  >
    <span i-tabler-x cursor-pointer absolute right-4 top-4 @click="close" />
    <h3 m-0 text-center text-base font-semibold text-gray-900>New Role</h3>
    <div flex flex-col gap-3>
      <div flex flex-col gap-1>
        <label text-gray-700 text-sm font-bold for="name"> Role Name </label>
        <input
          v-model="name"
          flex-grow
          shadow
          appearance-none
          border
          rounded
          px-2
          py-1
          text-gray-700
          leading-tight
          focus-outline-blue
          id="name"
          type="text"
          placeholder="Role name"
        />
      </div>
      
      <div flex flex-col gap-1>
        <label text-gray-700 text-sm font-bold for="description"> Description </label>
        <textarea
          v-model="description"
          flex-grow
          shadow
          appearance-none
          border
          rounded
          px-2
          py-1
          text-gray-700
          leading-tight
          focus-outline-blue
          id="description"
          type="text"
          placeholder="Role description (optional)"
          rows="3"
        />
      </div>
      
      <div flex flex-col gap-1>
        <label text-gray-700 text-sm font-bold> Associated Policies </label>
        <PolicySelectAll ref="policySelect" exclStatus="non-exclusive" />
      </div>
    </div>
    <div bg-gray-50 px-6 py-3 flex flex-row-reverse gap-2>
      <button
        type="button"
        rounded-md
        bg-blue-600
        px-3
        py-2
        text-white
        text-sm
        font-medium
        hover:bg-blue-700
        focus-outline-blue
        @click="createRole"
      >
        Create Role
      </button>
      <button
        type="button"
        rounded-md
        border
        border-gray-300
        bg-white
        px-3
        py-2
        text-gray-700
        text-sm
        font-medium
        hover-bg-gray-50
        focus-outline-blue
        @click="close"
      >
        Cancel
      </button>
    </div>
  </div>
</template>