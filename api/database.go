package api

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

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
	}

	database.database = db;
	// defer database.database.Close()
	return database
}


// Create an application
func (db *Database) CreateApplication (appName string) *Application {
	var application *Application = &Application{
		ID: uuid.New(), 
		Name: appName,
	} 

	var SQL string = "INSERT INTO Applications (ID, Name) VALUES ($1, $2)"

	// Omit the result return
	if _, err := db.database.Exec(SQL, application.ID, application.Name); err != nil {
		panic(err)
	}

	return application
}

// Create a user
func (db *Database) CreateUser (applicationID uuid.UUID, username string, password string) *User {
	var user *User = &User{
		ID: uuid.New(),
		Username: username,
		Password: password,
		ApplicationID: applicationID,
	}

	var SQL string = "INSERT INTO Users (ID, ApplicationID, Username, Password) VALUES ($1, $2, $3, $4)"

	// Omit the result return
	if _, err := db.database.Exec(SQL, user.ID, user.ApplicationID, user.Username, user.Password); err != nil {
		panic(err)
	}

	return user
}

// Verify a username and password
func (db *Database) VerifyUser (applicationID uuid.UUID, username string, password string) (*User, error) {
	var SQL string = "SELECT ApplicationID, ID, password FROM users WHERE ApplicationID = $1 AND Username = $2";

	rows, err := db.database.Query(SQL, applicationID, username)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			ApplicationID uuid.UUID
			ID uuid.UUID
			Password string
		)
		if err := rows.Scan(&ApplicationID, &ID, &Password); err != nil {
			return nil, err
		}

		valid, err := CompareString(password, Password)
			
		if err != nil {
			return nil, err
		}

		if valid {
			user := &User {
				ID: ID,
				ApplicationID: ApplicationID,
				Username: username,
				Password: Password,
			}

			return user, nil
		}
	}
	return nil, errors.New("User was not verified") 
}
