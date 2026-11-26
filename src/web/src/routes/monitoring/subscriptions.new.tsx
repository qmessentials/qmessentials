import { FormEvent, useState } from 'react'
import { createFileRoute, Link, useNavigate } from '@tanstack/react-router'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { subscriptionMutations } from '@/lib/subscription/mutations.ts'
import { subscriptionKeys } from '@/lib/subscription/queries.ts'
import { Button, buttonVariants } from '@/components/ui/button.tsx'
import { Field, FieldLabel } from '@/components/ui/field.tsx'

export const Route = createFileRoute('/monitoring/subscriptions/new')({
  component: NewSubscription,
})

function NewSubscription() {
  const [ruleText, setRuleText] = useState('')
  const queryClient = useQueryClient()
  const navigate = useNavigate()
  const createSubscription = useMutation({
    ...subscriptionMutations.create(),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: subscriptionKeys.root })
      await navigate({ to: '/monitoring/subscriptions' })
    },
  })

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!ruleText.trim()) return
    createSubscription.mutate({ ruleText })
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <header className="mb-8">
        <h1 className="text-4xl font-bold text-gray-900 dark:text-white">Create Subscription</h1>
      </header>

      <div className="bg-white dark:bg-gray-800 rounded-xl shadow p-8 border border-gray-200 dark:border-gray-700">
        <form className="space-y-6" onSubmit={handleSubmit}>
          <Field>
            <FieldLabel htmlFor="ruleText">Query</FieldLabel>
            <textarea
              id="ruleText"
              name="ruleText"
              value={ruleText}
              onChange={(event) => setRuleText(event.target.value)}
              rows={12}
              className="border-input bg-transparent dark:bg-input/30 focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive min-h-40 w-full rounded-md border px-3 py-2 text-sm shadow-xs outline-none transition-[color,box-shadow] focus-visible:ring-3 disabled:cursor-not-allowed disabled:opacity-50"
              required
            />
          </Field>

          {createSubscription.isError && (
            <p role="alert" className="text-sm text-destructive">
              Failed to create subscription. Please try again.
            </p>
          )}

          <div className="flex gap-3">
            <Button type="submit" disabled={!ruleText.trim() || createSubscription.isPending}>
              {createSubscription.isPending ? 'Saving…' : 'Save'}
            </Button>
            <Link to="/monitoring/subscriptions" className={buttonVariants({ variant: 'outline' })}>
              Cancel
            </Link>
          </div>
        </form>
      </div>
    </div>
  )
}
