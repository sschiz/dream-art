package main

import (
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sschiz/dream-art/pkg/actions"
	"github.com/sschiz/dream-art/pkg/shop"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	if err != nil {
		log.Panic(err)
	}

	shop := new(shop.Shop)
	err = shop.Syncer.Sync(shop)
	actionPool := make(map[int64]actions.Action)

	// TODO: add signal handler which sync the shop with database when the program finishes

	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		go handleUpdate(update, bot, shop, actionPool)
	}
}

func handleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI, store *shop.Shop, actionPool map[int64]actions.Action) {
	if update.CallbackQuery != nil {
		chatID := update.CallbackQuery.Message.Chat.ID
		messageID := update.CallbackQuery.Message.MessageID
		data := update.CallbackQuery.Data

		if strings.HasPrefix(data, "append-") {
			if strings.HasSuffix(data, "admin") {
				action, err := actions.NewAction("append", "admin", store)

				if err != nil {
					log.Panic(err)
				}

				_, _ = bot.Send(tgbotapi.NewEditMessageText(chatID, messageID, action.Next()))
				_, _ = bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Добавление админа"))

				actionPool[chatID] = action
			}
		} else {
			switch data {
			case "admin":
				msg := tgbotapi.NewEditMessageText(chatID, messageID, "Управление админами")
				msg.ReplyMarkup = &shop.AdminManagmentKeyboard
				_, _ = bot.Send(msg)
			case "product":
				msg := tgbotapi.NewEditMessageText(chatID, messageID, "Управление продуктами")
				msg.ReplyMarkup = &shop.ProductManagmentKeyboard
				_, _ = bot.Send(msg)
			case "category":
				msg := tgbotapi.NewEditMessageText(chatID, messageID, "Управление категориями")
				msg.ReplyMarkup = &shop.CategoryManagmentKeyboard
				_, _ = bot.Send(msg)
			case "order":
			}
			_, _ = bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Открыто"))
		}
	}

	if update.Message != nil {
		chatID := update.Message.Chat.ID

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(chatID, "")
			switch update.Message.Command() {
			case "admin":
				msg.Text = "Панель администратора"
				msg.ReplyMarkup = shop.AdminKeyboard
			case "buy", "start":
				msg.Text = "В разработке"
			default:
				msg.Text = "Я не знаю этой команды"
			}

			_, _ = bot.Send(msg)
		} else if action, ok := actionPool[chatID]; ok {
			msg := tgbotapi.NewMessage(chatID, "")

			err := action.AddChunk(update.Message.Text)

			if err != nil {
				log.Printf("An error has occurred: %s", err)
				_, _ = bot.Send(tgbotapi.NewMessage(chatID, "An error has occurred: "+err.Error()))
			}

			if action.IsDone() {
				delete(actionPool, chatID)

				msg.Text = "Панель администратора"
				msg.ReplyMarkup = shop.AdminKeyboard
			} else {
				msg.Text = action.Next()
			}

			_, _ = bot.Send(msg)
		}
	}
}
