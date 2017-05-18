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
	var balance int
	err := account.DB.QueryRow("SELECT balance FROM current_balances where account_id=$1", &id).Scan(&balance)
	
	switch {
		case err == sql.ErrNoRows:
        at = &Account{Id: id, Balance: 0}
		case err != nil:
		    log.Fatal(err)
		default:
				at = &Account{Id: id, Balance: balance}
	}
	return
}
