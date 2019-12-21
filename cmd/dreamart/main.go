/*
 * (c) 2019, Matyushkin Alexander <sav3nme@gmail.com>
 * GNU General Public License v3.0+ (see COPYING or https://www.gnu.org/licenses/gpl-3.0.txt)
 */

/*
 * (c) 2019, Matyushkin Alexander <sav3nme@gmail.com>
 * GNU General Public License v3.0+ (see COPYING or https://www.gnu.org/licenses/gpl-3.0.txt)
 */

package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sschiz/dream-art/pkg/action"
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

	newShop, err := shop.New(&shop.Syncer{ConnectionString: "user=sschiz password=60egozaz dbname=shop"})

	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if r := recover(); r != nil {
			err := newShop.Sync()
			if err != nil {
				log.Printf("Err while syncing: %s", err)
			}
			log.Println("Panicking", r)
		}
	}()

	actionPool := make(map[int64]action.Action)
	mu := new(sync.RWMutex)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			// sig is a ^C, handle it
			log.Printf("Called signal: %s", sig)
			err := newShop.Sync()
			if err != nil {
				log.Printf("error while syncing: %s", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
	}()

	for update := range updates {
		go handleUpdate(update, bot, newShop, actionPool, mu)
	}
}

func handleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI, store *shop.Shop, actionPool map[int64]action.Action, mu *sync.RWMutex) {
	if update.CallbackQuery != nil {
		chatID := update.CallbackQuery.Message.Chat.ID
		messageID := update.CallbackQuery.Message.MessageID
		data := update.CallbackQuery.Data

		if actionStrings := strings.Split(data, "-"); len(actionStrings) == 2 {
			if actionStrings[0] == "delete" || actionStrings[0] == "append" || actionStrings[0] == "change" {
				act, err := action.New(actionStrings[0], actionStrings[1], store)

				if err != nil {
					log.Panic(err)
				}

				msg := tgbotapi.NewEditMessageText(chatID, messageID, "")
				var markup interface{}
				msg.Text, markup = act.Next()

				if markup != nil {
					msg.ReplyMarkup = markup.(*tgbotapi.InlineKeyboardMarkup)
				}

				_, _ = bot.Send(msg)
				_, _ = bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, actionStrings[0]+" "+actionStrings[1]))

				mu.Lock()
				actionPool[chatID] = act
				mu.Unlock()
			} else if act, ok := actionPool[chatID]; ok {
				msg := tgbotapi.NewEditMessageText(chatID, messageID, "")

				err := act.AddChunk(data)

				if err != nil {
					log.Printf("An error has occurred: %s", err)
					_, _ = bot.Send(tgbotapi.NewMessage(chatID, "An error has occurred: "+err.Error()))
				}

				if act.IsDone() {
					if _, ok := act.(*action.Buy); !ok {
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
						_, _ = bot.Send(tgbotapi.NewPhotoShare(chatID, act.(*action.Buy).Photo()))
					} else {
						var markup interface{}
						msg.Text, markup = act.Next()

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
				if act, ok := actionPool[chatID]; ok {
					if _, ok := act.(*action.Buy); ok {
						act.SetDone()

						mu.Lock()
						delete(actionPool, chatID)
						mu.Unlock()

						msg := tgbotapi.NewEditMessageText(chatID, messageID, "Возвращайтесь! Чтобы открыть меню магазина снова, напишите /buy")
						_, _ = bot.Send(msg)
					} else {
						act.SetDone()

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
				if _, ok := store.Admins[update.Message.From.UserName]; ok {
					msg.Text = "Панель администратора"
					msg.ReplyMarkup = shop.AdminKeyboard
				}
			case "buy":
				mu.Lock()
				actionPool[chatID], _ = action.New("buy", "product", store)
				mu.Unlock()

				msg.Text, msg.ReplyMarkup = actionPool[chatID].Next()
			case "start":
				msg.Text = "Напишите /buy, а после выберите понравившийся стиль и размер будущего портрета. " +
					"По всем вопросам к @yakovlevpave1"
			default:
				msg.Text = "Я не знаю этой команды"
			}
			_, _ = bot.Send(msg)

		} else if act, ok := actionPool[chatID]; ok {
			msg := tgbotapi.NewMessage(chatID, "")

			err := act.AddChunk(update)

			if err != nil {
				log.Printf("An error has occurred: %s", err)

				act.SetDone()

				mu.Lock()
				delete(actionPool, chatID)
				mu.Unlock()

				msg := tgbotapi.NewMessage(chatID, "")
				msg.Text = "Панель администратора"
				msg.ReplyMarkup = shop.AdminKeyboard

				_, _ = bot.Send(tgbotapi.NewMessage(chatID, "An error has occurred: "+err.Error()))
				_, _ = bot.Send(msg)
			}

			if act.IsDone() {
				if _, ok := act.(*action.Buy); ok {
					msg.Text = "Спасибо за покупку!"
				} else {
					mu.Lock()
					delete(actionPool, chatID)
					mu.Unlock()

					msg.Text = "Панель администратора"
					msg.ReplyMarkup = shop.AdminKeyboard
				}
			} else {
				msg.Text, msg.ReplyMarkup = act.Next()
			}

			_, _ = bot.Send(msg)
		}
	}
}
