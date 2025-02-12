package main

import (
	"MUJ_AMG/pkg/database"
	"MUJ_AMG/pkg/middleware"
	"MUJ_AMG/portal_service/config"
	"MUJ_AMG/portal_service/controller"
	submission_controller "MUJ_AMG/submission_service/controller"
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
		AllowOrigins:     []string{"https://student-portal-nu-one.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "ngrok-skip-browser-warning"},
		AllowCredentials: true,
	}))

	// csvFile := "/home/ubuntu/MUJ_automated_mail_generation/Students_VIII.csv"
	// err = util.ImportCSVToPostgres(csvFile, database.DB)
	// if err != nil {
	// 	log.Fatalf("Failed to import CSV to PostgreSQL: %v", err)
	// }
	// log.Println("CSV data imported successfully!")

	// csvFile2 := "/home/ubuntu/MUJ_automated_mail_generation/Students_VI.csv"
	// err = util.ImportCSVToPostgres(csvFile2, database.DB)
	// if err != nil {
	// 	log.Fatalf("Failed to import second CSV to PostgreSQL: %v", err)
	// }
	// log.Println("Second CSV data imported successfully!")

	// Submission service
	r.POST("/generate-otp", submission_controller.GenerateOTPHandler)
	r.POST("/validate-otp", submission_controller.ValidateOTPHandler)
	r.POST("/submit", submission_controller.SubmitHandler)

	// Files serving
	r.Static("/files", "../uploads")

	r.POST("/admin", controller.CreateAdminHandler)
	r.POST("/admin/login", controller.LoginAdminHandler)

	authAdmin := r.Group("/admin")
	authAdmin.Use(middleware.AuthMiddleware(cfg.JwtSecretKey, "admin"))
	{
		authAdmin.POST("/fpc", controller.CreateFpCHandler)
		authAdmin.POST("/hod", controller.CreateHoDHandler)

		authAdmin.GET("/fpcs", controller.GetFpcsHandler)
		authAdmin.GET("/hods", controller.GetHoDsHandler)

		authAdmin.PATCH("/fpc", controller.UpdateFpCHandler)
		authAdmin.PATCH("/hod", controller.UpdateHoDHandler)

		authAdmin.DELETE("/fpc", controller.DeleteFpCHandler)
		authAdmin.DELETE("/hod", controller.DeleteHoDHandler)
	}

	r.POST("/fpc/login", controller.LoginFpcHandler)

	//[GIN] 2025/01/16 - 11:33:19 | 204 |            0s |             ::1 | OPTIONS  "/fpc/login"      asking cors for persmission
	//[GIN] 2025/01/16 - 11:33:19 | 200 |     71.4962ms |             ::1 | POST     "/fpc/login"		post request

	// (1) fpc login credentials will be handled here. Jwt token will be outputted here.
	//     {
	//			"email": "arjunsaxena04@gmail.com",
	//			"password": "secure"
	//		}
	// 		This is how it expects the body

	authFpc := r.Group("/fpc")
	authFpc.Use(middleware.AuthMiddleware(cfg.JwtSecretKey, "fpc"))
	{
		authFpc.GET("/submissions", controller.GetSubmissionscontroller)
		// (2) On successful login, JWT token should be placed here.
		//     fpc should be redirected to this URL.

		authFpc.POST("/fpc_reviews", controller.CreateFpcReviewHandler)
		// (3) When a fpc clicks "Approve", "Reject", or "Rework", send a JSON body like:
		//     {
		//         "submission_id": 5,
		//         "fpc_id": 1,
		//         "status": "Approved",   // Status changes based on button clicked.
		//         "comments": "Documents verified"
		//     }
		//     POST this JSON body to "/fpc_reviews". Note: Everything inside authFpc will require the jwt token at login to be carry forward.

		// authFpc.PATCH("/fpc_reviews", controller.UpdateFpcReviewHandler) // If fpc wants to reject or rework an approved submission
		// authFpc.GET("/fpc_reviews", controller.GetFpcReviewsHandler)     // For testing
	}

	r.POST("/hod/login", controller.LoginHoDHandler)
	// (4) Hod login credentials will be handled here. Jwt token will be outputted here.
	//     {
	//			"email": "arjunsaxena04@gmail.com",
	//			"password": "secure"
	//		}
	// 		This is how it expects the body

	authHoD := r.Group("/hod")
	authHoD.Use(middleware.AuthMiddleware(cfg.JwtSecretKey, "hod"))
	{
		authHoD.GET("/submissions", controller.GetSubmissionsForHoDcontroller)
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
		//     POST this JSON body to "/hod_reviews". Note: Everything inside authFpc will require the jwt token at login to be carry forward.
		//     Make sure the status in fpc or action in hod only posts status as
		//     "Approve", "Reject", or "Rework" according to the button clicked. This is case sensitive.

		// authHoD.GET("/hod_reviews", controller.GetHodReviewsHandler) // For testing
	}

	log.Printf("Server starting on port %s...", cfg.Port)
	r.Run("0.0.0.0:" + cfg.Port)
}
