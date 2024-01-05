package routes

import (
	"party-bot/controller"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/api/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from Gin!",
		})
	})

	r.GET("/api/images/list", controller.ListImages)

	r.POST("/api/images/mark", controller.MarkImage)

	r.POST("/callback", controller.LineCallback)
}
