package applications

import (
	"github.com/Azpect3120/AuthenticationServer/core/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Retrieves an application from the database. 'id'
// parameter should be parsed into a UUID to ensure
// validity. The 'int' return of this function is
// the HTTP status code that should be sent back
// to the user upon calling this function.
func Retrieve (db *model.Database, id uuid.UUID) (*model.Application, int, error) {
	stmt, err := db.Conn.Prepare("SELECT * FROM applications WHERE id = $1;")
	if err != nil {
		return nil, 500, err
	}
	defer stmt.Close()

	var app *model.Application = &model.Application{}
	var columns pq.StringArray

	if err = stmt.QueryRow(id.String()).Scan(
    &app.ID,
    &app.Name,
    &columns,
    &app.CreatedAt,
    &app.LastUpdatedAt,
  ); err != nil {
		return nil, 500, err
	}

	app.Columns = []string(columns)
	return app, 200, nil
}

// Retrieves all applications from the database. The 
// 'int' return of this function is the HTTP status
// code that should be sent back to the user upon
// requesting this data.
func RetrieveAll (db *model.Database) ([]*model.Application, int, error) {
	stmt, err := db.Conn.Prepare("SELECT * FROM applications;")
	if err != nil {
		return nil, 500, err
	}
  defer stmt.Close()

  var apps []*model.Application = make([]*model.Application, 0)
  var columns pq.StringArray

  rows, err := stmt.Query()
  if err != nil {
    return nil, 500, err
  }
  defer rows.Close()

  for rows.Next() {
    var app *model.Application = &model.Application{}
    rows.Scan(
      &app.ID,
      &app.Name,
      &columns,
      &app.CreatedAt,
      &app.LastUpdatedAt,
    )
    app.Columns = []string(columns)
    apps = append(apps, app)
  }

  return apps, 200, nil
}
