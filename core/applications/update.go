package applications

import (
	"errors"
	"fmt"
	"time"

	"database/sql"

	"github.com/Azpect3120/AuthenticationServer/core/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Updates an application in the database. The 'id' parameter
// should be parsed into a UUID to ensure validity. Only provided
// parameters will be updated. The 'int' return of this function
// is the HTTP status code that should be sent back to the user
// upon calling this function.
func Update (db *model.Database, id uuid.UUID, name string, columns []string) (*model.Application, string, int, error) {
  var sqlString string = "UPDATE applications SET"
  var params []interface{}
  var message string
  params = append(params, time.Now().UTC())

  var idIndex int = 2

  if name != "" {
    sqlString += fmt.Sprintf(" name = $%d,", idIndex)
    params = append(params, name)
    idIndex++
  }

  if len(columns) > 0 {
    sqlString += fmt.Sprintf(" columns = columns || $%d,", idIndex)
    message = MatchColumns(&columns)
    params = append(params, pq.Array(columns))
    idIndex++
  }

  if name == "" && len(columns) == 0 {
    return nil, "", 400, errors.New("Name or columns are required.")
  }

  params = append(params, id)
  sqlString += fmt.Sprintf(" lastupdatedat = $1 WHERE id = $%d RETURNING *;", idIndex)

  stmt, err := db.Conn.Prepare(sqlString)
  if err != nil {
    return nil, "", 500, err
  }
  defer stmt.Close()

  var app model.Application
  var updatedColumns pq.StringArray

  if err := stmt.QueryRow(params...).Scan(&app.ID, &app.Name, &updatedColumns, &app.CreatedAt, &app.LastUpdatedAt); err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      return nil, "", 404, errors.New("Application not found.")
    }
    return nil, "", 500, err
  }
  app.Columns = []string(updatedColumns)

  validatedApp, code, err := Validate(db, id)
  if err != nil {
    return nil, "", code, err
  }

  return validatedApp, message, code, nil
}

// Overwrites an application in the database. The 'id' parameter
// should be parsed into a UUID to ensure validity. All parameters
// should be provided in the request body. Other they will be set 
// to their default values. The 'int' return of this function is 
// the HTTP status code that should be sent back to the user upon 
// calling this function.
func Overwrite (db *model.Database, id uuid.UUID, name string, columns []string) (*model.Application, string, int, error) {
  var sqlString string = "UPDATE applications SET"
  var params []interface{}
  var message string
  params = append(params, time.Now().UTC())

  if name != "" {
    sqlString += " name = $2,"
    params = append(params, name)
  } else {
    return nil, "", 400, errors.New("Name is required.")
  }

  if len(columns) > 0 {
    sqlString += " columns = $3,"
    columns = append(columns, []string{"id", "applicationid", "createdat", "lastupdatedat"}...)
    message = MatchColumns(&columns)
    params = append(params, pq.Array(columns))
  } else {
    return nil, "", 400, errors.New("Columns are required.")
  }

  params = append(params, id)
  sqlString += " lastupdatedat = $1 WHERE id = $4 RETURNING *;"

  stmt, err := db.Conn.Prepare(sqlString)
  if err != nil {
    return nil, "", 500, err
  }
  defer stmt.Close()

  var app model.Application
  var updatedColumns pq.StringArray

  if err := stmt.QueryRow(params...).Scan(&app.ID, &app.Name, &updatedColumns, &app.CreatedAt, &app.LastUpdatedAt); err != nil {
    if errors.Is(err, sql.ErrNoRows) {
      return nil, "", 404, errors.New("Application not found.")
    }
    return nil, "", 500, err
  }
  app.Columns = []string(updatedColumns)

  validatedApp, code, err := Validate(db, id)
  if err != nil {
    return nil, "", code, err
  }

  return validatedApp, message, code, nil
}
