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

	r.POST("/spc", controller.CreateSpCHandler)
	r.POST("/spc/login", controller.LoginSpcHandler)

	//[GIN] 2025/01/16 - 11:33:19 | 204 |            0s |             ::1 | OPTIONS  "/spc/login"      asking cors for persmission
	//Current working directory:
	//[GIN] 2025/01/16 - 11:33:19 | 200 |     71.4962ms |             ::1 | POST     "/spc/login"		post request

	// (1) spc login credentials will be handled here. Jwt token will be outputted here.
	//     {
	//			"email": "arjunsaxena04@gmail.com",
	//			"password": "secure"
	//		}
	// 		This is how it expects the body

	r.GET("/spcs", controller.GetSpcsHandler)

	authSpc := r.Group("/spc")
	authSpc.Use(middleware.AuthMiddleware(cfg.JwtSecretKey, "spc"))
	{
		authSpc.GET("/submissions", controller.GetSubmissionscontroller)
		// (2) On successful login, JWT token should be placed here.
		//     spc should be redirected to this URL.

		authSpc.POST("/spc_reviews", controller.CreateSpcReviewHandler)
		// (3) When a spc clicks "Approve", "Reject", or "Rework", send a JSON body like:
		//     {
		//         "submission_id": 5,
		//         "spc_id": 1,
		//         "status": "Approved",   // Status changes based on button clicked.
		//         "comments": "Documents verified"
		//     }
		//     POST this JSON body to "/spc_reviews". Note: Everything inside authSpc will require the jwt token at login to be carry forward.

		authSpc.PATCH("/spc_reviews", controller.UpdateSpcReviewHandler) // If fpc wants to reject or rework an approved submission
		authSpc.GET("/spc_reviews", controller.GetSpcReviewsHandler)     // For testing
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
		//     POST this JSON body to "/hod_reviews". Note: Everything inside authspc will require the jwt token at login to be carry forward.
		//     Make sure the status in spc or action in hod only posts status as
		//     "Approve", "Reject", or "Rework" according to the button clicked. This is case sensitive.

		authHoD.GET("/hod_reviews", controller.GetHodReviewsHandler) // For testing
	}

	log.Printf("Server starting on port %s...", cfg.Port)
	r.Run(":" + cfg.Port)
}
