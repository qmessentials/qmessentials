# Architecture

This document describes the intended system shape and the stable boundaries
toward which implementation should move. It is not an implementation-status
report. Code, tests, migrations, and deployment configuration determine what
exists today.

The reasons for accepted architectural choices belong in
[architecture decision records](../architecture/README.md).

## System shape

QMEssentials is a service-oriented, event-driven application:

```text
Browser
   |
   v
React web application
   |
   v
API gateway
   +--------> Configuration service ----> configuration schema
   +--------> Intake service -----------> intake schema
   +--------> Subscription service ------> subscription schema
                         |                         ^
                         v                         |
Enriched test events -> Calculation broker        |
                         |                         |
                         v                         |
                 Calculation engines -------------+
                         |
                         v
                  Notification events
```

NATS JetStream transports events between asynchronous stages. PostgreSQL holds
durable application data and service-owned processing state.

## Browser boundary

The React application calls a dedicated API gateway rather than addressing
domain services directly. The gateway owns the browser-facing authentication,
authorization, routing, and cross-origin boundary. Domain services expose
contracts intended for the gateway or other trusted services.

## Service ownership

The configuration service owns products, tests, units, and product-test
configuration. The intake service owns samples, captured test results, and
result enrichment. The subscription service owns subscription definitions,
revisions, activation, and notification coordination.

The calculation broker owns compiled subscription rules, matches enriched
events, and routes matches by subscription. Calculation engines own
statistical processing state, anomaly detection, processed-event identities,
and match events.

Each data-owning service has exclusive ownership of its schema and migration
stream. Schemas may share a PostgreSQL server, but services communicate through
APIs or messages rather than cross-schema access. The rationale is recorded in
[ADR 0002](../architecture/0002-isolate-service-persistence.md).

## Result event flow

Intake resolves canonical test details, parses applicable sample metadata, and
publishes a versioned enriched result event. Enrichment happens once so
downstream calculation does not repeatedly query product configuration.

The calculation broker maintains an in-memory representation of active
subscription rules. It narrows candidate subscriptions, evaluates their rules,
and routes each match by subscription identity. Calculation engines consume
those routed readings, update statistical state atomically, and emit match
events. Subscription management coordinates notification delivery from those
events.

Messages have stable identities, and consumers are idempotent because
JetStream delivery can be repeated. Rebuildable processing state can be
reconstructed from retained events or authoritative result history.

## Persistence and deployment

Durable domain records use logged PostgreSQL tables. Derived statistical state
may use a different persistence policy when it is demonstrably rebuildable,
but that choice depends on recovery and performance requirements.

Components remain independently deployable where their workloads differ:
subscription management is primarily transactional, broker matching is
CPU- and memory-oriented, and calculation is driven by statistical complexity
and state I/O. Initial deployments may colocate them on one host without
erasing their ownership boundaries.
