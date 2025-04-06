package main

import (
	"MUJ_AMG/pkg/database"
	"MUJ_AMG/pkg/middleware"
	"MUJ_AMG/portal_service/config"
	"MUJ_AMG/portal_service/controller"
	"log"
	"path/filepath"

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
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "ngrok-skip-browser-warning"},
		AllowCredentials: true,
	}))

	// Submission Service main.go code here so that I dont need to make two seperate ngrok tunnels for two seperate ports during testing

	/////////////////////////////////////////////////////////////////////////////////////////////////////////

	// Importing student data in postgres table for extra validation

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

	// r.POST("/generate-otp", submission_controller.GenerateOTPHandler)
	// r.POST("/validate-otp", submission_controller.ValidateOTPHandler)
	// r.POST("/submit", submission_controller.SubmitHandle

	/////////////////////////////////////////////////////////////////////////////////////////////////////////

	authFiles := r.Group("/files")
	authFiles.Use(middleware.AuthMiddleware(cfg.JwtSecretKey, "fpc", "hod"))
	{
		authFiles.Static("/", filepath.Join("..", "uploads"))
	}

	r.POST("/admin", controller.CreateAdminHandler)
	r.POST("/admin/login", controller.LoginAdminHandler)

	authAdmin := r.Group("/admin")
	authAdmin.Use(middleware.AuthMiddleware(cfg.JwtSecretKey, "admin"))
	{
		authAdmin.POST("/logout", controller.LogoutAdminHandler)

		authAdmin.POST("/fpc", controller.CreateFpCHandler)
		authAdmin.POST("/hod", controller.CreateHoDHandler)

		authAdmin.GET("/fpcs", controller.GetFpcsHandler)
		authAdmin.GET("/hods", controller.GetHoDsHandler)

		authAdmin.DELETE("/fpc", controller.DeleteFpCHandler)
		authAdmin.DELETE("/hod", controller.DeleteHoDHandler)

		authAdmin.PATCH("/update", controller.UpdateAdminHandler)
	}

	r.POST("/fpc/login", controller.LoginFpcHandler)

	authFpc := r.Group("/fpc")
	authFpc.Use(middleware.AuthMiddleware(cfg.JwtSecretKey, "fpc"))
	{
		authFpc.PATCH("/", controller.UpdateFpCHandler)
		authFpc.POST("/logout", controller.LogoutFpcHandler)
		authFpc.GET("/submissions", controller.GetSubmissionscontroller)
		authFpc.POST("/fpc_reviews", controller.CreateFpcReviewHandler)
	}

	r.POST("/hod/login", controller.LoginHoDHandler)

	authHoD := r.Group("/hod")
	authHoD.Use(middleware.AuthMiddleware(cfg.JwtSecretKey, "hod"))
	{
		authHoD.PATCH("/", controller.UpdateHoDHandler)
		authHoD.POST("/logout", controller.LogoutHodHandler)
		authHoD.GET("/submissions", controller.GetSubmissionsForHoDcontroller)
		authHoD.POST("/hod_reviews", controller.CreateHodReviewHandler)
	}

	log.Printf("Server starting on port %s...", cfg.Port)
	r.Run("0.0.0.0:" + cfg.Port)
}
