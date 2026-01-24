package main

import (
	"context"
	"log"
	"os"

	_ "example/web-service-gin/docs"
	"example/web-service-gin/internal/di"
	"example/web-service-gin/internal/interfaces/http/server"
)

// @title           Gin Swagger Example
func main() {

	ctx := context.Background()
	app, err := di.Build(ctx)
	if err != nil {
		log.Fatal("DI build error:", err)
	}
	defer func() { _ = app.Close() }()

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	// 6. Запуск сервера
	if err := server.Start(addr, app.Router); err != nil {
		log.Fatal("Server error:", err)
	}
}
