<script setup lang="ts">
import { useRoute, useRouter } from "vue-router";

import { ref, watch } from "vue";
import { aclPolicyGetDetail, aclPolicyList } from "../common/alova";
import {
  alova,
  type ACLLink,
  type PolicyDetailInfo,
  debounce,
  b64Encode,
  b64Decode,
} from "../common/kz";
import PolicyList from "./acl/PolicyList.vue";
import PolicyDetail from "./acl/PolicyDetail.vue";
import SplitView from "./common/SplitView.vue";
import EmptyView from "./common/EmptyView.vue";
import PolicyForm from "./acl/PolicyForm.vue";
import FullScreenModal from "./common/FullScreenModal.vue";
import { toast } from "vue3-toastify";
import { invalidateCache } from "alova";
import emitter from "../common/mitt";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const route = useRoute();
const router = useRouter();

function initialName() {
  const b64name = route.params.b64name as string;
  return b64name ? b64Decode(b64name) || "" : "";
}

const currentPolicyName = ref(initialName());

const loading = ref(true);
const error = ref<Error | undefined>(undefined);
const policyList = ref<ACLLink[]>([]);

aclPolicyList()
  .then((data) => {
    policyList.value = data;
    loading.value = false;
  })
  .catch((e: Error) => {
    currentPolicyName.value = "";
    loading.value = false;
    error.value = e;
    toast.error(e);
  });

emitter.on("policyCreate", refresh);
emitter.on("policyDelete", refresh);

function refresh() {
  invalidateCache(alova.snapshots.match("aclPolicyList"));
  setTimeout(() => {
    aclPolicyList()
      .then((data) => {
        policyList.value = data;
        if (
          currentPolicyName.value &&
          !policyList.value.map(({ name }) => name).includes(currentPolicyName.value)
        ) {
          currentPolicyName.value = "";
          router.replace("/acl/policies");
        }
      })
      .catch((e: Error) => {
        toast.warn(t("policyPage.refreshFailed", { message: e.message }));
      });
  }, 300);
}

function onItemClick(policyName: string) {
  currentPolicyName.value = policyName;
  router.push(`/acl/policy/${b64Encode(policyName)}`);
}

const detail = ref<PolicyDetailInfo | undefined>(undefined);
function setDetail() {
  if (!currentPolicyName.value) {
    detail.value = undefined;
    return;
  }
  aclPolicyGetDetail(b64Encode(currentPolicyName.value))
    .then((data: PolicyDetailInfo) => {
      detail.value = data;
      loading.value = false;
    })
    .catch((e: Error) => {
      toast.error(e);
      loading.value = false;
      error.value = e;
    });
}
setDetail();

const debouncedSetDetail = debounce(setDetail, 50);

watch(currentPolicyName, debouncedSetDetail);
</script>

<template>
  <SplitView :loading :error :title="t('policyPage.title')" :current="currentPolicyName">
    <template #controls>
      <FullScreenModal>
        <template #trigger="{ open }">
          <button @click="open"
            class="flex items-center justify-center p-2 gap-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition-colors duration-200 shadow-sm"
            :title="t('policyPage.createNewPolicy')">
            <i class="w-4 h-4 i-tabler-plus" />
            <span class="text-sm font-medium hidden sm:inline">{{ t('policyPage.newPolicy') }}</span>
          </button>
        </template>
        <template #default="{ close }">
          <PolicyForm :close="close" />
        </template>
      </FullScreenModal>
    </template>
    <template #list>
      <PolicyList :current="currentPolicyName" :data="policyList" @item-click="onItemClick" />
    </template>

    <!-- Token Detail Section -->
    <div v-if="detail" class="flex flex-col">
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 flex-1 overflow-hidden">
        <PolicyDetail :data="detail" />
      </div>
    </div>

    <template #empty>
      <EmptyView :resource="t('policyPage.newPolicy')">
        <p class="text-sm">{{ t('policyPage.selectPrompt') }}</p>
      </EmptyView>
    </template>
  </SplitView>
</template>
