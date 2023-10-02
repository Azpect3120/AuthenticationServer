package main

import "github.com/Azpect3120/AuthenticationServer/internal/model"

func main () {
	var database *model.Database = model.CreateDatabase();

	if database == nil {
		panic("Database connection failure")
	}

	server := model.CreateServer("8080", *database);

	server.Listen()
}
