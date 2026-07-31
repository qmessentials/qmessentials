drop table product_test_configuration_allowed_modifier_combinations;

create table test_allowed_modifiers (
    test_id int not null references tests (id),
    modifier text not null references modifiers (modifier),
    is_active boolean not null default true,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now(),
    primary key (test_id, modifier)
);

alter table tests add column are_modifiers_allowed boolean not null default true;
