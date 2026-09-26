# 5. Store application configuration in database tables

Date: 2026-09-26

## Status

Accepted

## Context

Application-level configuration must be inspectable, queryable, and changeable
in a traceable way. File-based settings make it difficult to expose configured
values to users or other application components and do not provide a natural
GitOps workflow for installation-specific changes.

Some settings are fixed lists of values whose availability may change over
time. Those lists need an explicit data model rather than being embedded in
application configuration or code.

## Decision

Application-level configuration settings will be stored in database tables,
rather than configuration files.

Each fixed list of configurable values will have its own table. Where
appropriate, a value will have either an `active` flag or effective beginning
and ending dates. A null beginning or ending date means that bound is open.

Changes to installation-specific configuration will be applied through tracked,
installation-specific schema migrations in the owning service's migration
stream. This preserves a GitOps-compatible record of configuration changes.

This decision does not define the mechanism that selects and applies
installation-specific configuration. Forks and other deployment-specific
mechanisms remain open options.

## Consequences

Application components, user interfaces, and integrations can query a shared,
durable representation of configuration. Configuration changes can be reviewed,
versioned, and deployed with the same traceability as other migrations.

Each setting and fixed list requires schema design, ownership, and a migration
path. Consumers must apply `active` or effective-date rules consistently.
Installation-specific migrations also require a deployment process that selects
the correct configuration for the installation.
