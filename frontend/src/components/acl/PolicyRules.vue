<script setup lang="ts">
import { computed, ref, watch } from "vue";
import {
  policyRuleTypeList,
  policyRuleWithParamTypeSet,
  type PolicyFormRule,
  type PolicyFormRuleListElement,
} from "../../common/kz";
import shortid from "shortid";

interface Props {
  smallText?: boolean;
  rules?: PolicyFormRule[];
  readonly?: boolean;
}

const props = defineProps<Props>();

const preview = ref(false);
const newElement = ref(false);

watch(
  () => props.rules,
  (newVal) => {
    ruleList.value = initRuleList(newVal);
  }
);

function initRuleList(rules?: PolicyFormRule[]): PolicyFormRuleListElement[] {
  return rules ? rules.map((rule) => ({ id: shortid.generate(), rule })) : [];
}
const ruleList = ref<PolicyFormRuleListElement[]>(initRuleList(props.rules));

function ruleListElementCompare(
  a: PolicyFormRuleListElement,
  b: PolicyFormRuleListElement
): number {
  const rtypeCmp = a.rule.rtype < b.rule.rtype ? -1 : a.rule.rtype > b.rule.rtype ? 1 : 0;
  const matchCmp =
    a.rule.match && b.rule.match
      ? a.rule.match === b.rule.match
        ? 0
        : a.rule.match === "prefix"
        ? 1
        : -1
      : 0;
  const paramCmp =
    a.rule.param && b.rule.param
      ? a.rule.param < b.rule.param
        ? -1
        : a.rule.param > b.rule.param
        ? 1
        : 0
      : 0;
  const actionCmp =
    a.rule.access === b.rule.access
      ? 0
      : a.rule.access === "deny"
      ? 1
      : b.rule.access === "deny"
      ? -1
      : a.rule.access === "write"
      ? 1
      : -1;
  return rtypeCmp || matchCmp || paramCmp || actionCmp;
}

// for new element
function resetNewRuleElement() {
  newRuleElement.value.rtype = "";
  newRuleElement.value.match = undefined;
  newRuleElement.value.param = undefined;
  newRuleElement.value.access = "read";
}
const newRuleElement = ref<PolicyFormRule>({ rtype: "", access: "read" });
const withParam = computed(() => policyRuleWithParamTypeSet.has(newRuleElement.value.rtype));

function addNewRule() {
  if (!newRuleElement.value.rtype) {
    resetNewRuleElement();
    newElement.value = false;
    return;
  }

  ruleList.value.push({
    id: shortid.generate(),
    rule: {
      rtype: newRuleElement.value.rtype,
      match: withParam.value
        ? !newRuleElement.value.match || newRuleElement.value.match === "exact"
          ? "exact"
          : "prefix"
        : undefined,
      param: withParam.value ? newRuleElement.value.param || "" : undefined,
      access: newRuleElement.value.access,
    },
  });
  ruleList.value.sort(ruleListElementCompare);
  resetNewRuleElement();
}

function cancelAddNewRule() {
  resetNewRuleElement();
  newElement.value = false;
}

function removeRule(id: string) {
  ruleList.value = ruleList.value.filter((r) => r.id !== id);
}

function rule2String(rule: PolicyFormRule): string {
  if (!policyRuleWithParamTypeSet.has(rule.rtype)) {
    return `${rule.rtype} = "${rule.access}"`;
  }
  const header = rule.match && rule.match !== "exact" ? `${rule.rtype}_prefix` : rule.rtype;
  const param = rule.match !== "all" ? rule.param : "";
  return `${header} "${param}" {
  policy = "${rule.access}"
}`;
}

const rules = computed(() => ruleList.value.map((r) => rule2String(r.rule)).join("\n\n"));

defineExpose({ rules });
</script>

