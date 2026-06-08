import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { sampleQueries } from '../../lib/intake/queries'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

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
          <Table>
              <TableHeader>
              <TableRow>
                  <TableHead>Part Number</TableHead>
                  <TableHead>Serial Number</TableHead>
                  <TableHead>Status</TableHead>
              </TableRow>
              </TableHeader>
              <TableBody>
              {samples?.map(sample => (
                  <TableRow key={sample.id}>
                      <TableCell>{sample.partNumber}</TableCell>
                      <TableCell>{sample.serialNumber}</TableCell>
                      <TableCell>{sample.status}</TableCell>
                  </TableRow>
              ))}
              </TableBody>
          </Table>
      </div>
    </div>
  );
}
