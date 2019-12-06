package product

import "regexp"

import "errors"

// ErrWrongTitle is error which provides error for titles
var ErrWrongTitle = errors.New("title is wrong")

// ErrWrongText is error which provides error for texts
var ErrWrongText = errors.New("text is wrong")

// Product is typical struct of product
type Product struct {
	id          uint
	name        string
	description string
	price       uint
	photo       string
}

// NewProduct creates product checking the name and description
func NewProduct(name, description, photo string, price, id uint) (*Product, error) {
	matched, _ := regexp.MatchString(`^([a-zа-я0-9 _.])+$`, name)

	if !matched {
		return nil, ErrWrongTitle
	}

	matched, _ = regexp.MatchString(`^([\w\sА-Яа-я0-9.,!?:;])+$`, description)

	if !matched {
		return nil, ErrWrongText
	}

	return &Product{name: name, description: description, price: price, photo: photo, id: id}, nil
}
