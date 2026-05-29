create table schema_migrations (
    version_number int primary key,
    applied_at timestamp with time zone not null default now()
);

insert into schema_migrations (version_number) values (0);