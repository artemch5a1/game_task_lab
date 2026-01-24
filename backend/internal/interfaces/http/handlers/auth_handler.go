package handlers

import (
	"net/http"

	"example/web-service-gin/internal/application/dto"
	"example/web-service-gin/internal/application/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Login проверяет логин и пароль
// @Summary      Авторизация
// @Description  Принимает логин и пароль и возвращает true/false
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data body dto.LoginDto true "Логин и пароль"
// @Success      200 {boolean} boolean
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}

	ok, err := h.userService.Authenticate(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ok)
}

