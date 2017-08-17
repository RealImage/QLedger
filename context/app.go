package context

import (
	"database/sql"
)

// AppContext provides the context to the app components such as controllers, jobs, etc.,
type AppContext struct {
	DB *sql.DB
}
