package users

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Azpect3120/AuthenticationServer/core/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Validates the user objects column inputs based on
// the application's columns. Returns a message and
// status code for use in the response.
func Validate(db *model.Database, appId uuid.UUID, user *model.User) (string, int, error) {
	requiredColumns, err := GetApplicationColumns(db, appId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "Application not found", 404, err
		}
		return "", 500, err
	}

	providedColumns := getUserColumns(user)
	var missingColumns []string = make([]string, 0)

	for _, col := range requiredColumns {
		var found bool = false
		for _, pcol := range providedColumns {
			if pcol == col {
				found = true
				break
			}
		}
		if !found {
			missingColumns = append(missingColumns, col)
		}
	}

	if len(missingColumns) > 0 {
		return fmt.Sprintf("Missing required columns: %s", strings.Join(missingColumns, ", ")), 400, errors.New("Missing required columns: " + strings.Join(missingColumns, ", "))
	}

	return "", 200, nil
}

// Retrieves a list of the valid columns for an application based on the given ID.
func GetApplicationColumns(db *model.Database, appId uuid.UUID) ([]string, error) {
	stmt, err := db.Conn.Prepare("SELECT columns FROM applications WHERE id = $1;")
	if err != nil {
		return []string{}, err
	}
	defer stmt.Close()

	var columns pq.StringArray
	if err := stmt.QueryRow(appId).Scan(&columns); err != nil {
		return []string{}, err
	}

	return []string(columns), nil
}

// Returns an array of columns that are not empty
func getUserColumns(user *model.User) []string {
	var columns []string
	t := reflect.TypeOf(*user)
	v := reflect.ValueOf(*user)

	for i := 0; i < t.NumField(); i++ {
		name := t.Field(i).Name
		val := v.Field(i)
		if !val.IsZero() {
			columns = append(columns, strings.ToLower(name))
		}
	}

	return columns
}

// Hashes the provided string
func HashString(str string) (string, error) {
	hash := sha256.New()
	if _, err := hash.Write([]byte(str)); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// Compares the provided hash with the given string
func CompareHash(str, hash string) bool {
	if h, err := HashString(str); err != nil {
		return false
	} else if h != hash {
		return false
	}
	return true
}
