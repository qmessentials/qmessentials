import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/testing/samples/$serialNumber/edit')({
  component: EditSample,
})

function EditSample() {
  const { serialNumber } = Route.useParams()
  
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-8">
        Edit Sample: {serialNumber}
      </h1>
      <div className="bg-white dark:bg-gray-800 rounded-xl shadow p-8 border border-gray-200 dark:border-gray-700">
        <p className="text-gray-600 dark:text-gray-400">
          This is the edit page for sample with serial number: <strong>{serialNumber}</strong>
        </p>
      </div>
    </div>
  )
}
