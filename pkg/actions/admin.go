package actions

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sschiz/dream-art/pkg/shop"
)

type AdminAppendAction struct {
	isDone            bool
	isChunksCollected bool
	shop              *shop.Shop
	adminName         string
}

func (a *AdminAppendAction) SetDone() {
	a.isDone = true
}

func (a *AdminAppendAction) Execute() error {
	if !a.isChunksCollected {
		return ErrChunksIsNotCollected
	}

	if a.isDone {
		return ErrActionIsAlreadyDone
	}

	err := a.shop.AppendAdmin(a.adminName)
	if err != nil {
		return err
	}

	a.isDone = true

	return nil
}

func (a AdminAppendAction) IsDone() bool {
	return a.isDone
}

func (a AdminAppendAction) IsChunksCollected() bool {
	return a.isChunksCollected
}

func (a *AdminAppendAction) AddChunk(chunk interface{}) error {
	a.adminName = chunk.(string)
	a.isChunksCollected = true

	return a.Execute()
}

func (a AdminAppendAction) Next() (string, *tgbotapi.InlineKeyboardMarkup) {
	return "Введите ник нового администратора. Например, @kek123", nil
}
