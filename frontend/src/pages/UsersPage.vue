<template>
  <q-page class="column items-center justify-evenly">
    <q-table title="Пользователи" :columns="columns" :rows="rows" :loading="isLoading" >
      <template v-slot:top>
        <h5 style="margin: unset">Пользователи</h5>
        <q-space />
        <q-btn square color="primary" icon="refresh" size="sm" @click="refreshData">
          <q-tooltip :delay=400>Обновить список пользователей</q-tooltip>
        </q-btn>
        <q-btn class="q-ml-sm" square color="green" icon="person_add" size="sm" @click="createUser">
          <q-tooltip :delay=400>Создать нового пользователя</q-tooltip>
        </q-btn>
      </template>
      <template #body-cell-name="props">
        <q-td :props="props">
          <q-input v-model="props.row[props.col.field].value" />
        </q-td>
      </template>
      <template #body-cell-role="props">
        <q-td :props="props">
          <q-select v-model="props.row[props.col.field].value" :options="Object.values(Role)"></q-select>
        </q-td>
      </template>
      <template #body-cell-actions="props">
        <q-td :props="props"><q-btn flat round color="negative" icon="clear" /></q-td>
      </template>
    </q-table>
  </q-page>
</template>

<script setup lang="ts">
import { useApolloClient } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import { Notify } from 'quasar';
import { Query, Role } from 'src/gql/graphql';
import { UserData, UserDataKeys } from 'src/types/types';
import { onMounted, ref } from 'vue';

const columns = [{ name: 'name', label: 'Имя', field: 'login', sortable: true },
{ name: 'role', label: 'Роль', field: 'role', sortable: true },
{ name: 'actions', label: 'Действия', field: 'actions' }]

const rows = ref<UserData[]>([])
const isLoading = ref(true)

const { client } = useApolloClient()

function createUser() {
  rows.value.push({ login: { oldValue: 'biba', value: '' }, role: { oldValue: Role.User, value: Role.User }, isNew: true })
}

function refreshData(){
  client.query<Query>({
    query: gql`query GetUsers{
        users{
          id
          role
          login
        }
      }`}).then((result) => {
      rows.value = result.data.users!.map((item) => {
        let newItem: UserData = {isNew: false};
        // let key: keyof User
        for (const key of UserDataKeys) {
          Object.assign(newItem, {
            [key]: {
              oldValue: item[key],
              value: item[key],
            }
          })
        }
        return newItem;
      });
    },
      (reason) => {
        Notify.create({
          color: 'negative',
          message: reason,
          timeout: 60000,
          actions: [{ icon: 'close', color: 'white', handler: () => { /* ... */ } }]
        })
      }
    )
    .finally(() => {
      isLoading.value = false
    })
}

onMounted(() => {
  refreshData()
})



function saveCell(row, field) {
  const rowIndex = rows.value.indexOf(row);

  const newRow = { ...row, [field]: row[field] };

  rows.value[field] = newRow;
  //rows.value.toSpliced()
}

</script>

<style>
.buttons {
  align-self: stretch;
  margin: 10px;
  width: 100%;
}
</style>