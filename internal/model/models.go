package model

import (
	"errors"
	"reflect"

	"github.com/google/uuid"
)

// Email: Email Object
type Email struct {
	To		string		`json:"to"`
	Subject	string		`json:"subject"`
	Content string		`json:"content"`
}

// Channel Result: Application
type AppResult struct {
	Application *Application
	Error *Error
}

// Channel Result: User
type UserResult struct {
	User *User
	Error *Error
}

// Channel Result: Users
type UsersResult struct {
	Users []*User
	Error *Error
}

// Channel Result: String
type StringResult struct {
	String string
	Error *Error
}

// Channel Result: Error
type ErrorResult struct {
	Error *Error
}

// Error: Custom Error
type Error struct {
	Message			string		`json:"error"`
	Status			int			`json:"status"`
}

// Database table: Users
type User struct {
	ID				uuid.UUID   `json:"ID"` 
	Username		string		`json:"username"` 
	Password		string		`json:"password"`
	ApplicationID	uuid.UUID	`json:"applicationID"`
	Data			string		`json:"data"`
}

// Database table: Applications
type Application struct {
	ID			uuid.UUID	`json:"ID"`	
	Name		string		`json:"name"`
}

// Request map: /createApplication
type CreateApplicationRequest struct {
	Name	string		`json:"name"`
	Email   string		`json:"email"`
}

// Request map: /createUser
type CreateUserRequest struct {
	ApplicationID	uuid.UUID	`json:"applicationID"`
	Username		string		`json:"username"`
	Password		string		`json:"password"`
	Data			string		`json:"data"`
}

// Request map: /verifyUser
type VerifyUserRequest struct {
	ApplicationID	uuid.UUID	`json:"applicationID"`
	Username		string		`json:"username"`
	Password		string		`json:"password"`
}

// Request map: /setUsername
type SetUsernameRequest struct {
	ApplicationID	uuid.UUID	`json:"applicationID"`
	ID				uuid.UUID	`json:"ID"`
	Username		string		`json:"username"`
}

// Request map: /setPassword
type SetPasswordRequest struct {
	ApplicationID	uuid.UUID	`json:"applicationID"`
	ID				uuid.UUID	`json:"ID"`
	Password		string		`json:"password"`
}

// Request map: /setData
type SetDataRequest struct {
	ApplicationID	uuid.UUID	`json:"applicationID"`
	ID				uuid.UUID	`json:"ID"`
	Data			string		`json:"data"`
}

// Request map: /deleteUser
type DeleteUserReqest struct {
	ApplicationID	uuid.UUID	`json:"applicationID"`	
	ID				uuid.UUID   `json:"ID"`
}

// Validates a struct to ensure nothing is blank
func Validate (s interface{}) error {
	// Get type of interface passed
	sType := reflect.ValueOf(s)

	// Convert pointer types to values
	if sType.Kind() == reflect.Ptr {
		sType = sType.Elem()
	}

	// Ensure the interface is a struct/pStruct
	if sType.Kind() != reflect.Struct {
		return errors.New("Input is not a struct or pointer to a struct")
	}
	
	// Iterate over each field in the struct
	for i := 0; i < sType.NumField(); i++ {
		// Get each field and type
		f := sType.Field(i)
		fType := f.Type()
		fName := sType.Type().Field(i).Name


		// Check if field is an empty pointer
		if fType.Kind() == reflect.Ptr && f.IsNil() {
			return errors.New("Field '" + fName + "' is nil")
		}

		// Check if field is an empty string 
		if fType.Kind() == reflect.String && f.String() == "" {
			return errors.New("Field '" + fName + "' is an empty string")
		}
	}
	return nil
} 
