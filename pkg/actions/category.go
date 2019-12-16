package actions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sschiz/dream-art/pkg/shop"
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
