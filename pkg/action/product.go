/*
 * (c) 2019, Matyushkin Alexander <sav3nme@gmail.com>
 * GNU General Public License v3.0+ (see COPYING or https://www.gnu.org/licenses/gpl-3.0.txt)
 */

package action

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sschiz/dream-art/pkg/category"
	"github.com/sschiz/dream-art/pkg/product"
	"github.com/sschiz/dream-art/pkg/shop"
)

type ProductAppend struct {
	categoryId        int
	product           *product.Product
	shop              *shop.Shop
	isDone            bool
	isChunksCollected bool
	step              int
}

func (a *ProductAppend) SetDone() {
	a.isDone = true
}

func (a *ProductAppend) Execute(...interface{}) error {
	if !a.isChunksCollected {
		return ErrChunksIsNotCollected
	}

	if a.isDone {
		return ErrActionIsAlreadyDone
	}

	err := a.shop.Categories()[a.categoryId].AppendProduct(a.product.Name, a.product.Description, a.product.Photo, a.product.Price)

	if err != nil {
		return err
	}

	a.isDone = true

	return nil
}

func (a ProductAppend) IsDone() bool {
	return a.isDone
}

func (a ProductAppend) IsChunksCollected() bool {
	return a.isChunksCollected
}

func (a *ProductAppend) AddChunk(chunk interface{}) (err error) {
	switch a.step {
	case 0:
		data := chunk.(string)

		a.categoryId, err = strconv.Atoi(strings.Split(data, "-")[1])

		if err != nil {
			return err
		}
	case 1:
		update := chunk.(tgbotapi.Update)
		a.product.Name = update.Message.Text
	case 2:
		update := chunk.(tgbotapi.Update)
		a.product.Description = update.Message.Text
	case 3:
		update := chunk.(tgbotapi.Update)
		price, err := strconv.Atoi(update.Message.Text)

		if err != nil {
			return err
		}

		a.product.Price = uint(price)
	case 4:
		a.product.Photo = (*chunk.(tgbotapi.Update).Message.Photo)[0].FileID
		a.isChunksCollected = true

		return a.Execute()
	}

	a.step++
	return nil
}

func (a ProductAppend) Next() (string, interface{}) {
	switch a.step {
	case 0:
		if len(a.shop.Categories()) == 0 {
			a.isDone = true
			return "Категории отсутствуют", &shop.AdminKeyboard
		}

		var rows [][]tgbotapi.InlineKeyboardButton

		for i, c := range a.shop.Categories() {
			rows = append(
				rows, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(c.Name, c.Name+"-"+strconv.Itoa(i)),
				),
			)
		}

		rows = append(rows, shop.CancelRow.InlineKeyboard[0])

		return "Выберите категорию, в которую хотите добавить новый продукт", &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
	case 1:
		return "Отправьте имя нового продукта. Строго с маленькой буквы. Например, бумага", &shop.CancelRow
	case 2:
		return "Отправьте описание продукта", shop.CancelRow
	case 3:
		return "Отправьте цену продукта до копеек целым числом. Например, 1000 - это 100 рублей или 9005 - это 900 рублей и 5 копеек", &shop.CancelRow
	case 4:
		return "Отправьте фотографию продукта", &shop.CancelRow
	default:
		return "", nil
	}
}

type ProductDelete struct {
	shop              *shop.Shop
	category          *category.Category
	productID         int
	isDone            bool
	isChunksCollected bool
	step              int
}

func (a *ProductDelete) SetDone() {
	a.isDone = true
}

func (a *ProductDelete) Execute(...interface{}) error {
	if !a.isChunksCollected {
		return ErrChunksIsNotCollected
	}

	if a.isDone {
		return ErrActionIsAlreadyDone
	}

	a.category.DeleteProduct(a.productID)

	a.isDone = true

	return nil
}

func (a ProductDelete) IsDone() bool {
	return a.isDone
}

func (a ProductDelete) IsChunksCollected() bool {
	return a.isChunksCollected
}

func (a *ProductDelete) AddChunk(chunk interface{}) (err error) {
	data := chunk.(string)

	switch a.step {
	case 0:
		i, err := strconv.Atoi(strings.Split(data, "-")[1])

		if err != nil {
			return err
		}

		a.category = a.shop.Categories()[i]
	case 1:
		a.productID, err = strconv.Atoi(strings.Split(data, "-")[1])

		if err != nil {
			return err
		}

		a.isChunksCollected = true

		return a.Execute()
	}

	a.step++
	return nil
}

func (a ProductDelete) Next() (string, interface{}) {
	switch a.step {
	case 0:
		if len(a.shop.Categories()) == 0 {
			a.isDone = true
			return "Категории отсутствуют", &shop.AdminKeyboard
		}

		var rows [][]tgbotapi.InlineKeyboardButton

		for i, c := range a.shop.Categories() {
			rows = append(
				rows, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(c.Name, c.Name+"-"+strconv.Itoa(i)),
				),
			)
		}

		rows = append(rows, shop.CancelRow.InlineKeyboard[0])

		return "Выберите категорию, из которой хотите добавить новый продукт", &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
	case 1:
		if len(a.category.Products()) == 0 {
			a.isDone = true
			return "Продукты отсутствуют", &shop.AdminKeyboard
		}

		var rows [][]tgbotapi.InlineKeyboardButton

		for i, p := range a.category.Products() {
			rows = append(
				rows, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(p.Name, p.Name+"-"+strconv.Itoa(i)),
				),
			)
		}

		rows = append(rows, shop.CancelRow.InlineKeyboard[0])

		return "Выберите удаляемый продукт", &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
	default:
		return "", nil
	}
}
