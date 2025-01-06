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
	ID              int       `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"unique; not null"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	PopularityScore int       `json:"popularity_score"`
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

	var existingMemeCoin Memecoin
	err := db.Where("name = ?", memeCoin.Name).First(&existingMemeCoin).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Memecoin with this name already exists"})
		return
	}

	memeCoin.CreatedAt = time.Now()
	memeCoin.PopularityScore = 0

	if err := db.Create(&memeCoin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	if err := db.First(&memeCoin, id).Error; err != nil {
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
		Description     string     `json:"description"`
		Name            *string    `json:"name,omitempty"`
		CreatedAt       *time.Time `json:"created_at,omitempty"`
		PopularityScore *int       `json:"popularity_score,omitempty"`
	}

	var memeCoin Memecoin
	if err := db.First(&memeCoin, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != nil || input.CreatedAt != nil || input.PopularityScore != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can only update the description."})
		return
	}

	memeCoin.Description = input.Description
	db.Save(&memeCoin)
	c.JSON(http.StatusOK, memeCoin)
}

func deleteMemecoin(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var memeCoin Memecoin
	if err := db.First(&memeCoin, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.Delete(&memeCoin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
	c.JSON(http.StatusOK, gin.H{"message": "Memecoin deleted successfully!"})
}

func pokeMemecoin(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	var memeCoin Memecoin
	if err := db.First(&memeCoin, id).Error; err != nil {
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

func main() {
	initDB()
	r := gin.Default()

	r.POST("/memecoins", createMemecoin)
	r.GET("/memecoins/:id", getMemecoin)
	r.PUT("/memecoins/:id", updateMemecoin)
	r.DELETE("/memecoins/:id", deleteMemecoin)
	r.PATCH("/memecoins/:id", pokeMemecoin)

	r.Run(":8080")
}
