# QMEssentials Requirements

## Status

This document captures planned product behavior. Unless a requirement is also
described as a current capability in `PROJECT.md`, it should not be assumed to
be implemented.

## Result Subscriptions

Users will be able to subscribe to test results using four matching
dimensions:

1. Part number and sample metadata
2. Test details
3. Result characteristics
4. Time and statistical window

### Part Number and Sample Metadata

Tests are performed against samples. Each sample represents a specific product:

- Products are identified by part number.
- Samples are identified by serial number.
- A sample belongs to one product.

Serial numbers are assumed to encode information about the manufacturing
process that produced the item under test. For example:

```text
KE-TST-HE01-25012-CHA-00001
```

This could identify an instance of product `KE-TST-HE01`—a Kitchen Essentials
brand toaster heating element—manufactured on January 12, 2025, at the
Chattanooga, Tennessee plant.

#### Serial-number schemes

The system will support a product-specific scheme for extracting named
metadata from a sample's serial number.

The planned product configuration fields are:

- `sample_number_scheme_type`: the evaluator type, initially `regex` or `lua`
- `sample_number_scheme`: the expression or script used by that evaluator

A regular-expression scheme could use named capture groups:

```regex
^(?<line>[A-Z]+)-(?<product>[A-Z]+)-(?<part>[A-Z0-9]+)-(?<julian>\d{5})-(?<plant>[A-Z0-9]{3})-\d+$
```

For the example serial number, the resulting metadata would include values
such as `line=KE`, `product=TST`, `part=HE01`, `julian=25012`, and `plant=CHA`.

Lua, or another sandboxed scripting mechanism, may be used when a regular
expression cannot express the required parsing logic.

#### Subscription matching

When evaluating a result against a subscription, the system will:

1. Find the product associated with the result's part number.
2. Parse the sample serial number using that product's configured scheme.
3. Evaluate the extracted named values against the subscription's sample
   metadata rules.

An illustrative rule is:

```text
line = 'KE' AND plant = 'CHA' AND part = 'HE*'
```

The rule language remains to be selected. It may be a small purpose-built DSL
or an existing expression language.

### Test Details

A test result will be uniquely identified by the combination of:

- Sample ID
- Part number
- Product test sequence number
- Modifier or modifiers

For a product, each product test sequence number maps to one test name.

When evaluating a result against a subscription, the system will:

1. Use the result's part number and product test sequence number to resolve the
   test name.
2. Match the resolved test name against the subscription's test-name pattern,
   including wildcard support.
3. Include the result when the subscription has no modifier constraint, or
   when the result's modifiers match the subscription's modifier constraint.

### Result Characteristics

Product-test configuration may define minimum and maximum specification
values. The effective values are copied onto each test result when it is
submitted so that historical readings retain their original specification
context.

A subscription may select:

- All results
- Out-of-specification results
- Results more than a configured number of standard units from the mean

More advanced statistical process control is outside the current scope.

### Time and Statistical Window

Statistics used to compare a reading with its distribution will use a rolling
seven-day window by default.

The window will be configurable so users can control how quickly older data
rolls out of the distribution and, consequently, how quickly the statistics
adapt to process drift.

### Power-User Subscription Experience

The initial subscription interface will favor a text-oriented workflow for
power users.

Users will be able to:

- Author subscription filters using the supported text query syntax.
- Receive syntax and validation errors before a rule is activated.
- Dry-run a proposed rule against historical results before saving it.
- Inspect the live observation count, mean, and standard deviation used by a
  statistical subscription.

The UI and calculation broker must use the same query grammar and evaluation
semantics.

### Notification Output

When a reading meets a subscription's criteria, the system will create a
notification event associated with the subscription and source test result.

The user-facing delivery channels are not yet selected.

## Service and Isolation Requirements

- Subscription management, calculation brokering, and statistical calculation
  will be independently deployable services.
- Each service will own its persistence and data model.
- A service must not read or write another service's database schema.
- Service-owned schemas may share a PostgreSQL server initially, but each must
  be movable to a separate server without changing another service's storage
  access.
- The subscription-management service will own subscription definitions,
  revisions, activation state, notification coordination, and related durable
  records.
- The calculation-broker service will own subscription matching and routing.
  It will maintain compiled rules in memory but will not read the
  subscription-management database directly.
- The calculation-engine service will own statistical state, processed-event
  identities, and anomaly-calculation records.
- The services must be independently scalable so rule-matching CPU/memory load
  does not dictate calculation-engine capacity, and stateful calculation I/O
  does not dictate subscription API capacity.

### Broker Initialization and Rule Updates

Before accepting enriched result events, a calculation-broker instance must:

1. Obtain a consistent snapshot of all active subscriptions through the
   subscription-management API.
2. Parse and compile the snapshot's rules into its in-memory matching
   structures.
3. Apply subscription changes that occurred while the snapshot was loading.
4. Report ready only after its rule cache is current.

The snapshot and change stream must share a revision, cursor, or equivalent
ordering mechanism so a broker cannot miss a subscription change during
startup.

Subscription changes must be durably published and delivered to every active
broker instance. The broker must retry initialization with backoff when the
subscription-management service is unavailable.

The initial design will not add broker-owned persistent projections. A
broker-owned read model may be introduced if measured snapshot size, startup
latency, or replica churn makes API-based initialization too costly.

## Open Design Questions

- Does the serial-number scheme belong to a product, product family, plant, or
  another configuration entity?
- What should happen when a serial number cannot be parsed by its configured
  scheme?
- Is Lua the preferred script evaluator, and how will scripts be sandboxed,
  validated, versioned, and resource-limited?
- Which expression language should define sample-metadata subscription rules?
- What wildcard syntax and escaping rules should test-name and metadata
  patterns use?
- Are modifier comparisons ordered, set-based, case-sensitive, or normalized?
- Does "standard units from the mean" mean a z-score calculated from population
  or sample standard deviation?
- What minimum observation count is required before a statistical subscription
  can be evaluated?
- Which population dimensions define a result's comparison distribution—for
  example product, test, modifiers, plant, or line?
- Is the statistical window configured per subscription, per product/test, or
  globally?
- How should late-arriving, corrected, or voided readings affect statistics
  and previously emitted subscription matches?
- Through which channels will matching subscriptions be delivered?
