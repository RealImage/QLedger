ALTER TABLE ONLY transaction_tags
    ADD CONSTRAINT transaction_tags_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES transactions(id);