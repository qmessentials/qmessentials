import { createRootRoute, Outlet } from '@tanstack/react-router'
import Navbar from '../components/Navbar'

export const Route = createRootRoute({
  component: () => (
    <div className="min-h-screen bg-gray-100 dark:bg-gray-900">
      <Navbar />
      <main className="pt-20">
        <Outlet />
      </main>
    </div>
  ),
})
