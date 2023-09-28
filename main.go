package main

import (
	"authServer/api"
	"fmt"
)


func main () {
	var database *api.Database = api.CreateDatabase();

	if database != nil {
		fmt.Println("SUCCESS!!!!")
	}
	
	server := api.CreateServer("8080", *database);

	server.Listen()
}
