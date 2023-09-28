package api

import "github.com/google/uuid"

// Database table: Users
type User struct {
	ID				uuid.UUID   `json:"ID"` 
	Username		string		`json:"username"` 
	Password		string		`json:"password"`
	ApplicationID	uuid.UUID	`json:"applicationID"`
}

// Database table: Applications
type Application struct {
	ID			uuid.UUID	`json:"ID"`	
	Name		string		`json:"name"`
}

// Request map: /createApplication
type CreateApplicationRequest struct {
	Name	string		`json:"name"`
}

// Request map: /createUser
type CreateUserRequest struct {
	ApplicationID	uuid.UUID	`json:"applicationID"`
	Username		string		`json:"username"`
	Password		string		`json:"password"`
}

// Request map: /verifyUser
type VerifyUserRequest struct {
	ApplicationID	uuid.UUID	`json:"applicationID"`
	Username		string		`json:"username"`
	Password		string		`json:"password"`
}
