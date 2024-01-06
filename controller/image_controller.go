package controller

import (
	"net/http"
	"party-bot/service"

	"github.com/gin-gonic/gin"
)

// errorResponse represents the response format for an error
type errorResponse struct {
	Message string `json:"message"`
	// Add other fields as needed
}

// imageResponse represents the response format for an image
type imageResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

// markRequest represents the request payload for marking an image
type markRequest struct {
	ID uint `json:"id"`
}

// ListImages godoc
// @Summary List images
// @Description Get a list of images with optional filtering by name.
// @ID list-images
// @Accept  json
// @Produce  json
// @Param name query string false "Image name to filter by"
// @Success 200 {array} imageResponse
// @Failure 500 {object} errorResponse
// @Router /api/images/list [get]
func ListImages(c *gin.Context) {
	name := c.Query("name")
	images, err := service.ListImages(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse{Message: err.Error()})
		return
	}
	// Convert []models.Image to []imageResponse
	var imageResponses []imageResponse
	for _, img := range images {
		imageResponses = append(imageResponses, imageResponse{
			ID:   img.ID,
			Name: img.Name,
			Path: img.Path,
			// Add other fields as needed
		})
	}
	c.JSON(http.StatusOK, imageResponses)
}

// MarkImage godoc
// @Summary Mark an image
// @Description Mark an image with the specified ID.
// @ID mark-image
// @Accept  json
// @Produce  json
// @Param request body markRequest true "JSON payload with image ID"
// @Success 200
// @Failure 400 {object} errorResponse "Bad Request"
// @Router /api/images/mark [post]
func MarkImage(c *gin.Context) {
	var markRequest markRequest
	if err := c.BindJSON(&markRequest); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Message: err.Error()})
		return
	}

	err := service.MarkImage(int(markRequest.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Message: err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
