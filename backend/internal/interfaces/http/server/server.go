package server

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Start(addr string, router *gin.Engine) error {
	log.Printf("Server starting on %s", addr)
	return router.Run(addr)
}
