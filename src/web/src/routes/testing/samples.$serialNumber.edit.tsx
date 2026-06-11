import { createFileRoute } from '@tanstack/react-router'
import {useQuery} from "@tanstack/react-query";
import {sampleQueries} from "@/lib/intake/queries.ts";
import {Field, FieldGroup, FieldLabel} from "@/components/ui/field.tsx";
import {Input} from "@/components/ui/input.tsx";
import {Select, SelectTrigger, SelectValue, SelectItem, SelectContent} from "@/components/ui/select.tsx";
import { useState} from "react";
import {productQueries} from "@/lib/configuration/queries.ts";

export const Route = createFileRoute('/testing/samples/$serialNumber/edit')({
  component: EditSample,
})

function EditSample() {
  const { serialNumber } = Route.useParams()

    const {data: sample} = useQuery(sampleQueries.sample(serialNumber))
    const {data: product} = useQuery({
        ...productQueries.product(sample?.partNumber ?? ''),
        enabled: !!sample?.partNumber
    })
    const [status, setStatus] = useState<string | null>(sample?.status ?? null);

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <h1 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">
        Edit Sample: {serialNumber}
      </h1>
      <div className="bg-white dark:bg-gray-800 rounded-md p-8 border border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-400">
          <form className="space-y-6">
              <FieldGroup className="grid  grid-cols-2">
                  <Field>
                    <FieldLabel htmlFor="partNumber">Part #</FieldLabel>
                      <Input id="partNumber" type="text" value={sample?.partNumber ?? ''} readOnly/>
                  </Field>
                  <Field>
                      <FieldLabel htmlFor="status">Status</FieldLabel>
                      <Select id="status" value={status ?? sample?.status ?? ''} onValueChange={setStatus}>
                          <SelectTrigger>
                              <SelectValue placeholder="Select a status" />
                          </SelectTrigger>
                          <SelectContent>
                              <SelectItem value="TESTING">TESTING</SelectItem>
                              <SelectItem value="COMPLETED">COMPLETED</SelectItem>
                              <SelectItem value="INVALIDATED">INVALIDATED</SelectItem>
                              <SelectItem value="ARCHIVED">ARCHIVED</SelectItem>
                          </SelectContent>
                      </Select>
                  </Field>
              </FieldGroup>
          </form>
          <hr className="my-4"/>
          {product?.productTestConfigurations?.map((testConfiguration) => (
              <div key={testConfiguration.test.testName} className="mb-4">
                  <h3>
                      <span>{testConfiguration.test.testName}</span>
                      {testConfiguration?.specificModifiers.length > 0 ? <span className="text-sm text-gray-500 ml-2">({testConfiguration.specificModifiers.join(', ')})</span> : ''}
                  </h3>
              </div>
          ))}

      </div>
    </div>
  )
}
