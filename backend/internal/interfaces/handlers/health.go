package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler godoc
// @Summary      Health check
// @Description  Проверка состояния сервера
// @Tags         system
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /health [get]
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
