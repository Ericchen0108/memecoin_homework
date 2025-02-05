package service

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"memecoin_homework/internal/model"
)

var db *gorm.DB

func SetDB(database *gorm.DB) {
	db = database
}

func CreateMemecoin(c *gin.Context) {
	var memeCoin model.Memecoin
	if err := c.ShouldBindJSON(&memeCoin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if strings.TrimSpace(memeCoin.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required and cannot be empty."})
		return
	}

	memeCoin.CreatedAt = time.Now()

	result := db.Where("name = ?", memeCoin.Name).FirstOrCreate(&memeCoin)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 1 {
		c.JSON(http.StatusCreated, memeCoin)
	} else {
		c.JSON(http.StatusConflict, gin.H{"error": "Memecoin with this name already exists"})
	}
}

func GetMemecoin(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var memeCoin model.Memecoin
	if err := db.Where("id = ? AND deleted_at IS NULL", id).First(&memeCoin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Memecoin not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, memeCoin)
}

func UpdateMemecoin(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var input struct {
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Model(&model.Memecoin{}).Where("id = ? AND deleted_at IS NULL", id).Updates(model.Memecoin{Description: input.Description})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Memecoin not found"})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Memecoin updated successfully."})
}

func DeleteMemecoin(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	now := time.Now()
	result := db.Model(&model.Memecoin{}).Where("id = ? AND deleted_at IS NULL", id).Update("deleted_at", now)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Memecoin not found"})
		return
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Memecoin is deleted."})
}

func PokeMemecoin(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var memeCoin model.Memecoin
	if err := db.Where("id = ? AND deleted_at IS NULL", id).First(&memeCoin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Memecoin not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	memeCoin.PopularityScore++

	if err := db.Save(&memeCoin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, memeCoin)
}
