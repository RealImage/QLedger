package models

import (
	"database/sql"
	"log"
)

type Account struct {
	DB      *sql.DB `json:"-"`
	Id      string  `json:"id"`
	Balance int     `json:"balance"`
}

func (account *Account) GetByID(id string) (at *Account) {
	at = &Account{Id: id}

	err := account.DB.QueryRow("SELECT balance FROM current_balances where account_id=$1", &id).Scan(&at.Balance)

	switch {
	case err == sql.ErrNoRows:
		at.Balance = 0
	case err != nil:
		log.Fatal(err)
	}

	return
}
