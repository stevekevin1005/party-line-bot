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
	r.Static("/images", "./images")
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
