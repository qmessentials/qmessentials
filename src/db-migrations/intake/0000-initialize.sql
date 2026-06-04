create table if not exists schema_migrations (
    domain text not null,
    version_number int not null,
    applied_at timestamp with time zone not null default now(),
    primary key (domain, version_number)
);
