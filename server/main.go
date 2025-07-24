package main

import (
	"log"
	"main/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Server is starting...")

	r := gin.Default()

	// endpoints
	r.GET("/users", handlers.GetUsers)
	r.POST("/users", handlers.CreateUser)
	r.Run()
}
