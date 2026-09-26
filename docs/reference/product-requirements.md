# Product requirements

## Test-result capture

- The system must let a user submit a configured test result for a sample.
- Invalid, duplicate, or incomplete submissions
  must be prevented or clearly surfaced.
- Every stored result is traceable to its sample, product, configured test,
  submitting user, and effective specification context.

Result identity and the preserved interpretation fields are defined in the
[domain model](domain-model.md).

## Sample metadata

The system must extract named manufacturing metadata from sample serial numbers
using a product-appropriate parsing scheme. For example:

```text
KE-TST-HE01-25012-CHA-00001
```

could expose `line=KE`, `product=TST`, `part=HE01`, `julian=25012`, and
`plant=CHA`. Subscription rules must be able to match these named values.

## Result subscriptions

Users must be able to subscribe to test results using:

1. part number and extracted sample metadata
2. canonical test name and modifiers
3. result characteristics
4. a time or statistical window

A subscription may select all matching results, out-of-specification results,
or results more than a configured number of standard units
from the applicable mean.
Advanced statistical process control is outside the current scope.

Statistical comparisons must use a configurable rolling window
with a default of seven days.

## Subscription authoring

The initial experience favors power users. Users must be able to:

- author filters using a supported text-query syntax
- see syntax and validation errors before activation
- dry-run a rule against historical results before saving it
- inspect the live observation count, mean, and standard deviation used by a
  statistical subscription

Validation, dry runs, and live evaluation
must share a grammar and evaluation semantics.

## Notification output

When a reading meets an active subscription's criteria,
the system must create an idempotent notification event
associated with both the subscription and the source result.
