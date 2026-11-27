import { createFileRoute, Link } from '@tanstack/react-router'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useState, type FormEvent } from 'react'
import { getSession } from '@/lib/api'
import { adminUserKeys, adminUserQueries } from '@/lib/admin/queries.ts'
import { adminUserMutations } from '@/lib/admin/mutations.ts'
import { roles, type User } from '@/lib/admin/types.ts'

export const Route = createFileRoute('/users')({ component: UsersPage })

function UsersPage() {
  const client = useQueryClient()
  const session = useQuery({ queryKey: ['session'], queryFn: getSession })
  const users = useQuery({ ...adminUserQueries.all(), enabled: session.data?.roles?.includes('Administrator') })
  const [userId, setUserId] = useState('')
  const [password, setPassword] = useState('')
  const [selectedRoles, setSelectedRoles] = useState<string[]>([])
  const create = useMutation({
    ...adminUserMutations.create(),
    onSuccess: () => { setUserId(''); setPassword(''); setSelectedRoles([]); client.invalidateQueries({queryKey: adminUserKeys.all}) },
  })
  const updateRoles = useMutation({
    ...adminUserMutations.setRoles(),
    onSuccess: () => { client.invalidateQueries({queryKey: adminUserKeys.all}); client.invalidateQueries({queryKey: ['session']}) },
  })
  const updateDisabled = useMutation({
    ...adminUserMutations.setDisabled(),
    onSuccess: () => { client.invalidateQueries({queryKey: adminUserKeys.all}); client.invalidateQueries({queryKey: ['session']}) },
  })
  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    create.mutate({userId, password, roles: selectedRoles})
  }
  if (session.isLoading) return <div className="max-w-5xl mx-auto p-8">Loading…</div>
  if (!session.data) return <div className="max-w-5xl mx-auto p-8"><Link to="/login">Sign in</Link> to manage users.</div>
  if (!session.data.roles?.includes('Administrator')) return <div className="max-w-5xl mx-auto p-8">Administrator role required.</div>
  return <div className="max-w-5xl mx-auto px-4 py-12 space-y-8 text-gray-900 dark:text-white">
    <h1 className="text-4xl font-bold">Users</h1>
    <form onSubmit={submit} className="bg-white dark:bg-gray-800 rounded-xl border p-6 space-y-4">
      <h2 className="text-2xl font-semibold">Create user</h2>
      <div className="flex flex-wrap gap-4">
        <label>User ID <input className="block border rounded p-2 text-gray-900" required minLength={3} maxLength={64} value={userId} onChange={event => setUserId(event.target.value)} /></label>
        <label>Initial password <input className="block border rounded p-2 text-gray-900" type="text" required minLength={12} maxLength={72} value={password} onChange={event => setPassword(event.target.value)} /></label>
      </div>
      <fieldset><legend>Roles</legend><div className="flex flex-wrap gap-4">{roles.map(role => <label key={role}><input type="checkbox" checked={selectedRoles.includes(role)} onChange={event => setSelectedRoles(event.target.checked ? [...selectedRoles, role] : selectedRoles.filter(value => value !== role))} /> {role}</label>)}</div></fieldset>
      <button className="bg-blue-600 text-white px-4 py-2 rounded" disabled={create.isPending}>Create user</button>
      {create.isError && <p role="alert">{create.error.message}</p>}
    </form>
    <section className="bg-white dark:bg-gray-800 rounded-xl border p-6 space-y-4">
      <h2 className="text-2xl font-semibold">Accounts</h2>
      {users.isLoading && <p>Loading users…</p>}
      {users.isError && <p role="alert">{users.error.message}</p>}
      {updateRoles.isError && <p role="alert">{updateRoles.error.message}</p>}
      {updateDisabled.isError && <p role="alert">{updateDisabled.error.message}</p>}
      {users.data?.map((account: User) => <div key={account.userId} className="border-t py-4 flex flex-wrap gap-4 items-center justify-between">
        <div><strong>{account.userId}</strong> <span>{account.disabled ? '(disabled)' : '(active)'}</span></div>
        <div className="flex flex-wrap gap-4">{roles.map(role => <label key={role}><input type="checkbox" checked={account.roles.includes(role)} disabled={updateRoles.isPending || updateDisabled.isPending} onChange={event => updateRoles.mutate({userId: account.userId, roles: event.target.checked ? [...account.roles, role] : account.roles.filter(value => value !== role)})} /> {role}</label>)}</div>
        <button type="button" className="border rounded px-3 py-1" disabled={updateRoles.isPending || updateDisabled.isPending} onClick={() => updateDisabled.mutate({userId: account.userId, disabled: !account.disabled})}>{account.disabled ? 'Enable' : 'Disable'}</button>
      </div>)}
    </section>
  </div>
}
