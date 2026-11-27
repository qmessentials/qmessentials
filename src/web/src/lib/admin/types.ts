export const roles = ['Administrator', 'Quality Manager', 'Technician'] as const

export type User = {
  userId: string
  disabled: boolean
  roles: string[]
}

export type CreateUserInput = {
  userId: string
  password: string
  roles: string[]
}

export type SetUserRolesInput = {
  userId: string
  roles: string[]
}

export type SetUserDisabledInput = {
  userId: string
  disabled: boolean
}
