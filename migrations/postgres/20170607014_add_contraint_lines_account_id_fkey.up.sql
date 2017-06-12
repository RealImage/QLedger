ALTER TABLE ONLY lines
    ADD CONSTRAINT lines_account_id_fkey FOREIGN KEY (account_id) REFERENCES accounts(id);