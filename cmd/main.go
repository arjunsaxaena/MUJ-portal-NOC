package main

import (
	"MUJ_automated_mail_generation/pkg/database"
	"MUJ_automated_mail_generation/pkg/handler"
	"MUJ_automated_mail_generation/pkg/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	r := gin.Default()

	r.Use(cors.Default())

	// Student form data collection
	r.POST("/submit", handler.SubmitHandler)

	// Reviewer-related routes
	r.POST("/reviewer", handler.CreateReviewerHandler)
	r.POST("/reviewer/login", handler.LoginReviewerHandler)
	// (1) Reviewer login credentials will be handled here. Jwt token will be outputted here.
	//     {
	//			"email": "arjunsaxena04@gmail.com",
	//			"password": "secure"
	//		}
	// 		This is how it expects the body

	r.GET("/reviewer/:department", handler.GetReviewerDetailsHandler)
	r.GET("/reviewers", handler.GetAllReviewersHandler)

	// Authenticated Reviewer routes
	authReviewer := r.Group("/reviewer")
	authReviewer.Use(middleware.AuthMiddleware("reviewer")) // Middleware for authentication
	{
		authReviewer.GET("/submissions", handler.GetSubmissionsHandler)
		// (2) On successful login, JWT token should be placed here.
		//     Reviewer should be redirected to this URL.

		authReviewer.POST("/fpc_reviews", handler.CreateReviewHandler)
		// (3) When a reviewer clicks "Approve", "Reject", or "Rework", send a JSON body like:
		//     {
		//         "submission_id": 5,
		//         "reviewer_id": 1,
		//         "status": "Approved",   // Status changes based on button clicked.
		//         "comments": "Documents verified"
		//     }
		//     POST this JSON body to "/fpc_reviews". Note: Everything inside authReviewer will require the jwt token at login to be carry forward.

		authReviewer.PUT("/fpc_reviews/:id", handler.UpdateReviewHandler)                                 // Dont implement, only for testing
		authReviewer.GET("/fpc_reviews", handler.GetAllReviewerReviewsHandler)                            // Dont implement, only for testing
		authReviewer.GET("/fpc_reviews/submission/:submission_id", handler.GetReviewsBySubmissionHandler) // Dont implement, only for testing
		authReviewer.GET("/fpc_reviews/reviewer/:reviewer_id", handler.GetReviewsByReviewerHandler)       // Dont implement, only for testing
	}

	// HoD-related routes
	r.POST("/hod/login", handler.LoginHoDHandler)
	// (4) Hod login credentials will be handled here. Jwt token will be outputted here.
	//     {
	//			"email": "arjunsaxena04@gmail.com",
	//			"password": "secure"
	//		}
	// 		This is how it expects the body

	r.POST("/hod", handler.CreateHoDHandler)  // For creating Hod, most probably only for testing locally
	r.GET("/hods", handler.GetAllHoDsHandler) // Dont implement, only for testing

	r.GET("/hod/:department", handler.GetHoDsByDepartmentHandler) // Dont implement, only for testing

	// Authenticated HoD routes
	authHoD := r.Group("/hod")
	authHoD.Use(middleware.AuthMiddleware("hod")) // Middleware for authentication
	{
		authHoD.GET("/submissions/approved", handler.GetApprovedSubmissionsByDepartmentHandler)
		// (5) On successful login of Hod, JWT token should be placed here.
		//     Hod should be redirected to this URL.

		authHoD.POST("/hod_reviews", handler.CreateHodReviewHandler)
		// (3) When hod clicks "Approve", "Reject", or "Rework", send a JSON body like:
		//     {
		//         "submission_id": 5,
		//         "hod_id": 1,
		//         "action": "Approved",   // Action changes based on button clicked.
		//         "remarks": "Everything looks good, approved for NOC." // Remarks based on the status.
		//     }
		//     POST this JSON body to "/hod_reviews". Note: Everything inside authReviewer will require the jwt token at login to be carry forward.
		//     Make sure the status in reviewer or action in hod only posts status as
		//     "Approve", "Reject", or "Rework" according to the button clicked. This is case sensitive.

		authHoD.GET("/hod_reviews", handler.GetAllHodReviewsHandler) // Dont implement, only for testing
	}

	// Start the server
	r.Run(":8080")
}
