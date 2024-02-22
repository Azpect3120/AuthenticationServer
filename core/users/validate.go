package users

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Azpect3120/AuthenticationServer/core/applications"
	"github.com/Azpect3120/AuthenticationServer/core/model"
	"github.com/google/uuid"
)

// Validates a user's login credentials. The id parameter
// should be parsed to a UUID to ensure its validity. The
// columns parameter is a list of columns to be checked
// against the user's data. The user parameter is a pointer
// to a user model which contains the user's inputted data.
func ValidateLogin(db *model.Database, id uuid.UUID, columns []string, user *model.UserData) (*model.User, string, int, error) {
	if columns == nil || len(columns) == 0 {
		return nil, "No columns provided.", 400, errors.New("No columns provided.")
	}

	message := applications.MatchColumns(&columns)

	fieldValue := reflect.ValueOf(*user).FieldByName(COLUMNS[columns[0]])

	if columns[0] == "password" {
		fieldValue = reflect.ValueOf(HashString(fieldValue.String()))
	}
	fmt.Printf("SELECT * FROM users WHERE applicationid = $1 AND %s = '%s';\n", columns[0], fieldValue.String())

	stmt, err := db.Conn.Prepare(fmt.Sprintf("SELECT * FROM users WHERE applicationid = $1 AND %s = '%s';", columns[0], fieldValue.String()))
	if err != nil {
		return nil, message, 500, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, message, 500, err
	}
	defer rows.Close()

	var u model.User
	for rows.Next() {
		if err := rows.Scan(&u.ID, &u.ApplicationID, &u.Username, &u.First, &u.Last, &u.Full, &u.Email, &u.Password, &u.Data, &u.CreatedAt, &u.LastUpdatedAt); err != nil {
			return nil, message, 500, err
		}

		var valid bool = true
		for _, col := range columns {
			if col == "password" {
				dbVal := reflect.ValueOf(u).FieldByName(COLUMNS[col]).String()
				inpVal := reflect.ValueOf(*user).FieldByName(COLUMNS[col]).String()
				if !CompareHash(inpVal, dbVal) {
					valid = false
					break
				}

			} else {
				dbVal := reflect.ValueOf(u).FieldByName(COLUMNS[col]).String()
				inpVal := reflect.ValueOf(*user).FieldByName(COLUMNS[col]).String()
				if dbVal != inpVal {
					valid = false
					break
				}
			}
		}

		if valid {
			return &u, "", 200, nil
		}
	}

	return nil, message, 404, errors.New("Could not validate user.")
}
