# 4. Use ordered SQL migration files

Date: 2026-09-02

## Status

Accepted

## Context

Each data-owning service needs reviewable, deterministic PostgreSQL schema
changes. The project does not use an object-relational mapper whose migration
model should define the schema.

## Decision

Database changes will be expressed as explicitly ordered SQL files under
`src/db-migrations/<service>`. The migration utility will apply each service's
migration stream in filename order.

## Consequences

Schema changes remain visible as PostgreSQL SQL and are grouped by owning
service. Contributors must choose migration ordering deliberately and lint SQL
with `just lint-sql`.

Applied migrations should be treated as immutable. Roll-forward corrections
require new migration files, and database portability is not an objective of
this approach.
