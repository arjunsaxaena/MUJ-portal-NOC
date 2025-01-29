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

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/files", "./uploads")
	r.POST("/submit", controller.SubmitHandler)
	//r.GET("/submissions", controller.GetSubmissionsHandler)
	//r.PUT("/submissions", controller.UpdateSubmissionStatusHandler)

	log.Printf("Starting server on port %s", cfg.Port)
	r.Run("0.0.0.0:" + cfg.Port)
}
