package di

import (
	"context"
	"errors"
	"os"
	"strings"

	"example/web-service-gin/internal/application/services"
	"example/web-service-gin/internal/infrastructure/persistence/sqlite"
	"example/web-service-gin/internal/interfaces/http/handlers"
	"example/web-service-gin/internal/interfaces/http/router"

	"github.com/gin-gonic/gin"
)

type App struct {
	Router *gin.Engine
	Close  func() error
}

func Build(ctx context.Context) (*App, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}

	dbPath := strings.TrimSpace(os.Getenv("DB_PATH"))
	db, err := sqlite.Open(ctx, sqlite.Config{Path: dbPath})
	if err != nil {
		return nil, err
	}

	gameRepo := sqlite.NewGameRepository(db.SQL)
	genreRepo := sqlite.NewGenreRepository(db.SQL)
	userRepo := sqlite.NewUserRepository(db.SQL)

	gameService := services.NewGameService(gameRepo)
	genreService := services.NewGenreService(genreRepo)
	userService := services.NewUserService(userRepo)

	gameHandler := handlers.NewGameHandler(gameService)
	genreHandler := handlers.NewGenreHandler(genreService)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(userService)

	r := router.NewRouter(gameHandler, genreHandler, userHandler, authHandler)

	return &App{
		Router: r,
		Close:  db.Close,
	}, nil
}

