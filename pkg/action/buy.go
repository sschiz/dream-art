/*
 * (c) 2019, Matyushkin Alexander <sav3nme@gmail.com>
 * GNU General Public License v3.0+ (see COPYING or https://www.gnu.org/licenses/gpl-3.0.txt)
 */

package action

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sschiz/dream-art/pkg/product"
	"github.com/sschiz/dream-art/pkg/shop"
)

type Buy struct {
	isDone            bool
	isChunksCollected bool
	shop              *shop.Shop
	currentCategory   int
	currentProduct    int
	cart              []*product.Product
	orderText         string
	userName          string
}

func (a *Buy) SetDone() {
	a.isDone = true
}

func (a *Buy) Execute(args ...interface{}) error {
	if !a.isChunksCollected {
		return ErrChunksIsNotCollected
	}

	if a.isDone {
		return ErrActionIsAlreadyDone
	}

	if len(args) == 2 {
		nickname, bot := args[0].(string), args[1].(*tgbotapi.BotAPI)

		for _, chatID := range a.shop.Admins {
			if chatID != 0 {
				_, _ = bot.Send(tgbotapi.NewMessage(chatID, "#заказ\n"+a.orderText+"\n От @"+nickname))
			}
		}
	}

	a.isDone = true

	return nil
}

func (a Buy) IsDone() bool {
	return a.isDone
}

func (a Buy) IsChunksCollected() bool {
	return a.isChunksCollected
}

func (a *Buy) AddChunk(chunk interface{}) error {
	if data, ok := chunk.(string); ok {
		if strings.HasPrefix(data, "next-") {
			a.currentProduct++

			if a.currentProduct > len(a.shop.Categories()[a.currentCategory].Products())-1 {
				a.currentProduct = 0
			}
		} else if strings.HasPrefix(data, "back-") {
			a.currentProduct--

			if a.currentProduct < 0 {
				a.currentProduct = len(a.shop.Categories()[a.currentCategory].Products()) - 1
			}
		} else if strings.HasPrefix(data, "select-") {
			i, err := strconv.Atoi(data[7:])

			if err != nil {
				return err
			}

			a.cart = append(a.cart, a.shop.Categories()[a.currentCategory].Products()[i])

			a.currentCategory++
			a.currentProduct = 0
			if a.currentCategory > len(a.shop.Categories())-1 {
				a.isChunksCollected = true
				return nil
			}
		} else {
			a.userName = data
		}
	}

	return nil
}

func (a *Buy) Next() (text string, out interface{}) {
	if len(a.shop.Categories()) == 0 || len(a.shop.Categories()[a.currentCategory].Products()) == 0 {
		return "Магазин пуст", nil
	}

	p := a.shop.Categories()[a.currentCategory].Products()[a.currentProduct]

	text += "[" + strings.Title(p.Name) + "]" + "(" + p.Photo + ")" + "\n\n"
	text += p.Description + "\n\n"
	text += "Цена: " + fmt.Sprintf("%.2f", float32(p.Price)/10) + " руб"

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️", "back-back"),
			tgbotapi.NewInlineKeyboardButtonData("Выбрать", "select-"+strconv.Itoa(a.currentProduct)),
			tgbotapi.NewInlineKeyboardButtonData("➡️", "next-next"),
		),
		shop.CancelRow.InlineKeyboard[0],
	)

	return text, &keyboard
}
