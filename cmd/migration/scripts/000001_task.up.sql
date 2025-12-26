CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

CREATE TABLE tasks
(
    id          uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now(),
    deleted_at  timestamptz,

    title       varchar(255) NOT NULL,
    description text NOT NULL,

    status      varchar(32) NOT NULL,
    duration    bigint NOT NULL
);
