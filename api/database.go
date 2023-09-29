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

// Get a user from the database using its ID
func (db *Database) GetUser (applicationID string, userID string) (*User, error) {
	var SQL = "SELECT * FROM users WHERE ApplicationID = $1 AND ID = $2"

	appUUID, err := uuid.Parse(applicationID) 

	if err != nil {
		return nil, err
	}

	userUUID, err := uuid.Parse(userID)

	if err != nil {
		return nil, err
	}

	rows, err := db.database.Query(SQL, appUUID, userUUID)
	
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			ID uuid.UUID
			ApplicationID uuid.UUID
			Username string
			Password string
		)

		if err := rows.Scan(&ID, &ApplicationID, &Username, &Password); err != nil {
			return nil, err 
		}

		user := &User{
			ApplicationID: ApplicationID,
			ID: ID,
			Username: Username,
			Password: Password,
		}

		return user, nil
	}
	return nil, errors.New("User was not found") 
}


// Returns an array of users found in an application
func (db *Database) GetUsers (applicationID string) ([]*User, error) {
	var SQL string = "SELECT * FROM users WHERE ApplicationID = $1"

	appUUID, err := uuid.Parse(applicationID)

	if err != nil {
		return nil, err
	}

	rows, err := db.database.Query(SQL, appUUID)

	if err != nil {
		return nil, err
	}

	users := []*User{}

	for rows.Next() {
		var (
			ID uuid.UUID
			ApplicationID uuid.UUID
			Username string
			Password string
		)

		rows.Scan(&ID, &ApplicationID, &Username, &Password) 

		user := &User{
			ApplicationID: ApplicationID,
			ID: ID,
			Username: Username,
			Password: Password,
		}

		users = append(users, user)
	}

	return users, nil
}

// Updates a users username
func (db *Database) SetUsername (applicationID uuid.UUID, userID uuid.UUID, newUsername string) (*User, error) {
	var SQL string = "UPDATE Users SET Username = $1 WHERE ApplicationID = $2 AND ID = $3"

	_, err := db.database.Exec(SQL, newUsername, applicationID, userID)

	if err != nil {
		return nil, err
	}

	SQL = "SELECT * FROM Users WHERE ApplicationID = $1 AND ID = $2"

	rows, err := db.database.Query(SQL, applicationID, userID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			ID uuid.UUID
			ApplicationID uuid.UUID
			Username string
			Password string
		)

		rows.Scan(&ID, &ApplicationID, &Username, &Password)

		user := &User{
			ID: ID,
			ApplicationID: ApplicationID,
			Username: Username,
			Password: Password,
		}

		return user, nil
	}
	return nil, errors.New("The users username could not be changed.")
}





// Updates a users password
