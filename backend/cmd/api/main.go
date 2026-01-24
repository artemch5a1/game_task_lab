package main

import (
	"log"

	_ "example/web-service-gin/docs"
	"example/web-service-gin/internal/application/services"
	"example/web-service-gin/internal/infrastructure/persistence/data"
	"example/web-service-gin/internal/infrastructure/persistence/inmemory"
	"example/web-service-gin/internal/interfaces/http/handlers"
	"example/web-service-gin/internal/interfaces/http/router"
	"example/web-service-gin/internal/interfaces/http/server"
)

// @title           Gin Swagger Example
func main() {

	store := data.New()
	gameRepo := inmemory.NewGameRepository(store)
	genreRepo := inmemory.NewGenreRepository(store)

	// 3. Сервис
	gameService := services.NewGameService(gameRepo)
	genreService := services.NewGenreService(genreRepo)

	// 4. Хендлер
	gameHandler := handlers.NewGameHandler(gameService)
	genreHandler := handlers.NewGenreHandler(genreService)

	// 5. Роутер (Gin engine)
	ginRouter := router.NewRouter(gameHandler, genreHandler)

	// 6. Запуск сервера
	if err := server.Start(":8080", ginRouter); err != nil {
		log.Fatal("Server error:", err)
	}
}
