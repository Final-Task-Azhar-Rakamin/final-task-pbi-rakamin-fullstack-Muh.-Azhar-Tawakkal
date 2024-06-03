package main

import (
	"finalapi/database"
	"finalapi/router"
)

func main() {
	router := router.SetRouter()
	database.ConnectDB()
	router.Run("localhost:9090")

}
