package service

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"memecoin_homework/internal/model"
	"memecoin_homework/internal/utils"
)

var db *gorm.DB

func SetDB(database *gorm.DB) {
	db = database
}

func CreateMemecoin(c *gin.Context, db *gorm.DB) {
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

	result := db.Create(&memeCoin)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Memecoin with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, memeCoin)
}

func GetMemecoin(c *gin.Context, db *gorm.DB) {
	id, ok := utils.ValidateID(c)
	if !ok {
		return
	}

	var memeCoin model.Memecoin
	if err := db.Where("id = ? AND deleted IS NULL", id).First(&memeCoin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Memecoin not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, memeCoin)
}

func UpdateMemecoin(c *gin.Context, db *gorm.DB) {
	id, ok := utils.ValidateID(c)
	if !ok {
		return
	}

	var input struct {
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.Model(&model.Memecoin{}).Where("id = ? AND deleted IS NULL", id).Updates(model.Memecoin{Description: input.Description})

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

func DeleteMemecoin(c *gin.Context, db *gorm.DB) {
	id, ok := utils.ValidateID(c)
	if !ok {
		return
	}

	result := db.Delete(&model.Memecoin{}, id)

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

func PokeMemecoin(c *gin.Context, db *gorm.DB) {
	id, ok := utils.ValidateID(c)
	if !ok {
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		result := tx.Exec("UPDATE memecoins SET popularity_score = popularity_score + 1 WHERE id = ? AND deleted IS NULL", id)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return result.Error
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Memecoin not found"})
			return gorm.ErrRecordNotFound
		}

		c.JSON(http.StatusOK, gin.H{"message": "Memecoin popularity increased successfully"})
		return nil
	})

	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
