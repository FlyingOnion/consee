<script setup lang="ts">
import type { ACLLink } from "../../common/kz";

interface Props {
  t: string;
  data: ACLLink[];
  onItemClick: (token: string) => void;
}

defineProps<Props>();
</script>

<template>
  <div class="space-y-1 p-2">
    <div
      v-for="token in data"
      :key="token.id"
      @click="onItemClick(token.id)"
      :class="[
        'group cursor-pointer transition-all duration-200 rounded-lg border',
        token.id === t
          ? 'bg-blue-50 border-blue-200 shadow-sm'
          : 'bg-white border-gray-200 hover:bg-gray-50 hover:border-gray-300',
      ]"
    >
      <div class="p-3" :title="token.name">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-3">
            <div
              :class="[
                'w-2 h-2 rounded-full transition-colors duration-200',
                token.id === t ? 'bg-blue-500' : 'bg-gray-300 group-hover:bg-blue-400',
              ]"
            />
            <div class="min-w-0 flex-1">
              <p
                :class="[
                  'text-sm font-medium truncate',
                  token.id === t ? 'text-blue-900' : 'text-gray-900',
                ]"
              >
                {{ token.name.length > 32 ? token.name.slice(0, 32) + "..." : token.name }}
              </p>
              <p :class="['text-xs truncate', token.id === t ? 'text-blue-700' : 'text-gray-500']">
                ID: {{ token.id }}
              </p>
            </div>
          </div>
          <i
            v-if="token.id === t"
            class="w-4 h-4 i-tabler-chevron-right text-blue-500 transition-transform duration-200 group-hover:translate-x-0.5"
          />
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="!data.length" class="text-center py-8">
      <i class="w-12 h-12 i-tabler-list-search mx-auto mb-3 text-gray-300" />
      <p class="text-gray-500 text-sm">No tokens available</p>
    </div>
  </div>
</template>
