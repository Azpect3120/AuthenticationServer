package database

import (
	"database/sql"
	"log"

	"github.com/Azpect3120/AuthenticationServer/core/model"
	_ "github.com/lib/pq"
)

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
