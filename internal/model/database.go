package model

import (
	"database/sql"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type Database struct {
	connectionString string
	database         *sql.DB
}

func CreateDatabase() *Database {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	database := &Database{
		connectionString: os.Getenv("db_url"),
	}

	db, err := sql.Open("postgres", database.connectionString)

	if err != nil {
		log.Fatal(err)
	}

	database.database = db
	return database
}

// Create an application
func (db *Database) CreateApplication(ch chan *AppResult, appName string) {
	var application *Application = &Application{
		ID:   uuid.New(),
		Name: appName,
	}

	var SQL string = "INSERT INTO applications (ID, name) VALUES ($1, $2)"

	// Omit the result return
	if _, err := db.database.Exec(SQL, application.ID, application.Name); err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &AppResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	// return application, nil
	ch <- &AppResult{application, nil}
}

// Create a user
func (db *Database) CreateUser(ch chan *UserResult, applicationID uuid.UUID, username string, password string, data string) {
	var applicationCount int

	var countSQL string = "SELECT COUNT(*) FROM applications WHERE ID = $1"
	if err := db.database.QueryRow(countSQL, applicationID).Scan(&applicationCount); err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UserResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	if applicationCount == 0 {
		// return nil, &Error{ Message: "Invalid applicationID.", Status: 401 }
		ch <- &UserResult{nil, &Error{Message: "Invalid applicationID.", Status: 401}}
		return
	}

	var userCount int
	var userSQL string = "SELECT COUNT(*) FROM users WHERE applicationId = $1 AND username ILIKE $2"
	if err := db.database.QueryRow(userSQL, applicationID, username).Scan(&userCount); err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UserResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	if userCount > 0 {
		// return nil, &Error{ Message: "A user with that username already exists.", Status: 401 }
		ch <- &UserResult{nil, &Error{Message: "A user with that username already exists.", Status: 401}}
		return
	}

	var user *User = &User{
		ID:            uuid.New(),
		Username:      username,
		Password:      password,
		ApplicationID: applicationID,
		Data:          data,
	}

	var SQL string = "INSERT INTO users (ID, applicationID, username, password, data) VALUES ($1, $2, $3, $4, $5)"

	// Omit the result return
	if _, err := db.database.Exec(SQL, user.ID, user.ApplicationID, user.Username, user.Password, user.Data); err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UserResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	// return user, nil
	ch <- &UserResult{user, nil}
}

// Verify a username and password
func (db *Database) VerifyUser(ch chan *UserResult, applicationID uuid.UUID, username string, password string) {
	var SQL string = "SELECT applicationID, ID, password, data FROM users WHERE applicationID = $1 AND username = $2"

	rows, err := db.database.Query(SQL, applicationID, username)

	if err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UserResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	defer rows.Close()

	for rows.Next() {
		var (
			ApplicationID uuid.UUID
			ID            uuid.UUID
			Password      string
			Data		string
		)
		if err := rows.Scan(&ApplicationID, &ID, &Password, &Data); err != nil {
			// return nil, &Error{ Message: err.Error(), Status: 500 }
			ch <- &UserResult{nil, &Error{Message: err.Error(), Status: 500}}
			return
		}

		valid, err := CompareString(password, Password)

		if err != nil {
			// return nil, &Error{ Message: err.Message, Status: err.Status }
			ch <- &UserResult{nil, &Error{Message: err.Message, Status: err.Status}}
			return
		}

		if valid {
			user := &User{
				ID:            ID,
				ApplicationID: ApplicationID,
				Username:      username,
				Password:      Password,
				Data:          Data,
			}

			// return user, nil
			ch <- &UserResult{user, nil}
			return
		}
	}
	// return nil, &Error{ Message: "User was not verified", Status: 401 }
	ch <- &UserResult{nil, &Error{Message: "User was not verified", Status: 401}}
}

// Get a user from the database using its ID
func (db *Database) GetUser(ch chan *UserResult, applicationID string, userID string) {
	var SQL string = "SELECT * FROM users WHERE applicationID = $1 AND ID = $2"

	rows, err := db.database.Query(SQL, applicationID, userID)

	if err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UserResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	defer rows.Close()

	for rows.Next() {
		var (
			ID            uuid.UUID
			ApplicationID uuid.UUID
			Username      string
			Password      string
			Data          string
		)

		if err := rows.Scan(&ID, &ApplicationID, &Username, &Password, &Data); err != nil {
			// return nil, &Error{ Message: err.Error(), Status: 500 }
			ch <- &UserResult{nil, &Error{Message: err.Error(), Status: 500}}
			return
		}

		user := &User{
			ID:            ID,
			ApplicationID: ApplicationID,
			Username:      Username,
			Password:      Password,
			Data:          Data,
		}

		// return user, nil
		ch <- &UserResult{user, nil}
		return
	}

	if err := rows.Err(); err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500}
		ch <- &UserResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	// return nil, &Error{ Message: "User was not found.", Status: 404 }
	ch <- &UserResult{nil, &Error{Message: "User was not found.", Status: 404}}
}

