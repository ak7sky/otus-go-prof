-- +goose Up
CREATE TABLE calendar.events (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title         text NOT NULL,
    descr         text,
    event_start   timestamptz NOT NULL,
    event_end     timestamptz NOT NULL,
    notify_before bigint NOT NULL
);

-- +goose Down
drop table calendar.events;
