package main

import (
	"MUJ_automated_mail_generation/pkg/database"
	"MUJ_automated_mail_generation/pkg/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	r := gin.Default()

	// Routes
	r.POST("/submit", handler.SubmitHandler)

	r.Run(":8080")
}
