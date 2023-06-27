<template>
  <q-btn color="positive" icon="save" label="Сохранить изменения" @click="save"></q-btn>
</template>

<script setup lang="ts">
import { useMutation } from '@vue/apollo-composable'
import { UserData } from 'src/types/types'
import gql from 'graphql-tag';
import { Notify } from 'quasar';
import { Mutation } from 'src/gql/graphql';

const emit = defineEmits(['update:saveData'])
const props = defineProps<{ saveData: UserData[] }>()

async function save() {
  const result = await addUsers();
  await updateUsers();
 await deleteUsers();

  emit('update:saveData', result)
}

async function addUsers() {
  const newUsers = props.saveData.filter((value) => value.isNew).map((item) => ({
    login: item.fields.login.value,
    role: item.fields.role.value
  }))
  if (!newUsers.length) return props.saveData
  const { mutate: sendMessage,  } = useMutation<Mutation>(gql`
  mutation NewUser($users: [NewUserInput!]!){
      userNew(users: $users){
        id
        login
      }
    }
  `)
  try {
    const result = await sendMessage({
      users: newUsers
    })

    if (result?.errors) {
      result.errors.forEach(
        (currentValue) => {
          Notify.create({
            color: 'negative', message: currentValue.message
          })
        }
      )
      return props.saveData
    }
    return props.saveData.map(elem => {
      const user = result?.data?.userNew.find(val => elem.fields.login.value == val.login)
      if (user)
        return {
          ...elem, fields: {
            ...elem.fields, login: {
              oldValue: elem.fields.login.value,
              value: elem.fields.login.value
            }
          }, id: user.id, isNew: false
        }
      return elem;
    })
  }
  catch (e: any) {
    Notify.create({
      color: 'negative', message: e.message
    })
  }
  return
}

async function updateUsers() {
  const updateUsers = props.saveData.filter((value) => {
    for (const v of Object.values(value.fields))
      if (v.oldValue != v.value)
        return true;
    return false;
  })
    .map((item) => ({
      id: item.id,
      login: item.fields.login.value,
      role: item.fields.role.value
    }))

  if (updateUsers.length == 0) return props.saveData

  const { mutate: sendMessage } = useMutation<Mutation>(gql`
    mutation UpdateUser($users: [UpdateUserInput!]!){
        userUpdate(users: $users){
          id
          login
          role
        }
      }
  `)

  try {
    const result = await sendMessage({
      users: updateUsers
    })

    if (result?.errors) {
      result.errors.forEach(
        (currentValue) => {
          Notify.create({
            color: 'negative', message: currentValue.message
          })
        }
      )
      return props.saveData
    }
    return props.saveData.map(elem => {
      const user = result?.data?.userUpdate.find(val => elem.id == val.id)
      if (user)
        return {
          ...elem, fields: {
            ...elem.fields, login: {
              oldValue: elem.fields.login.value,
              value: elem.fields.login.value
            },
            role: {
              oldValue: elem.fields.role.value,
              value: elem.fields.role.value
            }
          },
        }
      return elem;
    })
  }
  catch (e: any) {
    Notify.create({
      color: 'negative', message: e.message
    })
  }
  return;
}

async function deleteUsers() {
  const deleteUsers = props.saveData.filter((value) => value.isDelete).map((value) => value.id)

  if (deleteUsers.length == 0) return

  const { mutate: sendMessage } = useMutation<Mutation>(gql`
          mutation userDelete($users: [ID!]!){
        userDelete(users: $users)
      }
  `)

  try {
    const result = await sendMessage({
      users: deleteUsers
    })

    if (result?.errors) {
      result.errors.forEach(
        (currentValue) => {
          Notify.create({
            color: 'negative', message: currentValue.message
          })
        }
      )
      return
    }
   return result?.data?.userDelete != props.saveData.length && 'Не все пользователи удалены'; 
  }
  catch (e: any) {
    Notify.create({
      color: 'negative', message: e.message
    })
  }
  return;
}
</script>
