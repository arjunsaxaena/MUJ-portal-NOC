package main

import (
	"MUJ_AMG/pkg/database"
	"MUJ_AMG/submission_service/config"
	"MUJ_AMG/submission_service/controller"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database.Connect(cfg)

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

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // frontend url
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/generate-otp", controller.GenerateOTPHandler)
	r.POST("/validate-otp", controller.ValidateOTPHandler)
	r.POST("/submit", controller.SubmitHandler)
	//r.GET("/submissions", controller.GetSubmissionsHandler)
	//r.PUT("/submissions", controller.UpdateSubmissionStatusHandler)

	log.Printf("Starting server on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
