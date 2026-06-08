import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/testing/samples')({
  component: () => <Outlet />,
})
