package routes

import (
	"party-bot/controller"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	apiV1Group := r.Group("/api/v1")

	apiV1Group.GET("/images/list", controller.ListImages)
	apiV1Group.POST("/images/mark", controller.MarkImage)
	apiV1Group.GET("/danmaku/ws", controller.GetDanmaku)

	r.POST("/callback", controller.LineCallback)

}
