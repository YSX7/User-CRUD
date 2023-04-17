import { useApolloClient } from '@vue/apollo-composable';
import { defineStore } from 'pinia';
import { Role, User } from 'src/gql/graphql';
import gql from 'graphql-tag';

interface State {
  IsLogged: boolean
  Role?: Role
}

export const useUserStore = defineStore('user',
 {
  state: () : State => {
    const { resolveClient } = useApolloClient<User>()
    const client = resolveClient()
    const response = await client.mutate<User>({mutation: gql`mutation Mutation{
      validate{
        id
        login
        role
      }
    }`})
    return {
      IsLogged: true,
      Role: response.data?.role
    }
  },
  actions: {
    Login(role: Role) {
      this.IsLogged = true;
      this.Role = role;
    },
    Logout(){
      this.IsLogged = false;
      this.Role = undefined;
    },
  },
});
