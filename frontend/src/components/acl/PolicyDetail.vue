<script setup lang="ts">
import { RouterLink } from "vue-router";
import { b64Encode, type PolicyDetailInfo } from "../../common/kz";
import PolicyRules from "./PolicyRules.vue";
import FullScreenModal from "../common/FullScreenModal.vue";
import DeleteConfirm from "../common/DeleteConfirm.vue";
import { computed } from "vue";
import { toast } from "vue3-toastify";
import { aclPolicyDelete } from "../../common/alova";
import emitter from "../../common/mitt";

interface Props {
  data: PolicyDetailInfo;
}
const props = defineProps<Props>();

const isExclusive = computed(
  () => props.data.tokens.length === 1 && props.data.name === `--${props.data.tokens[0].id}`
);
const isReadOnly = computed(
  () => props.data.name === "global-management" || props.data.name === "builtin/global-read-only"
);

const policyDeleteHints = [`The following resource(s) will be deleted:`];
const deletedElements = computed(() => [`ACL Policy ${props.data.name}`]);

function deletePolicy() {
  if (isReadOnly.value || isExclusive.value) {
    return;
  }
  aclPolicyDelete(b64Encode(props.data.name))
    .then(() => {
      toast.success(`policy deleted successfully`);
      emitter.emit("policyDelete");
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
          <h2 class="text-xl font-semibold text-gray-900">{{ data.name || "Unnamed Policy" }}</h2>
          <p class="text-sm text-gray-500 mt-1">Policy Details</p>
        </div>
        <div class="flex items-center space-x-2">
          <span
            class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
            <i class="w-3 h-3 i-tabler-shield mr-1" />
            ACL Policy
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
          <label class="block text-sm font-medium text-gray-700 mb-1">Policy ID</label>
          <div class="font-mono text-sm bg-gray-50 rounded px-3 py-2 border border-gray-200 break-all">
            {{ data.id }}
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Policy Name</label>
          <div class="font-mono text-sm bg-gray-50 rounded px-3 py-2 border border-gray-200 break-all">
            {{ data.name }}
          </div>
        </div>
      </div>

      <!-- Policy Rules Section -->
      <div class="bg-white rounded-lg border border-gray-200 shadow-sm">
        <div class="border-b border-gray-200">
          <h3 class="text-lg font-medium text-gray-900 flex items-center gap-2">
            <i class="w-5 h-5 i-tabler-file-text text-green-500" />
            Policy Rules
          </h3>
        </div>
        <div class="p-4">
          <PolicyRules :readonly="isReadOnly" :rules="data.parsed_rules" />
        </div>
      </div>

      <!-- Tokens Section -->
      <div class="bg-white rounded-lg border border-gray-200 shadow-sm">
        <div class="px-4 py-3 border-b border-gray-200">
          <h3 class="text-lg font-medium text-gray-900 flex items-center">
            <i class="w-5 h-5 i-tabler-key mr-2 text-purple-500" />
            Tokens
          </h3>
        </div>
        <div class="p-4">
          <div v-if="data.tokens.length" class="space-y-2">
            <RouterLink v-for="{ id, name } in data.tokens" :key="id" :to="`/acl/token/${id}`"
              class="flex items-center justify-between p-3 bg-gray-50 border border-gray-200 rounded-lg hover:bg-gray-100 transition-colors duration-200">
              <div class="flex items-center">
                <i class="w-4 h-4 i-tabler-key mr-2 text-gray-600" />
                <span class="font-medium text-gray-900">{{ name }}</span>
              </div>
              <i class="w-4 h-4 i-tabler-chevron-right text-gray-400" />
            </RouterLink>
          </div>
          <div v-else class="text-center py-8">
            <i class="w-12 h-12 i-tabler-key-off mx-auto mb-3 text-gray-300" />
            <p class="text-gray-500 text-sm">No tokens are using this policy</p>
          </div>
        </div>
      </div>

      <!-- Policy Status -->
      <div v-if="isReadOnly || isExclusive" class="bg-blue-50 border border-blue-200 rounded-lg p-4">
        <div class="flex items-start">
          <i class="w-5 h-5 i-tabler-info-circle text-blue-600 mt-0.5 mr-3 flex-shrink-0" />
          <div>
            <h4 class="text-sm font-medium text-blue-900 mb-1">
              {{ isReadOnly ? "Read-only Policy" : "Exclusive Policy" }}
            </h4>
            <p class="text-sm text-blue-700">
              {{
                isReadOnly
                  ? "This policy is readonly. You can't modify or delete it."
                  : `This policy is an exclusive policy. You can only modify its rules. To delete it, please go to token
              page and delete the token that uses this policy.`
              }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Action Buttons -->
    <div v-if="!isReadOnly && !isExclusive" class="bg-white border-t border-gray-200 px-6 py-4">
      <div class="flex flex-col sm:flex-row sm:justify-end gap-2">
        <FullScreenModal>
          <template #trigger="{ open }">
            <button type="button" @click="open"
              class="inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition-colors duration-200">
              <i class="w-4 h-4 i-tabler-trash mr-2" />
              Delete Policy
            </button>
          </template>
          <template #default="{ close }">
            <DeleteConfirm :close="close" :onDelete="deletePolicy" :hints="policyDeleteHints"
              :deleted-elements="deletedElements" />
          </template>
        </FullScreenModal>
      </div>
    </div>
  </div>
</template>
