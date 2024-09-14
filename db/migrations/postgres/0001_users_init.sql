-- +goose Up
create table if not exists users
(
    id      uuid primary key default gen_random_uuid(),
    email   text unique not null,
    name    text,
    surname text
);

create index if not exists users_email_idx on users (email);

-- +goose Down
drop table if exists users;