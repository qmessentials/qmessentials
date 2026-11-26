# QMEssentials

QMEssentials is an early-stage quality-management application for small and
medium-sized manufacturers. It helps teams configure product tests, record
results against manufacturing samples, and retain the context needed to
interpret those results later.

See the [product overview](docs/explanation/product-overview.md) for the
product's purpose and the
[architecture overview](docs/explanation/architecture.md) for the intended
system shape.

## Repository layout

- `src/web`: React browser application
- `src/apigw`: browser-facing Go API gateway
- `src/configuration`: product and test configuration service
- `src/intake`: sample and test-result intake service
- `src/subscription`: subscription CRUD service
- `src/calculation-broker`: calculation broker
- `src/calculation-engine`: calculation engine
- `src/db-migrations`: PostgreSQL migrations, grouped by owning service
- `src/utils`: database migration utility
- `test/data`: local development data
- `docs`: project documentation and architecture decisions

## Local development

Prerequisites and commands are documented in
[Run QMEssentials locally](docs/how-to/run-locally.md).

Common commands:

```sh
just run-services
just run-web
just lint-sql
```

## Documentation

The documentation follows the [Diataxis](https://diataxis.fr/) categories
where they are useful:

- [How-to guides](docs/how-to/)
- [Reference](docs/reference/)
- [Explanation](docs/explanation/)
- [Architecture decisions](docs/architecture/)
- [Work-tracker inbox](docs/work-tracker-inbox.md)

Repository instructions for coding agents are in [AGENTS.md](AGENTS.md).
