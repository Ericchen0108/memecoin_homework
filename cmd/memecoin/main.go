package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"memecoin_homework/api"
	"memecoin_homework/internal/model"
	"memecoin_homework/internal/service"
)

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

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbName + " port=5432 sslmode=disable"

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
	db.AutoMigrate(&model.Memecoin{})
}

func main() {
	initDB()

	r := gin.Default()

	service.SetDB(db)

	api.SetupRoutes(r)

	r.Run(":8080")
}
