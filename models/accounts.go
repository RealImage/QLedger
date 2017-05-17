package models

import (
	"log"

	"github.com/RealImage/QLedger/database"
)

type Account struct {
	Id      string `json:"id"`
	Balance int    `json:"balance"`
}

func GetAccounts() (accounts []*Account) {
	rows, err := database.Conn.Query("SELECT * FROM current_balances")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		account := &Account{}
		err := rows.Scan(&account.Id, &account.Balance)
		if err != nil {
			log.Fatal(err)
		}
		accounts = append(accounts, account)
	}
	return
}
