<script setup lang="ts">
import TokenForm from "./acl/TokenForm.vue";
import { useRoute, useRouter } from "vue-router";
import { aclTokenGetDetail, aclTokenList } from "../common/alova";
import { computed, ref, watch } from "vue";
import {
  alova,
  debounce,
  type ACLLink,
  type TokenDetailInfo
} from "../common/kz";
import FullScreenModal from "./common/FullScreenModal.vue";
import TokenList from "./acl/TokenList.vue";
import SplitView from "./common/SplitView.vue";
import EmptyView from "./common/EmptyView.vue";
import { toast } from "vue3-toastify";
import { invalidateCache } from "alova";
import TokenDetail from "./acl/TokenDetail.vue";
import emitter from "../common/mitt";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const route = useRoute();
const router = useRouter();
const currentToken = ref((route.params.id as string) || "");

const loading = ref(true);
const error = ref<Error | undefined>(undefined);
const tokenList = ref<ACLLink[]>([]);

aclTokenList()
  .then((data) => {
    tokenList.value = data;
    loading.value = false;
  })
  .catch((e: Error) => {
    currentToken.value = "";
    loading.value = false;
    error.value = e;
    toast.error(e);
  });

const currentTokenName = computed(() => {
  return currentToken.value && tokenList.value
    ? tokenList.value.find((item: ACLLink) => item.id === currentToken.value)?.name
    : "";
});

emitter.on("tokenCreate", refresh);
emitter.on("tokenDelete", refresh);

function refresh() {
  invalidateCache(alova.snapshots.match("aclTokenList"));
  setTimeout(() => {
    aclTokenList()
      .then((data) => {
        tokenList.value = data;
        if (
          currentToken.value &&
          !tokenList.value.map(({ id }) => id).includes(currentToken.value)
        ) {
          currentToken.value = "";
          router.replace("/ui/acl/tokens");
        }
      })
      .catch((e: Error) => {
        toast.warn(t("tokenPage.refreshFailed", { message: e.message }));
      });
  }, 300);
}

function onItemClick(token: string) {
  currentToken.value = token;
  router.push(`/acl/token/${token}`);
}

const detail = ref<TokenDetailInfo | undefined>(undefined);
function setDetail() {
  if (!currentToken.value) {
    detail.value = undefined;
    return;
  }
  aclTokenGetDetail(currentToken.value)
    .then((data: TokenDetailInfo) => {
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

watch(currentToken, debouncedSetDetail);

const splitview = ref<InstanceType<typeof SplitView> | null>(null);
</script>

<template>
  <SplitView ref="splitview" :loading :error :title="t('tokenPage.title')" :current="currentTokenName">
    <template #controls>
      <FullScreenModal>
        <template #trigger="{ open }">
          <button @click="open"
            class="flex items-center justify-center p-2 gap-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition-colors duration-200 shadow-sm"
            :title="t('tokenPage.createNewToken')">
            <i class="w-4 h-4 i-tabler-plus" />
            <span class="text-sm font-medium hidden sm:inline">{{ t('tokenPage.newToken') }}</span>
          </button>
        </template>
        <template #default="{ close }">
          <TokenForm :close="close" />
        </template>
      </FullScreenModal>
    </template>
    <template #list>
      <TokenList :mobile="splitview?.isMobile" et="token" :t="currentToken" :data="tokenList" @item-click="onItemClick" />
    </template>

    <!-- Token Detail Section -->
    <div v-if="detail" class="flex flex-col">
      <div class="bg-white rounded-lg shadow-sm border border-gray-200 flex-1 overflow-hidden">
        <TokenDetail :data="detail" />
      </div>
    </div>

    <template #empty>
      <EmptyView :resource="t('tokenPage.resource')">
        <p class="text-sm">{{ t('tokenPage.selectPrompt') }}</p>
      </EmptyView>
    </template>
  </SplitView>
</template>
