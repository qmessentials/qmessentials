# Product requirements

## Status and interpretation

This document captures intended product behavior. It is a living product
reference, not a claim that every capability is implemented and not a
prioritized roadmap. The [product overview](../explanation/product-overview.md)
describes the broader product intent. Unresolved design and prioritization work
is staged in the [work-tracker inbox](../work-tracker-inbox.md).

## Test-result capture

- A sample belongs to a product identified by part number and is identified by
  serial number.
- A captured result is identified by sample serial number, part number,
  canonical test name, and any applicable modifiers.
- Product-test sequence controls presentation and execution order; it is not
  part of result identity.
- The effective unit, precision, minimum, and maximum are copied onto a result
  at submission so historical interpretation survives configuration changes.
- Invalid, duplicate, or incomplete submissions are prevented or clearly
  surfaced.
- Every stored result is traceable to its sample, product, configured test,
  submitting user, and effective specification context.

## Sample metadata

The system will support extracting named manufacturing metadata from sample
serial numbers using a product-appropriate parsing scheme. For example:

```text
KE-TST-HE01-25012-CHA-00001
```

could expose `line=KE`, `product=TST`, `part=HE01`, `julian=25012`, and
`plant=CHA`. Subscription rules will be able to match these named values.

The ownership and evaluator for parsing schemes remain open design questions
in the work-tracker inbox.

## Result subscriptions

Users will be able to subscribe to test results using:

1. part number and extracted sample metadata;
2. canonical test name and modifiers;
3. result characteristics; and
4. a time or statistical window.

A subscription may select all matching results, out-of-specification results,
or results more than a configured number of standard units from the applicable
mean. Advanced statistical process control is outside the current scope.

Statistical comparisons will use a configurable rolling window with a default
of seven days. The statistical population, exact variance definition, minimum
observation count, and rolling algorithm remain to be decided.

## Subscription authoring

The initial experience will favor power users. Users will be able to:

- author filters using a supported text-query syntax;
- see syntax and validation errors before activation;
- dry-run a rule against historical results before saving it; and
- inspect the live observation count, mean, and standard deviation used by a
  statistical subscription.

Validation, dry runs, and live evaluation must share a grammar and evaluation
semantics. The grammar itself remains to be specified.

## Notification output

When a reading meets an active subscription's criteria, the system will create
an idempotent notification event associated with both the subscription and the
source result. User-facing delivery channels remain to be selected.

## Service qualities

- Configuration, intake, subscription management, matching, and statistical
  calculation respect service ownership boundaries.
- A service does not read or write another service's database schema.
- Event processing tolerates redelivery without changing statistics or
  notifying twice for the same event.
- Rebuildable state can be reconstructed from an authoritative result history
  or retained events.
- Production authentication, authorization, transport security,
  observability, backup, and recovery must be defined before a production
  release.
