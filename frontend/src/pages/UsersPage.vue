<template>
  <q-page class="q-mt-lg">
    <div class="column q-gutter-lg items-center justify-start">
      <ButtonSave  v-model:save-data="rows"></ButtonSave>
      <q-table bordered title="Пользователи" :columns="columns" :rows="rows" :loading="isLoading" table-class="table"
        :table-style="{ borderCollapse: 'collapse' }" :pagination="{rowsPerPage: 10}">
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
        <template #body="props">
          <q-tr :props="props" :class="{ newRow: props.row.isNew, deleteRow: props.row.isDelete }">
            <q-td key="name" :props="props" :class="{ newRow: CheckNewRow(props.row.fields.login, props.row.isNew) }">
              <q-input v-model="props.row.fields.login.value" :rules="[nameRule]" />
            </q-td>
            <q-td key="role" :props="props" :class="{ newRow: CheckNewRow(props.row.fields.role, props.row.isNew) }">
              <q-select v-model="props.row.fields.role.value" :options="Object.values(Role)"
                :rules="[fakeRule]"></q-select>
            </q-td>
            <q-td key="actions" :props="props">
              <q-btn v-if="CheckRowChanged(props.row.fields)" flat round color="warning" icon="refresh"
                @click="revertRow(props.rowIndex)"></q-btn>
              <q-btn flat round color="negative" icon="clear" @click="removeRow(props.rowIndex)" />
            </q-td>
          </q-tr>
        </template>
      </q-table>
    </div>
  </q-page>
</template>

<script setup lang="ts">
import { useApolloClient, useMutation } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import { Notify } from 'quasar';
import ButtonSave from 'src/components/users/ButtonSave.vue';
import { Query, Role, User } from 'src/gql/graphql';
import { UserData, UserDataFields, UserField, UserDataKey } from 'src/types/types';
import { onMounted, ref } from 'vue';

const columns = [{ name: 'name', label: 'Имя', field: 'login', sortable: true },
{ name: 'role', label: 'Роль', field: 'role', sortable: true },
{ name: 'actions', label: 'Действия', field: 'actions' }]

const rows = ref<UserData[]>([])
const isLoading = ref(true)

const { client } = useApolloClient()

onMounted(() => {
  refreshData()
})

function createUser() {
  rows.value.push({
    id: '',
    fields: {
      login:
        { oldValue: '', value: '' },
      role:
        { oldValue: Role.User, value: Role.User },
    },
    isNew: true,
    isDelete: false
  })
}

async function refreshData() {
  try {
    const result = await client.query<Query>({
      query: gql`query GetUsers{
        users{
          id
          login
          role
        }
      }`,fetchPolicy: 'network-only'})

    const users = result.data.users
    if (!users) {
      return;
    }
    rows.value = users.map((item) => {
      let newItem: UserData = {
        id: item.id,
        isNew: false, isDelete: false, fields: {
          login: {
            oldValue: '',
            value: ''
          },
          role: {
            oldValue: Role.User,
            value: Role.User
          }
        }
      };
      //let key: keyof SelectedUserFields & keyof User;
      Object.keys(item).forEach((key) => {
        const fieldKey = key as keyof UserDataKey
        if (fieldKey in newItem.fields) {
          const val = item[fieldKey];
          newItem.fields[fieldKey] = {
            oldValue: val,
            value: val,
          } as UserField<typeof val>
        }
      })
      return newItem;
    });
  }
  catch (e) {
    Notify.create({
      color: 'negative',
      message: 'Ошипрка',
      timeout: 60000,
      actions: [{ icon: 'close', color: 'white'/* , handler: () => {  }  */ }]
    })
  }
  isLoading.value = false

}


function revertRow(rowIndex: number) {
  const fields = rows.value[rowIndex].fields
  const revertedObj = Object.keys(fields).map((key) => {
    const fieldKey = key as keyof UserDataKey
    return {
      [fieldKey]: {
        oldValue: fields[fieldKey].oldValue,
        value: fields[fieldKey].oldValue,
      }
    }
  });

  rows.value[rowIndex].fields = Object.assign({}, ...revertedObj)
}

function removeRow(rowIndex: number) {
  const row = rows.value[rowIndex]
  if (row.isNew) {
    rows.value.splice(rowIndex, 1)
  } else {
    rows.value[rowIndex].isDelete = !rows.value[rowIndex].isDelete
  }
}

function CheckNewRow(field: UserField<string>, isNew: boolean) {
  return field && field.value && field.value !== field.oldValue && !isNew
}

function CheckRowChanged(fields: UserDataFields) {
  return Object.keys(fields).findIndex((key) => {
    const field = fields[key as keyof UserDataFields]
    return field.oldValue != field.value
  }
  ) > -1
}


function nameRule(value: string) {
  return (value && value.length >= 3) || 'Имя должно быть не меньше трёх букв'
}


function fakeRule() {
  return true
}

function onSaved(n:User[]){
  console.log(n)
}

</script>

<style lang="scss">
$new-row-bg: #122112;
$delete-row-bg: #211212;

.buttons {
  align-self: stretch;
  margin: 10px;
  width: 100%;
}

@mixin single-cell($color-vertical, $color-horizontal) {
  border-color: $color-vertical $color-horizontal
}

@mixin row-border($color) {
  @include single-cell($color, transparent);
  @content;

  &:first-child {
    border-left: 2px solid $color;
  }

  &:last-child {
    border-right: 2px solid $color;
  }
}

// Границы ячеек и строк

tr {
  td.q-td {
    @include row-border(transparent);
    //border: 2px solid;
    border-width: 2px 0;
    border-color: transparent;
  }

  td.newRow {
    background-color: $new-row-bg;
  }

  &.newRow td {
    @include row-border($positive) {
      background-color: $new-row-bg;
    }
  }

  &.deleteRow td {
    @include row-border($negative) {
      background-color: $delete-row-bg;
    }
  }
}
</style>