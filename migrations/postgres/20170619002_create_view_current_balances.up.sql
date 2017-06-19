CREATE VIEW current_balances AS
SELECT accounts.id, accounts.data,
    SUM(lines.delta) AS balance
  FROM accounts, lines
  WHERE accounts.id = lines.account_id
  GROUP BY accounts.id;