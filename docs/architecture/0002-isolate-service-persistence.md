# 2. Isolate service persistence

Date: 2026-09-02

## Status

Accepted

## Context

Configuration, sample intake, and subscriptions represent different service
concerns. They share a PostgreSQL server in the initial deployment to keep
operations simple, but direct access to one another's tables would couple
service code and deployment decisions.

## Decision

Each service will own its persistence and migration stream. Service-owned
schemas may share a PostgreSQL server, but a service will not read or mutate
another service's schema. Cross-service access will use an API or a message.

## Consequences

Each schema can later move to a separate PostgreSQL server without removing
storage-level dependencies from other services. Ownership and schema-change
responsibility are clearer.

Some workflows require network calls, messages, or duplicated snapshots of
data that a database join could obtain more simply. Cross-service consistency
must be handled explicitly.
