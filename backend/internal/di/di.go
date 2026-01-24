package di

import (
	"context"
	"errors"

	"example/web-service-gin/internal/application/services"
	"example/web-service-gin/internal/config"
	jwtinfra "example/web-service-gin/internal/infrastructure/auth/jwt"
	"example/web-service-gin/internal/infrastructure/persistence/sqlite"
	"example/web-service-gin/internal/interfaces/http/handlers"
	"example/web-service-gin/internal/interfaces/http/router"

	"github.com/gin-gonic/gin"
	"time"
)

type App struct {
	Router *gin.Engine
	Close  func() error
}

func Build(ctx context.Context) (*App, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}

	cfg := config.Load()
	db, err := sqlite.Open(ctx, sqlite.Config{Path: cfg.DBPath})
	if err != nil {
		return nil, err
	}

	gameRepo := sqlite.NewGameRepository(db.SQL)
	genreRepo := sqlite.NewGenreRepository(db.SQL)
	userRepo := sqlite.NewUserRepository(db.SQL)

	gameService := services.NewGameService(gameRepo)
	genreService := services.NewGenreService(genreRepo)
	userService := services.NewUserService(userRepo)
	jwtProvider := jwtinfra.NewProvider(cfg.JWTSecret, cfg.JWTIssuer, time.Duration(cfg.JWTTTLHours)*time.Hour)
	authService := services.NewAuthService(userRepo, jwtProvider)

	gameHandler := handlers.NewGameHandler(gameService)
	genreHandler := handlers.NewGenreHandler(genreService)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)

	r := router.NewRouter(gameHandler, genreHandler, userHandler, authHandler)

	return &App{
		Router: r,
		Close:  db.Close,
	}, nil
}

