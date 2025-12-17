<script setup lang="ts">
import type { ACLLink } from "../../common/kz";

interface Props {
  mobile?: boolean;
  current: string;
  data: ACLLink[];
  onItemClick: (current: string) => void;
}

defineProps<Props>();
</script>

<template>
  <div class="space-y-1 p-2">
    <div
      v-for="role in data"
      :key="role.name"
      @click="onItemClick(role.name)"
      :class="[
        'group cursor-pointer transition-all duration-200 rounded-lg border',
        role.name === current
          ? 'bg-blue-50 border-blue-200 shadow-sm'
          : 'bg-white border-gray-200 hover:bg-gray-50 hover:border-gray-300',
      ]"
    >
      <div class="p-3" :title="role.name">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-3">
            <div
              :class="[
                'w-2 h-2 rounded-full transition-colors duration-200',
                role.name === current ? 'bg-blue-500' : 'bg-gray-300 group-hover:bg-blue-400',
              ]"
            />
            <div class="min-w-0 flex-1">
              <p
                :class="[
                  'text-sm font-medium truncate',
                  role.name === current ? 'text-blue-900' : 'text-gray-900',
                ]"
              >
                {{ role.name.length > 32 && !mobile ? role.name.slice(0, 32) + "..." : role.name }}
              </p>
              <p
                :class="[
                  'text-xs truncate',
                  role.name === current ? 'text-blue-700' : 'text-gray-500',
                ]"
              >
                ID: {{ role.id }}
              </p>
            </div>
          </div>
          <i
            v-if="role.name === current"
            class="w-4 h-4 i-tabler-chevron-right text-blue-500 transition-transform duration-200 group-hover:translate-x-0.5"
          />
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-if="!data.length" class="text-center py-8">
      <i class="w-12 h-12 i-tabler-list-search mx-auto mb-3 text-gray-300" />
      <p class="text-gray-500 text-sm">No roles available</p>
    </div>
  </div>
</template>