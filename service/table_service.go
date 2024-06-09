package service

import (
	"party-bot/models"
	"party-bot/utils"
)

func FindTable(name string) (models.Table, error) {
	var table models.Table
	db := utils.GetDB()
	if name != "" {
		db = db.Where("name = ?", name)
	}
	if err := db.First(&table).Error; err != nil {
		return table, err
	}
	return table, nil
}
