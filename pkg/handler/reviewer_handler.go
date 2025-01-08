package handler

import (
	"MUJ_automated_mail_generation/pkg/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func CreateReviewerHandler(c *gin.Context) {
	var input struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		Department string `json:"department"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	id, err := database.CreateReviewer(input.Username, string(hashedPassword), input.Department)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reviewer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Reviewer created successfully"})
}

var jwtSecretKey = []byte("your_secret_key") // Secret key for signing the JWT token

// LoginReviewerHandler handles reviewer login.
func LoginReviewerHandler(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Fetch the reviewer from the database
	reviewer, err := database.GetReviewerByUsername(input.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Validate the password
	if err := bcrypt.CompareHashAndPassword([]byte(reviewer.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":   reviewer.Username,
		"department": reviewer.Department,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // Expire in 24 hours
	})

	// Sign the token
	tokenString, err := token.SignedString(jwtSecretKey)
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
	username := c.Param("username")

	reviewer, err := database.GetReviewerByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reviewer not found"})
		return
	}

	c.JSON(http.StatusOK, reviewer)
}

func GetAllReviewersHandler(c *gin.Context) {
	reviewers, err := database.GetAllReviewers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviewers"})
		return
	}

	c.JSON(http.StatusOK, reviewers)
}
