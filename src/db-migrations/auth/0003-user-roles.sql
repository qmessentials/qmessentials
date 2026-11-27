create table user_roles (
    user_id text not null references users (user_id),
    role text not null,
    primary key (user_id, role)
);
