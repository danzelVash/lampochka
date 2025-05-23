-- +goose Up
-- +goose StatementBegin
create table users
(
    tgID bigint primary key,
    state bigint not null default 0
)

create table devices
(
    user      bigint not null,
    device_id text   not null
)

create table commands
(
    user bigint   not null,
    device_id text not null ,
    command text ,
    action text,
    color text,
    done bool not null default false
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
