package models

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AccountsSuite struct {
	suite.Suite
	db *sql.DB
}

func (as *AccountsSuite) SetupTest() {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	assert.NotEmpty(as.T(), databaseURL)
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Panic("Unable to connect to Database:", err)
	} else {
		log.Println("Successfully established connection to database.")
		as.db = db
	}
}

func (as *AccountsSuite) TestAccountsInfoAPI() {
	t := as.T()

	accountsDB := AccountDB{DB: as.db}
	account, err := accountsDB.GetByID("100")
	assert.Equal(t, err, nil, "Error while getting acccount")
	assert.Equal(t, account.Id, "100", "Invalid account ID")
	assert.Equal(t, account.Balance, 0, "Invalid account balance")
}

func TestAccountsSuite(t *testing.T) {
	suite.Run(t, new(AccountsSuite))
}
