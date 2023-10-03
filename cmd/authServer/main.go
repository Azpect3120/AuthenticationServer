package main

import (
	"github.com/Azpect3120/AuthenticationServer/internal/model"
	"os"
)

func main () {
	var database *model.Database = model.CreateDatabase();

	if database == nil {
		panic("Database connection failure")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := model.CreateServer(port, *database);

	server.Listen()
}
