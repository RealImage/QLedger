ALTER TABLE ONLY accounts ADD COLUMN data jsonb DEFAULT '{}'::jsonb NOT NULL;
