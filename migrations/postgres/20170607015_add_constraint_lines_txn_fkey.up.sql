ALTER TABLE ONLY lines
    ADD CONSTRAINT lines_txn_fkey FOREIGN KEY (transaction_id) REFERENCES transactions(id);