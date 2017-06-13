CREATE VIEW current_balances AS
 SELECT lines.account_id,
    sum(lines.delta) AS balance
   FROM lines
  GROUP BY lines.account_id;