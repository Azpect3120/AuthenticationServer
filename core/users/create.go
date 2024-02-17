package users

import (
	"errors"
	"time"

	"github.com/Azpect3120/AuthenticationServer/core/model"
	"github.com/google/uuid"
)

// Creates a new user with the given data and application ID.
func New (appId uuid.UUID, data *model.UserData) *model.User {
  return &model.User{
    ID: uuid.New(),
    ApplicationID: appId,
    Username:      data.Username,
    FirstName:     data.FirstName,
    LastName:      data.LastName,
    FullName:      data.FullName,
    Email:         data.Email,
    Password:      data.Password,
    Data:          data.Data,
    CreatedAt:     time.Now().UTC(),
    LastUpdatedAt: time.Now().UTC(),
  }
}

// Inserts a new user into the database. The 'int' return
// of this function is the HTTP status code that should be
// sent back to the user upon calling this function.
func Insert (db *model.Database, user *model.User) (int, error) {
  var sqlStirng string = `
    INSERT INTO users 
    (id, applicationid, username, firstname, lastname, fullname, email, password, data, createdat, lastupdatedat)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
    `
  stmt, err := db.Conn.Prepare(sqlStirng)
  if err != nil {
    return 500, err
  }
  defer stmt.Close()

  res, err := stmt.Exec(
    user.ID,
    user.ApplicationID,
    user.Username,
    user.FirstName,
    user.LastName,
    user.FullName,
    user.Email,
    user.Password,
    user.Data,
    user.CreatedAt,
    user.LastUpdatedAt,
  )
  if err != nil {
    return 500, err
  }
  
  count, err := res.RowsAffected()
  if err != nil {
    return 500, err
  }

  if count < 1 {
    return 500, errors.New("Failed to insert user into database.")
  }

  return 201, nil
}
