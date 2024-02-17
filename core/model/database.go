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
  Username      string    `json:"username"`
  FirstName     string    `json:"firstname"`
  LastName      string    `json:"lastname"`
  FullName      string    `json:"fullname"`
  Email         string    `json:"email"`
  Password      string    `json:"password"`
  Data          string    `json:"data"`
	CreatedAt     time.Time `json:"createdat"`
	LastUpdatedAt time.Time `json:"lastupdatedat"`
}

type UserData struct {
  Username      string    `json:"username"`
  FirstName     string    `json:"firstname"`
  LastName      string    `json:"lastname"`
  FullName      string    `json:"fullname"`
  Email         string    `json:"email"`
  Password      string    `json:"password"`
  Data          string    `json:"data"`
}

type Application struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Columns       []string  `json:"columns"`
	CreatedAt     time.Time `json:"createdat"`
	LastUpdatedAt time.Time `json:"lastupdatedat"`
}
