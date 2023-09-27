package api

import "github.com/google/uuid"

type User struct {
	ID			uuid.UUID   `json:"ID"` 
	Username	string		`json:"username"` 
	Password	string		`json:"password"`
}

type Application struct {
	ID			uuid.UUID	`json:"ID"`	
	Name		string		`json:"name"`
	Key			uuid.UUID	`json:"key"`
}
