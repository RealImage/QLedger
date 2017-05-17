package context

import (
	"database/sql"

	"github.com/RealImage/QLedger/database"
)

type AppContext struct {
	DB *sql.DB
}

func (appContext *AppContext) Initialize() {
	appContext.DB = database.Dial()
	// Other contexts to initialize...
}

func (appContext *AppContext) Cleanup() {
	appContext.DB.Close()
	// Other contexts to cleanup...
}
