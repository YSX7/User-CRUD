export interface User {
  id: string
  login: string
}

export interface AuthPayload {
  token: string
  user: User
}