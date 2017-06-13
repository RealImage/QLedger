ALTER TABLE ONLY account_tags
    ADD CONSTRAINT account_tags_account_id_fkey FOREIGN KEY (account_id) REFERENCES accounts(id);