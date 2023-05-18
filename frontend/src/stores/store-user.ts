import { provideApolloClient, useApolloClient } from '@vue/apollo-composable';
import { defineStore } from 'pinia';
import { Role, User } from 'src/gql/graphql';
import gql from 'graphql-tag';
import { ref } from 'vue';
import { getDefaultApolloClient } from 'src/boot/apollo';
import { Mutation } from 'src/gql/graphql';

export const useUserStore = defineStore('user', () => {
  const IsLogged = ref<boolean>();
  const Role = ref<Role>();

  async function Init(): Promise<boolean> {
    let response;
    provideApolloClient(getDefaultApolloClient())

    const { resolveClient } = useApolloClient<Mutation>()
    const client = resolveClient()
    try {
      response = await client.mutate<Mutation>({
        mutation: gql`mutation Mutation{
      validate{
        id
        login
        role
      }
    }`})
    } catch (e) {
      console.log('not logged')
      IsLogged.value = false;
      return false
    }
    IsLogged.value = true;
    Role.value = response.data?.validate?.role;
    return true
  }

  function Login(role: Role) {
    IsLogged.value = true;
    Role.value = role
  }
  function Logout() {
    IsLogged.value = false;
    Role.value = undefined;
  }

  return { IsLogged, Role, Login, Logout, Init }
});
