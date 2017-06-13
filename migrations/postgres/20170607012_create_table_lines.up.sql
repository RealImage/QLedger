CREATE TABLE lines (
    id bigint NOT NULL,
    transaction_id character varying NOT NULL,
    account_id character varying NOT NULL,
    delta integer NOT NULL
);