package main

import (
	"MUJ_AMG/pkg/database"
	"MUJ_AMG/submission_service/config"
	"MUJ_AMG/submission_service/controller"
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

	// csvFile := "../Students_VIII.csv"
	// err = util.ImportCSVToPostgres(csvFile, database.DB)
	// if err != nil {
	// 	log.Fatalf("Failed to import CSV to PostgreSQL: %v", err)
	// }
	// log.Println("CSV data imported successfully!")

	// csvFile2 := "../Students_VI.csv"
	// err = util.ImportCSVToPostgres(csvFile2, database.DB)
	// if err != nil {
	// 	log.Fatalf("Failed to import second CSV to PostgreSQL: %v", err)
	// }
	// log.Println("Second CSV data imported successfully!")

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "ngrok-skip-browser-warning"},
		AllowCredentials: true,
	}))

	r.POST("/generate-otp", controller.GenerateOTPHandler)
	r.POST("/validate-otp", controller.ValidateOTPHandler)
	r.POST("/submit", controller.SubmitHandler)
	//r.GET("/submissions", controller.GetSubmissionsHandler)
	//r.PUT("/submissions", controller.UpdateSubmissionStatusHandler)

	log.Printf("Starting server on port %s", cfg.Port)
	r.Run("0.0.0.0:" + cfg.Port)
}
