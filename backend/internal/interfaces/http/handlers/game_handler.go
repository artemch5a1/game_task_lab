// internal/interfaces/http/handlers/game_handler.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"example/web-service-gin/internal/application/dto"
	"example/web-service-gin/internal/application/services"

	"github.com/google/uuid"
)

type GameHandler struct {
	gameService *services.GameService
}

func NewGameHandler(gameService *services.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

// CreateGame обрабатывает POST /games
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

// GetGame обрабатывает GET /games/:id
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

// GetAllGames обрабатывает GET /games
func (h *GameHandler) GetAllGames(c *gin.Context) {
	games, err := h.gameService.GetAllGames(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении игр"})
		return
	}

	c.JSON(http.StatusOK, games)
}

// UpdateGame обрабатывает PUT /games/:id
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

// DeleteGame обрабатывает DELETE /games/:id
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
