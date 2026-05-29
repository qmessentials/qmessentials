create table products (  
    id int generated always as identity primary key,  
    product_name text not null,  
    is_active bool not null default true,
    created_at timestamp with time zone not null default now(),  
    updated_at timestamp with time zone not null default now()  
);

create table test_unit_categories (
    category text primary key,
    is_active bool not null default true,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);

create table test_units (
    unit text primary key,
    unit_category text not null references test_unit_categories (category),
    unit_abbreviation text not null,
    measurement_system text not null, --examples; "metric", "imperial", "any"
    is_active bool not null default true,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);

create table tests (
    id int generated always as identity primary key,
    test_name text not null,
    test_unit_category text null references test_unit_categories (category),
    documentation_references text[] not null,
    is_active bool not null default true,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);

create table modifiers (
    modifier text primary key, --examples: "left", "top", "inside", "beginning", "outer"
    is_active bool not null default true,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);

--configures how a test can be used for a product; can be modifier-combo-specific or global
create table product_test_configurations (
    id int generated always as identity primary key,
    product_id int not null references products (id),
    test_id int not null references tests (id),
    specific_modifiers text[] null, --null means configuration applies to all modifiers; should match allowed modifier combinations for product test configuration, but nothing enforces that
    unit text not null references test_units (unit),
    decimal_places int null,
    min_value double precision null,
    max_value double precision null,
    documentation_references text[] null,
    is_critical bool not null,
    is_active bool not null default true,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);

--defines which modifier combinations are allowed for a given product test configuration
create table product_test_configuration_allowed_modifier_combinations (
    id int generated always as identity primary key,
    product_test_configuration_id int not null references product_test_configurations (id),
    modifier_combination text[] null, --null means a test result can be created without a modifier
    is_active bool not null default true,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);

--adds product-specific selections to the referenced test
create table product_test_configuration_added_selection_options (
    product_test_configuration_id int not null references product_test_configurations (id),
    selection_id int not null, --should not overlap with default test options
    label text not null,
    is_active bool not null default true,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),
    primary key (product_test_configuration_id, selection_id)
);

--"hides" specific selections for the referenced product and test
create table product_test_configuration_disabled_selection_options (
    product_test_configuration_id int not null references product_test_configurations (id),
    selection_id int not null, --should exist in default test options
    is_active bool not null default true,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),
    primary key (product_test_configuration_id, selection_id)
);

create table void_reasons (
    reason text primary key,
    is_active bool not null default true,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);

create table samples (
                         id int generated always as identity primary key,
                         sample_unique_identifier text not null unique,
                         product_id int not null references products (id),
                         status text not null, --examples: "in progress", "completed", "archived"
                         created_at timestamp with time zone not null default now(),
                         updated_at timestamp with time zone not null default now()
);

create table test_results (
    id int generated always as identity primary key,
    sample_id int not null references samples (id),
    product_test_configuration_id int not null references product_test_configurations (id),
    modifiers text[] null,
    test_result double precision not null,
    unit text not null references test_units (unit),
    decimal_places int not null,
    min_value double precision null, --at execution time (keeps from having to version test configurations)
    max_value double precision null,
    hash_value text not null,
    voided_at timestamp with time zone null,
    voided_by text null,
    voided_reason text null references void_reasons (reason),
    void_comment text null,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);