# Service catalog

This page is a catalog of the system's deployable components and their
responsibilities.

| Component          | Location                 | Responsibility                                               |
| ------------------ | ------------------------ | ------------------------------------------------------------ |
| Web                | `src/web`                | Browser workflows for configuration, testing, and monitoring |
| API gateway        | `src/apigw`              | Browser authentication and routing to backend services       |
| Configuration      | `src/configuration`      | Products, tests, and product-test configuration              |
| Intake             | `src/intake`             | Samples, result submission, and result enrichment            |
| Subscription       | `src/subscription`       | Subscription lifecycle and notification coordination         |
| Calculation broker | `src/calculation-broker` | Subscription-rule matching and event routing                 |
| Calculation engine | `src/calculation-engine` | Statistical state, anomaly detection, and match events       |
| Migration utility  | `src/utils`              | Ordered service-owned database migrations                    |

Component interactions are described in the
[architecture overview](../explanation/architecture.md). The browser-backend
and persistence-ownership decisions are recorded in
[ADR 0003](../architecture/0003-use-an-api-gateway-as-the-browser-backend.md)
and [ADR 0002](../architecture/0002-isolate-service-persistence.md).
