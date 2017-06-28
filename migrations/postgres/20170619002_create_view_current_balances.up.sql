CREATE VIEW current_balances AS
SELECT accounts.id, accounts.data,
    COALESCE(SUM(lines.delta), 0) AS balance
  FROM accounts LEFT OUTER JOIN lines
  ON (accounts.id = lines.account_id)
  GROUP BY accounts.id;