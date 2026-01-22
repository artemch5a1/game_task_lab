package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HelloHandler godoc
// @Summary      Hello by name
// @Description  Приветствие пользователя по имени
// @Tags         demo
// @Produce      json
// @Param        name path string true "Имя пользователя"
// @Success      200 {object} map[string]string
// @Router       /hello/{name} [get]
func HelloHandler(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, gin.H{"message": "Hello " + name + "!"})
}
