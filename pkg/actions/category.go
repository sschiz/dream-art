package actions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sschiz/dream-art/pkg/shop"
	"strconv"
	"strings"
)

type CategoryAppendAction struct {
	isDone            bool
	isChunksCollected bool
	shop              *shop.Shop
	categoryName      string
}

func (a *CategoryAppendAction) SetDone() {
	a.isDone = true
}

func (a *CategoryAppendAction) Execute() error {
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

func (a CategoryAppendAction) IsDone() bool {
	return a.isDone
}

func (a CategoryAppendAction) IsChunksCollected() bool {
	return a.isChunksCollected
}

func (a *CategoryAppendAction) AddChunk(chunk interface{}) error {
	a.categoryName = chunk.(string)
	a.isChunksCollected = true

	return a.Execute()
}

func (a CategoryAppendAction) Next() (string, *tgbotapi.InlineKeyboardMarkup) {
	return "Введите имя новой категории. Например, цвет", nil
}

type CategoryDeleteAction struct {
	isDone            bool
	isChunksCollected bool
	shop              *shop.Shop
	categoryId        int
}

func (a *CategoryDeleteAction) SetDone() {
	a.isDone = true
}

func (a *CategoryDeleteAction) Execute() error {
	a.shop.DeleteCategory(a.categoryId)

	a.isDone = true

	return nil
}

func (a CategoryDeleteAction) IsDone() bool {
	return a.isDone
}

func (a CategoryDeleteAction) IsChunksCollected() bool {
	return a.isChunksCollected
}

func (a *CategoryDeleteAction) AddChunk(chunk interface{}) (err error) {
	data := chunk.(string)

	a.categoryId, err = strconv.Atoi(strings.Split(data, "-")[1])

	if err != nil {
		return err
	}

	a.isChunksCollected = true

	return a.Execute()
}

func (a CategoryDeleteAction) Next() (string, *tgbotapi.InlineKeyboardMarkup) {
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

	return "Выберите категорию", &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
}
