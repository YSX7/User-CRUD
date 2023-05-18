<template>
  <router-view />
</template>


<script lang="ts">
import { useUserStore } from './stores/store-user';

export default {

  async preFetch({redirect}) { 
    const userStore = useUserStore()
    const initResult = await userStore.Init()
    if (!initResult){
      redirect({name: 'login'})
    } else{
      redirect({name: 'index'})
    }
   }
}

</script>
