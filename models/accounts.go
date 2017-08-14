package models

import (
	"database/sql"
	"encoding/json"
	"log"

	ledgerError "github.com/RealImage/QLedger/errors"
)

type Account struct {
	Id      string                 `json:"id"`
	Balance int                    `json:"balance"`
	Data    map[string]interface{} `json:"data"`
}

type AccountDB struct {
	DB *sql.DB `json:"-"`
}

func (a *AccountDB) GetByID(id string) (*Account, ledgerError.ApplicationError) {
	account := &Account{Id: id}

	err := a.DB.QueryRow("SELECT balance FROM current_balances WHERE id=$1", &id).Scan(&account.Balance)
	switch {
	case err == sql.ErrNoRows:
		account.Balance = 0
	case err != nil:
		return nil, DBError(err)
	}

	return account, nil
}

func (a *AccountDB) IsExists(id string) bool {
	var exists bool
	err := a.DB.QueryRow("SELECT EXISTS (SELECT * FROM accounts WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		log.Println("Error executing account exists query:", err)
	}
	return exists
}

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
	_, err = a.DB.Exec(q, account.Id, accountData)
	if err != nil {
		return DBError(err)
	}

	return nil
}

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
	_, err = a.DB.Exec(q, accountData, account.Id)
	if err != nil {
		return DBError(err)
	}

	return nil
}
