<template>
  <q-form @submit="onSubmit" @reset="onReset" class="q-gutter-md">
    <q-input
      filled
      v-model="name"
      label="Логин *"
      lazy-rules
      no-error-icon
      :rules="[(val) => (val && val.length > 0) || 'Поле не может быть пустым']"
    >
</q-input>

    <q-input
      filled
      :type="isPwd ? 'password' : 'text'"
      v-model="password"
      label="Пароль *"
      lazy-rules
      no-error-icon
      :rules="[(val) => (val && val.length > 0) || 'Поле не может быть пустым']"
    >    <template v-slot:append>
          <q-icon
            :name="isPwd ? 'visibility_off' : 'visibility'"
            class="cursor-pointer"
            @mousedown="isPwd = false"
            @mouseup="isPwd = true"
          />
        </template>
</q-input>

    <div>
      <q-btn label="Войти" type="submit" color="primary"></q-btn>
      <q-btn
        label="Сброс"
        type="reset"
        color="primary"
        flat
        class="q-ml-sm"
      ></q-btn>
      <q-btn label="Кукисы" color="primary" flat @click="cookiesCheck"></q-btn>
    </div>
  </q-form>
</template>

<script setup lang="ts">

import { Cookies, Notify } from 'quasar';
import { useApolloClient } from '@vue/apollo-composable'
import { ref } from 'vue';
import gql from 'graphql-tag';
import { AuthPayload } from 'src/apollo/models';
import { ApolloError } from '@apollo/client/errors';
import { useRouter } from 'vue-router';

const router = useRouter();

const isPwd = ref(true);
const name = ref('');
const password = ref('');

const {resolveClient} = useApolloClient<AuthPayload>()


const cookiesCheck =  () => {
  Notify.create({color: 'neutral', textColor: 'white', message: JSON.stringify(Cookies.getAll()), multiLine: true})
  router.push({name:'index'})
}

const onSubmit = () => {
  const client = resolveClient()
  client.mutate({mutation: gql`mutation Auth ($login: String!, $password: String!) {
    auth {
        login(login: $login, password: $password) {
            user {
                id
              	login
            }
        }
    }
}`, variables: {login: name.value, password: password.value}}).then((result)=> {
  Notify.create({color:'positive', textColor: 'white', icon: 'cloud_done',message:'Ok'})
  router.push({name: 'index'})
},(reason: ApolloError) => {
  Notify.create({color:'negative', textColor: 'white', icon: 'cloud_done',message: reason.message})
})
};

const onReset = () => {
  name.value = '';
  password.value = '';
};
</script>
