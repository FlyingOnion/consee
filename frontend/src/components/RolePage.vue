<script setup lang="ts">
import { useRoute, useRouter } from "vue-router";

import { ref, watch } from "vue";
import { aclRoleGetDetail, aclRoleList } from "../common/alova";
import {
  alova,
  type ACLLink,
  type RoleDetailInfo,
  debounce,
  b64Encode,
  b64Decode,
} from "../common/kz";
import RoleList from "./acl/RoleList.vue";
import RoleDetail from "./acl/RoleDetail.vue";
import SplitView from "./common/SplitView.vue";
import EmptyView from "./common/EmptyView.vue";
import RoleForm from "./acl/RoleForm.vue";
import FullScreenModal from "./common/FullScreenModal.vue";
import { toast } from "vue3-toastify";
import { invalidateCache } from "alova";
import emitter from "../common/mitt";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const route = useRoute();
const router = useRouter();

function initialName() {
  const b64name = route.params.id as string;
  return b64name ? b64Decode(b64name) || "" : "";
}

const currentRoleName = ref(initialName());

const loading = ref(true);
const error = ref<Error | undefined>(undefined);
const roleList = ref<ACLLink[]>([]);

aclRoleList()
  .then((data) => {
    roleList.value = data;
    loading.value = false;
  })
  .catch((e: Error) => {
    currentRoleName.value = "";
    loading.value = false;
    error.value = e;
    toast.error(e);
  });

emitter.on("roleCreate", refresh);
emitter.on("roleDelete", refresh);

function refresh() {
  invalidateCache(alova.snapshots.match("aclRoleList"));
  setTimeout(() => {
    aclRoleList()
      .then((data) => {
        roleList.value = data;
        if (
          currentRoleName.value &&
          !roleList.value.map(({ name }) => name).includes(currentRoleName.value)
        ) {
          currentRoleName.value = "";
          router.replace("/acl/roles");
        }
      })
      .catch((e: Error) => {
        toast.warn(t("rolePage.refreshFailed", { message: e.message }));
      });
  }, 300);
}

function onItemClick(roleName: string) {
  currentRoleName.value = roleName;
  router.push(`/acl/roles/${b64Encode(roleName)}`);
}

const detail = ref<RoleDetailInfo | undefined>(undefined);
function setDetail() {
  if (!currentRoleName.value) {
    detail.value = undefined;
    return;
  }
  aclRoleGetDetail(b64Encode(currentRoleName.value))
    .then((data: RoleDetailInfo) => {
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

watch(currentRoleName, debouncedSetDetail);

const splitview = ref<InstanceType<typeof SplitView> | null>(null);
</script>

<template>
  <SplitView ref="splitview" :loading :error :title="t('rolePage.title')" :current="currentRoleName">
    <template #controls>
      <FullScreenModal>
        <template #trigger="{ open }">
          <button @click="open"
            class="flex items-center justify-center p-2 gap-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition-colors duration-200 shadow-sm"
            :title="t('rolePage.createNewRole')">
            <i class="w-4 h-4 i-tabler-plus" />
            <span class="text-sm font-medium hidden sm:inline">{{ t('rolePage.newRole') }}</span>
          </button>
        </template>
        <template #default="{ close }">
          <RoleForm :close="close" />
        </template>
      </FullScreenModal>
    </template>
    <template #list>
      <RoleList :mobile="splitview?.isMobile" :current="currentRoleName" :data="roleList" @item-click="onItemClick" />
    </template>

    <!-- Role Detail Section -->
    <div v-if="detail" class="flex flex-col">
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 flex-1 overflow-hidden">
        <RoleDetail :data="detail" />
      </div>
    </div>

    <template #empty>
      <EmptyView :resource="t('rolePage.resource')">
        <p class="text-sm">{{ t('rolePage.selectPrompt') }}</p>
      </EmptyView>
    </template>
  </SplitView>
</template>
