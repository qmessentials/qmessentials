# Persistence model

This reference describes PostgreSQL tables by owning service. Domain concepts
and relationships are described in the [domain model](domain-model.md).

All data-owning services keep migration bookkeeping in `schema_migrations`:

| Column                     | Type                     | Notes                      |
| -------------------------- | ------------------------ | -------------------------- |
| `domain`, `version_number` | `text`, `int`            | Composite primary key      |
| `applied_at`               | timestamp with time zone | Migration application time |

## Configuration

### `products`

| Column                               | Type                     | Notes                                                 |
| ------------------------------------ | ------------------------ | ----------------------------------------------------- |
| `id`                                 | identity integer         | Primary key                                           |
| `part_number`, `product_name`        | text                     | Part number is unique                                 |
| `is_active`                          | boolean                  | Active product flag                                   |
| `created_at`, `updated_at`           | timestamp with time zone | Audit timestamps                                      |

### `serial_number_parsing_schemes`

| Column | Type | Notes |
| --- | --- | --- |
| `id` | identity integer | Primary key |
| `scheme_name` | text | Unique scheme name |
| `created_at`, `updated_at` | timestamp with time zone | Audit timestamps |

### `serial_number_parsing_scheme_versions`

| Column | Type | Notes |
| --- | --- | --- |
| `id` | identity integer | Primary key |
| `parsing_scheme_id`, `version_number` | integer | Unique version within a scheme |
| `parser_type` | text | References `serial_number_parser_types` |
| `parser_definition` | text | Regular expression, sandboxed Lua script, or external webhook configuration |
| `output_schema` | JSONB | JSON Schema that parser output must satisfy |
| `created_at`, `updated_at` | timestamp with time zone | Audit timestamps |

### `serial_number_parsing_scheme_version_retry_waits`

| Column | Type | Notes |
| --- | --- | --- |
| `parsing_scheme_version_id`, `attempt_number` | integer | Composite primary key and retry order |
| `wait_duration` | interval | Wait before this retry attempt |
| `created_at`, `updated_at` | timestamp with time zone | Audit timestamps |

### `product_parsing_scheme_versions`

| Column | Type | Notes |
| --- | --- | --- |
| `product_id`, `parsing_scheme_version_id` | integer | Product and parsing-scheme version references |
| `effective_from`, `effective_until` | timestamp with time zone | Sample creation time selects this span; null means an open bound |
| `created_at`, `updated_at` | timestamp with time zone | Audit timestamps |

Effective spans for one product must not overlap.

### `serial_number_parser_types`

| Column | Type | Notes |
| --- | --- | --- |
| `parser_type` | text | Primary key: `regex`, `lua`, or `webhook` |
| `is_active` | boolean | Active parser type flag |
| `created_at`, `updated_at` | timestamp with time zone | Audit timestamps |

### `test_unit_categories`

| Column                     | Type                     | Notes                |
| -------------------------- | ------------------------ | -------------------- |
| `category`                 | text                     | Primary key          |
| `is_active`                | boolean                  | Active category flag |
| `created_at`, `updated_at` | timestamp with time zone | Audit timestamps     |

### `test_units`

| Column                                    | Type                     | Notes                             |
| ----------------------------------------- | ------------------------ | --------------------------------- |
| `unit`                                    | text                     | Primary key                       |
| `unit_category`                           | text                     | References `test_unit_categories` |
| `unit_abbreviation`, `measurement_system` | text                     | Display and system metadata       |
| `is_active`                               | boolean                  | Active unit flag                  |
| `created_at`, `updated_at`                | timestamp with time zone | Audit timestamps                  |

### `tests`

| Column                               | Type                     | Notes                                                 |
| ------------------------------------ | ------------------------ | ----------------------------------------------------- |
| `id`                                 | identity integer         | Primary key                                           |
| `canonical_test_name`, `test_name`   | text                     | Canonical name is unique; `test_name` is display text |
| `test_unit_category`                 | text                     | Optional unit-category reference                      |
| `documentation_references`           | text array               | Supporting documentation                              |
| `are_modifiers_allowed`, `is_active` | boolean                  | Test behavior and availability                        |
| `created_at`, `updated_at`           | timestamp with time zone | Audit timestamps                                      |

### `modifiers` and `test_allowed_modifiers`

| Table                    | Columns                                                             | Purpose                      |
| ------------------------ | ------------------------------------------------------------------- | ---------------------------- |
| `modifiers`              | `modifier` (PK), `is_active`, audit timestamps                      | Catalog of modifier values   |
| `test_allowed_modifiers` | `test_id`, `modifier` (composite PK), `is_active`, audit timestamps | Modifiers allowed for a test |

### `product_test_configurations`

| Column                     | Type                     | Notes                                    |
| -------------------------- | ------------------------ | ---------------------------------------- |
| `id`                       | identity integer         | Primary key                              |
| `product_id`, `test_id`    | integer                  | References product and test              |
| `specific_modifiers`       | text array               | Optional modifier-specific configuration |
| `unit`, `decimal_places`   | text, integer            | Display and measurement format           |
| `min_value`, `max_value`   | double precision         | Optional specification limits            |
| `product_test_sequence`    | integer                  | Unique within a product                  |
| `is_critical`, `is_active` | boolean                  | Configuration behavior and availability  |
| `documentation_references` | text array               | Supporting documentation                 |
| `result_input_behavior`    | text                     | References `result_input_behaviors`      |
| `result_comment_enabled`   | boolean                  | Enables an optional result comment       |
| `created_at`, `updated_at` | timestamp with time zone | Audit timestamps                         |

