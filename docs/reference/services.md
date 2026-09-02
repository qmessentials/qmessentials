# Service catalog

This page defines the intended responsibilities of the system's deployable
components. The presence of a responsibility here does not indicate that it is
fully implemented; inspect code and tests when implementation status matters.

| Component          | Location                 | Responsibility                                               |
| ------------------ | ------------------------ | ------------------------------------------------------------ |
| Web                | `src/web`                | Browser workflows for configuration, testing, and monitoring |
| API gateway        | `src/apigw`              | Browser authentication and routing to backend services       |
| Configuration      | `src/configuration`      | Products, tests, and product-test configuration               |
| Intake             | `src/intake`             | Samples, result submission, and result enrichment             |
| Subscription       | `src/subscription`       | Subscription lifecycle and notification coordination         |
| Calculation broker | `src/calculation-broker` | Subscription-rule matching and event routing                  |
| Calculation engine | `src/calculation-engine` | Statistical state, anomaly detection, and match events        |
| Migration utility  | `src/utils`              | Ordered service-owned database migrations                     |

Browser API requests pass through the API gateway. Backend services interact
through explicit APIs and NATS JetStream messages. Each data-owning service has
an isolated PostgreSQL schema and migration stream, even when schemas share a
PostgreSQL server.
