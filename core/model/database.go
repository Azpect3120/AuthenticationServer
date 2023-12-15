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

type Application struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Columns       []string  `json:"columns"`
	CreatedAt     time.Time `json:"createdat"`
	LastUpdatedAt time.Time `json:"lastupdatedat"`
}
