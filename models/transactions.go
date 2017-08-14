package models

import (
	"database/sql"
	"encoding/json"
	"log"
	"reflect"
	"time"

	ledgerError "github.com/RealImage/QLedger/errors"
	"github.com/pkg/errors"
)

type Transaction struct {
	ID        string                 `json:"id"`
	Data      map[string]interface{} `json:"data"`
	Timestamp string                 `json:"timestamp"`
	Lines     []*TransactionLine     `json:"lines"`
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

func (t *TransactionDB) IsExists(id string) bool {
	var exists bool
	err := t.DB.QueryRow("SELECT EXISTS (SELECT * FROM transactions WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		log.Println("Error executing transaction exists query:", err)
	}
	return exists
}

func (t *TransactionDB) IsConflict(transaction *Transaction) bool {
	// Read existing lines
	rows, err := t.DB.Query("SELECT account_id, delta FROM lines WHERE transaction_id=$1", transaction.ID)
	if err != nil {
		log.Println("Error executing transaction lines query:", err)
		return false
	}
	defer rows.Close()
	var existingLines []*TransactionLine
	for rows.Next() {
		line := &TransactionLine{}
		if err := rows.Scan(&line.AccountID, &line.Delta); err != nil {
			log.Println("Error scanning transaction lines:", err)
			return false
		}
		existingLines = append(existingLines, line)
	}
	if err := rows.Err(); err != nil {
		log.Println("Error iterating transaction lines rows:", err)
		return false
	}

	// Compare new and existing transaction lines
	return !reflect.DeepEqual(transaction.Lines, existingLines)
}

func (t *TransactionDB) Transact(txn *Transaction) bool {
	// Start the transaction
	var err error
	tx, err := t.DB.Begin()
	if err != nil {
		log.Println("Error beginning transaction:", err)
		return false
	}

	// Rollback transaction on any failures
	handleTransactionError := func(tx *sql.Tx, err error) bool {
		log.Println(err)
		log.Println("Rolling back the transaction:", txn.ID)
		err = tx.Rollback()
		if err != nil {
			log.Println("Error rolling back transaction:", err)
		}
		return false
	}

	// Accounts do not need to be predefined
	// they are called into existence when they are first used.
	for _, line := range txn.Lines {
		_, err = tx.Exec("INSERT INTO accounts (id) VALUES ($1) ON CONFLICT (id) DO NOTHING", line.AccountID)
		if err != nil {
			return handleTransactionError(tx, errors.Wrap(err, "insert account failed"))
		}
	}

	// Add transaction
	data, err := json.Marshal(txn.Data)
	if err != nil {
		return handleTransactionError(tx, errors.Wrap(err, "transaction data parse error"))
	}
	transactionData := "{}"
	if txn.Data != nil && data != nil {
		transactionData = string(data)
	}

	_, err = tx.Exec("INSERT INTO transactions (id, timestamp, data) VALUES ($1, $2, $3)", txn.ID, time.Now().UTC(), transactionData)
	if err != nil {
		return handleTransactionError(tx, errors.Wrap(err, "insert transaction failed"))
	}

	// Add transaction lines
	for _, line := range txn.Lines {
		_, err = tx.Exec("INSERT INTO lines (transaction_id, account_id, delta) VALUES ($1, $2, $3)", txn.ID, line.AccountID, line.Delta)
		if err != nil {
			return handleTransactionError(tx, errors.Wrap(err, "insert lines failed"))
		}
	}

	// Commit the entire transaction
	err = tx.Commit()
	if err != nil {
		return handleTransactionError(tx, errors.Wrap(err, "commit transaction failed"))
	}

	return true
}

func (t *TransactionDB) UpdateTransaction(txn *Transaction) ledgerError.ApplicationError {
	data, err := json.Marshal(txn.Data)
	if err != nil {
		return JSONError(err)
	}
	tData := "{}"
	if txn.Data != nil && data != nil {
		tData = string(data)
	}

	q := "UPDATE transactions SET data = $1 WHERE id = $2"
	_, err = t.DB.Exec(q, tData, txn.ID)
	if err != nil {
		return DBError(err)
	}
	return nil
}
