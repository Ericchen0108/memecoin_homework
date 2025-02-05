package api

import (
	"memecoin_homework/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/memecoins", service.CreateMemecoin)
	r.GET("/memecoins/:id", service.GetMemecoin)
	r.PATCH("/memecoins/:id", service.UpdateMemecoin)
	r.DELETE("/memecoins/:id", service.DeleteMemecoin)
	r.POST("/memecoins/:id/poke", service.PokeMemecoin)
}
