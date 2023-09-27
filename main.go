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

	//application := database.CreateApplication("Testing")

	//var uuidStr string = "9ee5c706-a722-47f1-9049-dc486f8641e2"
	//uuid, err := uuid.Parse(uuidStr)
	//user := database.CreateUser(uuid, "Azpect", "root")

	server.Listen()
}
