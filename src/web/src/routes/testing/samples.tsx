import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { sampleQueries } from '../../lib/intake/queries'

export const Route = createFileRoute('/testing/samples')({
  component: TestingSamples,
})

function TestingSamples() {
  const { data: samples } = useQuery(sampleQueries.all())

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <header className="mb-8">
        <h1 className="text-4xl font-bold text-gray-900 dark:text-white">
          Samples
        </h1>
      </header>
      
      <div className="bg-white dark:bg-gray-800 rounded-xl shadow p-8 border border-gray-200 dark:border-gray-700">
          <table>
              <thead>
              <tr>
                  <th>Part Number</th>
                  <th>Serial Number</th>
                  <th>Status</th>
              </tr>
              </thead>
              <tbody>
              {samples?.map(sample => (
                  <tr key={sample.id}>
                      <td>{sample.partNumber}</td>
                      <td>{sample.serialNumber}</td>
                      <td>{sample.status}</td>
                  </tr>
              ))}
              </tbody>
          </table>
      </div>
    </div>
  );
}
