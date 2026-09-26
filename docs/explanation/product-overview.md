# Product overview

## Purpose

QMEssentials is a quality-management application for small and medium-sized
manufacturers. It makes product-quality testing easier to configure, perform,
monitor, and trace without requiring a large enterprise quality-management
system.

Its core value is reliable capture of the right quality-test result for the
right product and sample, with enough preserved context to interpret and audit
the result later. It provides a system of record for that quality information,
which customers can use to support their own compliance activities.

## Market fit

QMEssentials is intended to be applicable across manufacturing segments. A
prospective customer's process must be able to assign meaningful part numbers
to products, identify samples with serial numbers, and test products as part of
manufacturing or quality control.

QMEssentials does not seek to satisfy a particular regulatory or accrediting
standard. It records and makes quality information accessible; customers remain
responsible for applying that information to their own compliance obligations.

## Interoperability

QMEssentials is intended to interoperate with customers' existing systems
through HTTP and message-based interfaces.

## Deployment model

QMEssentials is designed for single-tenant deployments. If it is hosted, each
tenant must use an isolated environment, such as a Kubernetes namespace.

## Target users

- Quality technicians performing tests on manufactured samples
- Quality managers defining products, tests, limits, units, and procedures
- Manufacturing teams monitoring results and responding to quality signals
- Manufacturers that need practical quality traceability
  without enterprise-system complexity

## Product shape

**Quality Managers** define products and their required tests, including units,
precision, specification limits, modifiers, execution order, and supporting
documentation.

**Technicians** identify a manufacturing sample, perform the tests configured
for its product, and submit results with their effective specification context
preserved.

**Subscribers** create text-based subscriptions to select results by product,
sample metadata, test details, result characteristics, and statistical window.
The system evaluates new results against those subscriptions, maintains any
needed statistical state, and produces notifications when a rule matches.

The [product requirements](../reference/product-requirements.md) specify
intended behavior in more detail.
