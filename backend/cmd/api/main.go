package main

import (
	"context"
	"log"
	"os"

	_ "example/web-service-gin/docs"
	"example/web-service-gin/internal/application/services"
	"example/web-service-gin/internal/infrastructure/persistence/sqlite"
	"example/web-service-gin/internal/interfaces/http/handlers"
	"example/web-service-gin/internal/interfaces/http/router"
	"example/web-service-gin/internal/interfaces/http/server"
)

// @title           Gin Swagger Example
func main() {

	ctx := context.Background()
	dbPath := os.Getenv("DB_PATH")

	db, err := sqlite.Open(ctx, sqlite.Config{Path: dbPath})
	if err != nil {
		log.Fatal("DB init error:", err)
	}
	defer func() { _ = db.Close() }()

	gameRepo := sqlite.NewGameRepository(db.SQL)
	genreRepo := sqlite.NewGenreRepository(db.SQL)

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
