<script setup lang="ts">
import { computed, ref, watch } from "vue";
import { type ACLLink } from "../../common/kz";

interface Props {
  allPolicies: ACLLink[];
  selectedPolicies?: string[];
}

const props = defineProps<Props>();

const policyIdNameMap = ref<Map<string, string>>(
  new Map(props.allPolicies.map((item) => [item.id, item.name]))
);
const selectedPolicies = ref<string[]>(props.selectedPolicies || []);
const availablePolicies = ref<string[]>(
  props.allPolicies.map(({ id }) => id).filter((id) => !selectedPolicies.value.includes(id))
);

const selected = computed(() =>
  selectedPolicies.value
    .map((id) => ({ id, name: policyIdNameMap.value.get(id) || `unknown policy ${id}` }))
    .sort((a, b) => a.name.localeCompare(b.name))
);
const options = computed(() =>
  availablePolicies.value
    .map((id) => ({ id, name: policyIdNameMap.value.get(id) || `unknown policy ${id}` }))
    .sort((a, b) => a.name.localeCompare(b.name))
);

function addPolicy(id: string) {
  selectedPolicies.value.push(id);
  availablePolicies.value = availablePolicies.value.filter((item) => item !== id);
  pSelect.value = "";
}

function removePolicy(id: string) {
  availablePolicies.value.push(id);
  selectedPolicies.value = selectedPolicies.value.filter((item) => item !== id);
}

const pSelect = ref("");

defineExpose({ selected });

function resetOnSelectedChanged(newVal?: string[]) {
  selectedPolicies.value = newVal || [];
  availablePolicies.value = props.allPolicies
    .map(({ id }) => id)
    .filter((id) => !selectedPolicies.value.includes(id));
}

watch(() => props.selectedPolicies, resetOnSelectedChanged);

function resetOnAllPoliciesChanged(newVal: ACLLink[]) {
  policyIdNameMap.value = new Map(newVal.map((item) => [item.id, item.name]));
  availablePolicies.value = newVal
    .map(({ id }) => id)
    .filter((id) => !selectedPolicies.value.includes(id));
}

watch(() => props.allPolicies, resetOnAllPoliciesChanged);
</script>

<template>
  <select
    :disabled="!availablePolicies.length"
    v-model="pSelect"
    flex-grow
    shadow
    border
    rounded
    px-1
    py-1
    text-gray-500
    leading-tight
    focus:outline-blue
    @change="addPolicy(pSelect)"
  >
    <option v-if="!availablePolicies.length" value="" hidden>No more policies available</option>
    <template v-else>
      <option value="" hidden>Select a policy</option>
      <option v-for="{ id, name } in options" :key="id" :value="id">{{ name }}</option>
    </template>
  </select>
  <div flex flex-wrap gap-2>
    <div v-for="{ id, name } in selected" flex items-center bg-gray-2 hover:bg-gray-3 rounded>
      <p my-0 px-1 py-1 text-sm text-bluegray-5>
        {{ name }}
      </p>
      <span i-tabler-x cursor-pointer text-bluegray-5 @click="removePolicy(id)"></span>
    </div>
  </div>
</template>
