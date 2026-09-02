# Product overview

## Purpose

QMEssentials is a quality-management application for small and medium-sized
manufacturers. It makes product-quality testing easier to configure, perform,
monitor, and trace without requiring a large enterprise quality-management
system.

Its core value is reliable capture of the right quality-test result for the
right product and sample, with enough preserved context to interpret and audit
the result later.

## Intended users

- Quality technicians performing tests on manufactured samples
- Quality managers defining products, tests, limits, units, and procedures
- Manufacturing teams monitoring results and responding to quality signals
- Manufacturers that need practical quality traceability without
  enterprise-system complexity

Detailed personas and assumptions about manufacturing segments and regulatory
environments still require product discovery.

## Product shape

Quality managers define products and their required tests, including units,
precision, specification limits, modifiers, execution order, and supporting
documentation. Technicians identify a manufacturing sample, perform the tests
configured for its product, and submit results with their effective
specification context preserved.

Users create text-based subscriptions to select results by product, sample
metadata, test details, result characteristics, and statistical window. The
system evaluates new results against those subscriptions, maintains any needed
statistical state, and produces notifications when a rule matches.

The [product requirements](../reference/product-requirements.md) specify the
intended behavior in more detail. Open questions and potential work are staged
in the [work-tracker inbox](../work-tracker-inbox.md) before transfer to the
authoritative work tracker.
