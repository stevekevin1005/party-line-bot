package main

import (
	"os"
	"party-bot/routes"
	"path"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	routes.SetupRoutes(r)
	r.NoRoute(func(c *gin.Context) {
		filepath := path.Join("./frontend/dist", c.Request.URL.Path)
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			c.File("./frontend/dist/index.html")
		} else {
			c.File(filepath)
		}
	})
	r.Run(":8080") // 啟動伺服器在 :8080 port
}
