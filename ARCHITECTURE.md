# QMEssentials Architecture

## Status and Scope

This document describes the target architecture for test-result enrichment,
subscription routing, statistical calculation, and notification triggering.
It extends the current system described in `PROJECT.md`; most components in
this document are not yet implemented.

The state-storage addendum in the design discussion supersedes the earlier
Redis proposal. PostgreSQL is the planned statistical-state store for the
initial single-VM deployment model.

## Service Boundaries

The target system uses separate services for workloads that need independent
scaling:

- **Subscription-management service:** Owns subscription CRUD, rule
  validation, dry runs, revisions, activation, notification coordination, and
  the Subscription API.
- **Calculation-broker service:** Owns the compiled rule cache, matches
  enriched events against subscriptions, and routes matches by subscription
  ID.
- **Calculation-engine service:** Owns statistical calculation, rebuildable
  calculation state, anomaly detection, and notification-event emission.

Notification channel delivery may initially run as a worker within the
subscription-management boundary. It can become a separate service if channel
credentials, retries, rate limits, or workload require independent ownership
and scaling.

The broker is deliberately separate from subscription management. Rule
matching is expected to be primarily CPU- and memory-bound, while subscription
management is an API and transactional workload. Calculation engines are
separate because their capacity is driven by matched-event volume, database
I/O, statistical complexity, and state contention.

### Persistence ownership

Each service exclusively owns its persistence:

- Intake owns samples, raw results, and ingestion publication state.
- Subscription management owns subscription definitions, revisions,
  activation, notification configuration, and delivery coordination records.
- Calculation engines own statistical state, processed-event identities, and
  anomaly records.
- The calculation broker initially has no durable application store; its
  compiled cache is disposable and reconstructed from subscription management.

Service schemas may reside on the same PostgreSQL server for initial
deployments. No service may directly query or mutate another service's schema.
This preserves the option to place each schema on a separate PostgreSQL server
without changing storage-level dependencies.

## Domain Segmentation

Test-result subscriptions operate over three domain layers plus a time window:

1. **Part number and sample metadata — where and what:** A serial number maps a
   sample to a product and can encode manufacturing attributes such as plant,
   line, and date.
2. **Test details — context:** A result is identified by sample, part number,
   product test sequence, and optional modifiers. The product and sequence
   resolve to a canonical test name.
3. **Result characteristics — trigger:** A subscription can select all
   results, configured out-of-specification results, or statistical anomalies.
4. **Statistical window — time:** Stateful calculations operate over a
   configurable window that defaults to seven days.

## Target Event Flow

```text
Raw result
    |
    v
Ingestion and enrichment
    |  parse serial number; resolve canonical test details
    v
Enriched test event in NATS JetStream
    |
    v
Calculation broker
    |  rules initialized from Subscription API
    |  evaluate compiled subscription rules
    +----> subscription.<subscription_id>.readings
                                      |
                                      v
                          Calculation engine queue group
                                      |
                                      +----> PostgreSQL statistical state
                                      |
                                      +----> notification event when matched
                                                    |
                                                    v
                                      Subscription management /
                                        notification delivery
```

### 1. Ingestion and Enrichment

The ingestion layer will:

- Accept a raw test result.
- Parse the sample serial number once using the configured regular expression
  or sandboxed script.
- Resolve the canonical test name and other required test context.
- Publish an enriched, versioned test event to the primary NATS JetStream.

The enriched event separates parsing and configuration lookups from
subscription calculation. Calculation workers therefore do not need direct
knowledge of product configuration for each message.

### 2. Calculation Broker

The calculation broker will:

- Consume enriched events from the primary stream.
- Evaluate each event against active, compiled subscription criteria.
- Route a matching event by subscription ID to:

```text
subscription.<subscription_id>.readings
```

Rules are evaluated independently per subscription. Equivalent subscriptions
may therefore duplicate downstream work. This is an intentional initial
tradeoff in favor of simpler routing, isolation, and auditability.

The broker will use inexpensive indexed attributes to narrow the candidate
rule set before evaluating compiled expression trees. The indexing strategy
will evolve from measured subscription volume and event distributions.

### 3. Calculation Engines

Calculation engine instances will:

- Join a shared NATS consumer/queue group.
- Consume subscription reading events so one engine processes each delivery.
- Atomically read and update the subscription's statistical state.
- Determine whether the reading meets the subscription trigger.
- Emit a notification event when the trigger matches.

Workers will not use sticky subscription assignments or process-local state as
the authoritative statistical state. This permits horizontal scaling and
worker replacement.

## Subscription Rule Evaluation

### Source of truth

PostgreSQL will store a subscription's filter criteria as one raw query string,
for example:

```text
line:KE AND plant:CHA AND part:HE*
```

The raw query is the authoritative representation. Compiled SQL fragments or a
second structured rule representation will not be persisted initially.

This PostgreSQL data belongs exclusively to the subscription-management
service. Calculation brokers obtain it through a service contract, not by
accessing the schema.

### Query syntax and execution

The planned syntax is a supported subset of Solr/Lucene query syntax.

The calculation broker will:

1. Subscribe to or begin buffering durable subscription-change events.
2. Request a consistent, paginated or compressed snapshot of active
   subscriptions from the Subscription API.
3. Parse and compile the snapshot into in-memory abstract syntax trees.
4. Apply changes after the snapshot's revision or stream position.
5. Report ready and begin consuming enriched test events.
6. Evaluate enriched events directly against the compiled trees.
7. Recompile an affected rule when its subscription changes.

