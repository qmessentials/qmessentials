# Architecture decision records

This directory contains accepted architecture decision records (ADRs).
They explain why consequential choices were made;
they are not a description of the entire current architecture
or a backlog of possible work.

The repository uses `adr-tools`.
The root `.adr-dir` file configures the CLI to store records here.
Create the next record from the repository root with:

```sh
adr new "Decision title"
```

ADR status describes the decision, not implementation completeness.
Supersede an ADR instead of silently rewriting a decision
whose context or consequences have materially changed.

## Decisions

- [0001: Record architecture decisions](0001-record-architecture-decisions.md)
- [0002: Isolate service persistence](0002-isolate-service-persistence.md)
- [0003: Use an API gateway as the browser backend](0003-use-an-api-gateway-as-the-browser-backend.md)
- [0004: Use ordered SQL migration files](0004-use-ordered-sql-migration-files.md)
- [0005: Store application configuration in database tables](0005-store-application-configuration-in-database-tables.md)
- [0006: Use append-only test-result history](0006-use-append-only-test-result-history.md)
