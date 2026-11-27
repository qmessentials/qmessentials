create table users (
    user_id text primary key,
    password_hash text not null,
    disabled boolean not null default false,
    created_at timestamp with time zone not null default now(),
    updated_at timestamp with time zone not null default now()
);
