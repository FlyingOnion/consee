<script setup lang="ts">
import { aclPolicyList } from "../../common/alova";
import type { ACLLink } from "../../common/kz";
import PolicySelect from "./PolicySelect.vue";
import { computed, ref } from "vue";
import { toast } from "vue3-toastify";

interface Props {
  exclStatus?: "" | "exclusive" | "non-exclusive" | "all";
  // allPolicies: ACLLink[];
  selectedPolicies?: string[];
}

const props = defineProps<Props>();
const exclusiveParam =
  props.exclStatus === "exclusive"
    ? { exclusive: "1" }
    : props.exclStatus === "non-exclusive"
      ? { exclusive: "0" }
      : undefined;

const loading = ref(true);
const list = ref<ACLLink[]>([]);

aclPolicyList(exclusiveParam).then((data) => {
  list.value = data;
}).catch((e: Error) => {
  toast.error(e);
}).finally(() => {
  loading.value = false;
});

const policyselect = ref<InstanceType<typeof PolicySelect> | null>(null);
const selected = computed(() => policyselect.value?.selected);

defineExpose({ selected });
</script>

<template>
  <select v-if="loading" disabled flex-grow shadow border rounded px-1 py-1 text-gray-500 leading-tight
    focus:outline-blue>
    <option hidden value="">Fetching policies ...</option>
  </select>
  <PolicySelect v-else-if="list.length" ref="policyselect" :all-policies="list" :selected-policies />
  <div v-else>
    <p text-bluegray-5>No policies found</p>
  </div>
</template>
