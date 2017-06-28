CREATE INDEX accounts_data_idx ON accounts USING GIN (data jsonb_path_ops);
