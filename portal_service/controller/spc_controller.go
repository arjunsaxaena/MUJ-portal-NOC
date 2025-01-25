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

func CreateSpCHandler(c *gin.Context) {
	var input struct {
		Name       string `json:"name"`
		Email      string `json:"email"`
		Password   string `json:"password"`
		Department string `json:"department"`
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

	id, err := repository.CreateSpC(input.Name, input.Email, string(hashedPassword), input.Department)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create spc"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "spc created successfully"})
}

func LoginSpcHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	filters := model.GetSpCFilters{Email: input.Email}
	spcs, err := repository.GetSpCs(filters)
	if err != nil || len(spcs) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	spc := spcs[0] // incase multiple, unlikely

	if err := bcrypt.CompareHashAndPassword([]byte(spc.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	cfg, _ := config.LoadConfig()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         spc.ID,
		"email":      spc.Email,
		"department": spc.Department,
		"role":       "spc",
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

func GetSpcsHandler(c *gin.Context) {
	department := c.DefaultQuery("department", "")

	var filters model.GetSpCFilters
	if department != "" {
		filters.Department = department
	}

	spcs, err := repository.GetSpCs(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch spcs"})
		return
	}

	if len(spcs) == 0 {
		if department != "" {
			c.JSON(http.StatusNotFound, gin.H{"error": "No spcs found for the specified department"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "No spcs found"})
		}
		return
	}

	c.JSON(http.StatusOK, spcs)
}

func GetSubmissionscontroller(c *gin.Context) {
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

	submissionURL := cfg.SubmissionServiceURL + "/submissions?department=" + department.(string) + "&status=Pending"

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
