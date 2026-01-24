package handlers

import (
	"net/http"

	"example/web-service-gin/internal/application/dto"
	"example/web-service-gin/internal/application/services"
	"example/web-service-gin/internal/constants"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login проверяет логин и пароль
// @Summary      Авторизация
// @Description  Принимает логин и пароль и возвращает JWT токен
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data body dto.LoginDto true "Логин и пароль"
// @Success      200 {object} dto.AuthTokenDto
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if err.Error() == constants.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": constants.ErrUnauthorized})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.AuthTokenDto{Token: token})
}

// Register регистрирует обычного пользователя
// @Summary      Регистрация
// @Description  Создает пользователя с ролью user и возвращает JWT токен
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data body dto.RegisterDto true "Логин и пароль"
// @Success      201 {object} dto.AuthTokenDto
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}

	token, err := h.authService.Register(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.AuthTokenDto{Token: token})
}

