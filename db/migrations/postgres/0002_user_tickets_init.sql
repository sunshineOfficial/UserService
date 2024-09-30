-- +goose Up
create table if not exists user_tickets
(
    user_id   uuid references users (id) on delete cascade,
    ticket_id varchar(48) not null,
    constraint pk_user_tickets primary key (user_id, ticket_id)
);

create index if not exists user_tickets_user_id_idx on user_tickets (user_id);

-- +goose Down
drop table if exists user_tickets;