-- +goose Up
-- +goose StatementBegin
create table users
(
    tg_id bigint primary key,
    state bigint not null default 0
);

create table devices
(
    tg_id      bigint not null,
    device_id text   not null
);

create table command_list (
   action text
);

create table commands
(
    tg_id bigint not null,
    device_id text ,
    command text ,
    action text,
    color text,
    done bool not null default false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
