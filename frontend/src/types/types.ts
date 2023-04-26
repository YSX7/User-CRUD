import { User } from 'src/gql/graphql'

export type UserField<T> = {
    oldValue: T
    value: T
}

//export const UserDataKeys: Array<keyof User> = ['login', 'role']

export type UserDataKey = Omit<User,'__typename'>

export type UserDataFields = {[K in keyof UserDataKey]?: UserField<User[K]>};

export type UserData = {
    fields: UserDataFields
    isNew: boolean
    isDelete: boolean 
}