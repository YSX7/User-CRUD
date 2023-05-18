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
  let changedUsers = await addUsers()
  emit('update:saveData', changedUsers)
}

async function addUsers() {
  const newUsers = props.saveData.filter((value) => value.isNew).map((item) => ({
    login: item.fields.login.value,
    role: item.fields.role.value
  }))
  if (!newUsers.length) return props.saveData
  const { mutate: sendMessage } = useMutation<Mutation>(gql`
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
}

</script>
