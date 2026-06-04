import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/testing/samples')({
  component: TestingSamples,
})

function TestingSamples() {
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <header className="mb-8">
        <h1 className="text-4xl font-bold text-gray-900 dark:text-white">
          Samples
        </h1>
      </header>
      
      <div className="bg-white dark:bg-gray-800 rounded-xl shadow p-8 border border-gray-200 dark:border-gray-700">
        <p className="text-gray-500 dark:text-gray-400 italic">
          No samples to display yet.
        </p>
      </div>
    </div>
  );
}
