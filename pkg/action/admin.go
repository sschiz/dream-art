/*
 * (c) 2019, Matyushkin Alexander <sav3nme@gmail.com>
 * GNU General Public License v3.0+ (see COPYING or https://www.gnu.org/licenses/gpl-3.0.txt)
 */

package action

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sschiz/dream-art/pkg/shop"
)

type AdminAppend struct {
	isDone            bool
	isChunksCollected bool
	shop              *shop.Shop
	adminName         string
}

func (a *AdminAppend) SetDone() {
	a.isDone = true
}

func (a *AdminAppend) Execute(...interface{}) error {
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

func (a AdminAppend) IsDone() bool {
	return a.isDone
}

func (a AdminAppend) IsChunksCollected() bool {
	return a.isChunksCollected
}

func (a *AdminAppend) AddChunk(chunk interface{}) error {
	a.adminName = chunk.(tgbotapi.Update).Message.Text
	a.isChunksCollected = true

	return a.Execute()
}

func (a AdminAppend) Next() (string, interface{}) {
	return "Введите ник нового администратора. Например, @kek123", &shop.CancelRow
}

type AdminDelete struct {
	isDone            bool
	isChunksCollected bool
	shop              *shop.Shop
	adminName         string
}

func (a *AdminDelete) SetDone() {
	a.isDone = true
}

func (a *AdminDelete) Execute(...interface{}) error {
	if !a.isChunksCollected {
		return ErrChunksIsNotCollected
	}

	if a.isDone {
		return ErrActionIsAlreadyDone
	}

	err := a.shop.DeleteAdmin(a.adminName)
	if err != nil {
		return err
	}

	a.isDone = true

	return nil
}

func (a AdminDelete) IsDone() bool {
	return a.isDone
}

func (a AdminDelete) IsChunksCollected() bool {
	return a.isChunksCollected
}

func (a *AdminDelete) AddChunk(chunk interface{}) error {
	a.adminName = chunk.(tgbotapi.Update).Message.Text
	a.isChunksCollected = true

	return a.Execute()
}

func (a AdminDelete) Next() (string, interface{}) {
	return "Введите ник удаляемого администратора. Например, @kek123", &shop.CancelRow
}
