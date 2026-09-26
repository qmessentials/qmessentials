# Work-tracker inbox

This is a temporary staging area for work discovered while developing or
documenting QMEssentials. It is not the authoritative backlog and does not
imply priority or commitment.

Move an item to the external work tracker when it can be prioritized, then
remove it here or replace it with a tracker reference. Accepted architectural
decisions belong in `docs/architecture`.

## Product discovery

- Validate target manufacturing segments, regulatory environments, and roles.
- Decide whether QMEssentials supplements or replaces QMS, ERP, or MES tools.
- Define a complete sample and its lifecycle transitions.
- Define review, correction, voiding, audit history, retention, and signatures.
- Select result types beyond numeric measurements.
- Decide whether deployments are single-tenant or multi-tenant.
- Establish measurable outcomes for a first production-ready release.

## Serial-number metadata

- Decide the scope that owns a parsing scheme.
- Define behavior for an unparseable serial number.
- Decide whether regular expressions suffice or sandboxed scripting is needed.

## Subscription rules and user experience

- Specify the query grammar, data types, comparisons, escaping, and normalization.
- Ensure validation, dry runs, and broker matching share one evaluator.
- Design dry runs and statistical diagnostics.
- Select notification channels and delivery, retry, and deduplication guarantees.

## Subscription event pipeline

- Define and version the enriched result-event contract and stable event ID.
- Define the active-subscription snapshot and subscription-change contract.
- Deliver and reconcile subscription changes at every broker instance.
- Define JetStream retention, consumer, acknowledgement, replay, and
  poison-message policies.
- Establish throughput, latency, recovery-time, and state-size targets.

## Statistical processing

- Define comparison-population dimensions and rolling-window configuration.
- Define sample-size, variance, and update-order semantics.
- Design atomic, idempotent state updates.
- Decide the persistence policy for rebuildable calculation state.
- Define replay, reconciliation, and handling of late, corrected, and voided results.