// Returns an array of users found in an application
func (db *Database) GetUsers(ch chan *UsersResult, applicationID string) {
	appUUID, err := uuid.Parse(applicationID)
	if err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UsersResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	var appCount int

	var countSQL string = "SELECT COUNT(*) FROM applications WHERE ID = $1"
	if err := db.database.QueryRow(countSQL, appUUID).Scan(&appCount); err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UsersResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	if appCount == 0 {
		// return nil, &Error{ Message: "Application with the provided ID does not exist.", Status: 404 }
		ch <- &UsersResult{nil, &Error{Message: "Application with the provided ID does not exist.", Status: 404}}
		return
	}

	var SQL string = "SELECT * FROM users WHERE applicationID = $1"
	rows, err := db.database.Query(SQL, appUUID)
	if err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UsersResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	users := []*User{}

	for rows.Next() {
		var (
			ID            uuid.UUID
			ApplicationID uuid.UUID
			Username      string
			Password      string
			Data          string
		)

		rows.Scan(&ID, &ApplicationID, &Username, &Password, &Data)

		user := &User{
			ApplicationID: ApplicationID,
			ID:            ID,
			Username:      Username,
			Password:      Password,
			Data:          Data,
		}

		users = append(users, user)
	}

	// return users, nil
	ch <- &UsersResult{users, nil}
}