### `result_input_behaviors`

| Column                     | Type                     | Notes                                                                                    |
| -------------------------- | ------------------------ | ---------------------------------------------------------------------------------------- |
| `behavior`                 | text                     | Primary key: `single`, `average`, `minimum`, `maximum`, `median`, `count`, or `separate` |
| `is_active`                | boolean                  | Active behavior flag                                                                     |
| `created_at`, `updated_at` | timestamp with time zone | Audit timestamps                                                                         |

### `product_metadata_definitions`

The selected parsing-scheme version's JSON Schema is authoritative for metadata
structure and types. These rows are the product's queryable-field catalog.

| Column                       | Type                     | Notes                           |
| ---------------------------- | ------------------------ | ------------------------------- |
| `product_id`, `metadata_key` | integer, text            | Composite primary key           |
| `value_type`                 | text                     | Derived from the parser version's JSON Schema |
| `is_active`                  | boolean                  | Makes a schema-defined field available for querying |
| `created_at`, `updated_at`   | timestamp with time zone | Audit timestamps                |

### `void_reasons`

| Column                     | Type                     | Notes              |
| -------------------------- | ------------------------ | ------------------ |
| `reason`                   | text                     | Primary key        |
| `is_active`                | boolean                  | Active reason flag |
| `created_at`, `updated_at` | timestamp with time zone | Audit timestamps   |

## Intake

### `samples`

| Column                     | Type                     | Notes                            |
| -------------------------- | ------------------------ | -------------------------------- |
| `id`                       | identity integer         | Primary key                      |
| `serial_number`            | text                     | Unique external identity         |
| `part_number`              | text                     | Product identity snapshot        |
| `is_available_for_testing` | boolean                  | Whether results may be submitted |
| `created_at`, `updated_at` | timestamp with time zone | Audit timestamps                 |

### `test_result_events`

| Column                                                             | Type                                     | Notes                                                            |
| ------------------------------------------------------------------ | ---------------------------------------- | ---------------------------------------------------------------- |
| `id`                                                               | UUID                                     | Primary key and event identity                                   |
| `result_id`                                                        | UUID                                     | Identifies the result history; the recorded event establishes it |
| `event_type`                                                       | text                                     | `recorded` or `voided`                                           |
| `serial_number`, `part_number`, `canonical_test_name`, `modifiers` | text, text, text, text array             | Result identity snapshot on a recorded event                     |
| `parsing_scheme_version_id`                                      | integer                                  | Version that interpreted the serial number                       |
| `result_value`, `input_values`                                     | double precision, double precision array | Numeric measurement and original input values                    |
| `unit`, `decimal_places`, `min_value`, `max_value`                 | measurement fields                       | Interpretation snapshot on a recorded event                      |
| `result_comment`                                                   | text                                     | Optional comment on a recorded event                              |
| `acting_user`, `void_reason`, `void_comment`                       | text                                     | Void-event actor, required reason, optional comment              |
| `occurred_at`                                                      | timestamp with time zone                 | User-facing history order                                        |
| `digest_algorithm`, `digest_value`                                 | text                                     | Digest of canonical JSON using HMAC-SHA-256 or SHA-256           |

Rows are authoritative, append-only history. The current effective-results
projection is derived from these rows for efficient queries.

### `current_test_results` materialized view

| Column | Source | Notes |
| --- | --- | --- |
| result fields | `test_result_events` | One row for each non-voided result |

The view is derived from authoritative history and is not itself authoritative.

## Subscription

### `subscriptions`

| Column                       | Type                     | Notes                          |
| ---------------------------- | ------------------------ | ------------------------------ |
| `id`                         | identity integer         | Primary key                    |
| `owner_user_id`, `rule_text` | text                     | Rule ownership and source text |
| `version_id`, `is_active`    | integer, boolean         | Revision and activation state  |
| `created_at`, `updated_at`   | timestamp with time zone | Audit timestamps               |

### `subscription_notification_enrollments`

| Column                           | Type                     | Notes                        |
| -------------------------------- | ------------------------ | ---------------------------- |
| `id`                             | identity integer         | Primary key                  |
| `subscription_id`                | integer                  | References `subscriptions`   |
| `notification_type`, `is_active` | enum, boolean            | Channel and activation state |
| `created_at`, `updated_at`       | timestamp with time zone | Audit timestamps             |

### Channel enrollment tables

| Table                                         | Columns                             | Purpose           |
| --------------------------------------------- | ----------------------------------- | ----------------- |
| `subscription_notification_email_enrollments` | enrollment ID (PK), `email_address` | Email destination |
| `subscription_notification_sms_enrollments`   | enrollment ID (PK), `phone_number`  | SMS destination   |

### `subscription_hits`

| Column                              | Type                     | Notes                                   |
| ----------------------------------- | ------------------------ | --------------------------------------- |
| `id`                                | identity bigint          | Primary key                             |
| `subscription_id`, `test_result_id` | integer, UUID            | Matching subscription and source result |
| `created_at`                        | timestamp with time zone | Match time                              |
| unique key                          | subscription and result  | Prevents duplicate hits                 |

### `subscription_notifications`

| Column                                                           | Type                     | Notes                                 |
| ---------------------------------------------------------------- | ------------------------ | ------------------------------------- |
| `id`                                                             | identity bigint          | Primary key                           |
| `subscription_hit_id`, `subscription_notification_enrollment_id` | bigint, integer          | Source hit and destination enrollment |
| `notification_status`                                            | text                     | Delivery-work status                  |
| `created_at`, `updated_at`                                       | timestamp with time zone | Audit timestamps                      |
