package handler

import (
	"MUJ_automated_mail_generation/pkg/config"
	"MUJ_automated_mail_generation/pkg/database"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// CreateHoDHandler handles creating a new HoD
func CreateHoDHandler(c *gin.Context) {
	var input struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		Department string `json:"department"`
	}

	// Bind input data from the request body
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

	// Create HoD in the database
	id, err := database.CreateHoD(input.Email, string(hashedPassword), input.Department)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create HoD"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "HoD created successfully"})
}

// LoginHoDHandler handles HoD login and token generation
func LoginHoDHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Bind input data from the request body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Fetch HoD from the database
	hod, err := database.GetHoDByEmail(input.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Validate the password
	if err := bcrypt.CompareHashAndPassword([]byte(hod.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":      hod.Email,
		"department": hod.Department,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // Expire in 24 hours
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

// GetHoDDetailsHandler retrieves HoD details based on department
func GetHoDDetailsHandler(c *gin.Context) {
	department := c.Param("department")

	// Fetch HoDs for a specific department
	hods, err := database.GetHoDsByDepartment(department)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "HoDs not found"})
		return
	}

	// Return HoDs as JSON response
	c.JSON(http.StatusOK, hods)
}

// GetAllHoDsHandler retrieves all HoDs from the database
func GetAllHoDsHandler(c *gin.Context) {
	// Fetch all HoDs from the database
	hods, err := database.GetAllHoDs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch HoDs"})
		return
	}

	// Return all HoDs as JSON response
	c.JSON(http.StatusOK, hods)
}

// GetHoDsByDepartmentHandler retrieves HoDs by department.
func GetHoDsByDepartmentHandler(c *gin.Context) {
	department := c.Param("department")

	hoDs, err := database.GetHoDsByDepartment(department)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "HoDs not found in department"})
		return
	}

	c.JSON(http.StatusOK, hoDs)
}
