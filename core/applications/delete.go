package applications

import (
	"errors"
	"fmt"

	"github.com/Azpect3120/AuthenticationServer/core/model"
	"github.com/google/uuid"
)

// Deletes an application from the database
// id should be parsed to a UUID to ensure
// its validity. The 'int' return on this
// function is the HTTP status code which
// should be sent back to the user upon
// calling this function.
func Delete (db *model.Database, id uuid.UUID) (int, error) {
	stmt, err := db.Conn.Prepare("DELETE FROM applications WHERE id = $1;")
	if err != nil {
		return 500, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(id.String())
	if err != nil {
		return 500, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 500, err
	}
	if count < 1 {
		return 404, errors.New(fmt.Sprintf("Could not find an application with the id: '%s'", id))
	}
	return 204, nil
}
