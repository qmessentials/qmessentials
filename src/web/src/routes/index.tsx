import { createFileRoute, Link } from '@tanstack/react-router'
import { buttonVariants } from '@/components/ui/button'

export const Route = createFileRoute('/')({
  component: Home,
})

function Home() {
  const cards = [
    {
      title: 'Testing',
      actions: [
        { label: 'Samples', path: '/testing/samples' },
      ],
    },
  ];

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        {cards.map((card) => (
          <div
            key={card.title}
            className="bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-200 dark:border-gray-700 overflow-hidden hover:shadow-2xl transition-shadow"
          >
            <div className="p-6">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">
                {card.title}
              </h2>
              <div className="flex flex-wrap gap-4">
                {card.actions.map((action) => (
                  <Link
                    key={action.path}
                    to={action.path}
                    className={buttonVariants({ variant: "outline" })}
                  >
                    {action.label}
                  </Link>
                ))}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
