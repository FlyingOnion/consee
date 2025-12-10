<script setup lang="ts">
import { ref } from "vue";
import { uuidv7 } from "uuidv7";
import { uuidRegexp, type CreateTokenRequest } from "../../common/kz";
import { toast } from "vue3-toastify";
import PolicySelectAll from "./PolicySelectAll.vue";
import PolicyRules from "./PolicyRules.vue";
import emitter from "../../common/mitt";
import { aclTokenCreate } from "../../common/alova";

interface Props {
  close: () => void;
}

const props = defineProps<Props>();

const accessorId = ref("");
const token = ref("");
const name = ref("");

type PolicyMode = "common" | "exclusive";

const applyPolicyMode = ref<PolicyMode>("common");
const policyselect = ref<InstanceType<typeof PolicySelectAll> | null>(null);
const policyrules = ref<InstanceType<typeof PolicyRules> | null>(null);

function createToken() {
  const id = accessorId.value.trim();
  if (id && !uuidRegexp.test(id)) {
    toast.error("Invalid AccessorID");
    return;
  }

  const secretId = token.value.trim();
  if (secretId && !uuidRegexp.test(secretId)) {
    toast.error("Invalid SecretID");
    return;
  }

  // const data = {
  //   accessor_id: accessorId.value.trim(),
  //   secret_id: token.value.trim(),
  //   name: name.value.trim(),
  // };
  const basicData = {
    accessor_id: id,
    secret_id: secretId,
    name: name.value.trim(),
    policy_mode: applyPolicyMode.value || "common",
  };
  let data: CreateTokenRequest;
  switch (applyPolicyMode.value) {
    case "exclusive":
      data = {
        ...basicData,
        rules: policyrules.value!.rules,
      };
      break;
    case "common":
      data = {
        ...basicData,
        policies: policyselect.value?.selected?.map(({ id }) => id) || [],
      };
      break;
    default:
      toast.error("Invalid apply mode");
      return;
  }
  aclTokenCreate(data)
    .then(() => {
      toast.success("Token created successfully");
      props.close();
      emitter.emit("tokenCreate");
    })
    .catch((error) => {
      toast.error(error);
    });
}
</script>

<template>
  <div
    class="relative bg-white p-4 rounded-lg box-border w-150 max-[600px]:w-full flex flex-col gap-2"
  >
    <span i-tabler-x cursor-pointer absolute right-4 top-4 @click="close" />
    <h3 m-0 text-center text-base font-semibold text-gray-900>New Token</h3>

    <div flex flex-col gap-1>
      <label flex-grow text-gray-700 text-sm font-bold for="accessor_id"> Token AccessorID </label>
      <label text-gray-500 text-xs for="accessor_id">
        Token identifier, an 8-4-4-4-12 UUID. You can leave it empty, and consee will generate an id
        for you.
      </label>
    </div>

    <div flex items-center gap-2>
      <input
        v-model="accessorId"
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
        id="accessor_id"
        type="text"
        placeholder="Token accessor id"
      />
      <button
        px-2
        py-1
        rounded
        cursor-pointer
        not-disabled:hover:bg-gray-2
        @click="accessorId = uuidv7()"
      >
        Random
      </button>
    </div>

    <div flex flex-col gap-1>
      <label flex-grow text-gray-700 text-sm font-bold for="token"> Token SecretID </label>
      <label text-gray-500 text-xs for="token">
        The access token itself, an 8-4-4-4-12 UUID. You can leave it empty, and consee will
        generate a token for you.
      </label>
    </div>

    <div flex items-center gap-2>
      <input
        v-model="token"
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
        id="token"
        type="text"
        placeholder="Token secret id"
      />
      <button
        px-2
        py-1
        rounded
        cursor-pointer
        not-disabled:hover:bg-gray-2
        @click="token = uuidv7()"
      >
        Random
      </button>
    </div>

    <div flex flex-col gap-1>
      <label text-gray-700 text-sm font-bold for="name"> Token Name </label>
      <label text-gray-500 text-xs for="name">
        <strong>(Suggested)</strong> Human readable name for the token, like "kv-database-write" or
        "node-foo-identity".
      </label>
      <label text-gray-500 text-xs for="name">
        You can leave it empty, and consee will generate a name for you.
      </label>
    </div>

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
      placeholder="Token name"
    />

    <p my-0 flex-grow text-gray-700 text-sm font-bold>Policy</p>

    <select
      v-model="applyPolicyMode"
      flex-grow
      shadow
      border
      rounded
      px-1
      py-1
      text-gray-500
      leading-tight
      focus:outline-blue
    >
      <option value="common">Apply existing policies</option>
      <option value="exclusive">Create an exclusive policy</option>
    </select>

    <keep-alive>
      <PolicySelectAll
        v-if="applyPolicyMode === 'common'"
        ref="policyselect"
        excl-status="non-exclusive"
      />
    </keep-alive>
    <keep-alive>
      <PolicyRules v-if="applyPolicyMode === 'exclusive'" ref="policyrules" small-text />
    </keep-alive>

    <label flex-grow text-gray-700 text-sm font-bold> Role </label>
    <select
      flex-grow
      shadow
      border
      rounded
      px-1
      py-1
      text-gray-500
      leading-tight
      focus:outline-blue
    >
      <option value="">TODO: Select a role</option>
      <option value="role1">Role 1</option>
      <option value="role2">Role 2</option>
      <option value="role3">Role 3</option>
    </select>

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
        @click="createToken"
      >
        Save
      </button>
    </div>
  </div>
</template>
