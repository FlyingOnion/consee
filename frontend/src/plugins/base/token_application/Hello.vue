<script setup lang="ts">
import { useRouter } from 'vue-router';
import { authenticate } from '../../../common/alova';
import { conseeTokenKey } from '../../../common/const';
import emitter from '../../../common/mitt';

const router = useRouter();

if (localStorage.getItem(conseeTokenKey)) {
  authenticate().then((data) => {
    emitter.emit("login", data);
  }).catch((_: Error) => {
    localStorage.removeItem(conseeTokenKey);
    router.push("/home");
  });
} else {
  router.push("/home");
}

</script>

<template></template>