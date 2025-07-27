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
	r.PATCH("/users/:user_id", handlers.UpdateUser)
	r.GET("/messages", handlers.GetMessages)
	r.POST("/messages", handlers.CreateMessage)
	r.GET("/users/:user_id/messages", handlers.GetMessagesByUser)
	r.Run()
}
