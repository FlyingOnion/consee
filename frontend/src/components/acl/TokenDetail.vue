<script setup lang="ts">
import { computed, ref } from "vue";
import { type TokenDetailInfo, b64Encode } from "../../common/kz";
import DeleteConfirm from "../common/DeleteConfirm.vue";
import Drawer from "../common/Drawer.vue";
import FullScreenModal from "../common/FullScreenModal.vue";
import { conseeTokenKey } from "../../common/const";
import { toast } from "vue3-toastify";
import PolicySelectAll from "./PolicySelectAll.vue";
import { aclTokenDelete, aclTokenUpdate } from "../../common/alova";
import emitter from "../../common/mitt";

interface Props {
  data: TokenDetailInfo;
}

const loginToken = localStorage.getItem(conseeTokenKey) || "";

const props = defineProps<Props>();

const hasExclusivePolicy = computed(() => {
  return (
    props.data.policies.length === 1 && props.data.policies[0].name == `--${props.data.accessor_id}`
  );
});

const tokenDeleteHints = [
  `If a token is created with an exclusive policy, deleting it will also delete the policy.`,
  `The following resource(s) will be deleted:`,
];

const deletedElements = computed(() =>
  hasExclusivePolicy.value
    ? [`ACL Token  ${props.data.name}`, `ACL Policy ${props.data.policies[0].name}`]
    : [`ACL Token ${props.data.name}`]
);

const policyselect = ref<InstanceType<typeof PolicySelectAll> | null>(null);

function saveToken() {
  const selected = policyselect.value?.selected || [];
  aclTokenUpdate(
    props.data.accessor_id,
    selected.map(({ name }) => name)
  )
    .then(() => {
      toast.success(`token updated successfully`);
    })
    .catch((e: Error) => {
      toast.error(e);
    });
}

function deleteToken() {
  aclTokenDelete(props.data.accessor_id)
    .then(() => {
      toast.success(`token deleted successfully`);
      emitter.emit("tokenDelete");
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
          <h2 class="text-xl font-semibold text-gray-900">{{ data.name || "Unnamed Token" }}</h2>
          <p class="text-sm text-gray-500 mt-1">Token Details</p>
        </div>
        <div class="flex items-center space-x-2">
          <span
            class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
            <i class="w-3 h-3 i-tabler-key mr-1" />
            ACL Token
          </span>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div class="flex-1 overflow-y-auto p-4 space-y-6">
      <!-- Basic Information -->

      <h3 class="text-lg font-medium text-gray-900 mb-4 flex items-center">
        <i class="w-5 h-5 i-tabler-info-circle mr-2 text-blue-500" />
        Basic Information
      </h3>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Accessor ID</label>
          <div class="font-mono text-sm bg-gray-50 rounded px-3 py-2 border border-gray-200 break-all">
            {{ data.accessor_id }}
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Secret ID</label>
          <div class="font-mono text-sm bg-gray-50 rounded px-3 py-2 border border-gray-200 break-all">
            {{ data.secret_id }}
          </div>
        </div>
      </div>

      <!-- Policies Section -->
      <Drawer title="Policies" open>
        <template v-if="hasExclusivePolicy">
          <div class="space-y-2">
            <RouterLink v-for="policy in data.policies" :key="policy.name"
              :to="`/acl/policy/${b64Encode(policy.name)}`"
              class="flex items-center justify-between p-3 bg-yellow-50 border border-yellow-200 rounded-lg hover:bg-yellow-100 transition-colors duration-200">
              <div class="flex items-center">
                <i class="w-4 h-4 i-tabler-shield-check mr-2 text-yellow-600" />
                <span class="font-medium text-yellow-800">{{ policy.name }}</span>
              </div>
              <span
                class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
                Exclusive
              </span>
            </RouterLink>
          </div>
        </template>
        <div v-else flex flex-col gap-2>
          <PolicySelectAll ref="policyselect" excl-status="non-exclusive"
            :selected-policies="data.policies.map(({ id }) => id)" />
        </div>
      </Drawer>

      <!-- Roles Section -->
      <Drawer title="Roles" open>
        <div v-if="!data.roles.length" class="text-center py-8">
          <i class="w-12 h-12 i-tabler-users-group mx-auto mb-3 text-gray-300" />
          <p class="text-gray-500 text-sm">No roles are adopted to this token.</p>
        </div>
      </Drawer>

      <!-- Metadata Section -->
      <Drawer title="Metadata" open>
        <div v-if="Object.keys(data.metadata).length" class="space-y-3">
          <div v-for="(value, key) in data.metadata" :key="key"
            class="flex flex-col sm:flex-row sm:items-center justify-between p-3 bg-gray-50 rounded-lg border border-gray-200">
            <span class="text-sm font-medium text-gray-700 capitalize mb-1 sm:mb-0 sm:w-32">
              {{ key.replace(/_/g, " ") }}
            </span>
            <span class="text-sm text-gray-900 break-all">{{ value || "Not specified" }}</span>
          </div>
        </div>
        <div v-else class="text-center py-8">
          <i class="w-12 h-12 i-tabler-tags mx-auto mb-3 text-gray-300" />
          <p class="text-gray-500 text-sm">No metadata available</p>
        </div>
      </Drawer>
    </div>

    <!-- Action Buttons -->
    <div v-if="!(loginToken === data.secret_id)" class="bg-white border-t border-gray-200 px-6 py-4">
      <div class="flex flex-col sm:flex-row sm:justify-end gap-2">
        <FullScreenModal>
          <template #trigger="{ open }">
            <button type="button" @click="open"
              class="inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition-colors duration-200">
              <i class="w-4 h-4 i-tabler-trash mr-2" />
              Delete Token
            </button>
          </template>
          <template #default="{ close }">
            <DeleteConfirm :close="close" :hints="tokenDeleteHints" :deleted-elements="deletedElements"
              :onDelete="deleteToken" />
          </template>
        </FullScreenModal>
        <button v-if="!hasExclusivePolicy" type="button" @click="saveToken"
          class="inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors duration-200">
          <i class="w-4 h-4 i-tabler-device-floppy mr-2" />
          Save Changes
        </button>
      </div>
    </div>
  </div>
</template>
