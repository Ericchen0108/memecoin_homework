package api

import (
	"memecoin_homework/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/memecoins", func(c *gin.Context) {
		service.CreateMemecoin(c, db)
	})
	r.GET("/memecoins/:id", func(c *gin.Context) {
		service.GetMemecoin(c, db)
	})
	r.PATCH("/memecoins/:id", func(c *gin.Context) {
		service.UpdateMemecoin(c, db)
	})
	r.DELETE("/memecoins/:id", func(c *gin.Context) {
		service.DeleteMemecoin(c, db)
	})
	r.POST("/memecoins/:id/poke", func(c *gin.Context) {
		service.PokeMemecoin(c, db)
	})
}
