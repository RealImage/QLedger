ALTER TABLE ONLY transactions ADD COLUMN data jsonb DEFAULT '{}'::jsonb NOT NULL;