# Work-tracker inbox

This is a temporary staging area for work discovered while developing or
documenting QMEssentials. It is not the authoritative backlog and does not
imply priority or commitment.

Move an item to the external work tracker when it can be prioritized, then
remove it here or replace it with a tracker reference. Accepted architectural
decisions belong in `docs/architecture`.

## Potential customer-specific enhancements

- Assess workflow, record, and reporting enhancements for a specific compliance
  regimen when a customer need justifies them.
- Assess a blockchain-backed event-integrity mechanism in place of HMAC when a
  customer requires irrefutable history.
- Assess in-application result-history retention policies when a customer
  requires them.

## Possible future directions

- Build and maintain customer-specific integrations with QMS, ERP, and MES
  systems.
- Support QMS result exchange, including ingesting exported readings and
  publishing readings to systems that can import them.
- Assess custom ODBC adapters for IBM i (AS/400) systems with the required
  middleware; consider native Java-based tooling running on IBM i later.

## Application configuration

- Decide how installation-specific configuration is isolated and applied,
  including whether it uses forks or another deployment-specific mechanism.

## Asynchronous intake feedback

- Design rejected-result feedback for browser producers and programmatic
  producers. Browser input may use a BFF WebSocket hub; programmatic producers
  need a monitorable outcome queue.

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
