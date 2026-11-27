import { apiUrl } from '../api.ts'
import type { CreateUserInput, SetUserDisabledInput, SetUserRolesInput } from './types.ts'

async function adminRequest(path: string, method: string, body: object): Promise<void> {
  const response = await fetch(`${apiUrl}/admin${path}`, {
    method,
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })
  if (!response.ok) {
    const result = await response.json().catch(() => ({}))
    throw new Error(result.error || 'User request failed')
  }
}

async function createUser(input: CreateUserInput): Promise<void> {
  await adminRequest('/users', 'POST', input)
}

async function setUserRoles(input: SetUserRolesInput): Promise<void> {
  await adminRequest(`/users/${encodeURIComponent(input.userId)}/roles`, 'PUT', { roles: input.roles })
}

async function setUserDisabled(input: SetUserDisabledInput): Promise<void> {
  await adminRequest(`/users/${encodeURIComponent(input.userId)}/disabled`, 'PUT', { disabled: input.disabled })
}

export const adminUserMutations = {
  create: () => ({ mutationFn: createUser }),
  setRoles: () => ({ mutationFn: setUserRoles }),
  setDisabled: () => ({ mutationFn: setUserDisabled }),
}
