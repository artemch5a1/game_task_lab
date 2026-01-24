package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"example/web-service-gin/internal/interfaces/http/handlers"
)

func NewRouter(
	gameHandler *handlers.GameHandler,
	genreHandler *handlers.GenreHandler,
	userHandler *handlers.UserHandler,
	authHandler *handlers.AuthHandler,
	adminOnly gin.HandlerFunc,
) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"Accept",
			"X-Requested-With",
			"Cache-Control",
			// Можно добавить любые кастомные заголовки
		},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 часов
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if adminOnly != nil {
		r.POST("/games", adminOnly, gameHandler.CreateGame)
	} else {
		r.POST("/games", gameHandler.CreateGame)
	}
	r.GET("/games", gameHandler.GetAllGames)
	r.GET("/games/:id", gameHandler.GetGame)
	if adminOnly != nil {
		r.PUT("/games/:id", adminOnly, gameHandler.UpdateGame)
		r.DELETE("/games/:id", adminOnly, gameHandler.DeleteGame)
	} else {
		r.PUT("/games/:id", gameHandler.UpdateGame)
		r.DELETE("/games/:id", gameHandler.DeleteGame)
	}

	if adminOnly != nil {
		r.POST("/genres", adminOnly, genreHandler.CreateGenre)
	} else {
		r.POST("/genres", genreHandler.CreateGenre)
	}
	r.GET("/genres", genreHandler.GetAllGenres)
	r.GET("/genres/:id", genreHandler.GetGenre)
	if adminOnly != nil {
		r.PUT("/genres/:id", adminOnly, genreHandler.UpdateGenre)
		r.DELETE("/genres/:id", adminOnly, genreHandler.DeleteGenre)
	} else {
		r.PUT("/genres/:id", genreHandler.UpdateGenre)
		r.DELETE("/genres/:id", genreHandler.DeleteGenre)
	}

	if adminOnly != nil {
		r.POST("/users", adminOnly, userHandler.CreateUser)
	} else {
		r.POST("/users", userHandler.CreateUser)
	}
	r.GET("/users", userHandler.GetAllUsers)
	r.GET("/users/:id", userHandler.GetUser)
	if adminOnly != nil {
		r.PUT("/users/:id", adminOnly, userHandler.UpdateUser)
		r.DELETE("/users/:id", adminOnly, userHandler.DeleteUser)
	} else {
		r.PUT("/users/:id", userHandler.UpdateUser)
		r.DELETE("/users/:id", userHandler.DeleteUser)
	}

	r.POST("/auth/login", authHandler.Login)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":    "Game API",
			"version": "1.0.0",
		})
	})

	return r
}
