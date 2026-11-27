import { apiUrl } from '../api.ts'
import type { User } from './types.ts'

export const adminUserKeys = {
  all: ['admin', 'users'] as const,
}

async function getUsers(): Promise<User[]> {
  const response = await fetch(`${apiUrl}/admin/users`, { credentials: 'include' })
  if (!response.ok) {
    throw new Error('Failed to fetch users')
  }
  return response.json()
}

export const adminUserQueries = {
  all: () => ({
    queryKey: adminUserKeys.all,
    queryFn: getUsers,
  }),
}
