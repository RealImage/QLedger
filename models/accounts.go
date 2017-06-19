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

func (adb *AccountDB) GetByID(id string) *Account {
	account := &Account{Id: id}

	err := adb.DB.QueryRow("SELECT balance FROM current_balances WHERE id=$1", &id).Scan(&account.Balance)
	switch {
	case err == sql.ErrNoRows:
		account.Balance = 0
	case err != nil:
		log.Fatal(err)
	}

	return account
}

func (adb *AccountDB) IsExists(id string) bool {
	var exists bool
	err := adb.DB.QueryRow("SELECT EXISTS (SELECT * FROM accounts WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		log.Println("Error executing account exists query:", err)
	}
	return exists
}

func (adb *AccountDB) CreateAccount(account *Account) ledgerError.ApplicationError {
	data, jerr := json.Marshal(account.Data)
	if jerr != nil {
		return JSONError(jerr)
	}
	accountData := "{}"
	if account.Data != nil && data != nil {
		accountData = string(data)
	}

	sql := "INSERT INTO accounts (id, data)  VALUES ($1, $2)"
	_, derr := adb.DB.Exec(sql, account.Id, accountData)
	if derr != nil {
		return DBError(derr)
	}
	return nil
}

func (adb *AccountDB) UpdateAccount(account *Account) ledgerError.ApplicationError {
	data, jerr := json.Marshal(account.Data)
	if jerr != nil {
		return JSONError(jerr)
	}
	accountData := "{}"
	if account.Data != nil && data != nil {
		accountData = string(data)
	}

	sql := "UPDATE accounts SET data = $1 WHERE id = $2"
	_, derr := adb.DB.Exec(sql, accountData, account.Id)
	if derr != nil {
		return DBError(derr)
	}
	return nil
}
