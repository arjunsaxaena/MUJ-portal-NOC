package main

import (
	"MUJ_AMG/pkg/database"
	"MUJ_AMG/pkg/middleware"
	"MUJ_AMG/portal_service/config"
	"MUJ_AMG/portal_service/controller"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database.Connect(cfg)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},         // Frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS"}, // Allowed methods
		AllowHeaders:     []string{"Content-Type", "Authorization"}, // Allowed headers
		AllowCredentials: true,                                      // Allow cookies or credentials
	}))

	r.POST("/reviewer", controller.CreateReviewerHandler)
	r.POST("/reviewer/login", controller.LoginReviewerHandler)

	//[GIN] 2025/01/16 - 11:33:19 | 204 |            0s |             ::1 | OPTIONS  "/reviewer/login"      asking cors for persmission
	//Current working directory:
	//[GIN] 2025/01/16 - 11:33:19 | 200 |     71.4962ms |             ::1 | POST     "/reviewer/login"		post request

	// (1) Reviewer login credentials will be handled here. Jwt token will be outputted here.
	//     {
	//			"email": "arjunsaxena04@gmail.com",
	//			"password": "secure"
	//		}
	// 		This is how it expects the body

	r.GET("/reviewers", controller.GetReviewersHandler)

	authReviewer := r.Group("/reviewer")
	authReviewer.Use(middleware.AuthMiddleware(cfg.JwtSecretKey, "reviewer"))
	{
		authReviewer.GET("/submissions", controller.GetSubmissionscontroller)
		// (2) On successful login, JWT token should be placed here.
		//     Reviewer should be redirected to this URL.

		authReviewer.POST("/fpc_reviews", controller.CreateReviewHandler)
		// (3) When a reviewer clicks "Approve", "Reject", or "Rework", send a JSON body like:
		//     {
		//         "submission_id": 5,
		//         "reviewer_id": 1,
		//         "status": "Approved",   // Status changes based on button clicked.
		//         "comments": "Documents verified"
		//     }
		//     POST this JSON body to "/fpc_reviews". Note: Everything inside authReviewer will require the jwt token at login to be carry forward.

		authReviewer.PUT("/fpc_reviews", controller.UpdateReviewHandler) // If fpc wants to reject or rework an approved submission
		authReviewer.GET("/fpc_reviews", controller.GetReviewsHandler)   // For testing
	}

	r.POST("/hod/login", controller.LoginHoDHandler)
	// (4) Hod login credentials will be handled here. Jwt token will be outputted here.
	//     {
	//			"email": "arjunsaxena04@gmail.com",
	//			"password": "secure"
	//		}
	// 		This is how it expects the body

	r.POST("/hod", controller.CreateHoDHandler) // For creating HoD (used during testing)
	r.GET("/hods", controller.GetHoDsHandler)   // For testing

	authHoD := r.Group("/hod")
	authHoD.Use(middleware.AuthMiddleware(cfg.JwtSecretKey, "hod"))
	{
		authHoD.GET("/approved_submissions", controller.GetSubmissionsForHoDcontroller)
		// (5) On successful login of Hod, JWT token should be placed here.
		//     Hod should be redirected to this URL.

		authHoD.POST("/hod_reviews", controller.CreateHodReviewHandler)
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

		authHoD.GET("/hod_reviews", controller.GetHodReviewsHandler) // For testing
	}

	log.Printf("Server starting on port %s...", cfg.Port)
	r.Run(":" + cfg.Port)
}
