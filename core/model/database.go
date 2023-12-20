package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Database struct {
	ConnectionString string
	Conn             *sql.DB
}

type User struct {
	ID            uuid.UUID `json:"id"`
	ApplicationID uuid.UUID `json:"applicationID"`
	// data here, columns defined in the
	// parent application ...
	CreatedAt     time.Time `json:"createdat"`
	LastUpdatedAt time.Time `json:"lastupdatedat"`
}

type Application struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Columns       []string  `json:"columns"`
	CreatedAt     time.Time `json:"createdat"`
	LastUpdatedAt time.Time `json:"lastupdatedat"`
}
