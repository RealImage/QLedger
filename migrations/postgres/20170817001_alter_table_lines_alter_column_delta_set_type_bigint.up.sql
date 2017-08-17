BEGIN;

DROP VIEW IF EXISTS current_balances;
DROP VIEW IF EXISTS invalid_transactions;

ALTER TABLE lines ALTER COLUMN delta SET DATA TYPE bigint;

CREATE VIEW current_balances AS
  SELECT accounts.id, accounts.data,
    COALESCE(SUM(lines.delta), 0) AS balance
  FROM accounts LEFT OUTER JOIN lines
  ON (accounts.id = lines.account_id)
  GROUP BY accounts.id;
CREATE VIEW invalid_transactions AS
  SELECT lines.transaction_id,
    sum(lines.delta) AS sum
   FROM lines
  GROUP BY lines.transaction_id
 HAVING (sum(lines.delta) > 0);

COMMIT;
