package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "example/web-service-gin/internal/interfaces/docs"
	"example/web-service-gin/internal/interfaces/handlers"
)

// @title           Gin Swagger Example
func main() {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Маршруты
	// Роуты
	r.GET("/", handlers.RootHandler)
	r.GET("/health", handlers.HealthHandler)
	r.GET("/hello/:name", handlers.HelloHandler)
	r.POST("/echo", handlers.EchoHandler)

	// Запуск сервера
	fmt.Println("Server starting on http://localhost:8080")
	r.Run(":8080") // слушает на порту 8080
}
