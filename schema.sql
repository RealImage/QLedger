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
CREATE TABLE account_tags (
    account_id character varying NOT NULL,
    key character varying NOT NULL,
    value character varying
);
CREATE TABLE accounts (
    id character varying NOT NULL,
    data jsonb DEFAULT '{}'::jsonb NOT NULL
);
CREATE TABLE lines (
    id bigint NOT NULL,
    transaction_id character varying NOT NULL,
    account_id character varying NOT NULL,
    delta integer NOT NULL
);
CREATE VIEW current_balances AS
 SELECT lines.account_id,
    sum(lines.delta) AS balance
   FROM lines
  GROUP BY lines.account_id;
CREATE VIEW invalid_transactions AS
 SELECT lines.transaction_id,
    sum(lines.delta) AS sum
   FROM lines
  GROUP BY lines.transaction_id
 HAVING (sum(lines.delta) > 0);
CREATE SEQUENCE lines_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE lines_id_seq OWNED BY lines.id;
CREATE TABLE transaction_tags (
    transaction_id character varying NOT NULL,
    key character varying NOT NULL,
    value character varying
);
CREATE TABLE transactions (
    id character varying NOT NULL,
    "timestamp" timestamp without time zone NOT NULL
);
ALTER TABLE ONLY lines ALTER COLUMN id SET DEFAULT nextval('lines_id_seq'::regclass);
ALTER TABLE ONLY accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (id);
ALTER TABLE ONLY lines
    ADD CONSTRAINT lines_pkey PRIMARY KEY (id);
ALTER TABLE ONLY transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);
CREATE UNIQUE INDEX account_tags_lookup_idx ON account_tags USING btree (value, key, account_id);
CREATE INDEX timestamp_idx ON transactions USING brin ("timestamp");
CREATE UNIQUE INDEX transaction_tags_lookup_idx ON transaction_tags USING btree (value, key, transaction_id);
ALTER TABLE ONLY account_tags
    ADD CONSTRAINT account_tags_account_id_fkey FOREIGN KEY (account_id) REFERENCES accounts(id);
ALTER TABLE ONLY lines
    ADD CONSTRAINT lines_account_id_fkey FOREIGN KEY (account_id) REFERENCES accounts(id);
ALTER TABLE ONLY lines
    ADD CONSTRAINT lines_txn_fkey FOREIGN KEY (transaction_id) REFERENCES transactions(id);
ALTER TABLE ONLY transaction_tags
    ADD CONSTRAINT transaction_tags_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES transactions(id);
