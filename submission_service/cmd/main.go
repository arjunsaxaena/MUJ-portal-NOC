package main

import (
	"MUJ_AMG/pkg/database"
	"MUJ_AMG/submission_service/config"
	"MUJ_AMG/submission_service/controller"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database.Connect(cfg)

	r := gin.Default()

	r.POST("/submit", controller.SubmitHandler)
	r.GET("/submissions", controller.GetSubmissionsHandler)
	r.PUT("/submissions", controller.UpdateSubmissionStatusHandler)

	log.Printf("Starting server on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
