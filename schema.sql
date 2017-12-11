SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;
CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';
SET search_path = public, pg_catalog;
SET default_tablespace = '';
SET default_with_oids = false;
CREATE TABLE accounts (
    id character varying NOT NULL,
    data jsonb DEFAULT '{}'::jsonb NOT NULL
);
CREATE TABLE current_balances (
    id character varying,
    data jsonb,
    balance numeric
);
ALTER TABLE ONLY current_balances REPLICA IDENTITY NOTHING;
CREATE TABLE lines (
    id bigint NOT NULL,
    transaction_id character varying NOT NULL,
    account_id character varying NOT NULL,
    delta bigint NOT NULL
);
CREATE VIEW invalid_transactions AS
 SELECT lines.transaction_id,
    sum(lines.delta) AS sum
   FROM lines
  GROUP BY lines.transaction_id
 HAVING (sum(lines.delta) > (0)::numeric);
CREATE SEQUENCE lines_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE lines_id_seq OWNED BY lines.id;
CREATE TABLE schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);
CREATE TABLE transactions (
    id character varying NOT NULL,
    "timestamp" timestamp without time zone NOT NULL,
    data jsonb DEFAULT '{}'::jsonb NOT NULL
);
ALTER TABLE ONLY lines ALTER COLUMN id SET DEFAULT nextval('lines_id_seq'::regclass);
ALTER TABLE ONLY accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);
ALTER TABLE ONLY lines
    ADD CONSTRAINT lines_pkey PRIMARY KEY (id);
ALTER TABLE ONLY schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);
ALTER TABLE ONLY transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);
CREATE INDEX accounts_data_idx ON accounts USING gin (data jsonb_path_ops);
CREATE INDEX lines_account_id_idx ON lines USING btree (account_id);
CREATE INDEX lines_transaction_id_idx ON lines USING btree (transaction_id);
CREATE INDEX timestamp_idx ON transactions USING brin ("timestamp");
CREATE INDEX transactions_data_idx ON transactions USING gin (data jsonb_path_ops);
CREATE RULE "_RETURN" AS
    ON SELECT TO current_balances DO INSTEAD  SELECT accounts.id,
    accounts.data,
    COALESCE(sum(lines.delta), (0)::numeric) AS balance
   FROM (accounts
     LEFT JOIN lines ON (((accounts.id)::text = (lines.account_id)::text)))
  GROUP BY accounts.id;
ALTER TABLE ONLY lines
    ADD CONSTRAINT lines_account_id_fkey FOREIGN KEY (account_id) REFERENCES accounts(id);
ALTER TABLE ONLY lines
    ADD CONSTRAINT lines_txn_fkey FOREIGN KEY (transaction_id) REFERENCES transactions(id);
