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
    </div>
  </q-form>
</template>

<script setup lang="ts">
import { useQuasar } from 'quasar';
import { ref } from 'vue';

const $q = useQuasar();

const isPwd = ref(true);
const name = ref(null);
const password = ref(null);

const onSubmit = () => {
    $q.notify({
      color: 'green-4',
      textColor: 'white',
      icon: 'cloud_done',
      message: 'Ок',
    });
};

const onReset = () => {
  name.value = null;
  password.value = null;
};
</script>
