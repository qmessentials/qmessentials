# Test-result processing

A result begins as a value for a test configured for a sample's product.
Its identity and relationships are defined
in the [domain model](../reference/domain-model.md).

When Intake captures the result,
it records the result identity
and a snapshot of the fields needed to interpret the value:

- sample serial number and product part number
- canonical test name and applicable modifiers
- result value, unit, and decimal precision
- effective minimum and maximum specification values

This snapshot preserves what the result meant when it was recorded.
Later configuration changes therefore
do not silently change the interpretation of a historical result.

After Intake persists the result, it publishes a versioned enriched event.
The [architecture's result-event flow](architecture.md#result-event-flow)
describes how the broker, calculation engine, and subscription service process it.

Stable event identities and idempotent consumers
prevent redelivery from altering statistics or notifying twice.
Retained events support recovery of rebuildable processing state.
