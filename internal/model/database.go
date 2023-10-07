package model

import (
	"database/sql"
	"errors"
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
func (db *Database) CreateApplication (appName string) (*Application, error) {
	var application *Application = &Application{
		ID: uuid.New(), 
		Name: appName,
	} 

	var SQL string = "INSERT INTO applications (ID, name) VALUES ($1, $2)"

	// Omit the result return
	if _, err := db.database.Exec(SQL, application.ID, application.Name); err != nil {
		return nil, err
	}
		
	return application, nil
}

// Create a user
func (db *Database) CreateUser (applicationID uuid.UUID, username string, password string) (*User, error) {
	var user *User = &User{
		ID: uuid.New(),
		Username: username,
		Password: password,
		ApplicationID: applicationID,
	}

	var SQL string = "INSERT INTO users (ID, applicationID, username, password) VALUES ($1, $2, $3, $4)"

	// Omit the result return
	if _, err := db.database.Exec(SQL, user.ID, user.ApplicationID, user.Username, user.Password); err != nil {
		return nil, err
	}

	return user, nil
}

// Verify a username and password
func (db *Database) VerifyUser (applicationID uuid.UUID, username string, password string) (*User, error) {
	var SQL string = "SELECT applicationID, ID, password FROM users WHERE applicationID = $1 AND username = $2";

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
	var SQL string = "SELECT * FROM users WHERE applicationID = $1 AND ID = $2"

	rows, err := db.database.Query(SQL, applicationID, userID)

	if err != nil {
		return nil, err
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
			return nil, err 
		}

		user := &User{
			ID: ID,
			ApplicationID: ApplicationID,
			Username: Username,
			Password: Password,
		}

		return user, nil
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nil, errors.New("User was not found") 
}


// Returns an array of users found in an application
func (db *Database) GetUsers (applicationID string) ([]*User, error) {
	appUUID, err := uuid.Parse(applicationID)
	if err != nil {
		return nil, err
	}
	
	var appCount int

	var countSQL string = "SELECT COUNT(*) FROM applications WHERE ID = $1"
	if err := db.database.QueryRow(countSQL, appUUID).Scan(&appCount); err != nil {
		return nil, err
	}

	if appCount == 0 {
		return nil, errors.New("Application with the provided ID does not exist.")
	}

	var SQL string = "SELECT * FROM users WHERE applicationID = $1"
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
	var userCount int

	var countSQL string = "SELECT COUNT(*) FROM users WHERE applicationID = $1 AND ID = $2"
	if err := db.database.QueryRow(countSQL, applicationID, userID).Scan(&userCount); err != nil {
		return nil, err
	}

	if userCount == 0 {
		return nil, errors.New("A user with the provided ID and applicationID does not exist.")
	}

	var SQL string = "UPDATE users SET username = $1 WHERE applicationID = $2 AND ID = $3"

	_, err := db.database.Exec(SQL, newUsername, applicationID, userID)

	if err != nil {
		return nil, err
	}

	SQL = "SELECT * FROM users WHERE applicationID = $1 AND ID = $2"

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
func (db *Database) SetPassword (applicationID uuid.UUID, userID uuid.UUID, newPassword string) (*User, error) {
	var userCount int

	var countSQL string = "SELECT COUNT(*) FROM users WHERE applicationID = $1 AND ID = $2"
	if err := db.database.QueryRow(countSQL, applicationID, userID).Scan(&userCount); err != nil {
		return nil, err
	}

	if userCount == 0 {
		return nil, errors.New("A user with the provided ID and applicationID does not exist.")
	}

	var SQL string = "UPDATE users SET password = $1 WHERE applicationID = $2 AND ID = $3"

	hashedPassword, err := HashString(newPassword)

	if err != nil {
		return nil, err
	}

	_, err2 := db.database.Exec(SQL, hashedPassword, applicationID, userID)

	if err2 != nil {
		return nil, err
	}

	SQL = "SELECT * FROM users WHERE applicationID = $1 AND ID = $2"

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
	return nil, errors.New("The users password could not be changed.")
}

// Deletes a user from the database
func (db *Database) DeleteUser (applicationID uuid.UUID, userID uuid.UUID) error {
	var userCount int

	var countSQL string = "SELECT COUNT(*) FROM users WHERE applicationID = $1 AND ID = $2"
	if err := db.database.QueryRow(countSQL, applicationID, userID).Scan(&userCount); err != nil {
		return err
	}

	if userCount == 0 {
		return errors.New("A user with the provided ID and applicationID does not exist.")
	}

	var SQL string = "DELETE FROM users WHERE applicationID = $1 AND ID = $2"

	result, err := db.database.Exec(SQL, applicationID, userID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("The user was not found")	
	}

	return nil
}
