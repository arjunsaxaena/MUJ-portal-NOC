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

func CreateFpCHandler(c *gin.Context) {
	var input struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		AppPassword string `json:"app_password"`
		RoleType    string `json:"role_type"`
		Department  string `json:"department"`
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

	id, err := repository.CreateFpC(input.Name, input.Email, string(hashedPassword), input.AppPassword, input.RoleType, input.Department)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create fpc"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "fpc created successfully"})
}

func LoginFpcHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	filters := model.GetFpCFilters{Email: input.Email}
	fpcs, err := repository.GetFpCs(filters)
	if err != nil || len(fpcs) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	fpc := fpcs[0]

	if err := bcrypt.CompareHashAndPassword([]byte(fpc.PasswordHash), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	cfg, _ := config.LoadConfig()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         fpc.ID,
		"email":      fpc.Email,
		"department": fpc.Department,
		"role":       "fpc",
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

func GetFpcsHandler(c *gin.Context) {
	id := c.Query("id")
	department := c.DefaultQuery("department", "")
	email := c.Query("email")

	var filters model.GetFpCFilters
	if id != "" {
		filters.ID = id
	}
	if department != "" {
		filters.Department = department
	}
	if email != "" {
		filters.Email = email
	}

	fpcs, err := repository.GetFpCs(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch fpcs"})
		return
	}

	if len(fpcs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No fpcs found matching the query"})
		return
	}

	c.JSON(http.StatusOK, fpcs)
}

func GetSubmissionscontroller(c *gin.Context) {
	department, exists := c.Get("department")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Department not found in context"})
		return
	}

	var filters model.GetSubmissionFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		log.Printf("Invalid filters: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filters"})
		return
	}

	filters.Department = department.(string)
	filters.NocType = "Specific"
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

func UpdateFpCHandler(c *gin.Context) {
	var input struct {
		Password *string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "FPC ID is required"})
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

	err = repository.UpdateFpC(id, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func DeleteFpCHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fpc ID is required"})
		return
	}

	err := repository.DeleteFpC(id)
	if err != nil {
		if err.Error() == "no reviewer found with the given ID" {
			c.JSON(http.StatusNotFound, gin.H{"error": "No fpc found with the given ID"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete fpc", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "fpc deleted successfully"})
}

func LogoutFpcHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}

	middleware.BlacklistToken(tokenString)

	c.JSON(http.StatusOK, gin.H{"message": "FPC logged out successfully"})
}
