import { createFileRoute, Link } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { subscriptionQueries } from '@/lib/subscription/queries.ts'
import { buttonVariants } from '@/components/ui/button.tsx'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

export const Route = createFileRoute('/monitoring/subscriptions/')({
  component: MonitoringSubscriptions,
})

function MonitoringSubscriptions() {
  const { data: subscriptions } = useQuery(subscriptionQueries.all())

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <header className="mb-8 flex items-center justify-between gap-4">
        <h1 className="text-4xl font-bold text-gray-900 dark:text-white">
          Subscriptions
        </h1>
        <Link
          to="/monitoring/subscriptions/new"
          className={buttonVariants()}
        >
          Create Subscription
        </Link>
      </header>

      <div className="bg-white dark:bg-gray-800 rounded-xl shadow p-8 border border-gray-200 dark:border-gray-700">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Subscription ID</TableHead>
              <TableHead>Rule</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {subscriptions?.map(subscription => (
              <TableRow
                key={subscription.id}
                className={subscription.isActive ? undefined : 'text-muted-foreground opacity-60'}
              >
                <TableCell>
                  <Link
                    to="/monitoring/subscriptions/$subscriptionId/edit"
                    params={{ subscriptionId: subscription.id.toString() }}
                    className="text-primary hover:underline font-medium"
                  >
                    {subscription.id}
                  </Link>
                </TableCell>
                <TableCell>{subscription.ruleText}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </div>
  )
}
