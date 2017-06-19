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

func (tdb *TransactionDB) IsExists(id string) bool {
	var exists bool
	err := tdb.DB.QueryRow("SELECT EXISTS (SELECT * FROM transactions WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		log.Println("Error executing transaction exists query:", err)
	}
	return exists
}

func (tdb *TransactionDB) IsConflict(transaction *Transaction) bool {
	// Read existing lines
	rows, err := tdb.DB.Query("SELECT account_id, delta FROM lines WHERE transaction_id=$1", transaction.ID)
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

func (tdb *TransactionDB) Transact(t *Transaction) bool {
	// Start the transaction
	var err error
	txn, err := tdb.DB.Begin()
	if err != nil {
		log.Println("Error beginning transaction:", err)
		return false
	}

	// Rollback transaction on any failures
	handleTransactionError := func(txn *sql.Tx, err error) bool {
		log.Println(err)
		log.Println("Rolling back the transaction:", t.ID)
		err = txn.Rollback()
		if err != nil {
			log.Println("Error rolling back transaction:", err)
		}
		return false
	}

	// Accounts do not need to be predefined
	// they are called into existence when they are first used.
	for _, line := range t.Lines {
		_, err = txn.Exec("INSERT INTO accounts (id) VALUES ($1) ON CONFLICT (id) DO NOTHING", line.AccountID)
		if err != nil {
			return handleTransactionError(txn, errors.Wrap(err, "insert account failed"))
		}
	}

	// Add transaction
	data, err := json.Marshal(t.Data)
	if err != nil {
		return handleTransactionError(txn, errors.Wrap(err, "transaction data parse error"))
	}
	transactionData := "{}"
	if t.Data != nil && data != nil {
		transactionData = string(data)
	}
	_, err = txn.Exec("INSERT INTO transactions (id, timestamp, data) VALUES ($1, $2, $3)", t.ID, time.Now().UTC(), transactionData)
	if err != nil {
		return handleTransactionError(txn, errors.Wrap(err, "insert transaction failed"))
	}

	// Add transaction lines
	for _, line := range t.Lines {
		_, err = txn.Exec("INSERT INTO lines (transaction_id, account_id, delta) VALUES ($1, $2, $3)", t.ID, line.AccountID, line.Delta)
		if err != nil {
			return handleTransactionError(txn, errors.Wrap(err, "insert lines failed"))
		}
	}

	// Commit the entire transaction
	err = txn.Commit()
	if err != nil {
		return handleTransactionError(txn, errors.Wrap(err, "commit transaction failed"))
	}

	return true
}

func (tdb *TransactionDB) UpdateTransaction(t *Transaction) ledgerError.ApplicationError {
	data, jerr := json.Marshal(t.Data)
	if jerr != nil {
		return JSONError(jerr)
	}
	tData := "{}"
	if t.Data != nil && data != nil {
		tData = string(data)
	}

	sql := "UPDATE transactions SET data = $1 WHERE id = $2"
	_, derr := tdb.DB.Exec(sql, tData, t.ID)
	if derr != nil {
		return DBError(derr)
	}
	return nil
}
