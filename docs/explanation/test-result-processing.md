# Test-result processing

A result begins as a value for one of the tests configured for a sample's
product.

The web application sends the result through the API gateway to the intake
service. Intake obtains the relevant configuration and builds a result payload
that includes external domain identifiers and a snapshot of interpretation
fields, including:

- sample serial number;
- product part number;
- canonical test name;
- modifiers;
- result value;
- unit and decimal precision; and
- effective minimum and maximum specification values.

Intake persists the captured result and publishes a versioned, enriched event
to NATS JetStream. Persisting interpretation fields means later configuration
changes do not silently change what a historical result meant.

The calculation broker evaluates enriched events against active subscription
rules and routes matches by subscription. Calculation engines update any
required statistical state and emit match events. Subscription management
turns those matches into notification delivery work.

Stable event identities and idempotent consumers prevent redelivery from
altering statistics or notifying twice. Retained events support recovery of
rebuildable processing state.
