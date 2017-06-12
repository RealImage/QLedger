CREATE VIEW invalid_transactions AS
 SELECT lines.transaction_id,
    sum(lines.delta) AS sum
   FROM lines
  GROUP BY lines.transaction_id
 HAVING (sum(lines.delta) > 0);