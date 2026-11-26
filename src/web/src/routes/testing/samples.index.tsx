import { createFileRoute, Link } from '@tanstack/react-router'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { sampleQueries } from '@/lib/intake/queries.ts'
import { sampleMutations } from '@/lib/intake/mutations.ts'
import { Input } from '@/components/ui/input.tsx'
import { Button } from '@/components/ui/button.tsx'
import { useState } from 'react'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"

export const Route = createFileRoute('/testing/samples/')({
  component: TestingSamples,
})

function TestingSamples() {
  const { data: samples } = useQuery(sampleQueries.all())
  const queryClient = useQueryClient()
  const [serialNumber, setSerialNumber] = useState('')
  const [partNumber, setPartNumber] = useState('')
  const [search, setSearch] = useState('')
  const create = useMutation({
    ...sampleMutations.create(),
    onSuccess: () => {
      setSerialNumber(''); setPartNumber('');
      queryClient.invalidateQueries({queryKey: sampleQueries.all().queryKey})
    },
  })

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <header className="mb-8">
        <h1 className="text-4xl font-bold text-gray-900 dark:text-white">
          Samples
        </h1>
      </header>
      <form className="mb-8 flex flex-wrap gap-3" onSubmit={e => {
        e.preventDefault()
        create.mutate({serialNumber: serialNumber.trim(), partNumber: partNumber.trim()})
      }}>
        <Input aria-label="New serial number" placeholder="Serial number" value={serialNumber} onChange={e => setSerialNumber(e.target.value)} required className="max-w-xs" />
        <Input aria-label="Product part number" placeholder="Product part number" value={partNumber} onChange={e => setPartNumber(e.target.value)} required className="max-w-xs" />
        <Button type="submit" disabled={create.isPending}>Register sample</Button>
      </form>
      {create.error && <p role="alert" className="mb-4 text-red-600">{create.error.message}</p>}
      <Input aria-label="Find sample by serial number" placeholder="Find by serial number" value={search} onChange={e => setSearch(e.target.value)} className="mb-4 max-w-sm" />
      
      <div className="bg-white dark:bg-gray-800 rounded-xl shadow p-8 border border-gray-200 dark:border-gray-700">
          <Table>
              <TableHeader>
              <TableRow>
                  <TableHead>Serial Number</TableHead>
                  <TableHead>Part Number</TableHead>
                  <TableHead>Status</TableHead>
              </TableRow>
              </TableHeader>
              <TableBody>
              {samples?.filter(sample => sample.serialNumber.toLowerCase().includes(search.trim().toLowerCase())).map(sample => (
                  <TableRow key={sample.id}>
                      <TableCell>
                        <Link 
                          to="/testing/samples/$serialNumber/edit"
                          params={{ serialNumber: sample.serialNumber }}
                          className="text-primary hover:underline font-medium"
                        >
                          {sample.serialNumber}
                        </Link>
                      </TableCell>
                      <TableCell>{sample.partNumber}</TableCell>
                      <TableCell>{sample.status}</TableCell>
                  </TableRow>
              ))}
              </TableBody>
          </Table>
      </div>
    </div>
  );
}
