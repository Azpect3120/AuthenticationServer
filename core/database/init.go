package database

import (
	"database/sql"
	"log"

	"github.com/Azpect3120/AuthenticationServer/core/model"
	_ "github.com/lib/pq"
)

// Creates a database obejct and attempts to connect to the
// database using the provided connection string. If the
// connection fails, the program will exit with a fatal.
func NewDatabase(connectionString string) *model.Database {
	db := &model.Database{
		ConnectionString: connectionString,
	}

	conn, err := sql.Open("postgres", db.ConnectionString)
	if err != nil {
		log.Fatalln("Error connecting to database:", err)
	}

	db.Conn = conn

	return db
}
