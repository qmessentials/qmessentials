# 6. Keep business logic in the application

Date: 2026-11-27

## Status

Accepted

## Context

Database constraints and triggers can enforce rules close to stored data, but
they can also hide business behavior from the application code that owns it.
Rules split between application code and migrations are harder to find, test,
and change.

Some rules need database enforcement to preserve relational integrity across
concurrent writes or to meet a demonstrated performance need.

## Decision

Business logic lives in the application by default. We encode a rule in the
database only when relational integrity or performance provides a specific
reason to do so.

Every new trigger, constraint, or comparable database-enforced rule that
encodes business policy requires an adjacent SQL comment explaining why it
belongs in the database rather than solely in application code. Ordinary
primary keys, foreign keys, and nullability that define table structure do
not need explanatory comments. One comment may cover related rules that share
the same justification.

Application validation remains responsible for clear errors, even when a
database rule also protects integrity. Database rules must not become an
undocumented substitute for application behavior.

This decision applies to new migrations and to database rules changed in later
migrations. Applied migrations remain immutable.

## Consequences

Business behavior remains visible in the owning application service. Database
enforcement is available where it protects relationships, concurrent writes,
or demonstrated performance, with its reason recorded beside the SQL.

Schema reviews require an explicit justification for database-enforced
business rules. Some rules may be checked in both the application and database,
so those checks must stay consistent.
