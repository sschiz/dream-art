package action

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sschiz/dream-art/pkg/shop"
	"strconv"
	"strings"
)

type CategoryAppend struct {
	isDone            bool
	isChunksCollected bool
	shop              *shop.Shop
	categoryName      string
}

func (a *CategoryAppend) SetDone() {
	a.isDone = true
}

func (a *CategoryAppend) Execute(...interface{}) error {
	if !a.isChunksCollected {
		return ErrChunksIsNotCollected
	}

	if a.isDone {
		return ErrActionIsAlreadyDone
	}

	a.shop.AppendCategory(a.categoryName)

	a.isDone = true

	return nil
}

func (a CategoryAppend) IsDone() bool {
	return a.isDone
}

func (a CategoryAppend) IsChunksCollected() bool {
	return a.isChunksCollected
}

func (a *CategoryAppend) AddChunk(chunk interface{}) error {
	a.categoryName = chunk.(tgbotapi.Update).Message.Text
	a.isChunksCollected = true

	return a.Execute()
}

func (a CategoryAppend) Next() (string, interface{}) {
	return "Введите имя новой категории. Например, цвет", &shop.CancelRow
}

type CategoryDelete struct {
	isDone            bool
	isChunksCollected bool
	shop              *shop.Shop
	categoryId        int
}

func (a *CategoryDelete) SetDone() {
	a.isDone = true
}

func (a *CategoryDelete) Execute(...interface{}) error {
	if !a.isChunksCollected {
		return ErrChunksIsNotCollected
	}

	if a.isDone {
		return ErrActionIsAlreadyDone
	}

	a.shop.DeleteCategory(a.categoryId)

	a.isDone = true

	return nil
}

func (a CategoryDelete) IsDone() bool {
	return a.isDone
}

func (a CategoryDelete) IsChunksCollected() bool {
	return a.isChunksCollected
}

func (a *CategoryDelete) AddChunk(chunk interface{}) (err error) {
	data := chunk.(string)

	a.categoryId, err = strconv.Atoi(strings.Split(data, "-")[1])

	if err != nil {
		return err
	}

	a.isChunksCollected = true

	return a.Execute()
}

func (a CategoryDelete) Next() (string, interface{}) {
	if len(a.shop.Categories()) == 0 {
		a.isDone = true
		return "Категории отсутствуют", &shop.AdminKeyboard
	}

	var rows [][]tgbotapi.InlineKeyboardButton

	for i, category := range a.shop.Categories() {
		rows = append(
			rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(category.Name, category.Name+"-"+strconv.Itoa(i)),
			),
		)
	}

	rows = append(rows, shop.CancelRow.InlineKeyboard[0])

	return "Выберите категорию", &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}
