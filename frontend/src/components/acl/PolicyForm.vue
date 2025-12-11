<script setup lang="ts">
import { inject, ref } from "vue";
import PolicyRules from "./PolicyRules.vue";
import { conseeErrorKey, conseeTokenKey } from "../../common/const";
import { alova } from "../../common/kz";
import { toast } from "vue3-toastify";

interface Props {
  close: () => void;
}
const props = defineProps<Props>();
const name = ref("");
const policyrules = ref<InstanceType<typeof PolicyRules> | null>(null);

const refreshPolicyList = inject<() => void>("refreshPolicyList")!;

async function createPolicy() {
  if (!name.value) {
    toast.error("Policy name cannot be empty");
    return;
  }
  if (name.value.startsWith("--")) {
    toast.error("Policy name cannot start with --. It's a reserved prefix for exclusive policies.");
    return;
  }
  const resp = await alova.Post<Response>(
    "/acl/policy",
    {
      name: name.value,
      rules: policyrules.value?.rules || "",
    },
    {
      headers: {
        [conseeTokenKey]: localStorage.getItem(conseeTokenKey) || "",
      },
    }
  );
  if (resp.status !== 201) {
    toast.error(resp.headers.get(conseeErrorKey) || "Failed to create policy");
    return;
  }
  toast.success("Policy created");
  props.close();
  refreshPolicyList();
}
</script>

<template>
  <div
    class="relative bg-white p-4 rounded-lg box-border w-150 max-[600px]:w-full flex flex-col gap-2"
  >
    <span i-tabler-x cursor-pointer absolute right-4 top-4 @click="close" />
    <h3 m-0 text-center text-base font-semibold text-gray-900>New Policy</h3>
    <div flex flex-col gap-1>
      <label text-gray-700 text-sm font-bold for="name"> Policy Name </label>
      <input
        v-model="name"
        flex-grow
        shadow
        appearance-none
        border
        rounded
        px-2
        py-1
        text-gray-700
        leading-tight
        focus-outline-blue
        id="name"
        type="text"
        placeholder="Policy name"
      />
      <PolicyRules ref="policyrules" small-text />
    </div>
    <div bg-gray-50 px-6 py-3 flex flex-row-reverse gap-2>
      <button
        type="button"
        rounded-md
        bg-white
        px-3
        py-2
        text-sm
        font-semibold
        text-gray-900
        ring-1
        shadow-xs
        ring-gray-300
        ring-inset
        hover:bg-gray-50
        w-20
        @click="close"
      >
        Cancel
      </button>
      <button
        type="button"
        rounded-md
        bg-blue-5
        px-3
        py-2
        text-sm
        font-semibold
        text-white
        shadow-xs
        hover:bg-blue-4
        w-20
        @click="createPolicy"
      >
        Save
      </button>
    </div>
  </div>
</template>
