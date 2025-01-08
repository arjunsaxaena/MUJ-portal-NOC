package handler

import (
	"MUJ_automated_mail_generation/pkg/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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

	// Create the reviewer
	id, err := database.CreateReviewer(input.Username, string(hashedPassword), input.Department)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reviewer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Reviewer created successfully"})
}

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

	// Fetch the reviewer
	reviewer, err := database.GetReviewerByUsername(input.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	fmt.Printf("Fetched reviewer for login: %+v\n", reviewer) // Debug log

	// Validate the password
	if err := bcrypt.CompareHashAndPassword([]byte(reviewer.PasswordHash), []byte(input.Password)); err != nil {
		fmt.Printf("Password comparison failed: %v\n", err) // Debug log
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// GetReviewerDetailsHandler handles fetching reviewer details.
func GetReviewerDetailsHandler(c *gin.Context) {
	username := c.Param("username")

	// Call the database function
	reviewer, err := database.GetReviewerByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reviewer not found"})
		return
	}

	c.JSON(http.StatusOK, reviewer)
}
