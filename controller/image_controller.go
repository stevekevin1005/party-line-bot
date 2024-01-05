package controller

import (
	"net/http"
	"party-bot/service"

	"github.com/gin-gonic/gin"
)

func ListImages(c *gin.Context) {
	name := c.Query("name")
	images, err := service.ListImages(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, images)
}

func MarkImage(c *gin.Context) {
	var markRequest struct {
		ID uint `json:"id"`
	}

	if err := c.BindJSON(&markRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := service.MarkImage(int(markRequest.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