// Updates a users username
func (db *Database) SetUsername(ch chan *UserResult, applicationID uuid.UUID, userID uuid.UUID, newUsername string) {
	var userCount int

	var countSQL string = "SELECT COUNT(*) FROM users WHERE applicationID = $1 AND ID = $2"
	if err := db.database.QueryRow(countSQL, applicationID, userID).Scan(&userCount); err != nil {
		ch <- &UserResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	if userCount == 0 {
		ch <- &UserResult{nil, &Error{Message: "A user with the provided ID and applicationID does not exist.", Status: 404}}
		return
	}

	var usernameCount int

	var usernameCountSQL string = "SELECT COUNT(*) FROM users WHERE username = $1"
	if err := db.database.QueryRow(usernameCountSQL, newUsername).Scan(&usernameCount); err != nil {
		ch <- &UserResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	if usernameCount > 0 {
		ch <- &UserResult{nil, &Error{Message: "That username is taken, please try again.", Status: 401}}
		return
	}

	var SQL string = "UPDATE users SET username = $1 WHERE applicationID = $2 AND ID = $3"

	_, err := db.database.Exec(SQL, newUsername, applicationID, userID)

	if err != nil {
		ch <- &UserResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	SQL = "SELECT * FROM users WHERE applicationID = $1 AND ID = $2"

	rows, err := db.database.Query(SQL, applicationID, userID)

	if err != nil {
		ch <- &UserResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	for rows.Next() {
		var (
			ID            uuid.UUID
			ApplicationID uuid.UUID
			Username      string
			Password      string
			Data          string
		)

		rows.Scan(&ID, &ApplicationID, &Username, &Password, &Data)

		user := &User{
			ID:            ID,
			ApplicationID: ApplicationID,
			Username:      Username,
			Password:      Password,
			Data:          Data,
		}

		ch <- &UserResult{user, nil}
		return
	}
	ch <- &UserResult{nil, &Error{Message: "The users username could not be changed.", Status: 401}}
}

// Updates a users password
func (db *Database) SetPassword(ch chan *UserResult, applicationID uuid.UUID, userID uuid.UUID, newPassword string) {
	var userCount int

	var countSQL string = "SELECT COUNT(*) FROM users WHERE applicationID = $1 AND ID = $2"
	if err := db.database.QueryRow(countSQL, applicationID, userID).Scan(&userCount); err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UserResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	if userCount == 0 {
		// return nil, &Error{ Message: "A user with the provided ID and applicationID does not exist.", Status: 404 }
		ch <- &UserResult{nil, &Error{Message: "A user with the provided ID and applicationID does not exist.", Status: 404}}
		return
	}

	var SQL string = "UPDATE users SET password = $1 WHERE applicationID = $2 AND ID = $3"

	strCh := make(chan *StringResult)
	go HashString(strCh, newPassword)
	result := <-strCh

	if result.Error != nil {
		// return nil, &Error{ Message: result.Error.Message, Status: result.Error.Status }
		ch <- &UserResult{nil, &Error{Message: result.Error.Message, Status: result.Error.Status}}
		return
	}

	_, err := db.database.Exec(SQL, result.String, applicationID, userID)

	if err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UserResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	SQL = "SELECT * FROM users WHERE applicationID = $1 AND ID = $2"

	rows, err := db.database.Query(SQL, applicationID, userID)

	if err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UserResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	for rows.Next() {
		var (
			ID            uuid.UUID
			ApplicationID uuid.UUID
			Username      string
			Password      string
			Data          string
		)

		rows.Scan(&ID, &ApplicationID, &Username, &Password, &Data)

		user := &User{
			ID:            ID,
			ApplicationID: ApplicationID,
			Username:      Username,
			Password:      Password,
			Data:          Data,
		}

		// return user, nil
		ch <- &UserResult{user, nil}
		return
	}
	// return nil, &Error{ Message: "The users password could not be changed.", Status: 401 }
	ch <- &UserResult{nil, &Error{Message: "The users password could not be changed.", Status: 401}}
}

// Updates a users data
func (db *Database) SetData(ch chan *UserResult, applicationID uuid.UUID, userID uuid.UUID, newData string) {
	var userCount int

	var countSQL string = "SELECT COUNT(*) FROM users WHERE applicationID = $1 AND ID = $2"
	if err := db.database.QueryRow(countSQL, applicationID, userID).Scan(&userCount); err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UserResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	if userCount == 0 {
		// return nil, &Error{ Message: "A user with the provided ID and applicationID does not exist.", Status: 404 }
		ch <- &UserResult{nil, &Error{Message: "A user with the provided ID and applicationID does not exist.", Status: 404}}
		return
	}

	var SQL string = "UPDATE users SET data = $1 WHERE applicationID = $2 AND ID = $3"

	_, err := db.database.Exec(SQL, newData, applicationID, userID)

	if err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UserResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	SQL = "SELECT * FROM users WHERE applicationID = $1 AND ID = $2"

	rows, err := db.database.Query(SQL, applicationID, userID)

	if err != nil {
		// return nil, &Error{ Message: err.Error(), Status: 500 }
		ch <- &UserResult{nil, &Error{Message: err.Error(), Status: 500}}
		return
	}

	for rows.Next() {
		var (
			ID            uuid.UUID
			ApplicationID uuid.UUID
			Username      string
			Password      string
			Data          string
		)

		rows.Scan(&ID, &ApplicationID, &Username, &Password, &Data)

		user := &User{
			ID:            ID,
			ApplicationID: ApplicationID,
			Username:      Username,
			Password:      Password,
			Data:          Data,
		}

		// return user, nil
		ch <- &UserResult{user, nil}
		return
	}
	// return nil, &Error{ Message: "The users username could not be changed.", Status: 401 }
	ch <- &UserResult{nil, &Error{Message: "The users username could not be changed.", Status: 401}}
}

// Deletes a user from the database
func (db *Database) DeleteUser(ch chan *ErrorResult, applicationID uuid.UUID, userID uuid.UUID) {
	var userCount int

	var countSQL string = "SELECT COUNT(*) FROM users WHERE applicationID = $1 AND ID = $2"
	if err := db.database.QueryRow(countSQL, applicationID, userID).Scan(&userCount); err != nil {
		// return &Error { Message: err.Error(), Status: 500 }
		ch <- &ErrorResult{&Error{Message: err.Error(), Status: 500}}
		return
	}

	if userCount == 0 {
		// return &Error { Message: "A user with the provided ID and applicationID does not exist.", Status: 404 }
		ch <- &ErrorResult{&Error{Message: "A user with the provided ID and applicationID does not exist.", Status: 404}}
		return
	}

	var SQL string = "DELETE FROM users WHERE applicationID = $1 AND ID = $2"

	result, err := db.database.Exec(SQL, applicationID, userID)

	if err != nil {
		// return &Error { Message: err.Error(), Status: 500 }
		ch <- &ErrorResult{&Error{Message: err.Error(), Status: 500}}
		return
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		// return &Error { Message: err.Error(), Status: 500 }
		ch <- &ErrorResult{&Error{Message: err.Error(), Status: 500}}
		return
	}

	if rowsAffected == 0 {
		// return &Error { Message: "The user was not found.", Status: 404 }
		ch <- &ErrorResult{&Error{Message: "The user was not found.", Status: 404}}
		return
	}

	// return nil
	ch <- &ErrorResult{nil}
}
