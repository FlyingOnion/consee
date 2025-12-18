<script setup lang="ts">
import { computed, ref } from 'vue';
import { wImport, type ImportResponse } from '../common/alova';
import { toast } from 'vue3-toastify';

// Import page - content placeholder
const fileinput = ref<HTMLInputElement | null>(null);
const selectedFile = ref<File | undefined>(undefined);

const response = ref<ImportResponse>();
const nErrors = computed(() => {
  return response.value ? response.value.errors.length : 0;
})
const nConflicts = computed(() => {
  return response.value ? response.value.conflicts.length : 0;
})

function handleFileChange(event: Event) {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  if (!file) {
    return;
  }
  selectedFile.value = file;
  response.value = undefined;
}

function handleDrop(event: DragEvent) {
  event.preventDefault();
  const file = event.dataTransfer?.files[0];
  if (file && (file.type === 'application/json' || file.name.endsWith('.zip'))) {
    selectedFile.value = file;
    response.value = undefined;
  }
}

function clearFile() {
  selectedFile.value = undefined;
  response.value = undefined;
  if (fileinput.value) {
    fileinput.value.value = '';
  }
}

function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

function importFile() { doImport() }
function importDryrun() { doImport(true) }

function doImport(dryrun?: boolean) {
  if (!selectedFile.value) {
    return;
  }
  wImport(selectedFile.value, dryrun).then((r) => {
    response.value = r
  }).catch((e: Error) => {
    toast.error(e);
  })
}

</script>

<template>
  <div class="flex-grow flex flex-col p-6 gap-3 overflow-y-auto">
    <h1 class="text-2xl font-semibold text-gray-900">Import Resources</h1>
    <p class="text-gray-600">Import your data from JSON or archived file</p>
    <input ref="fileinput" type="file" class="hidden" accept=".json,.zip" @change="handleFileChange" />

    <div v-if="selectedFile" class="flex flex-col p-4 gap-4 bg-blue-50 rounded-lg border border-blue-200">
      <div class="flex items-center justify-between">
        <div class="flex items-center">
          <svg class="w-8 h-8 text-blue-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z">
            </path>
          </svg>
          <div>
            <p class="font-medium text-gray-900">{{ selectedFile.name }}</p>
            <p class="text-sm text-gray-600">{{ formatFileSize(selectedFile.size) }}</p>
          </div>
        </div>
        <button @click="clearFile" class="text-red-500 hover:text-red-700 transition-colors cursor-pointer">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>
      <div class="flex flex-col md:flex-row gap-2">
        <button
          class="px-4 py-2 bg-blue-500 text-white font-medium rounded hover:bg-blue-6 transition-all duration-200 shadow-lg hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed"
          :disabled="response !== undefined" @click="importDryrun">
          <span class="flex items-center">
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z">
              </path>
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z">
              </path>
            </svg>
            Import Preview
          </span>
        </button>
        <button
          class="px-4 py-2 bg-green-500 text-white font-medium rounded hover:bg-green-600 transition-all duration-200 shadow-lg hover:shadow-xl disabled:opacity-50 disabled:cursor-not-allowed"
          @click="importFile">
          <span class="flex items-center">
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
            </svg>
            Import Data
          </span>
        </button>
      </div>
    </div>
    <div v-else
      class="border-2 border-dashed border-gray-300 rounded-xl p-12 text-center hover:border-blue-400 transition-colors duration-200 cursor-pointer bg-gray-50 hover:bg-blue-50"
      @click="fileinput?.click()" @dragover.prevent @drop.prevent="handleDrop">

      <div class="flex flex-col items-center">
        <svg class="w-16 h-16 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path>
        </svg>
        <h3 class="text-lg font-semibold text-gray-700 mb-2">
          Click to browse or drag and drop
        </h3>
        <p class="text-sm text-gray-500">Support for .json and .zip files</p>
      </div>
    </div>



    <div v-if="response" class="bg-gray-50 rounded-lg p-8">
      <h3 class="text-2xl font-bold text-gray-900 mb-6">Import Results</h3>

      <!-- Summary -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div class="bg-green-50 rounded-lg p-4 border border-green-200">
          <div class="flex items-center">
            <svg class="w-8 h-8 text-green-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            <div>
              <p class="text-2xl font-bold text-green-700">{{
                response.successes.keys + response.successes.policies + response.successes.tokens || 0 }}</p>
              <p class="text-sm text-green-600">Imported Successfully</p>
            </div>
          </div>
        </div>

        <div v-if="nConflicts > 0" class="bg-yellow-50 rounded-lg p-4 border border-yellow-200">
          <div class="flex items-center">
            <svg class="w-8 h-8 text-yellow-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z">
              </path>
            </svg>
            <div>
              <p class="text-2xl font-bold text-yellow-700">{{ nConflicts }}</p>
              <p class="text-sm text-yellow-600">Conflicts Found</p>
            </div>
          </div>
        </div>

        <div v-if="nErrors > 0" class="bg-red-50 rounded-lg p-4 border border-red-200">
          <div class="flex items-center">
            <svg class="w-8 h-8 text-red-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            <div>
              <p class="text-2xl font-bold text-red-700">{{ nErrors }}</p>
              <p class="text-sm text-red-600">Errors Found</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Conflicts -->
      <div v-if="nConflicts > 0" class="mb-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
          <svg class="w-5 h-5 text-yellow-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z">
            </path>
          </svg>
          Conflicts
        </h3>
        <div class="bg-yellow-50 rounded-lg border border-yellow-200">
          <div v-for="(conflict, index) in response.conflicts" :key="index"
            class="p-3 border-b border-yellow-200 last:border-b-0">
            <div class="flex items-center justify-between">
              <div>
                <span
                  class="inline-block px-2 py-1 bg-yellow-200 text-yellow-800 text-xs font-medium rounded-full mr-2">
                  {{ conflict.kind }}
                </span>
                <span class="text-gray-900 font-medium">{{ conflict.param }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Errors -->
      <div v-if="nErrors > 0">
        <h3 class="text-lg font-semibold text-gray-900 mb-3 flex items-center">
          <svg class="w-5 h-5 text-red-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
          Errors
        </h3>
        <div class="bg-red-50 rounded-lg border border-red-200">
          <div v-for="(error, index) in response.errors" :key="index"
            class="p-3 border-b border-red-200 last:border-b-0">
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <div class="flex items-center mb-1">
                  <span class="inline-block px-2 py-1 bg-red-200 text-red-800 text-xs font-medium rounded-full mr-2">
                    {{ error.kind }}
                  </span>
                  <span class="text-gray-900 font-medium">{{ error.param }}</span>
                </div>
                <p class="text-gray-700 text-sm">{{ error.cause }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Success Message -->
      <div v-if="nErrors === 0 && nConflicts === 0"
        class="bg-green-50 rounded-lg border border-green-200 p-6 text-center">
        <svg class="w-16 h-16 text-green-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        <h3 class="text-xl font-bold text-green-800 mb-2">Import Successful!</h3>
        <p class="text-green-700">Your data has been imported successfully without any issues.</p>
      </div>
    </div>
  </div>
</template>