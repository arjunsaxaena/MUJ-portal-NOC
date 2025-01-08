package main

import (
	"MUJ_automated_mail_generation/pkg/database"
	"MUJ_automated_mail_generation/pkg/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// Establish database connection
	database.Connect()

	// Initialize Gin router
	r := gin.Default()

	// Routes for student submission (existing)
	r.POST("/submit", handler.SubmitHandler)

	// Routes for reviewer management
	r.POST("/reviewer", handler.CreateReviewerHandler)              // Create reviewer
	r.POST("/reviewer/login", handler.LoginReviewerHandler)         // Reviewer login
	r.GET("/reviewer/:username", handler.GetReviewerDetailsHandler) // Get reviewer details

	// Start the server on port 8080
	r.Run(":8080")
}
