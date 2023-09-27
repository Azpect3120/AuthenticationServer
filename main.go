package main

import (
	"authServer/api"
	"fmt"
)


func main () {
	server := api.CreateServer("8080");

	database := api.CreateDatabase();

	if database != nil {
		fmt.Println("SUCCESS!!!!")
	}

	server.Listen()
}
