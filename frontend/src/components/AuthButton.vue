<template>
  <q-btn v-if="!userStore.IsLogged" color="white" text-color="primary" label="Login" @click="Login"></q-btn>
  <q-btn v-else color="white" text-color="primary" label="Logout" @click="Logout"></q-btn>
</template>

<script setup lang="ts">
import { Cookies } from 'quasar';
import { useUserStore } from 'src/stores/store-user';
import { useRouter } from 'vue-router';

  const userStore = useUserStore()
  const router = useRouter()

  function Login(){
    router.push({name: 'login'})
  }

  function Logout(){
    userStore.Logout();
    Cookies.remove(process.env.AUTH_COOKIE)
    router.push({name: 'login'})
  }
</script>