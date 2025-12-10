<!--
  Copyright (c) 2025 The Consee Authors. All rights reserved.
  SPDX-License-Identifier: MulanPSL-2.0
-->

<script setup lang="ts">
import PageHeader from "./components/PageHeader.vue";

</script>

<template>
  <div class="flex flex-col relative w-100vw h-100vh mx-auto max-w-full lg:max-w-1024px xl:max-w-80vw">
    <PageHeader />
    <router-view />
    <Hello />
  </div>
</template>
