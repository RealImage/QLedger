package models

import (
	"database/sql"
	"log"
)

type Account struct {
	Id      string `json:"id"`
	Balance int    `json:"balance"`
}

type AccountDB struct {
	DB *sql.DB `json:"-"`
}

func (adb *AccountDB) GetByID(id string) *Account {
	account := &Account{Id: id}

	err := adb.DB.QueryRow("SELECT balance FROM current_balances WHERE account_id=$1", &id).Scan(&account.Balance)
	switch {
	case err == sql.ErrNoRows:
		account.Balance = 0
	case err != nil:
		log.Fatal(err)
	}

	return account
}
