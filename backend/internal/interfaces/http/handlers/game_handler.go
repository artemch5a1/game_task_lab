// internal/interfaces/http/handlers/game_handler.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"example/web-service-gin/internal/application/dto"
	"example/web-service-gin/internal/application/services"
)

type GameHandler struct {
	gameService *services.GameService
}

func NewGameHandler(gameService *services.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

// CreateGame создает новую игру
// @Summary      Создать игру
// @Description  Создает новую игру в системе
// @Tags         games
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        data body dto.CreateGameDto true "Данные для создания игры"
// @Success      201 {object} dto.GameDto
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /games [post]
func (h *GameHandler) CreateGame(c *gin.Context) {
	var req dto.CreateGameDto

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}

	game, err := h.gameService.CreateGame(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, game)
}

// GetGame получает игру по ID
// @Summary      Получить игру
// @Description  Получает информацию об игре по её идентификатору
// @Tags         games
// @Accept       json
// @Produce      json
// @Param        id path string true "ID игры"
// @Success      200 {object} dto.GameDto
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /games/{id} [get]
func (h *GameHandler) GetGame(c *gin.Context) {
	gameID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID игры"})
		return
	}

	game, err := h.gameService.GetGameByID(c.Request.Context(), gameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Игра не найдена"})
		return
	}

	c.JSON(http.StatusOK, game)
}

// GetAllGames получает список всех игр
// @Summary      Получить все игры
// @Description  Возвращает список всех игр в системе
// @Tags         games
// @Accept       json
// @Produce      json
// @Success      200 {array} dto.GameDto
// @Failure      500 {object} map[string]string
// @Router       /games [get]
func (h *GameHandler) GetAllGames(c *gin.Context) {
	games, err := h.gameService.GetAllGames(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении игр"})
		return
	}

	c.JSON(http.StatusOK, games)
}

// UpdateGame обновляет игру
// @Summary      Обновить игру
// @Description  Обновляет информацию об игре
// @Tags         games
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "ID игры"
// @Param        data body dto.UpdateGameDto true "Данные для обновления игры"
// @Success      200 {object} dto.GameDto
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /games/{id} [put]
func (h *GameHandler) UpdateGame(c *gin.Context) {
	gameID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID игры"})
		return
	}

	var req dto.UpdateGameDto

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}

	req.ID = gameID

	game, err := h.gameService.UpdateGame(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, game)
}

// DeleteGame удаляет игру
// @Summary      Удалить игру
// @Description  Удаляет игру из системы по её идентификатору
// @Tags         games
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        id path string true "ID игры"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /games/{id} [delete]
func (h *GameHandler) DeleteGame(c *gin.Context) {
	gameID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID игры"})
		return
	}

	err = h.gameService.DeleteGame(c.Request.Context(), gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении игры"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Игра успешно удалена"})
}

// HealthCheck проверяет состояние сервиса
// @Summary      Проверка здоровья
// @Description  Проверяет, что сервис работает
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200 {object} map[string]string
// @Router       /health [get]
func (h *GameHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
