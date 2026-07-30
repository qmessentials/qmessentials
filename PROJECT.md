# QMEssentials

## What This Is

QMEssentials is a quality-management application for small and medium-sized
manufacturers. It is intended to make product-quality testing easier to
configure, perform, and trace without requiring a large enterprise quality
management system.

The current implementation is an early, working vertical slice. A user can
view manufacturing samples, open a sample, retrieve the tests configured for
its product, and submit numeric test results. The system preserves the
effective test limits and units with each result and processes submitted
results asynchronously.

## Core Value

Manufacturing teams can reliably capture the right quality-test result for the
right product and sample with enough context to understand and audit it later.

## Target Users

- Quality technicians performing tests on manufactured samples
- Quality managers defining products, tests, limits, units, and procedures
- Small and medium-sized manufacturers that need practical quality
  traceability without enterprise-system complexity

The repository does not yet contain detailed personas or validated user
research, so these roles are inferred from the implemented domain and
interfaces.

## Current Capabilities

### Product and test configuration

- Store products identified by part number
- Define tests, measurement units, and unit categories
- Configure product-specific tests, sequence, modifiers, limits, precision,
  criticality, and documentation references
- Define selection options, permitted modifier combinations, and void reasons
- Retrieve a product and its test configuration by part number

### Sample testing

- Store and list samples by serial number, part number, and status
- Retrieve an individual sample and its recorded test results
- Present the product's configured tests for a sample
- Accept one numeric result per configured product-test sequence
- Record the effective unit, precision, minimum, and maximum with the result
- Hash submitted result payloads and persist them through a NATS JetStream
  work queue

### Application access

- Route browser API requests through a dedicated API gateway
- Require an internal shared secret between the gateway and backend services
- Provide a temporary cookie-based development login

## Architecture

QMEssentials is a containerized, service-oriented application:

- **Web:** React 19, TypeScript, Vite, TanStack Router, TanStack Query,
  Tailwind CSS, and Base UI
- **API gateway:** Go and Gin; handles browser-facing authentication, CORS,
  and reverse proxying
- **Configuration service:** Go, Gin, and PostgreSQL; owns products and test
  configuration
- **Intake service:** Go, Gin, PostgreSQL, and NATS JetStream; owns samples and
  test-result intake
- **Utilities:** Go-based database migration command
- **Infrastructure:** Docker Compose, PostgreSQL 18, and NATS

The services share a PostgreSQL server but use separate database schemas and
migration streams for configuration and intake concerns.

## Current Constraints

- Local orchestration is defined in `docker-compose.yml`.
- PostgreSQL is exposed to the host on port `5433` and remains available to
  containers on port `5432`.
- Backend services are configured through environment variables.
- Database migrations are ordered SQL files under `src/db-migrations`.
- The web application and Go services are separate build units.
- Authentication and authorization are placeholders and are not suitable for
  production use.
- Database connections currently disable TLS.
- The current result-entry UI supports numeric results; the schema anticipates
  richer configuration than the UI currently exposes.

## Near-Term Product Gaps

These gaps are derived from the current code and should not be treated as an
approved roadmap:

- Production authentication, user management, roles, and authorization
- UI and APIs for creating and maintaining product/test configuration
- Sample creation and lifecycle/status updates
- Result correction or voiding workflows with audit history
- Validation and clear pass/fail feedback against configured limits
- Support for non-numeric or selection-based tests
- Reporting, search, filtering, exports, and quality trends
- Automated test coverage and deployment environments
- Production security, observability, backup, and recovery practices

## Success Signals

Success criteria have not yet been formally defined. Useful initial measures
would include:

- Technicians can complete all required tests for a sample without referring
  to a separate source for limits or units.
- Every stored result can be traced to a sample, product, configured test
  sequence, submitting user, and effective specification.
- Invalid, duplicate, or incomplete test submissions are prevented or clearly
  surfaced.
- Product and test configuration changes do not corrupt the historical meaning
  of existing results.

## Decisions Reflected in the Code

- Separate configuration data from sample/result intake concerns.
- Use an API gateway as the browser-facing backend boundary.
- Snapshot result interpretation fields such as limits and units at submission
  time instead of relying only on mutable configuration.
- Queue result persistence through NATS JetStream.
- Use ordered, explicit SQL migrations rather than an ORM migration system.

## Target Architectural Direction

Planned result subscriptions will be implemented as an event-driven pipeline:
ingestion enriches results, a broker routes matches by subscription, and
generic calculation workers maintain statistical state and emit notification
events.

Subscription management, calculation brokering, and statistical calculation
will be separate services so their distinct API, CPU/memory, and stateful I/O
workloads can scale independently. Notification coordination initially belongs
to the subscription-management service; channel delivery may be extracted
later if its retry and scaling needs justify another service.

The initial deployment will use NATS JetStream for event transport and
PostgreSQL for both durable application data and rebuildable real-time
statistical state. Subscription filters will remain text-based sources of truth
and will be compiled into in-memory rules for high-throughput matching.

Each service owns isolated persistence. Separate PostgreSQL schemas may share
one server initially, but services do not read or write one another's schemas,
allowing each schema to move to a dedicated PostgreSQL instance later.

Detailed target design and unresolved implementation concerns are tracked in
`ARCHITECTURE.md`.

## Open Questions

- Which manufacturing segments and regulatory environments are in scope?
- Is QMEssentials intended to supplement or replace existing QMS/ERP/MES
  systems?
- What constitutes a complete sample, and how should sample status transition?
- Which test-result types beyond numeric measurements are required?
- What audit, electronic-signature, retention, and data-integrity requirements
  apply?
- How should failed or out-of-specification results be reviewed and resolved?
- Will deployments be single-tenant per manufacturer or multi-tenant?
- Which product outcomes should define the first production-ready release?

Detailed planned behavior is tracked in `REQUIREMENTS.md`.
