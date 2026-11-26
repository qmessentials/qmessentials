import { createFileRoute } from "@tanstack/react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { testResultMutations } from "@/lib/intake/mutations.ts";
import { sampleQueries } from "@/lib/intake/queries.ts";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field.tsx";
import { Input } from "@/components/ui/input.tsx";
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectItem,
  SelectContent,
} from "@/components/ui/select.tsx";
import { useState, useMemo } from "react";
import { productQueries } from "@/lib/configuration/queries.ts";
import { Test } from "@/lib/configuration/types.ts";
import { Sample } from "@/lib/intake/types.ts";
import { Button } from "@/components/ui/button.tsx";

export const Route = createFileRoute("/testing/samples/$serialNumber/edit")({
  component: EditSample,
});

function EditSample() {
  const { serialNumber } = Route.useParams();

  const { data: sample } = useQuery(sampleQueries.sample(serialNumber));
  const { data: product } = useQuery({
    ...productQueries.product(sample?.partNumber ?? ""),
    enabled: !!sample?.partNumber,
  });
  const [status, setStatus] = useState<string | null>(sample?.status ?? null);
  const queryClient = useQueryClient();
  const { mutate: submitTestResult } = useMutation({
    ...testResultMutations.submit(),
    onSuccess: (_, variables) => {
      queryClient.setQueryData<Sample>(
        sampleQueries.sample(serialNumber).queryKey,
        (old) =>
          old && {
            ...old,
            testResults: [
              ...(old.testResults ?? []).filter(
                (tr) =>
                  tr.canonicalTestName !== variables.canonicalTestName,
              ),
              {
                id: "",
                serialNumber: variables.serialNumber,
                partNumber: variables.partNumber,
                canonicalTestName: variables.canonicalTestName,
                modifiers: variables.modifiers,
                testResult: variables.testResult,
                unit: variables.unit,
                decimalPlaces: variables.decimalPlaces,
                minValue: variables.minValue,
                maxValue: variables.maxValue,
                hashValue: "",
                voidedAt: null,
                voidedBy: null,
                voidedReason: null,
                voidComment: null,
                createdAt: new Date(),
                updatedAt: new Date(),
              },
            ],
          },
      );
    },
  });

  type mergedTest = {
    partNumber: string;
    productTestSequence: number;
    test: Test;
    specificModifiers: string[];
    testResult: number | null;
    unit: string;
    minValue: number | null;
    maxValue: number | null;
    decimalPlaces: number;
  };
  const mergedTests = useMemo<mergedTest[]>(() => {
    if (
      !product ||
      !product.productTestConfigurations ||
      !sample ||
      !sample.testResults
    ) {
      return [];
    }
    return product.productTestConfigurations.map((testConfiguration) => {
      const testResult = sample.testResults!.find(
        (tr) =>
          tr.partNumber === sample.partNumber &&
          tr.canonicalTestName === testConfiguration.test.canonicalTestName,
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
        decimalPlaces:
          testResult?.decimalPlaces ?? testConfiguration.decimalPlaces,
      };
    });
  }, [
    product?.productTestConfigurations,
    sample?.testResults,
    sample?.partNumber,
    product?.partNumber,
  ]);

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
              <Input
                id="partNumber"
                type="text"
                value={sample?.partNumber ?? ""}
                readOnly
              />
            </Field>
            <Field>
              <FieldLabel htmlFor="status">Status</FieldLabel>
              <Select
                id="status"
                value={status ?? sample?.status ?? ""}
                onValueChange={setStatus}
              >
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
      </div>
      <div className="mt-8 space-y-4">
        {mergedTests.map((test) => (
          <div
            key={test.productTestSequence}
            className="bg-white dark:bg-gray-800 rounded-md p-4 border border-gray-200 dark:border-gray-700"
          >
            <form
              className="grid grid-cols-1 md:grid-cols-12 gap-4 items-center"
              onSubmit={(e) => {
                e.preventDefault();
                const value = new FormData(e.currentTarget).get("testResult");
                if (value === null || value === "") return;
                submitTestResult({
                  serialNumber: sample!.serialNumber,
                  partNumber: test.partNumber,
                  canonicalTestName: test.test.canonicalTestName,
                  modifiers: test.specificModifiers,
                  testResult: Number(value),
                  unit: test.unit,
                  decimalPlaces: test.decimalPlaces,
                  minValue: test.minValue,
                  maxValue: test.maxValue,
                });
              }}
            >
              <div className="md:col-span-5 min-w-0">
                <div className="font-semibold text-gray-900 dark:text-white">
                  {test.test.testName}
                  {test.specificModifiers.length > 0 && (
                    <span className="text-sm font-normal text-gray-500 dark:text-gray-400 ml-2">
                      ({test.specificModifiers.join(", ")})
                    </span>
                  )}
                </div>
                <div className="text-sm text-gray-500 dark:text-gray-400 whitespace-nowrap">
                  <span className="font-medium mr-1">Range:</span>
                  {test.minValue?.toFixed(test.decimalPlaces) ?? "N/A"} -{" "}
                  {test.maxValue?.toFixed(test.decimalPlaces) ?? "N/A"}{" "}
                  {test.unit}
                </div>
              </div>
              <div className="md:col-span-3">
                <Input
                  type="number"
                  placeholder="Result"
                  name="testResult"
                  defaultValue={test.testResult ?? ""}
                  step={
                    test.decimalPlaces > 0
                      ? 1 / Math.pow(10, test.decimalPlaces)
                      : "1"
                  }
                  className="h-9 w-full"
                  disabled={test.testResult ? true : false}
                />
              </div>
              <div className="md:col-span-2">
                <span className="text-sm text-gray-500 dark:text-gray-400 truncate block">
                  {test.unit}
                </span>
              </div>
              <Button
                type="submit"
                size="sm"
                className="w-full md:w-auto md:col-span-2 md:justify-self-end"
                disabled={test.testResult ? true : false}
              >
                Submit
              </Button>
            </form>
          </div>
        ))}
      </div>
    </div>
  );
}
