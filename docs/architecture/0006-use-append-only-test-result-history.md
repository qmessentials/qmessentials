# 6. Use append-only test-result history

Date: 2026-09-26

## Status

Accepted

## Context

Test-result values must remain traceable after they are captured, while later
actions such as voiding must also be recorded. Updating a result row to add
void information makes the original record and its audit history harder to
distinguish.

The system must support both a queryable sample history and asynchronous event
consumers. Replaying retained events is appropriate for recovery and rebuilding
derived state, but ordinary sample-history queries should not require a message
replay.

## Decision

The intake service will store authoritative test-result history as append-only
records in PostgreSQL. Recording a test result and voiding a result are separate
history events. The original recorded result is not updated when it is voided.
Voiding is terminal; a replacement is recorded as a new result.
Replacements do not link explicitly to voided results.

The intake application will only insert into and select from its authoritative
history table. Database-level enforcement of this restriction is not required
at this stage; service ownership and application behavior enforce it.

Queryable current-state projections may be separately maintained or rebuilt
from the history. Those projections are not authoritative and may use updates
or replacement as their implementation requires.
The current effective-results projection is a materialized view with one row
per non-voided result.

When a result is voided, Intake will publish a `sample-history-invalidated`
event to NATS with the sample identity. Downstream consumers subscribe to this
subject, reload the authoritative sample history, and rebuild affected derived
state.
This subject is emitted only for voids; new results use the normal incremental
result-event path.

History-event timestamps order sample history for users. Calculations must not
depend on a particular input-event order.

Each history event will carry a digest computed by the intake service from its
canonical JSON event payload. When an HMAC key is configured, the digest will
be an HMAC-SHA-256 value; otherwise, it will be a SHA-256 hash. The service
must verify the digest whenever it reads or replays an event. When HMAC is
configured, Intake uses one application-wide key supplied by an
environment-appropriate secret-management mechanism rather than an application
database table.

## Consequences

The full history for a sample is queryable directly from PostgreSQL. Consumers
and projections can replay the same history for recovery without making replay
the ordinary query path. A void action does not overwrite the original result.

Canonical JSON ensures that identical event values produce the same digest. A
SHA-256 hash detects accidental changes but, unlike an HMAC, does not protect
against an actor who can change both the event and its hash.
