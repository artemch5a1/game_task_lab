package handlers

import (
	"net/http"

	"example/web-service-gin/internal/application/abstraction/repository"
	"example/web-service-gin/internal/application/dto"
	"example/web-service-gin/internal/application/services"
	"example/web-service-gin/internal/constants"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// CreateUser создает нового пользователя
// @Summary      Создать пользователя
// @Description  Создает нового пользователя в системе
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        data body dto.CreateUserDto true "Данные для создания пользователя"
// @Success      201 {object} dto.UserDto
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}

	created, err := h.userService.CreateUser(c.Request.Context(), req)
	if err != nil {
		if err == repository.ErrUserAlreadyExists || err.Error() == constants.ErrUserAlreadyExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrUserAlreadyExists})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetUser получает пользователя по ID
// @Summary      Получить пользователя
// @Description  Получает информацию о пользователе по его идентификатору
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "ID пользователя"
// @Success      200 {object} dto.UserDto
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID пользователя"})
		return
	}

	u, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if err == repository.ErrUserNotFound || err.Error() == constants.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrUserNotFound})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, u)
}

// GetAllUsers получает список всех пользователей
// @Summary      Получить всех пользователей
// @Description  Возвращает список всех пользователей в системе
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200 {array} dto.UserDto
// @Failure      500 {object} map[string]string
// @Router       /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении пользователей"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// UpdateUser обновляет пользователя
// @Summary      Обновить пользователя
// @Description  Обновляет информацию о пользователе
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "ID пользователя"
// @Param        data body dto.UpdateUserDto true "Данные для обновления пользователя"
// @Success      200 {object} dto.UserDto
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID пользователя"})
		return
	}

	var req dto.UpdateUserDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}
	req.ID = userID

	updated, err := h.userService.UpdateUser(c.Request.Context(), req)
	if err != nil {
		if err == repository.ErrUserNotFound || err.Error() == constants.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrUserNotFound})
			return
		}
		if err == repository.ErrUserAlreadyExists || err.Error() == constants.ErrUserAlreadyExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": constants.ErrUserAlreadyExists})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteUser удаляет пользователя
// @Summary      Удалить пользователя
// @Description  Удаляет пользователя из системы по его идентификатору
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "ID пользователя"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID пользователя"})
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), userID); err != nil {
		if err == repository.ErrUserNotFound || err.Error() == constants.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrUserNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении пользователя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно удален"})
}

