package handler

import (
	"MUJ_automated_mail_generation/pkg/config"
	"MUJ_automated_mail_generation/pkg/database"
	"MUJ_automated_mail_generation/pkg/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func CreateReviewerHandler(c *gin.Context) {
	var input struct {
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

	id, err := database.CreateReviewer(input.Email, string(hashedPassword), input.Department)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reviewer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Reviewer created successfully"})
}

func LoginReviewerHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	reviewer, err := database.GetReviewerByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(reviewer.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         reviewer.ID, // Include reviewer ID for identification
		"email":      reviewer.Email,
		"department": reviewer.Department, // Include department in claims
		"role":       "reviewer",
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // Expiration time (24 hours)
	})

	// Sign the token
	tokenString, err := token.SignedString(config.JwtSecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Respond with the token
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   tokenString,
	})
}

func GetReviewerDetailsHandler(c *gin.Context) {
	department := c.Param("department")

	reviewers, err := database.GetReviewerByDepartment(department)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reviewers not found"})
		return
	}

	c.JSON(http.StatusOK, reviewers)
}

func GetAllReviewersHandler(c *gin.Context) {
	reviewers, err := database.GetAllReviewers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviewers"})
		return
	}

	c.JSON(http.StatusOK, reviewers)
}

func GetReviewerByEmail(email string) (model.Reviewer, error) {
	var reviewer model.Reviewer
	err := database.DB.Get(&reviewer, "SELECT * FROM reviewers WHERE email = $1", email)
	if err != nil {
		return reviewer, fmt.Errorf("unable to fetch reviewer: %w", err)
	}
	return reviewer, nil
}
