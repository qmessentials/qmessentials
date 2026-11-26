alter table products
add column metadata_type text null;

alter table products
add column metadata_selector text null;

create table product_metadata_definitions (
    product_id int not null references products (id),
    metadata_key text not null,
    value_type text not null,
    is_active boolean not null default true,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),
    primary key (product_id, metadata_key)
);
