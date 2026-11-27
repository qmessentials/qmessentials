import { createFileRoute } from '@tanstack/react-router'
import {useMutation, useQueryClient} from "@tanstack/react-query";
import {useState, type FormEvent} from "react";
import {logIn} from "../lib/api.ts";

export const Route = createFileRoute('/login')({
  component: Login,
})

function Login() {
    const queryClient = useQueryClient()
    const [userId, setUserId] = useState('')
    const [password, setPassword] = useState('')
    const loginMutation = useMutation({
        mutationFn: logIn,
        onSuccess: () => queryClient.invalidateQueries({queryKey: ['session']})
    })
    function submit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        loginMutation.mutate({userId, password})
    }
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-6">
        Login
      </h1>
        <form onSubmit={submit} className="max-w-sm space-y-4">
          <label className="block">
            <span className="block mb-1">User ID</span>
            <input className="w-full rounded border p-2" autoComplete="username" required
              value={userId} onChange={event => setUserId(event.target.value)} />
          </label>
          <label className="block">
            <span className="block mb-1">Password</span>
            <input className="w-full rounded border p-2" type="password" autoComplete="current-password" required
              value={password} onChange={event => setPassword(event.target.value)} />
          </label>
          <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
            type="submit" disabled={loginMutation.isPending}>Log In</button>
          {loginMutation.isError && <p role="alert">Invalid user ID or password.</p>}
          {loginMutation.isSuccess && <p role="status">Signed in.</p>}
        </form>
    </div>
  )
}
