package controller

import (
	"MUJ_AMG/pkg/model"
	"MUJ_AMG/portal_service/config"
	"MUJ_AMG/portal_service/repository"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func CreateHoDHandler(c *gin.Context) {
	var input struct {
		Name       string `json:"name" binding:"required"`
		Email      string `json:"email" binding:"required,email"`
		Password   string `json:"password" binding:"required,min=6"`
		Department string `json:"department" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	id, err := repository.CreateHoD(input.Name, input.Email, string(hashedPassword), input.Department)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create HoD"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "HoD created successfully"})
}

func LoginHoDHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	filters := model.GetHoDFilters{Email: input.Email}
	hods, err := repository.GetHoDs(filters)
	if err != nil || len(hods) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	hod := hods[0] // Assume single match (unlikely to have duplicates)

	if err := bcrypt.CompareHashAndPassword([]byte(hod.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	cfg, _ := config.LoadConfig()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         hod.ID,
		"email":      hod.Email,
		"department": hod.Department,
		"role":       "hod",
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.JwtSecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})
}

func GetHoDsHandler(c *gin.Context) {
	department := c.DefaultQuery("department", "")

	var filters model.GetHoDFilters
	if department != "" {
		filters.Department = department
	}

	hods, err := repository.GetHoDs(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch HoDs"})
		return
	}

	if len(hods) == 0 {
		if department != "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "No HoDs found for the specified department"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "No HoDs found"})
		}
		return
	}

	c.JSON(http.StatusOK, hods)
}

func GetSubmissionsForHoDcontroller(c *gin.Context) {
	department, exists := c.Get("department")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Department not found in context"})
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load configuration"})
		return
	}

	submissionURL := cfg.SubmissionServiceURL + "/submissions?department=" + department.(string) + "&status=Approved"

	resp, err := http.Get(submissionURL)
	if err != nil {
		log.Printf("Error calling submission service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submissions"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %d", resp.StatusCode)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submissions"})
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}
	// log.Printf("Response from submission service: %s", body)

	type SubmissionsResponse struct {
		Data []model.StudentSubmission `json:"submissions"`
	}

	var result SubmissionsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error decoding response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"submissions": result.Data})
}
