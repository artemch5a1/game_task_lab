package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"example/web-service-gin/internal/application/services"
	"example/web-service-gin/internal/infrastructure/persistence/inmemory"
	_ "example/web-service-gin/internal/interfaces/docs"
	"example/web-service-gin/internal/interfaces/http/handlers"
	"example/web-service-gin/internal/interfaces/http/router"
	"example/web-service-gin/internal/interfaces/http/server"
)

// @title           Gin Swagger Example
func main() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	gameRepo := inmemory.NewGameRepository()

	// 3. Сервис
	gameService := services.NewGameService(gameRepo)

	// 4. Хендлер
	gameHandler := handlers.NewGameHandler(gameService)

	// 5. Роутер (Gin engine)
	ginRouter := router.NewRouter(gameHandler)

	// 6. Запуск сервера
	if err := server.Start(":8080", ginRouter); err != nil {
		log.Fatal("Server error:", err)
	}
}
