<template>
  <div v-if="ll" p-4>Connecting to consee backend</div>
  <div v-else-if="lerr" p-4>{{ lerr }}</div>
  <div v-else grid grid-cols-4 gap-4 p-4 box-border h-full>
    <div class="col-span-1 max-[960px]:hidden">
      <slot name="ltitle">
        <p my-2 text-lg font-800>{{ ltitle }}</p>
      </slot>
      <slot name="list"></slot>
    </div>

    <div class="col-span-3 max-[960px]:col-span-full relative">
      <div v-if="rl">Connecting to consee backend</div>
      <div v-else-if="rerr">{{ rerr }}</div>
      <template v-else>
        <div flex gap-2 items-center>
          <p v-if="!rtitle" class="my-2 text-bluegray-6 max-[960px]:hidden">
            👈 Choose one to see details.
          </p>
          <FullScreenModal>
            <template #trigger="{ open }">
              <span
                class="i-tabler-baseline-density-medium cursor-pointer min-[961px]:hidden"
                @click="open"
              >
              </span>
            </template>
            <template #default="{ close }">
              <div flex-grow relative bg-white p-4 rounded-lg max-w-150 box-border>
                <span i-tabler-x cursor-pointer absolute right-4 top-4 @click="close" />
                <div flex justify-center>
                  <div max-w-120px>
                    <slot name="ltitle">
                      <h3 text-center text-base font-semibold text-gray-900>{{ ltitle }}</h3>
                    </slot>
                  </div>
                </div>

                <slot name="list"></slot>
              </div>
            </template>
          </FullScreenModal>
          <slot v-if="rtitle" name="rtitle">
            <p my-2 text-lg font-800>{{ rtitle }}</p>
          </slot>
          <p v-if="!rtitle" class="my-2 text-bluegray-6 min-[961px]:hidden">
            👈 Show {{ ltitle.toLowerCase() }}
          </p>
        </div>
        <slot v-if="rtitle"></slot>

        <p v-if="!rtitle" v-for="hint in hints" my-2 text-bluegray-6>{{ hint }}</p>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defaultEmptyListHints } from "../../common/kz";
import FullScreenModal from "./FullScreenModal.vue";

interface Props {
  lerr: Error | undefined;
  ll: boolean;
  rerr: Error | undefined;
  rl: boolean;
  ltitle: string;
  rtitle?: string;
  hints?: string[];
}

const props = defineProps<Props>();

const hints = props.hints || defaultEmptyListHints;
</script>
