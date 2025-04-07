package controller

import (
	"MUJ_AMG/pkg/middleware"
	"MUJ_AMG/pkg/model"
	"MUJ_AMG/portal_service/config"
	"MUJ_AMG/portal_service/repository"
	submissionRepository "MUJ_AMG/submission_service/repository"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func CreateHoDHandler(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required"`
		AppPassword string `json:"app_password" binding:"required"`
		Department  string `json:"department" binding:"required"`
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

	id, err := repository.CreateHoD(input.Name, input.Email, string(hashedPassword), input.AppPassword, input.Department)
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

	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in context"})
		return
	}

	emailStr, ok := email.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	var filters model.GetSubmissionFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		log.Printf("Invalid filters: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filters"})
		return
	}

	filters.Department = department.(string)
	if emailStr == "shusheelavishnoi@gmail.com" { // All generic noc visible to shushila ma'am
		filters.NocType = "Generic"
	} else {
		filters.NocType = "Specific"
	}

	log.Printf("Filters received: %+v", filters)

	submissions, err := submissionRepository.GetSubmissions(filters)
	if err != nil {
		log.Printf("Error fetching submissions with filters %v: %v", filters, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submissions", "details": err.Error()})
		return
	}

	log.Printf("Submissions retrieved: %d records", len(submissions))
	c.JSON(http.StatusOK, gin.H{"submissions": submissions})
}

func UpdateHoDHandler(c *gin.Context) {
	var input struct {
		Password *string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "HoD ID is required"})
		return
	}

	if input.Password == nil || *input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	hashedPassword := string(hashedBytes)

	err = repository.UpdateHoD(id, hashedPassword)
	if err != nil {
		log.Printf("Error updating HoD password with ID %s: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func DeleteHoDHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "HoD ID is required"})
		return
	}

	err := repository.DeleteHoD(id)
	if err != nil {
		if err.Error() == "no HoD found with the given ID" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			log.Printf("Error deleting HoD with ID %s: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete HoD"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "HoD deleted successfully"})
}

func LogoutHodHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}

	middleware.BlacklistToken(tokenString)

	c.JSON(http.StatusOK, gin.H{"message": "HOD logged out successfully"})
}
