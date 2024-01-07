package controller

import (
	"party-bot/handlers"

	"github.com/gin-gonic/gin"
)

// GetDanmaku godoc
// @Summary Get danmaku message
// @Description Get danmaku by websocket
// @ID handle-websocket-danmaku
// @Router /api/v1/danmaku/ws [get]
// @tags Danmaku
func GetDanmaku(c *gin.Context) {
	handlers.WebSocketHandler(c)
}
