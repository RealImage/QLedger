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

func (account *Account) GetByID() (accounts []*Account) {
	//TODO: Get a single account by ID
	rows, err := account.DB.Query("SELECT * FROM current_balances")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&account.Id, &account.Balance)
		if err != nil {
			log.Fatal(err)
		}
		accounts = append(accounts, account)
	}
	return
}
