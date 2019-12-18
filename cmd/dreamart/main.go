package main

import (
	"log"
	"os"
	"strings"
	"sync"

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

	shop, err := shop.NewShop(new(shop.Syncer))

	if err != nil {
		log.Panic(err)
	}

	actionPool := make(map[int64]actions.Action)
	mu := new(sync.RWMutex)

	// TODO: add signal handler which sync the shop with database when the program finishes

	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		go handleUpdate(update, bot, shop, actionPool, mu)
	}
}

func handleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI, store *shop.Shop, actionPool map[int64]actions.Action, mu *sync.RWMutex) {
	if update.CallbackQuery != nil {
		chatID := update.CallbackQuery.Message.Chat.ID
		messageID := update.CallbackQuery.Message.MessageID
		data := update.CallbackQuery.Data

		if actionStrings := strings.Split(data, "-"); len(actionStrings) == 2 {
			if actionStrings[0] == "delete" || actionStrings[0] == "append" || actionStrings[0] == "change" {
				action, err := actions.NewAction(actionStrings[0], actionStrings[1], store)

				if err != nil {
					log.Panic(err)
				}

				msg := tgbotapi.NewEditMessageText(chatID, messageID, "")
				var markup interface{}
				msg.Text, markup = action.Next()

				if markup != nil {
					msg.ReplyMarkup = markup.(*tgbotapi.InlineKeyboardMarkup)
				}

				_, _ = bot.Send(msg)
				_, _ = bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, actionStrings[0]+" "+actionStrings[1]))

				mu.Lock()
				actionPool[chatID] = action
				mu.Unlock()
			} else if action, ok := actionPool[chatID]; ok {
				msg := tgbotapi.NewEditMessageText(chatID, messageID, "")

				err := action.AddChunk(data)

				if err != nil {
					log.Printf("An error has occurred: %s", err)
					_, _ = bot.Send(tgbotapi.NewMessage(chatID, "An error has occurred: "+err.Error()))
				}

				if action.IsDone() {
					if _, ok := action.(*actions.BuyAction); ok {
						mu.Lock()
						delete(actionPool, chatID)
						mu.Unlock()

						msg.Text = "Панель администратора"
						msg.ReplyMarkup = &shop.AdminKeyboard
					} else {
						msg.Text = "Спасибо за покупку!"
					}
				} else {
					if actionStrings[0] == "photo" {
						_, _ = bot.Send(tgbotapi.NewPhotoShare(chatID, action.(*actions.BuyAction).Photo()))
					} else {
						var markup interface{}
						msg.Text, markup = action.Next()

						if markup != nil {
							msg.ReplyMarkup = markup.(*tgbotapi.InlineKeyboardMarkup)
						}
					}
				}

				_, _ = bot.Send(msg)
				_, _ = bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "выбрано"))
			}
		} else if actionStrings[0] == "yes" || actionStrings[0] == "no" {
			switch actionStrings[0] {
			case "yes":
				mu.RLock()
				_ = actionPool[chatID].Execute(update.CallbackQuery.From.UserName, bot)
				mu.RUnlock()

				mu.Lock()
				delete(actionPool, chatID)
				mu.Unlock()

				_, _ = bot.Send(tgbotapi.NewEditMessageText(chatID, messageID, "С вами свяжется один из свободных администраторов"))
				_, _ = bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "готово"))
			case "no":
				mu.Lock()
				delete(actionPool, chatID)
				mu.Unlock()

				_, _ = bot.Send(tgbotapi.NewEditMessageText(chatID, messageID, "Чтобы вернуться в меню покупки напишите /buy"))
				_, _ = bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "отмена покупки"))
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
			case "cancel":
				if _, ok := actionPool[chatID].(*actions.BuyAction); ok {
					mu.RLock()
					actionPool[chatID].SetDone()
					mu.RUnlock()

					mu.Lock()
					delete(actionPool, chatID)
					mu.Unlock()

					msg := tgbotapi.NewEditMessageText(chatID, messageID, "Возвращайтесь! Чтобы открыть меню магазина снова, напишитн /buy")
					_, _ = bot.Send(msg)
				} else {
					mu.RLock()
					actionPool[chatID].SetDone()
					mu.RUnlock()

					mu.Lock()
					delete(actionPool, chatID)
					mu.Unlock()

					msg := tgbotapi.NewEditMessageText(chatID, messageID, "Панель администратора")
					msg.ReplyMarkup = &shop.AdminKeyboard
					_, _ = bot.Send(msg)
					_, _ = bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Действие отменено"))
					return
				}
			}
			_, _ = bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Открыто"))
		}
	}

	if update.Message != nil {
		chatID := update.Message.Chat.ID

		if _, ok := store.Admins[update.Message.From.UserName]; ok && store.Admins[update.Message.From.UserName] == 0 {
			store.AddChatID(update.Message.From.UserName, chatID)
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(chatID, "")

			switch update.Message.Command() {
			case "admin":
				msg.Text = "Панель администратора"
				msg.ReplyMarkup = shop.AdminKeyboard
			case "buy", "start":
				mu.Lock()
				actionPool[chatID], _ = actions.NewAction("buy", "product", store)
				mu.Unlock()

				msg.Text, msg.ReplyMarkup = actionPool[chatID].Next()
			default:
				msg.Text = "Я не знаю этой команды"
			}
			_, _ = bot.Send(msg)

		} else if action, ok := actionPool[chatID]; ok {
			msg := tgbotapi.NewMessage(chatID, "")

			err := action.AddChunk(update)

			if err != nil {
				log.Printf("An error has occurred: %s", err)
				_, _ = bot.Send(tgbotapi.NewMessage(chatID, "An error has occurred: "+err.Error()))
			}

			if action.IsDone() {
				if _, ok := action.(*actions.BuyAction); ok {
					msg.Text = "Спасибо за покупку!"
				} else {
					mu.Lock()
					delete(actionPool, chatID)
					mu.Unlock()

					msg.Text = "Панель администратора"
					msg.ReplyMarkup = shop.AdminKeyboard
				}
			} else {
				msg.Text, msg.ReplyMarkup = action.Next()
			}

			_, _ = bot.Send(msg)
		}
	}
}
