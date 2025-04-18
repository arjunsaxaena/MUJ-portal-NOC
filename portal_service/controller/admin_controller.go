package controller

import (
	"MUJ_AMG/pkg/middleware"
	"MUJ_AMG/pkg/model"
	"MUJ_AMG/portal_service/config"
	"MUJ_AMG/portal_service/repository"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func CreateAdminHandler(c *gin.Context) {
	var input struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		AppPassword string `json:"app_password" binding:"required"`
	}

	log.Printf("Received request to create admin with email: %s", input.Email)

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	log.Printf("Input validated: Name: %s, Email: %s", input.Name, input.Email)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err) // Log the password hashing error
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	log.Println("Password hashed successfully")

	id, err := repository.CreateAdmin(input.Name, input.Email, string(hashedPassword), input.AppPassword)
	if err != nil {
		log.Printf("Error creating admin: %v", err) // Log error while creating admin
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin"})
		return
	}

	log.Printf("Admin created successfully with ID: %s", id)

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Admin created successfully"})
}

func LoginAdminHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	filters := model.GetAdminFilters{Email: input.Email}
	admins, err := repository.GetAdmins(filters)
	if err != nil || len(admins) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	admin := admins[0] // Assuming there's only one admin with the email

	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	cfg, _ := config.LoadConfig()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    admin.ID,
		"email": admin.Email,
		"role":  "admin",
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
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

func GetAdminsHandler(c *gin.Context) {
	email := c.DefaultQuery("email", "")

	var filters model.GetAdminFilters
	if email != "" {
		filters.Email = email
	}

	admins, err := repository.GetAdmins(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch admins"})
		return
	}

	if len(admins) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No admins found"})
		return
	}

	c.JSON(http.StatusOK, admins)
}

func LogoutAdminHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}

	middleware.BlacklistToken(tokenString)

	c.JSON(http.StatusOK, gin.H{"message": "Admin logged out successfully"})
}

func UpdateAdminHandler(c *gin.Context) {
	var input struct {
		Name     *string `json:"name"`
		Password *string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	id := c.Query("id") // Use string directly for UUID
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required in the query parameter"})
		return
	}

	if input.Name == nil && input.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one field (name or password) must be provided"})
		return
	}

	var hashedPassword *string
	if input.Password != nil {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		hashedStr := string(hashed)
		hashedPassword = &hashedStr
	}

	err := repository.UpdateAdmin(id, input.Name, hashedPassword) // Pass string ID
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin updated successfully"})
}
