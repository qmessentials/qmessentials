create table samples
(
    id            int generated always as identity primary key,
    serial_number text                     not null unique,
    part_number   text                     not null,
    status        text                     not null, --examples: "in progress", "completed", "archived"
    created_at    timestamp with time zone not null default now(),
    updated_at    timestamp with time zone not null default now()
);

create table test_results
(
    id                            uuid                              default uuidv7() primary key,
    sample_id                     int                      not null references samples (id),
    product_test_configuration_id int                      not null references product_test_configurations (id),
    modifiers                     text[]                   null,
    test_result                   double precision         not null,
    unit                          text                     not null references test_units (unit),
    decimal_places                int                      not null,
    min_value                     double precision         null, --at execution time (keeps from having to version test configurations)
    max_value                     double precision         null,
    hash_value                    text                     not null,
    voided_at                     timestamp with time zone null,
    voided_by                     text                     null,
    voided_reason                 text                     null references void_reasons (reason),
    void_comment                  text                     null,
    created_at                    timestamp with time zone not null default now(),
    updated_at                    timestamp with time zone not null default now()
);
