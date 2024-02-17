package applications

import (
	"errors"

	"github.com/Azpect3120/AuthenticationServer/core/model"
	"github.com/google/uuid"
  "github.com/lib/pq"
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
	res, err := stmt.Exec(id)
	if err != nil {
    pqErr, ok := err.(*pq.Error)
    if ok && pqErr.Code == "23503" && pqErr.Constraint == "users_applicationid_fkey" {
      return 409, errors.New("Cannot delete applications with users.")
    }
		return 500, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 500, err
	}
	if count < 1 {
		return 404, errors.New("Application not found.")
	}
	return 204, nil
}
