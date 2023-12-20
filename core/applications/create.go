package applications

import (
	"errors"
	"time"

	"github.com/Azpect3120/AuthenticationServer/core/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Create an application using the provided
// input. The columns slice is not validated
// in this function.
func New(name string, columns []string) *model.Application {
	app := &model.Application{
		ID:            uuid.New(),
		Name:          name,
		Columns:       columns,
		CreatedAt:     time.Now().UTC(),
		LastUpdatedAt: time.Now().UTC(),
	}
	app.Columns = append([]string{"id", "applicationid"}, app.Columns...)
	app.Columns = append(app.Columns, []string{"createdat", "lastupdatedat"}...)
	return app
}

// Insert an application into the database.
// Object(id) should not exist in the database
// already.
func Insert(db *model.Database, app *model.Application) (int, error) {
	stmt, err := db.Conn.Prepare("INSERT INTO applications (id, name, columns, createdat, lastupdatedat) VALUES ($1, $2, $3, $4, $5);")
	if err != nil {
		return 500, err
	}
	defer stmt.Close()
	a := pq.StringArray(app.Columns)
	res, err := stmt.Exec(app.ID, app.Name, a, app.CreatedAt, app.LastUpdatedAt)
	if err != nil {
		return 500, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 500, err
	}
	if count < 1 {
		return 500, errors.New("Could not insert application into database.")
	}
	return 201, nil
}
