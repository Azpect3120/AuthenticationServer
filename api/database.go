package api

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	db_host		= "bubble.db.elephantsql.com"
	db_port		= 5432
	db_user		= "cihrecbo"
	db_password	= "u6hQwF7ceKcHeuu6I4uM4ewaE8MCJjqs"
	db_name		= "cihrecbo"	
)

type Database struct {
	connectionString string
	database *sql.DB
}

func CreateDatabase () *Database {
	database := &Database {
		connectionString: fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_password, db_name),
	}

	db, err := sql.Open("postgres", database.connectionString)

	if err != nil {
		panic(err)
		return nil
	}

	database.database = db;
	defer database.database.Close()
	return database
}


// Create an application
func (db *Database) CreateApplication (appName string) *Application {
	return nil
}

// Create a user
func (db *Database) CreateUser (user User) *User {
	return nil
}
