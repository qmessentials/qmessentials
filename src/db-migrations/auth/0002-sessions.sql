create table sessions (
    token_hash bytea primary key,
    user_id text not null references users (user_id),
    created_at timestamp with time zone not null default now(),
    expires_at timestamp with time zone not null,
    revoked_at timestamp with time zone
);
