CREATE TABLE accounts (
    id character varying NOT NULL,
    data jsonb DEFAULT '{}'::jsonb NOT NULL
);