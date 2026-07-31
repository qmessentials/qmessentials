create table subscriptions
(
    id int generated always as identity primary key,
    owner_user_id text not null,
    rule_text text not null,
    version_id int not null,
    is_active boolean not null,
    created_at    timestamp with time zone not null default now(),
    updated_at    timestamp with time zone not null default now()
);

create type notification_type as enum ('email', 'sms');

create table subscription_notification_enrollments (
    id int generated always as identity primary key,
    subscription_id int not null references subscriptions (id),
    notification_type notification_type not null,
    is_active boolean not null,
    created_at    timestamp with time zone not null default now(),
    updated_at    timestamp with time zone not null default now()
);

create table subscription_notification_email_enrollments (
    subscription_notification_enrollment_id int primary key
        references subscription_notification_enrollments (id),
    email_address text not null
);

create table subscription_notification_sms_enrollments (
    subscription_notification_enrollment_id int primary key
        references subscription_notification_enrollments (id),
    phone_number text not null
);

create table subscription_hits (
    id bigint generated always as identity primary key,
    subscription_id int not null references subscriptions (id),
    test_result_id uuid not null,
    created_at    timestamp with time zone not null default now(),
    unique(subscription_id, test_result_id)
);

create table subscription_notifications (
    id bigint generated always as identity primary key,
    subscription_hit_id int not null references subscription_hits (id),
    subscription_notification_enrollment_id int not null
        references subscription_notification_enrollments (id),
    notification_status text not null,
    created_at    timestamp with time zone not null default now(),
    updated_at    timestamp with time zone not null default now(),
    unique(subscription_hit_id, subscription_notification_enrollment_id)
);
