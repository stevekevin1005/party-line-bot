package service

import (
	"errors"
	"io"
	"log"
	"os"
	"party-bot/models"
	"party-bot/utils"
	"path/filepath"
)

func SaveImageLocally(content io.Reader, fileName string) string {
	// 指定保存的目錄，這裡假設是當前工作目錄的 "images" 子目錄
	dir := "images"

	// 確保目錄存在，如果不存在則創建
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("Failed to create directory: %v", err)
		return ""
	}

	// 構造本地文件的完整路徑
	filePath := filepath.Join(dir, fileName+".jpg")

	// 創建本地文件
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return ""
	}
	defer file.Close()

	// 將下載的圖片內容寫入本地文件
	_, err = io.Copy(file, content)
	if err != nil {
		log.Printf("Failed to write file content: %v", err)
		return ""
	}

	return filePath
}

func ListImages(name string) ([]models.Image, error) {
	var images []models.Image
	db := utils.GetDB()
	if name != "" {
		db = db.Where("name Like ?", "%"+name+"%")
	}
	if err := db.Find(&images).Error; err != nil {

		return nil, err
	}
	return images, nil
}

func MarkImage(id int) error {
	var image models.Image
	db := utils.GetDB()

	if err := db.First(&image, id).Error; err != nil {
		return errors.New("failed to find the image with specified ID")
	}

	image.Status = true

	if err := db.Save(&image).Error; err != nil {
		return errors.New("failed to mark the image")
	}

	return nil
}
