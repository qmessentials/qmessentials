import { createFileRoute } from '@tanstack/react-router'
import {useMutation} from "@tanstack/react-query";
import {logIn} from "../lib/api.ts";

export const Route = createFileRoute('/login')({
  component: Login,
})

function Login() {
    const loginMutation = useMutation({
        mutationFn: logIn
    })
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-6">
        Login
      </h1>
        <hr className="my-3"/>
        <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded" onClick={() => loginMutation.mutate()}>Log In</button>
    </div>
  )
}
