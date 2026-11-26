# Product requirements

## Authentication

All users must be authenticated.

## Interoperability

The system must support interoperability with customers' existing systems
through HTTP and message-based interfaces.

## Deployment

Deployments must be single-tenant. A hosted deployment must isolate each tenant
in a separate environment, such as a Kubernetes namespace.

## Test-result capture

- The system must let a user submit a configured test result for a sample.
- Test results must be numeric measurements.
- A product-test configuration must be able to allow an unbounded set of input
  values and select one of these outcomes: one result derived from the average,
  minimum, maximum, median, or count of the inputs; or a separate result for
  each input value.
- The system must retain the source input values with each recorded result for
  later display.
- In separate-result mode, each result must retain a one-element input array;
  the entry session does not need to be retained as a group.
- A product-test configuration must be able to enable an optional free-text
  comment on result entry. The system must retain a submitted comment with the
  recorded result.
- A sample must have an availability state: available to test or not available
  to test.
- The system must prevent result submission for a sample that is not available
  for testing.
- Invalid, duplicate, or incomplete submissions
  must be prevented or clearly surfaced.
- The system must retain an append-only result history for each sample.
- Result history must be retained indefinitely unless an external process
  manages retention.
- Voiding a result must create a separate history event without changing the
  original recorded result.
- Voiding is terminal. A replacement must be recorded as a new result.
- A replacement must not explicitly link to a voided result. Current effective
  result queries must exclude voided results.
- Current effective results must be exposed through a materialized view with
  one row per non-voided result.
- Once the user navigates away from a voided result, the user interface must
  not make it viewable.
- A voided result must be excluded from future subscription matching and
  statistical calculations.
- When a result is voided, Intake must publish a NATS
  `sample-history-invalidated` event with the sample identity. Downstream
  consumers must reload the sample history and rebuild affected derived state.
- New results must use the normal incremental result-event path rather than a
  sample-history invalidation event.
- Sample-history timestamps must provide user-facing order. Statistical
  calculations must not depend on input-event order.
- Only an authorized user may void a result. A void event must record the
  acting user, timestamp, and reason, and may include an optional comment.
- The system must support both complete sample-history queries and queries of
  current effective results.
- The intake service must compute and verify a digest for every result-history
  event when it is read or replayed. The digest must use HMAC-SHA-256 when an
  HMAC key is configured and SHA-256 otherwise, over a deterministic canonical
  JSON representation of the event.
- When HMAC is configured, the Intake application must use one key supplied
  through an environment-appropriate secret-management mechanism.
- Every stored result is traceable to its sample, product, configured test,
  submitting user, and effective specification context.

Result identity and the preserved interpretation fields are defined in the
[domain model](domain-model.md).

## Sample metadata

The system must extract named manufacturing metadata from sample serial numbers
using the parsing-scheme version mapped to the product at the sample's creation
time. For example:

```text
KE-TST-HE01-25012-CHA-00001
```

could expose `line=KE`, `product=TST`, `part=HE01`, `julian=25012`, and
`plant=CHA`. Subscription rules must be able to match these named values.

The system must support regular expressions, Lua scripts, and webhooks as
serial-number parsing mechanisms.
Lua scripts must execute in a QMEssentials sandbox.
Lua parsers must be pure functions that accept a serial number and return
metadata. They must not access the filesystem, network, database, or system.

Each parsing-scheme version must define a JSON Schema for parser output. The
system must reject a parse whose output does not satisfy that schema and notify
the result producer through the rejection-outcome process.

Parsing schemes must be versioned. A product must map to a parsing-scheme
version for an effective time span, whose beginning or ending may be open.
Each result must retain the parsing-scheme version that interpreted its serial
number.

Webhook parsing must use the parsing-scheme version's configurable, ordered
retry waits. A sample that exhausts them must be sent to the webhook process's
dead-letter queue. The webhook process owns monitoring and time-to-live
management for that queue and is external to QMEssentials.
Webhook output that fails JSON Schema validation must be rejected immediately
without consuming a retry wait.
Only an HTTP 200 webhook response is successful. A webhook 4xx response must
be rejected and notify the result producer. Any other HTTP response, including
non-200 2xx, 3xx, and 5xx responses, must consume configured retry waits and
send the sample to the webhook process's dead-letter queue when they are
exhausted.
Timeouts must follow the same retry and dead-letter-queue path.
Initial webhook authentication must be transport-level, such as mutual TLS. It
does not require application-level credentials or request signing.

The sample-creation user interface must validate serial-number parsability as
the user enters it. Intake must reject a test result whose sample serial number
cannot be parsed by its product's configured scheme.

Result producers must be able to observe asynchronous rejection outcomes and
take appropriate action.

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
