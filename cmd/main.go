package main

import (
	"MUJ_automated_mail_generation/pkg/database"
	"MUJ_automated_mail_generation/pkg/handler"
	"MUJ_automated_mail_generation/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	database.Connect()

	// Create a new Gin router
	r := gin.Default()

	// Routes for student submissions
	r.POST("/submit", handler.SubmitHandler)

	// Routes for reviewer management
	r.POST("/reviewer", handler.CreateReviewerHandler)
	r.POST("/reviewer/login", handler.LoginReviewerHandler)
	r.GET("/reviewer/:department", handler.GetReviewerDetailsHandler)

	// Authenticated routes for reviewers
	authReviewer := r.Group("/reviewer")
	authReviewer.Use(middleware.AuthMiddleware()) // Middleware for authentication
	{
		authReviewer.POST("/reviews", handler.CreateReviewHandler)
		authReviewer.PUT("/reviews/:id", handler.UpdateReviewHandler)
		authReviewer.GET("/reviews", handler.GetAllReviewsHandler)
		authReviewer.GET("/reviews/submission/:submission_id", handler.GetReviewsBySubmissionHandler)
		authReviewer.GET("/reviews/reviewer/:reviewer_id", handler.GetReviewsByReviewerHandler)
	}

	// Routes for HoD management
	r.POST("/hod/login", handler.LoginHoDHandler)                 // HoD login
	r.POST("/hod", handler.CreateHoDHandler)                      // Create HoD
	r.GET("/hods", handler.GetAllHoDsHandler)                     // Get all HoDs
	r.GET("/hod/:department", handler.GetHoDsByDepartmentHandler) // Get HoDs by department

	// Authenticated routes for HoDs
	authHoD := r.Group("/hod")
	authHoD.Use(middleware.AuthMiddleware()) // Middleware for authentication
	{
		// HoD review actions (To be added later)
	}

	// Routes for student submissions and reviewers
	r.GET("/submissions", handler.GetAllSubmissionsHandler) // Get all student submissions
	r.GET("/reviewers", handler.GetAllReviewersHandler)     // Get all reviewers

	// Start the server on port 8080
	r.Run(":8080")
}
