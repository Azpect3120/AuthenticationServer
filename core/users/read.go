package users

import (
	"reflect"
	"time"

	"github.com/Azpect3120/AuthenticationServer/core/model"
	"github.com/google/uuid"
)

// Define a map to store column names
var COLUMNS = map[string]string{
    "id":            "ID",
    "applicationid": "ApplicationID",
    "username":      "Username",
    "first":         "First",
    "last":          "Last",
    "full":          "Full",
    "email":         "Email",
    "password":      "Password",
    "data":          "Data",
    "createdat":     "CreatedAt",
    "lastupdatedat": "LastUpdatedAt",
}

// Retrieves a user from the database based on
// the given ID.
func Retrieve (db *model.Database, id uuid.UUID, uid uuid.UUID) (*map[string]string, int, error) {
  stmt, err := db.Conn.Prepare("SELECT * FROM users WHERE id = $1 AND applicationid = $2;")
  if err != nil {
    return nil, 500, err
  }
  defer stmt.Close()

  var user model.User = model.User{}
  if err := stmt.QueryRow(uid, id).Scan(
    &user.ID,
    &user.ApplicationID,
    &user.Username,
    &user.First,
    &user.Last,
    &user.Full,
    &user.Email,
    &user.Password,
    &user.Data,
    &user.CreatedAt,
    &user.LastUpdatedAt,
  ); err != nil {
    return nil, 404, err
  }

  appColumns, err := getApplicationColumns(db, user.ApplicationID)
  if err != nil {
    return nil, 500, err
  }

  providedColumns := make(map[string]string)

  for _, col := range appColumns {
    val := reflect.ValueOf(user).FieldByName(COLUMNS[col])
    var fieldValue string
    switch val.Interface().(type) {
    case uuid.UUID:
      fieldValue = val.Interface().(uuid.UUID).String()
    case time.Time:
      fieldValue = val.Interface().(time.Time).String()
    default:
      fieldValue = val.String()
  }
    providedColumns[col] = fieldValue
  }

  return &providedColumns, 200, nil
}

// Retrieves all users from the database based on 
// the given application ID.
func RetrieveAll (db *model.Database, id uuid.UUID) ([]map[string]string, int, error) {
  stmt, err := db.Conn.Prepare("SELECT * FROM users WHERE applicationid = $1;")
  if err != nil {
    return nil, 500, err
  }
  defer stmt.Close()

  appColumns, err := getApplicationColumns(db, id)
  if err != nil {
    return nil, 500, err
  }

  rows, err := stmt.Query(id)
  if err != nil {
    return nil, 404, err
  }
  defer rows.Close()

  providedColumns := make([]map[string]string, 0)
  
  for rows.Next() {
    var user model.User = model.User{}
    if err := rows.Scan(&user.ID, &user.ApplicationID, &user.Username, &user.First, &user.Last, &user.Full, &user.Email, &user.Password, &user.Data, &user.CreatedAt, &user.LastUpdatedAt); err != nil {
      return nil, 404, err
    }

    mapColumns := make(map[string]string)
    for _, col := range appColumns {
      val := reflect.ValueOf(user).FieldByName(COLUMNS[col])
      var fieldValue string
      switch val.Interface().(type) {
        case uuid.UUID:
          fieldValue = val.Interface().(uuid.UUID).String()
        case time.Time:
          fieldValue = val.Interface().(time.Time).String()
        default:
          fieldValue = val.String()
      }
      mapColumns[col] = fieldValue
    }
    providedColumns = append(providedColumns, mapColumns)
  }

  return providedColumns, 200, nil
}
