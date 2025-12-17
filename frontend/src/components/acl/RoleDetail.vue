<script setup lang="ts">
import { RouterLink } from "vue-router";
import { computed } from "vue";
import { b64Encode, type RoleDetailInfo } from "../../common/kz";
import FullScreenModal from "../common/FullScreenModal.vue";
import DeleteConfirm from "../common/DeleteConfirm.vue";
import { toast } from "vue3-toastify";
import { aclRoleDelete } from "../../common/alova";
import emitter from "../../common/mitt";

interface Props {
  data: RoleDetailInfo;
}
const props = defineProps<Props>();

const roleDeleteHints = [`The following resource(s) will be deleted:`];
const deletedElements = computed(() => [`ACL Role ${props.data.name}`]);

function deleteRole() {
  aclRoleDelete(b64Encode(props.data.name))
    .then(() => {
      toast.success(`role deleted successfully`);
      emitter.emit("roleDelete");
    })
    .catch((e: Error) => {
      toast.error(e);
    });
}
</script>

<template>
  <div class="flex flex-col">
    <!-- Header -->
    <div class="bg-white border-b border-gray-200 px-4">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-xl font-semibold text-gray-900">{{ data.name || "Unnamed Role" }}</h2>
          <p class="text-sm text-gray-500 mt-1">Role Details</p>
        </div>
        <div class="flex items-center space-x-2">
          <span
            class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
            <i class="w-3 h-3 i-tabler-users mr-1" />
            ACL Role
          </span>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex flex-col overflow-y-auto p-4 gap-4">
      <!-- Basic Information -->
      <h3 class="text-lg font-medium text-gray-900 flex items-center gap-2">
        <i class="w-5 h-5 i-tabler-info-circle text-blue-500" />
        Basic Information
      </h3>
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Role ID</label>
          <div class="font-mono text-sm bg-gray-50 rounded px-3 py-2 border border-gray-200 break-all">
            {{ data.id }}
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Role Name</label>
          <div class="font-mono text-sm bg-gray-50 rounded px-3 py-2 border border-gray-200 break-all">
            {{ data.name }}
          </div>
        </div>
      </div>

      <!-- Description -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
        <div class="text-sm bg-gray-50 rounded px-3 py-2 border border-gray-200">
          {{ data.description || "No description provided" }}
        </div>
      </div>

      <!-- Policies Section -->
      <div class="bg-white rounded-lg border border-gray-200 shadow-sm">
        <div class="border-b border-gray-200">
          <h3 class="text-lg font-medium text-gray-900 flex items-center">
            <i class="w-5 h-5 i-tabler-shield mr-2 text-green-500" />
            Associated Policies
          </h3>
        </div>
        <div class="py-4">
          <div v-if="data.policies.length" class="space-y-2">
            <RouterLink v-for="{id, name} in data.policies" :key="id" :to="`/acl/policy/${b64Encode(name)}`"
              class="flex items-center justify-between p-3 bg-gray-50 border border-gray-200 rounded-lg hover:bg-gray-100 transition-colors duration-200">
              <div class="flex items-center">
                <i class="w-4 h-4 i-tabler-shield mr-2 text-gray-600" />
                <span class="font-medium text-gray-900">{{ name }}</span>
              </div>
              <i class="w-4 h-4 i-tabler-chevron-right text-gray-400" />
            </RouterLink>
          </div>
          <div v-else class="text-center py-8">
            <i class="w-12 h-12 i-tabler-shield-off mx-auto mb-3 text-gray-300" />
            <p class="text-gray-500 text-sm">No policies are associated with this role</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Action Buttons -->
    <div class="bg-white border-t border-gray-200 px-6 py-4">
      <div class="flex flex-col sm:flex-row sm:justify-end gap-2">
        <FullScreenModal>
          <template #trigger="{ open }">
            <button type="button" @click="open"
              class="inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition-colors duration-200">
              <i class="w-4 h-4 i-tabler-trash mr-2" />
              Delete Role
            </button>
          </template>
          <template #default="{ close }">
            <DeleteConfirm :close="close" :onDelete="deleteRole" :hints="roleDeleteHints"
              :deleted-elements="deletedElements" />
          </template>
        </FullScreenModal>
      </div>
    </div>
  </div>
</template>