import { createFileRoute } from '@tanstack/react-router'
import {useQuery} from "@tanstack/react-query";
import {sampleQueries} from "@/lib/intake/queries.ts";
import {Field, FieldGroup, FieldLabel} from "@/components/ui/field.tsx";
import {Input} from "@/components/ui/input.tsx";
import {Select, SelectTrigger, SelectValue, SelectItem, SelectContent} from "@/components/ui/select.tsx";
import { useState, useMemo} from "react";
import {productQueries} from "@/lib/configuration/queries.ts";
import {Test} from "@/lib/configuration/types.ts";

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

    type mergedTest = {
        partNumber: string,
        productTestSequence: number,
        test: Test,
        specificModifiers: string[],
        testResult: number | null,
        unit: string,
        minValue: number | null,
        maxValue: number | null
    }
    const mergedTests = useMemo<mergedTest[]>(() => {
        if (!product || !product.productTestConfigurations || !sample || !sample.testResults) {
            return [];
        }
        return product.productTestConfigurations.map((testConfiguration) => {
            const testResult = sample.testResults!.find(tr =>
                tr.partNumber === sample.partNumber &&
                tr.productTestSequence === testConfiguration.productTestSequence
            );
            return {
                partNumber: product.partNumber,
                productTestSequence: testConfiguration.productTestSequence,
                test: testConfiguration.test,
                specificModifiers: testConfiguration.specificModifiers,
                testResult: testResult?.testResult ?? null,
                unit: testResult?.unit ?? testConfiguration.unit,
                minValue: testResult?.minValue ?? testConfiguration.minValue,
                maxValue: testResult?.maxValue ?? testConfiguration.maxValue,
            };
        });
    }, [product?.productTestConfigurations, sample?.testResults, sample?.partNumber, product?.partNumber]);

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
          {mergedTests.map((test) => (
              <div key={test.productTestSequence} className="mb-4">
                  <h3>
                      <span>{test.test.testName}</span>
                      {test.specificModifiers.length > 0 ? <span className="text-sm text-gray-500 ml-2">({test.specificModifiers.join(', ')})</span> : ''}
                  </h3>
                  {test.testResult && (
                      <div className="mt-2 text-sm">
                          Result: {test.testResult} {test.unit}
                      </div>
                  )}
              </div>
          ))}

      </div>
    </div>
  )
}
