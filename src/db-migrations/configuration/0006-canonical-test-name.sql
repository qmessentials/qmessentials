alter table tests
add column canonical_test_name text;

update tests
set canonical_test_name = trim(
    both '_' from lower(regexp_replace(test_name, '[^a-zA-Z0-9]+', '_', 'g'))
);

alter table tests
alter column canonical_test_name set not null;

alter table tests
add constraint canonical_test_name_unique unique (canonical_test_name);
