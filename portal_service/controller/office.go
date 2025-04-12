package controller

import (
	"MUJ_AMG/pkg/model"
	"MUJ_AMG/portal_service/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateOfficeHandler(c *gin.Context) {
	var req struct {
		Department string `json:"department" binding:"required"`
		Email      string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := repository.CreateOffice(req.Department, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create office: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Office created successfully"})
}

func GetOfficesHandler(c *gin.Context) {
	var filters model.GetOfficeFilters

	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offices, err := repository.GetOffices(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch offices: " + err.Error()})
		return
	}

	if len(offices) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No offices found", "offices": []model.Office{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"offices": offices})
}

func UpdateOfficeHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id query parameter is required"})
		return
	}

	var req struct {
		Department *string `json:"department"`
		Email      *string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Department == nil && req.Email == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one field must be provided for update"})
		return
	}

	err := repository.UpdateOffice(id, req.Department, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update office: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Office updated successfully"})
}

func DeleteOfficeHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id query parameter is required"})
		return
	}

	err := repository.DeleteOffice(id)
	if err != nil {
		if err.Error() == "no office found with the given ID" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete office: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Office deleted successfully"})
}
