import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/monitoring/subscriptions')({
  component: () => <Outlet />,
})
