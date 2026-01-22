package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type EchoRequest struct {
	Message string `json:"message" example:"Hello!"`
}

type EchoResponse struct {
	Echo string `json:"echo"`
}

// EchoHandler godoc
// @Summary      Echo message
// @Description  Повторяет сообщение из тела запроса
// @Tags         demo
// @Accept       json
// @Produce      json
// @Param        data body EchoRequest true "Message payload"
// @Success      200 {object} EchoResponse
// @Failure      400 {object} map[string]string
// @Router       /echo [post]
func EchoHandler(c *gin.Context) {
	var json EchoRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, EchoResponse{Echo: json.Message})
}
