package handlers

import (
	"net/http"

	"example/web-service-gin/internal/application/dto"
	"example/web-service-gin/internal/application/services"
	"example/web-service-gin/internal/application/abstraction/repository"
	"example/web-service-gin/internal/constants"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GenreHandler struct {
	genreService *services.GenreService
}

func NewGenreHandler(genreService *services.GenreService) *GenreHandler {
	return &GenreHandler{genreService: genreService}
}

// CreateGenre создает новый жанр
// @Summary      Создать жанр
// @Description  Создает новый жанр в системе
// @Tags         genres
// @Accept       json
// @Produce      json
// @Param        data body dto.CreateGenreDto true "Данные для создания жанра"
// @Success      201 {object} dto.GenreDto
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /genres [post]
func (h *GenreHandler) CreateGenre(c *gin.Context) {
	var req dto.CreateGenreDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}

	genre, err := h.genreService.CreateGenre(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, genre)
}

// GetGenre получает жанр по ID
// @Summary      Получить жанр
// @Description  Получает информацию о жанре по его идентификатору
// @Tags         genres
// @Accept       json
// @Produce      json
// @Param        id path string true "ID жанра"
// @Success      200 {object} dto.GenreDto
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /genres/{id} [get]
func (h *GenreHandler) GetGenre(c *gin.Context) {
	genreID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID жанра"})
		return
	}

	genre, err := h.genreService.GetGenreByID(c.Request.Context(), genreID)
	if err != nil {
		if err == repository.ErrGenreNotFound || err.Error() == constants.ErrGenreNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrGenreNotFound})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, genre)
}

// GetAllGenres получает список всех жанров
// @Summary      Получить все жанры
// @Description  Возвращает список всех жанров в системе
// @Tags         genres
// @Accept       json
// @Produce      json
// @Success      200 {array} dto.GenreDto
// @Failure      500 {object} map[string]string
// @Router       /genres [get]
func (h *GenreHandler) GetAllGenres(c *gin.Context) {
	genres, err := h.genreService.GetAllGenres(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении жанров"})
		return
	}

	c.JSON(http.StatusOK, genres)
}

// UpdateGenre обновляет жанр
// @Summary      Обновить жанр
// @Description  Обновляет информацию о жанре
// @Tags         genres
// @Accept       json
// @Produce      json
// @Param        id path string true "ID жанра"
// @Param        data body dto.UpdateGenreDto true "Данные для обновления жанра"
// @Success      200 {object} dto.GenreDto
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /genres/{id} [put]
func (h *GenreHandler) UpdateGenre(c *gin.Context) {
	genreID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID жанра"})
		return
	}

	var req dto.UpdateGenreDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат JSON"})
		return
	}

	req.ID = genreID

	genre, err := h.genreService.UpdateGenre(c.Request.Context(), req)
	if err != nil {
		if err == repository.ErrGenreNotFound || err.Error() == constants.ErrGenreNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrGenreNotFound})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, genre)
}

// DeleteGenre удаляет жанр
// @Summary      Удалить жанр
// @Description  Удаляет жанр из системы по его идентификатору
// @Tags         genres
// @Accept       json
// @Produce      json
// @Param        id path string true "ID жанра"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /genres/{id} [delete]
func (h *GenreHandler) DeleteGenre(c *gin.Context) {
	genreID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID жанра"})
		return
	}

	err = h.genreService.DeleteGenre(c.Request.Context(), genreID)
	if err != nil {
		if err == repository.ErrGenreNotFound || err.Error() == constants.ErrGenreNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrGenreNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении жанра"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Жанр успешно удален"})
}

