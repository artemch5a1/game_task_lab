package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RootHandler godoc
// @Summary      Root endpoint
// @Description  Проверка работы сервера
// @Tags         system
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       / [get]
func RootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello Gin!"})
}
