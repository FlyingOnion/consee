<template>
  <a v-if="!doNotShow" relative flex items-center p-2 bg-gray-50 hover:bg-gray-2 rounded cursor-pointer
    @click="showPremiumToast">
    <i i-tabler-bell w-4 h-4 />
  </a>
</template>

<script setup lang="ts">
import { h, ref } from "vue";
import { nShowNotification } from "../../../common/const";
import { toast } from "vue3-toastify";
import ShowPremiumToastContent from "./ShowPremiumToastContent.vue";

const doNotShow = ref(localStorage.getItem(nShowNotification) === "1");
const toastId = "premium-toast";

function ok() {
  toast.remove(toastId);
}

function hide() {
  console.log("hide");
  doNotShow.value = true;
  localStorage.setItem(nShowNotification, "1");
  toast.remove(toastId);
}

function showPremiumToast() {
  toast(h(ShowPremiumToastContent, { ok, hide }), {
    closeOnClick: false,
    toastId,
  });
}
</script>
