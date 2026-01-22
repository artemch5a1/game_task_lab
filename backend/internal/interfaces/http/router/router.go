package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"example/web-service-gin/internal/interfaces/http/handlers"
)

func NewRouter(gameHandler *handlers.GameHandler) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/games", gameHandler.CreateGame)
	r.GET("/games", gameHandler.GetAllGames)
	r.GET("/games/:id", gameHandler.GetGame)
	r.PUT("/games/:id", gameHandler.UpdateGame)
	r.DELETE("/games/:id", gameHandler.DeleteGame)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":    "Game API",
			"version": "1.0.0",
		})
	})

	return r
}
