package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Memecoin struct {
	ID              int        `json:"id" gorm:"primaryKey"`
	Name            string     `json:"name" gorm:"unique; not null"`
	Description     string     `json:"description"`
	CreatedAt       time.Time  `json:"created_at"`
	PopularityScore int        `json:"popularity_score"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
}

var db *gorm.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func initDB() {
	var err error
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")

	dsn := " host=" + host + " user=" + user + " password=" + password + " dbname=" + dbName + " port=5432 sslmode=disable"

	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("failed to connect database, retrying in 2 seconds... (%d/5)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("failed to connect database after retries: %v", err)
	}
	db.AutoMigrate(&Memecoin{})
}

func createMemecoin(c *gin.Context) {
	var memeCoin Memecoin
	if err := c.ShouldBindJSON(&memeCoin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if memeCoin.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	memeCoin.CreatedAt = time.Now()

	if err := db.Where("name = ?", memeCoin.Name).FirstOrCreate(&memeCoin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if memeCoin.ID != 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Memecoin with this name already exists"})
		return
	}

	c.JSON(http.StatusCreated, memeCoin)
}

func getMemecoin(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var memeCoin Memecoin
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

func updateMemecoin(c *gin.Context) {
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

	result := db.Model(&Memecoin{}).Where("id = ? AND deleted_at IS NULL", id).Updates(Memecoin{Description: input.Description})

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

func deleteMemecoin(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	now := time.Now()
	result := db.Model(&Memecoin{}).Where("id = ? AND deleted_at IS NULL", id).Update("deleted_at", now)

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

func pokeMemecoin(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var memeCoin Memecoin
	tx := db.Begin()

	if err := tx.Where("id = ? AND deleted_at IS NULL", id).First(&memeCoin).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Memecoin not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	memeCoin.PopularityScore++

	if err := tx.Save(&memeCoin).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, memeCoin)
}

func main() {
	initDB()
	r := gin.Default()

	r.POST("/memecoins", createMemecoin)
	r.GET("/memecoins/:id", getMemecoin)
	r.PATCH("/memecoins/:id", updateMemecoin)
	r.DELETE("/memecoins/:id", deleteMemecoin)
	r.POST("/memecoins/:id/poke", pokeMemecoin)

	r.Run(":8080")
}
