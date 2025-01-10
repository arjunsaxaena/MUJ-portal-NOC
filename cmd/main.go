package main

import (
	"MUJ_automated_mail_generation/pkg/database"
	"MUJ_automated_mail_generation/pkg/handler"
	"MUJ_automated_mail_generation/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	r := gin.Default()

	r.POST("/submit", handler.SubmitHandler)

	r.POST("/reviewer", handler.CreateReviewerHandler)
	r.POST("/reviewer/login", handler.LoginReviewerHandler)
	r.GET("/reviewer/:department", handler.GetReviewerDetailsHandler) // just for me useful while doing testing not to implement.

	authReviewer := r.Group("/reviewer")
	authReviewer.Use(middleware.AuthMiddleware()) // Middleware for authentication
	{
		authReviewer.GET("/submissions", handler.GetSubmissionsHandler)
		authReviewer.POST("/reviews", handler.CreateReviewHandler)
		authReviewer.PUT("/reviews/:id", handler.UpdateReviewHandler)
		authReviewer.GET("/reviews", handler.GetAllReviewsHandler)
		authReviewer.GET("/reviews/submission/:submission_id", handler.GetReviewsBySubmissionHandler)
		authReviewer.GET("/reviews/reviewer/:reviewer_id", handler.GetReviewsByReviewerHandler)
	}

	// Routes for HoD
	r.POST("/hod/login", handler.LoginHoDHandler)
	r.POST("/hod", handler.CreateHoDHandler)
	r.GET("/hods", handler.GetAllHoDsHandler)
	r.GET("/hod/:department", handler.GetHoDsByDepartmentHandler)

	authHoD := r.Group("/hod")
	authHoD.Use(middleware.AuthMiddleware()) // Middleware for authentication
	{
		// HoD review actions (To be added later)
	}

	r.GET("/reviewers", handler.GetAllReviewersHandler)

	r.Run(":8080")
}
