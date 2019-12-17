package actions

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sschiz/dream-art/pkg/product"
	"github.com/sschiz/dream-art/pkg/shop"
)

type ProductAppendAction struct {
	categoryId        int
	product           *product.Product
	shop              *shop.Shop
	isDone            bool
	isChunksCollected bool
	step              int
}

func (a *ProductAppendAction) SetDone() {
	a.isDone = true
}

func (a *ProductAppendAction) Execute() error {
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

func (a ProductAppendAction) IsDone() bool {
	return a.isDone
}

func (a ProductAppendAction) IsChunksCollected() bool {
	return a.isChunksCollected
}

func (a *ProductAppendAction) AddChunk(chunk interface{}) (err error) {
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

func (a ProductAppendAction) Next() (string, interface{}) {
	switch a.step {
	case 0:
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

		return "Выберите категорию, в которую хотите добавить новый продукт", &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: rows}
	case 1:
		return "Отправьте имя нового продукта. Строго с маленькой буквы. Например, бумага", nil
	case 2:
		return "Отправьте описание продукта", nil
	case 3:
		return "Отправьте цену продукта до копеек целым числом. Например, 1000 - это 100 рублей или 9005 - это 900 рублей и 5 копеек", nil
	case 4:
		return "Отправьте фотографию продукта", nil
	default:
		return "", nil
	}
}