<template>
  <div flex items-center gap-2>
    <p class="my-0 text-gray-700 font-bold" :class="{ 'text-sm': smallText }">Policy Rule</p>
    <p v-if="preview" class="my-0 text-bluegray-5" :class="{ 'text-sm': smallText }">
      (Preview Mode)
    </p>
    <div ml-auto flex gap-2>
      <span
        v-show="preview"
        i-tabler-eye-off
        cursor-pointer
        title="Exit preview mode"
        @click="preview = false"
      />
      <span v-if="!(preview || readonly)" cursor-pointer i-tabler-plus @click="newElement = true" />
      <span
        v-show="!preview"
        i-tabler-eye
        cursor-pointer
        title="Preview mode"
        @click="preview = true"
      />
    </div>
  </div>

  <div
    class="relative grid grid-items-baseline gap-2"
    :class="
      readonly
        ? 'grid-cols-[repeat(4,minmax(100px,1fr))]'
        : 'grid-cols-[repeat(4,minmax(100px,1fr))_44px]'
    "
  >
    <p class="my-0 text-gray-700 font-bold" :class="{ 'text-sm': smallText }">Resource</p>
    <p class="my-0 text-gray-700 font-bold" :class="{ 'text-sm': smallText }">Match</p>
    <p class="my-0 text-gray-700 font-bold" :class="{ 'text-sm': smallText }">Param</p>
    <p class="my-0 text-gray-700 font-bold" :class="{ 'text-sm': smallText }">Access</p>
    <p v-if="!readonly" my-0></p>
    <p
      v-if="!ruleList.length"
      class="my-0 text-bluegray-5 col-span-full"
      :class="{ 'text-sm': smallText }"
    >
      No rules defined yet.
    </p>
    <template v-for="r in ruleList" :key="r.id">
      <p class="my-0 text-gray-5" :class="{ 'text-sm': smallText }">{{ r.rule.rtype }}</p>
      <p class="my-0 text-gray-5" :class="{ 'text-sm': smallText }">{{ r.rule.match || "" }}</p>
      <p v-if="r.rule.match === 'prefix' && !r.rule.param" my-0 text-sm text-bluegray-5>
        (match all
        {{
          r.rule.rtype === "identity"
            ? "identities"
            : r.rule.rtype === "query"
            ? "queries"
            : r.rule.rtype === "mesh"
            ? "meshes"
            : r.rule.rtype + "s"
        }})
      </p>
      <p v-else class="my-0 text-gray-5" :class="{ 'text-sm': smallText }">
        {{ r.rule.param || "" }}
      </p>
      <p class="my-0 text-gray-5" :class="{ 'text-sm': smallText }">{{ r.rule.access }}</p>
      <span v-if="!readonly" ml-auto i-tabler-trash @click="removeRule(r.id)" />
    </template>
    <template v-if="newElement">
      <select v-model="newRuleElement.rtype" px-1 py-1>
        <option value="" hidden></option>
        <option v-for="ruleType in policyRuleTypeList" :value="ruleType">{{ ruleType }}</option>
      </select>

      <select :disabled="!withParam" v-model="newRuleElement.match" px-1 py-1>
        <option value="exact">exact</option>
        <option value="prefix">prefix</option>
        <option value="all">all</option>
      </select>
      <input
        :disabled="!withParam || newRuleElement.match === 'all'"
        px-1
        py-1
        v-model="newRuleElement.param"
        type="text"
        placeholder="key or prefix"
      />
      <select v-model="newRuleElement.access" px-1 py-1>
        <option value="read">read</option>
        <option value="write">write</option>
        <option value="deny">deny</option>
      </select>
      <div self-stretch flex flex-row-reverse items-center gap-2>
        <span i-tabler-x @click="cancelAddNewRule" />
        <span v-show="newRuleElement.rtype" i-tabler-check @click="addNewRule" />
      </div>
    </template>

    <div v-if="preview" absolute top-0 bottom-0 left-0 right-0 bg-white overflow-y-scroll>
      <pre v-if="ruleList.length" my-0 p-1><code>{{ rules }}</code></pre>
      <p v-else class="my-0 text-bluegray-5" :class="{ 'text-sm': smallText }">
        No rules defined yet.
      </p>
    </div>
  </div>
</template>