The subscription API will publish change events on a NATS subject such as:

```text
subscriptions.events
```

Brokers will consume these events and replace the affected compiled rule
without interrupting event ingestion. Every broker instance needs the complete
active rule set, so subscription changes must be delivered to every broker
cache rather than distributed across the brokers as a shared queue-group
workload.

The snapshot response and change stream will share a monotonic revision,
cursor, or equivalent position. This closes the race in which a subscription
could change while a broker is loading its snapshot. Broker readiness will
remain false until the snapshot is compiled and all subsequent changes are
applied.

Initialization retries will use exponential backoff and jitter to avoid a
restart surge overwhelming subscription management. A broker-owned persistent
read model is deferred until measured snapshot volume, startup latency, or
replica churn demonstrates a need for it. If introduced, it will be populated
through events and owned solely by the broker service; it will not expose or
share the subscription-management schema.

The exact supported syntax, parser library, data types, wildcard behavior, and
case-sensitivity remain to be specified before implementation.

## Statistical Calculation

### State representation

The initial anomaly calculation will use Welford-style online statistics. Its
basic state consists of:

- Observation count (`n`)
- Running mean (`M`)
- Aggregate squared distance from the mean (`S`, commonly named `M2`)

These values allow mean and variance to be updated without querying all
historical readings for every event.

### State storage

For the initial single-VM deployment model, PostgreSQL will store real-time
statistical state in a dedicated `UNLOGGED` table, tentatively named
`subscription_state`.

Updates will use atomic `INSERT ... ON CONFLICT DO UPDATE` operations. The
active state is expected to remain in PostgreSQL shared buffers. PostgreSQL
memory configuration will be deployment-specific; a fixed percentage of
system memory is not an application-level invariant.

An unlogged table reduces write-ahead-log overhead but is not durable across
certain crashes and is not replicated through normal physical replication.
It is therefore a rebuildable cache, not the system of record.

### Rolling windows

Basic Welford state describes a cumulative population; by itself, it does not
implement the configurable rolling window required by `REQUIREMENTS.md`.
Implementation must add a defined expiration strategy, such as time buckets or
a numerically stable removal algorithm backed by retained observations.

The window algorithm must be decided before statistical anomaly subscriptions
are considered implementable.

## Delivery, Recovery, and Reconciliation

JetStream will retain enriched events long enough to support redelivery and
state rebuilding. Calculation engines must use durable consumer semantics,
explicit acknowledgements, and a defined acknowledgement timeout.

Because event delivery can be repeated, state updates and notification
emission must be idempotent. Each enriched result event needs a stable event
identity, and processing records or equivalent deduplication must prevent a
redelivered event from changing statistics or notifying twice.

If rebuildable PostgreSQL state is lost, a recovery process will replay the
applicable retained events to reconstruct it. A low-priority reconciliation
service may compare calculated state with the transactional system of record
and repair drift.

Recovery time depends on retained event volume and cannot be assumed to take
only seconds without measured capacity targets.

The system will emit a notification event after a subscription trigger
matches. User-facing delivery channels and their retry policies remain outside
this architecture until product requirements select them.

## Power-User Experience Supported by the Architecture

The target architecture supports:

- A syntax-highlighted editor for raw subscription queries
- Immediate syntax validation using the same grammar as the broker
- Dry runs against historical results before saving a subscription
- Diagnostics showing live statistical count, mean, and standard deviation

These are product requirements rather than prerequisites for the event
pipeline, but shared parsing and diagnostics interfaces should be designed so
the UI and broker cannot silently interpret a rule differently.

## Key Decisions

- Enrich each result once before subscription evaluation.
- Use NATS JetStream to decouple ingestion, routing, and calculation.
- Deploy subscription management, calculation brokering, and statistical
  calculation as separate services.
- Give each service exclusive ownership of its persistence and prohibit
  cross-schema access.
- Scale rule matching independently from subscription APIs and stateful
  calculation.
- Route matched readings by subscription ID.
- Scale generic calculation workers through a shared durable consumer.
- Keep authoritative calculation state outside worker memory.
- Use PostgreSQL, not Redis, for initial real-time statistical state.
- Treat unlogged statistical tables as rebuildable caches.
- Store the raw subscription query as its source of truth.
- Compile rules into in-memory ASTs and invalidate them through NATS events.
- Initialize brokers through a versioned Subscription API snapshot plus a
  durable change stream.
- Defer broker-owned persistent projections until operational measurements
  justify them.
- Use online statistics rather than querying full history per reading.

## Required Design Follow-ups

- Define the enriched event schema, versioning policy, and stable event ID.
- Specify the supported Solr/Lucene syntax subset and evaluation semantics.
- Select a parser that can run consistently in API validation and the broker.
- Define the Subscription API snapshot contract, pagination, compression, and
  revision semantics.
- Define durable change-event fan-out and rule-cache reconciliation behavior
  for every broker instance.
- Define cache indexing and candidate-selection strategies.
- Design atomic, idempotent state updates under concurrent delivery.
- Choose and validate the rolling-window statistics algorithm.
- Define anomaly comparison timing: before or after including the new reading.
- Define minimum sample size and sample-versus-population variance.
- Set JetStream retention, consumer, acknowledgement, and replay policies.
- Define how corrected, voided, and late readings alter state.
- Define reconciliation authority, frequency, and repair behavior.
- Select notification delivery channels and retry/deduplication guarantees.
- Establish throughput, latency, recovery-time, and state-size targets.
