package main

import (
	"os"
	"party-bot/docs"
	"party-bot/routes"
	"party-bot/utils"
	"path"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func imagesHeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/images" || len(c.Request.URL.Path) > 7 && c.Request.URL.Path[:8] == "/images/" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
			// Add other headers if needed
		}
		c.Next()
	}
}

// @title party line bot
// @version 1.0
// @description this is a line bot for party usage
func main() {
	utils.AutoMigrate()
	r := gin.Default()
	r.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     []string{"*"},
				AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"*"},
				ExposeHeaders:    []string{"*"},
				AllowCredentials: true,
			},
		),
	)
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	imagesGroup := r.Group("/images")
	imagesGroup.Use(imagesHeaderMiddleware())
	imagesGroup.Static("/", "./images")

	routes.SetupRoutes(r)
	r.NoRoute(func(c *gin.Context) {
		filepath := path.Join("./frontend/dist", c.Request.URL.Path)
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			c.File("./frontend/dist/index.html")
		} else {
			c.File(filepath)
		}
	})
	r.Run(":8080")
}
