package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"party-bot/service"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var (
	bot *linebot.Client
)

func init() {
	var err error
	bot, err = linebot.New(
		os.Getenv("LINE_CHANNEL_SECRET"),
		os.Getenv("LINE_ACCESS_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}
}

// LineBotHandler 處理 Line Bot 訊息的 Handler
func LineBotHandler(c *gin.Context) {
	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.String(http.StatusBadRequest, "Bad Request")
		} else {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
		return
	}
	for _, event := range events {
		var userId string
		switch event.Source.Type {
		case linebot.EventSourceTypeUser:
			userId = event.Source.UserID
		case linebot.EventSourceTypeGroup:
			userId = event.Source.GroupID
		case linebot.EventSourceTypeRoom:
			userId = event.Source.RoomID
		}
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				handleDanmakuMessage(message.Text, userId)
			case *linebot.ImageMessage:
				handleImageMessage(bot, event.ReplyToken, message, userId)
			}
		}
	}

	c.String(http.StatusOK, "OK")
}

func handleImageMessage(bot *linebot.Client, replyToken string, message *linebot.ImageMessage, userId string) {

	// 下載圖片
	content, err := bot.GetMessageContent(message.ID).Do()
	if err != nil {
		log.Print(err)
		return
	}
	defer content.Content.Close()

	// 將圖片保存到本地文件系統
	filePath := service.SaveImageFileLocally(content.Content, message.ID)
	if filePath == "" {
		log.Println("Failed to save the image locally")
		return
	}

	// 在這裡你可以對本地保存的圖片進行進一步的處理
	senderProfile, err := bot.GetProfile(userId).Do()
	if err != nil {
		log.Printf("Error getting sender's profile: %v", err)
		// 錯誤處理...
		return
	}
	senderName := senderProfile.DisplayName
	newImage := service.SaveImage(senderName, filePath)
	// 回覆用戶
	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(fmt.Sprintf("已收到您的圖片，您的照片序號是：%d", newImage.ID)),
	).Do(); err != nil {
		log.Print(err)
	}
}

func handleDanmakuMessage(message string, userId string) {
	// 在這裡你可以對本地保存的圖片進行進一步的處理
	senderProfile, err := bot.GetProfile(userId).Do()
	if err != nil {
		log.Printf("Error getting sender's profile: %v", err)
		// 錯誤處理...
		return
	}
	senderName := senderProfile.DisplayName
	BroadcastMessage(message)
	service.SaveMessage(message, senderName)
}
