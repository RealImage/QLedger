CREATE INDEX transactions_data_idx ON transactions USING GIN (data jsonb_path_ops);
