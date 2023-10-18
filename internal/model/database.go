package model

import (
	"database/sql"
	"fmt"
	"os"
	"log"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type Database struct {
	connectionString string
	database *sql.DB
}

func CreateDatabase () *Database {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	db_host := os.Getenv("db_host")
	db_port	:= os.Getenv("db_port")
	db_user	:= os.Getenv("db_user")
	db_password	:= os.Getenv("db_password")
	db_name := os.Getenv("db_name")

	database := &Database {
		connectionString: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_password, db_name),
	}

	db, err := sql.Open("postgres", database.connectionString)

	if err != nil {
		log.Fatal(err)
	}

	database.database = db;
	return database
}


// Create an application
func (db *Database) CreateApplication (ch chan *AppResult, appName string) {
	var application *Application = &Application{
		ID: uuid.New(), 
		Name: appName,
	} 

	var SQL string = "INSERT INTO applications (ID, name) VALUES ($1, $2)"

	// Omit the result return
	if _, err := db.database.Exec(SQL, application.ID, application.Name); err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &AppResult{nil, &Error{ Message: err.Error(), Status: 500 }}
		return
	}
		
	// return application, nil
	ch <- &AppResult{application, nil}
}

// Create a user
func (db *Database) CreateUser (ch chan *UserResult, applicationID uuid.UUID, username string, password string) {
	var applicationCount int

	var countSQL string = "SELECT COUNT(*) FROM applications WHERE ID = $1"
	if err := db.database.QueryRow(countSQL, applicationID).Scan(&applicationCount); err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UserResult{nil, &Error{ Message: err.Error(), Status: 500 }}
	}


	if applicationCount == 0 {
		// return nil, &Error{ Message: "Invalid applicationID.", Status: 401 }
		ch <- &UserResult{nil, &Error{ Message: "Invalid applicationID.", Status: 401 }}
	}

	var user *User = &User{
		ID: uuid.New(),
		Username: username,
		Password: password,
		ApplicationID: applicationID,
	}

	var SQL string = "INSERT INTO users (ID, applicationID, username, password) VALUES ($1, $2, $3, $4)"

	// Omit the result return
	if _, err := db.database.Exec(SQL, user.ID, user.ApplicationID, user.Username, user.Password); err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UserResult{nil, &Error{ Message: err.Error(), Status: 500 }}
	}

	// return user, nil
	ch <- &UserResult{user, nil}
}

// Verify a username and password
func (db *Database) VerifyUser (ch chan *UserResult, applicationID uuid.UUID, username string, password string) {
	var SQL string = "SELECT applicationID, ID, password FROM users WHERE applicationID = $1 AND username = $2";

	rows, err := db.database.Query(SQL, applicationID, username)

	if err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UserResult{nil, &Error{ Message: err.Error(), Status: 500 }}
	}

	defer rows.Close()

	for rows.Next() {
		var (
			ApplicationID uuid.UUID
			ID uuid.UUID
			Password string
		)
		if err := rows.Scan(&ApplicationID, &ID, &Password); err != nil {
			// return nil, &Error{ Message: err.Error(), Status: 500 }
			ch <- &UserResult{nil, &Error{ Message: err.Error(), Status: 500 }}
		}

		valid, err := CompareString(password, Password)
			
		if err != nil {
			// return nil, &Error{ Message: err.Message, Status: err.Status }
			ch <- &UserResult{nil, &Error{ Message: err.Message, Status: err.Status }}
		}

		if valid {
			user := &User {
				ID: ID,
				ApplicationID: ApplicationID,
				Username: username,
				Password: Password,
			}

			// return user, nil
			ch <- &UserResult{user, nil}
		}
	}
	// return nil, &Error{ Message: "User was not verified", Status: 401 }
	ch <- &UserResult{nil, &Error{ Message: "User was not verified", Status: 401 }}
}

// Get a user from the database using its ID
func (db *Database) GetUser (ch chan *UserResult, applicationID string, userID string) {
	var SQL string = "SELECT * FROM users WHERE applicationID = $1 AND ID = $2"

	rows, err := db.database.Query(SQL, applicationID, userID)

	if err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UserResult{nil, &Error{ Message: err.Error(), Status: 500 }}
	}

	defer rows.Close()


	for rows.Next() {
		var (
			ID uuid.UUID
			ApplicationID uuid.UUID
			Username string
			Password string
		)

		if err := rows.Scan(&ID, &ApplicationID, &Username, &Password); err != nil {
			// return nil, &Error{ Message: err.Error(), Status: 500 }
			ch <- &UserResult{nil, &Error{ Message: err.Error(), Status: 500 }}
		}

		user := &User{
			ID: ID,
			ApplicationID: ApplicationID,
			Username: Username,
			Password: Password,
		}

		// return user, nil
		ch <- &UserResult{user, nil}
	}

	if err := rows.Err(); err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500}
		ch <- &UserResult{nil, &Error{ Message: err.Error(), Status: 500 }}
	}

	// return nil, &Error{ Message: "User was not found.", Status: 404 }
	ch <- &UserResult{nil, &Error{ Message: "User was not found.", Status: 404 }}
}


