package controller

import (
	"party-bot/handlers"

	"github.com/gin-gonic/gin"
)

func LineCallback(c *gin.Context) {
	handlers.LineBotHandler(c)
}
