# Work-tracker inbox

This is a staging area for work discovered while developing or documenting
QMEssentials. It is not the authoritative backlog, does not imply priority or
commitment, and must not be treated as an implementation specification.

When an item is added to the external work tracker, remove it from this file or
replace it briefly with a tracker reference. Accepted architectural decisions
belong in `docs/architecture`, and implemented behavior belongs in the
appropriate explanation or reference document.

## Product discovery

- Validate target manufacturing segments, regulatory environments, and user
  roles.
- Decide whether QMEssentials supplements or replaces QMS, ERP, or MES tools.
- Define a complete sample and its lifecycle transitions.
- Define review, correction, voiding, audit history, retention, and electronic
  signature requirements.
- Select result types beyond numeric measurements.
- Decide whether deployments are single-tenant or multi-tenant.
- Establish measurable outcomes for a first production-ready release.

## Serial-number metadata

- Decide whether a parsing scheme belongs to a product, product family, plant,
  or another configuration entity.
- Define failure behavior for serial numbers that cannot be parsed.
- Decide whether sandboxed scripting is needed in addition to regular
  expressions; if so, select an evaluator and define validation, versioning,
  resource limits, and sandboxing.

## Subscription rules and user experience

- Specify the supported text-query grammar, data types, wildcard behavior,
  escaping, normalization, and case sensitivity.
- Decide how modifier comparisons behave.
- Use one parser and evaluation semantics for validation, dry runs, and broker
  matching.
- Design dry-run behavior and diagnostics for observation count, mean, and
  standard deviation.
- Select notification channels and define delivery, retry, and deduplication
  guarantees.

## Subscription event pipeline

- Define and version an enriched result-event schema with a stable event ID.
- Implement result enrichment and publish enriched events to JetStream.
- Define the active-subscription snapshot contract, pagination or compression,
  and a revision/cursor shared with subscription change events.
- Deliver durable subscription changes to every broker instance and reconcile
  broker rule caches.
- Compile and index rules in the calculation broker, then route matching
  readings by subscription ID.
- Define JetStream retention, consumer, acknowledgement, replay, and poison
  message policies for the pipeline.
- Establish throughput, latency, recovery-time, and state-size targets before
  choosing optimization strategies.

## Statistical processing

- Define the comparison population dimensions for each result.
- Choose a rolling-window algorithm and decide where the window is configured.
- Define minimum sample size, sample-versus-population variance, and whether a
  new reading is compared before or after it updates state.
- Design atomic and idempotent state updates for repeated or concurrent event
  delivery.
- Decide whether PostgreSQL unlogged tables are appropriate for rebuildable
  calculation state after measuring operational requirements.
- Define replay and reconciliation for lost state.
- Define how late, corrected, and voided results affect statistics and earlier
  notifications.
