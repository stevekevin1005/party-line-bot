package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"party-bot/service"
	"regexp"
	"time"
	"unicode/utf8"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"gorm.io/gorm"
)

var (
	bot   *linebot.Client
	cache *service.Cache
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
	cache = service.NewCache()
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
				if message.Text == "電子喜帖" || message.Text == "婚宴交通" || message.Text == "停車資訊" {
					return
				} else if message.Text == "［愛的留言］" {
					cache.Set(userId+"Danmaku", true, 60*time.Second)
					if _, err := bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTextMessage("感謝您使用此功能：\n請在接下來的5分鐘內，將您想告訴新人的話傳給我～\n\n愛的留言就會投射至大螢幕\n～趕快留言給新人吧(๑ ◡ ๑)"),
					).Do(); err != nil {
						log.Print(err)
					}
				} else if message.Text == " [拍立得列印] " {
					var imageCount, _ = service.CountImages()
					if imageCount >= 200 {
						if _, err := bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewTextMessage("WOW（*＾-＾*）太搶手了！\n拍立得~已經全數印完囉！\n謝謝你們大家的喜愛 ♡"),
						).Do(); err != nil {
							log.Print(err)
						}
					} else {
						cache.Set(userId+"Photo", true, 300*time.Second)
						if _, err := bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewTextMessage("感謝您使用此功能：\n請在接下來的5分鐘內，將希望列印的照片上傳給我～\n數量有限，印完為止，可以到入口處看看您的照片有沒有印出來唷(๑•̀ㅂ•́)و✧~"),
						).Do(); err != nil {
							log.Print(err)
						}
					}
				} else if message.Text == "［座位查詢］" {
					senderProfile, err := bot.GetProfile(userId).Do()
					if err != nil {
						log.Printf("Error getting sender's profile: %v", err)
						// 錯誤處理...
						return
					}
					senderName := senderProfile.DisplayName
					table, err := service.FindTable(senderName)
					if err != nil {
						if errors.Is(err, gorm.ErrRecordNotFound) {
							cache.Set(userId+"findTable", true, 300*time.Second)
							if _, err := bot.ReplyMessage(
								event.ReplyToken,
								linebot.NewTextMessage("感謝您使用此功能：\n請在此輸入「您的名字」\n將會告知您的桌位"),
							).Do(); err != nil {
								log.Print(err)
							}
						}
					} else {
						if _, err := bot.ReplyMessage(
							event.ReplyToken,
							linebot.NewTextMessage("您的桌位是： "+table.TableName),
						).Do(); err != nil {
							log.Print(err)
						}
					}
				} else {
					if _, ok := cache.Get(userId + "findTable"); ok {
						handleFindTable(message.Text, event.ReplyToken)
						cache.Delete(userId + "findTable")
						return
					}
					if _, ok := cache.Get(userId + "Danmaku"); ok {
						handleDanmakuMessage(message.Text, userId, event.ReplyToken)
						return
					}
				}
			case *linebot.ImageMessage:
				if _, ok := cache.Get(userId + "Photo"); ok {
					handleImageMessage(bot, event.ReplyToken, message, userId)
				}
			}
		}
	}

	c.String(http.StatusOK, "OK")
}

func handleImageMessage(bot *linebot.Client, replyToken string, message *linebot.ImageMessage, userId string) {
	senderProfile, err := bot.GetProfile(userId).Do()
	if err != nil {
		log.Printf("Error getting sender's profile: %v", err)
		// 錯誤處理...
		return
	}
	senderName := senderProfile.DisplayName
	images, _ := service.ListImages(senderName)
	if len(images) >= 2 {
		if _, err := bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage("Oopsヾ(≧O≦)〃 拍立得數量有限! 1人只能上傳2張喔～"),
		).Do(); err != nil {
			log.Print(err)
		}
	}
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

	newImage := service.SaveImage(senderName, filePath)
	// 回覆用戶
	if _, err := bot.ReplyMessage(
		replyToken,
		linebot.NewTextMessage(fmt.Sprintf("已收到您的圖片，您的照片序號是：%s", newImage.Serial)),
	).Do(); err != nil {
		log.Print(err)
	}
}

func handleDanmakuMessage(message string, userId string, replyToken string) {
	senderProfile, err := bot.GetProfile(userId).Do()
	if err != nil {
		log.Printf("Error getting sender's profile: %v", err)
		// 錯誤處理...
		return
	}
	senderName := senderProfile.DisplayName
	regex := regexp.MustCompile(`^[\p{Han}\p{Katakana}\p{Hiragana}\p{Hangul}a-zA-Z0-9\s]+$`)
	if !regex.MatchString(message) || utf8.RuneCountInString(message) > 20 {
		if _, err := bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage("只能傳20個以內的訊息唷~"),
		).Do(); err != nil {
			log.Print(err)
		}
		return
	}
	BroadcastMessage(message)
	service.SaveMessage(message, senderName)
}

func handleFindTable(name string, replyToken string) {
	table, err := service.FindTable(name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if _, err := bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage("不好意思(*꧆▽꧆*)\n未查詢到您得名字，請您至婚宴入口處詢問招待呦！"),
		).Do(); err != nil {
			log.Print(err)
		}
	} else {
		if _, err := bot.ReplyMessage(
			replyToken,
			linebot.NewTextMessage("您的座位是: "+table.TableName),
		).Do(); err != nil {
			log.Print(err)
		}
	}
}
