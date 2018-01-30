package models

import (
	"database/sql"
	"encoding/json"
	"log"

	ledgerError "github.com/RealImage/QLedger/errors"
)

// Account represents the ledger account with information such as ID, balance and JSON data
type Account struct {
	ID      string                 `json:"id"`
	Balance int                    `json:"balance"`
	Data    map[string]interface{} `json:"data"`
}

// AccountDB provides all functions related to ledger account
type AccountDB struct {
	db *sql.DB
}

// NewAccountDB provides instance of `AccountDB`
func NewAccountDB(db *sql.DB) AccountDB {
	return AccountDB{db: db}
}

// GetByID returns an acccount with the given ID
func (a *AccountDB) GetByID(id string) (*Account, ledgerError.ApplicationError) {
	account := &Account{ID: id}

	err := a.db.QueryRow("SELECT balance FROM current_balances WHERE id=$1", &id).Scan(&account.Balance)
	switch {
	case err == sql.ErrNoRows:
		account.Balance = 0
	case err != nil:
		return nil, DBError(err)
	}

	return account, nil
}

// GetBalanceOnTime returns the balance of the account as on given time.
func (a *AccountDB) GetBalanceOnTime(accountID string, timestamp string) (*Account, ledgerError.ApplicationError) {
	account := &Account{ID: accountID}

	query := `SELECT COALESCE(sum(lines.delta), 0) AS balance
				FROM lines LEFT OUTER JOIN transactions ON (transactions.id = lines.transaction_id)
				WHERE lines.account_id = $1 AND transactions.timestamp <= $2
				GROUP  BY lines.account_id`
	err := a.db.QueryRow(query, &accountID, &timestamp).Scan(&account.Balance)
	switch {
	case err == sql.ErrNoRows:
		account.Balance = 0
	case err != nil:
		return nil, DBError(err)
	}

	return account, nil
}

// IsExists says whether an account exists or not
func (a *AccountDB) IsExists(id string) (bool, ledgerError.ApplicationError) {
	var exists bool
	err := a.db.QueryRow("SELECT EXISTS (SELECT id FROM accounts WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		log.Println("Error executing account exists query:", err)
		return false, DBError(err)
	}
	return exists, nil
}

// CreateAccount creates a new account in the ledger
func (a *AccountDB) CreateAccount(account *Account) ledgerError.ApplicationError {
	data, err := json.Marshal(account.Data)
	if err != nil {
		return JSONError(err)
	}

	accountData := "{}"
	if account.Data != nil && data != nil {
		accountData = string(data)
	}

	q := "INSERT INTO accounts (id, data)  VALUES ($1, $2)"
	_, err = a.db.Exec(q, account.ID, accountData)
	if err != nil {
		return DBError(err)
	}

	return nil
}

// UpdateAccount updates the account with new data
func (a *AccountDB) UpdateAccount(account *Account) ledgerError.ApplicationError {
	data, err := json.Marshal(account.Data)
	if err != nil {
		return JSONError(err)
	}
	accountData := "{}"
	if account.Data != nil && data != nil {
		accountData = string(data)
	}

	q := "UPDATE accounts SET data = $1 WHERE id = $2"
	_, err = a.db.Exec(q, accountData, account.ID)
	if err != nil {
		return DBError(err)
	}

	return nil
}