// Returns an array of users found in an application
func (db *Database) GetUsers (ch chan *UsersResult, applicationID string) {
	appUUID, err := uuid.Parse(applicationID)
	if err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UsersResult{nil, &Error{ Message: err.Error(), Status: 500 }}
	}
	
	var appCount int

	var countSQL string = "SELECT COUNT(*) FROM applications WHERE ID = $1"
	if err := db.database.QueryRow(countSQL, appUUID).Scan(&appCount); err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UsersResult{nil, &Error{ Message: err.Error(), Status: 500 }}
	}

	if appCount == 0 {
		// return nil, &Error{ Message: "Application with the provided ID does not exist.", Status: 404 }
		ch <- &UsersResult{nil, &Error{ Message: "Application with the provided ID does not exist.", Status: 404 }}
	}

	var SQL string = "SELECT * FROM users WHERE applicationID = $1"
	rows, err := db.database.Query(SQL, appUUID)
	if err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UsersResult{nil, &Error{ Message: err.Error(), Status: 500 }}
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

	// return users, nil
	ch <- &UsersResult{users, nil}
}

// Updates a users username
func (db *Database) SetUsername (applicationID uuid.UUID, userID uuid.UUID, newUsername string) (*User, *Error) {
	var userCount int

	var countSQL string = "SELECT COUNT(*) FROM users WHERE applicationID = $1 AND ID = $2"
	if err := db.database.QueryRow(countSQL, applicationID, userID).Scan(&userCount); err != nil {
		return nil, &Error{ Message: err.Error(), Status: 500 }
	}

	if userCount == 0 {
		return nil, &Error{ Message: "A user with the provided ID and applicationID does not exist.", Status: 404 }
	}

	var SQL string = "UPDATE users SET username = $1 WHERE applicationID = $2 AND ID = $3"

	_, err := db.database.Exec(SQL, newUsername, applicationID, userID)

	if err != nil {
		return nil, &Error{ Message: err.Error(), Status: 500 }
	}

	SQL = "SELECT * FROM users WHERE applicationID = $1 AND ID = $2"

	rows, err := db.database.Query(SQL, applicationID, userID)

	if err != nil {
		return nil, &Error{ Message: err.Error(), Status: 500 }
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
	return nil, &Error{ Message: "The users username could not be changed.", Status: 401 }
}

// Updates a users password
func (db *Database) SetPassword (applicationID uuid.UUID, userID uuid.UUID, newPassword string) (*User, *Error) {
	var userCount int

	var countSQL string = "SELECT COUNT(*) FROM users WHERE applicationID = $1 AND ID = $2"
	if err := db.database.QueryRow(countSQL, applicationID, userID).Scan(&userCount); err != nil {
		return nil, &Error{ Message: err.Error(), Status: 500 }
	}

	if userCount == 0 {
		return nil, &Error{ Message: "A user with the provided ID and applicationID does not exist.", Status: 404 }
	}

	var SQL string = "UPDATE users SET password = $1 WHERE applicationID = $2 AND ID = $3"

	strCh := make(chan *StringResult)
	go HashString(strCh, newPassword)
	result := <- strCh

	if result.Error != nil {
		return nil, &Error{ Message: result.Error.Message, Status: result.Error.Status }
	}

	_, err := db.database.Exec(SQL, result.String, applicationID, userID)

	if err != nil {
		return nil, &Error{ Message: err.Error(), Status: 500 }
	}

	SQL = "SELECT * FROM users WHERE applicationID = $1 AND ID = $2"

	rows, err := db.database.Query(SQL, applicationID, userID)

	if err != nil {
		return nil, &Error{ Message: err.Error(), Status: 500 }
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
	return nil, &Error{ Message: "The users password could not be changed.", Status: 401 }
}

// Deletes a user from the database
func (db *Database) DeleteUser (applicationID uuid.UUID, userID uuid.UUID) *Error {
	var userCount int

	var countSQL string = "SELECT COUNT(*) FROM users WHERE applicationID = $1 AND ID = $2"
	if err := db.database.QueryRow(countSQL, applicationID, userID).Scan(&userCount); err != nil {
		return &Error { Message: err.Error(), Status: 500 }
	}

	if userCount == 0 {
		return &Error { Message: "A user with the provided ID and applicationID does not exist.", Status: 404 }
	}

	var SQL string = "DELETE FROM users WHERE applicationID = $1 AND ID = $2"

	result, err := db.database.Exec(SQL, applicationID, userID)

	if err != nil {
		return &Error { Message: err.Error(), Status: 500 }
	}

	rowsAffected, err := result.RowsAffected()
	
	if err != nil {
		return &Error { Message: err.Error(), Status: 500 }
	}

	if rowsAffected == 0 {
		return &Error { Message: "The user was not found.", Status: 404 }
	}

	return nil
}
