package main

import (
	"log"

	_ "example/web-service-gin/docs"
	"example/web-service-gin/internal/application/services"
	"example/web-service-gin/internal/infrastructure/persistence/inmemory"
	"example/web-service-gin/internal/interfaces/http/handlers"
	"example/web-service-gin/internal/interfaces/http/router"
	"example/web-service-gin/internal/interfaces/http/server"
)

// @title           Gin Swagger Example
func main() {

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
