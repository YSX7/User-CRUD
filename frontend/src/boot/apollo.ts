import { ApolloClient /*, createHttpLink */ } from '@apollo/client/core'
import { ApolloClients } from '@vue/apollo-composable'
import { boot } from 'quasar/wrappers'
import { getClientOptions } from 'src/apollo'

export function getDefaultApolloClient(){
  const options = /* await */ getClientOptions(/* {app, router ...} */)
  return new ApolloClient(options)
}

export default boot(
  ({ app }) => {
    // Default client.

    const apolloClient = getDefaultApolloClient()


    const apolloClients: Record<string, ApolloClient<unknown>> = {
      default: apolloClient,
      // clientA,
      // clientB,
    }

    app.provide(ApolloClients, apolloClients)
  }
)
