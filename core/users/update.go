package users

import (
	"database/sql"
	"time"

	"github.com/Azpect3120/AuthenticationServer/core/model"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// Updates a user in the database. The 'id' parameter
// should be parsed into a UUID to ensure validity. Only provided
// parameters will be updated. The 'int' return of this function
// is the HTTP status code that should be sent back to the user
// upon calling this function.
func Update (db *model.Database, id uuid.UUID, uid uuid.UUID, data *model.UserData) (*model.User, int, error) {
  stmt, err := db.Conn.Prepare("SELECT * FROM users WHERE id = $1 AND applicationid = $2;")
  if err != nil {
    return nil, 500, err
  }
  defer stmt.Close()

  var user model.User
  var (
    username sql.NullString
    first sql.NullString
    last sql.NullString
    full sql.NullString
    email sql.NullString
    password sql.NullString
    userData sql.NullString
  )
  if err := stmt.QueryRow(uid, id).Scan(&user.ID, &user.ApplicationID, &username, &first, &last, &full, &email, &password, &userData, &user.CreatedAt, &user.LastUpdatedAt); err != nil {
    return nil, 404, err
  }

  user.Username = username.String
  user.First = first.String
  user.Last = last.String
  user.Full = full.String
  user.Email = email.String
  user.Password = password.String
  user.Data = userData.String

  cols, err := GetApplicationColumns(db, id)
  if err != nil {
    return nil, 500, err
  }

  for _, col := range cols {
    switch col {
    case "username":
      if data.Username != "" { user.Username = data.Username }
    case "first":
      if data.First != "" { user.First = data.First }
    case "last":
      if data.Last != "" { user.Last = data.Last }
    case "full":
      if data.Full != "" { user.Full = data.Full }
    case "email":
      if data.Email != "" { user.Email = data.Email }
    case "password":
      if data.Password != "" { user.Password = data.Password }
    case "data":
      if data.Data != "" { user.Data = data.Data }
    }
  }

  stmt, err = db.Conn.Prepare("UPDATE users SET username = $1, firstname = $2, lastname = $3, fullname = $4, email = $5, password = $6, data = $7, lastupdatedat = $8 WHERE id = $9 AND applicationid = $10 RETURNING *;")
  if err != nil {
    return nil, 500, err
  }
  defer stmt.Close()

  if err := stmt.QueryRow(user.Username, user.First, user.Last, user.Full, user.Email, user.Password, user.Data, time.Now().UTC(), uid, id).Scan(&user.ID, &user.ApplicationID, &user.Username, &user.First, &user.Last, &user.Full, &user.Email, &user.Password, &user.Data, &user.CreatedAt, &user.LastUpdatedAt); err != nil {
    return nil, 500, err
  }

  return &user, 200, nil
}
