package users

import (
	"errors"

	"github.com/Azpect3120/AuthenticationServer/core/model"
	"github.com/google/uuid"
)

// Deletes a user from the database, ID should be parsed 
// to a UUID to ensure its validity. The 'int' return on
// this function is the HTTP status code which should be
// sent back to the user upon calling this function.
func Delete(db *model.Database, id uuid.UUID, uid uuid.UUID) (int, error) {
	stmt, err := db.Conn.Prepare("DELETE FROM users WHERE id = $1 AND applicationid = $2;")
	if err != nil {
		return 500, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(uid, id)
	if err != nil {
		return 500, err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return 500, err
	}
	if count < 1 {
		return 404, errors.New("User not found.")
	}
	return 204, nil
}
