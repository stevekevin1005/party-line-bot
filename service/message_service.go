package service

import (
	"party-bot/models"
	"party-bot/utils"
)

func SaveMessage(name string, message string) models.Message {
	newImage := models.Message{
		Name: name,
		Text: message,
	}
	utils.GetDB().Create(&newImage)
	return newImage
}
