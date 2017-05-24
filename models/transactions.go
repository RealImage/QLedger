package models

import (
	"database/sql"
	"log"
	"time"
)

type Transaction struct {
	ID    string `json:"id"`
	Lines []*TransactionLine
}

type TransactionLine struct {
	AccountID string `json:"account"`
	Delta     int    `json:"delta"`
}

func (t *Transaction) IsValid() bool {
	sum := 0
	for _, line := range t.Lines {
		sum += line.Delta
	}
	return sum == 0
}

type TransactionDB struct {
	DB *sql.DB
}

func (tdb *TransactionDB) IsExists(id string) bool {
	var exists bool
	err := tdb.DB.QueryRow("SELECT exists (SELECT * FROM transactions WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		log.Println("Error executing transaction exists query:", err)
	}
	return exists
}

func (tdb *TransactionDB) DoTransaction(t *Transaction) bool {
	// Start the transaction
	var err error
	txn, err := tdb.DB.Begin()
	if err != nil {
		log.Println("Error beginning transaction:", err)
		return false
	}

	// Rollback transaction on any failures
	defer func() {
		if r := recover(); r != nil {
			log.Fatal("Recovering while transaction", r)

			log.Println("Rolling back the transaction:", t.ID)
			err = txn.Rollback()
			if err != nil {
				log.Println("Error rolling back transaction:", err)
			}
		} else if err != nil {
			log.Println("Rolling back the transaction:", t.ID)
			err = txn.Rollback()
			if err != nil {
				log.Println("Error rolling back transaction:", err)
			}
		}
	}()

	// Accounts do not need to be predefined
	// they are called into existence when they are first used.
	accStmt, err := txn.Prepare(`INSERT INTO accounts (id) VALUES ($1) ON CONFLICT (id) DO NOTHING`)
	if err != nil {
		log.Println("Error preparing statement of accounts:", err)
		return false
	}
	for _, line := range t.Lines {
		_, err = accStmt.Exec(line.AccountID)
		if err != nil {
			log.Println("Error executing prepared statement of accounts:", err)
			return false
		}
	}
	err = accStmt.Close()
	if err != nil {
		log.Println("Error while closing prepared statement of accounts:", err)
		return false
	}

	// Add transaction
	_, err = txn.Exec("INSERT INTO transactions (id, timestamp) VALUES ($1, $2)", t.ID, time.Now().UTC())
	if err != nil {
		log.Println("Error inserting transaction:", err)
		return false
	}

	// Add transaction lines
	linesStmt, err := txn.Prepare(`INSERT INTO lines (transaction_id, account_id, delta) VALUES ($1, $2, $3)`)
	if err != nil {
		log.Println("Error preparing statement of lines:", err)
		return false
	}
	for _, line := range t.Lines {
		_, err = linesStmt.Exec(t.ID, line.AccountID, line.Delta)
		if err != nil {
			log.Println("Error executing prepared statement of lines:", err)
			return false
		}
	}
	err = linesStmt.Close()
	if err != nil {
		log.Println("Error while closing prepared statement of lines:", err)
		return false
	}

	// Commit the entire transaction
	err = txn.Commit()
	if err != nil {
		log.Println("Error while committing the transaction:", err)
		return false
	}

	return true
}
