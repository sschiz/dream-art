/*
 * (c) 2019, Matyushkin Alexander <sav3nme@gmail.com>
 * GNU General Public License v3.0+ (see COPYING or https://www.gnu.org/licenses/gpl-3.0.txt)
 */

package product

import "errors"

// ErrWrongText is error which provides error for ids
var ErrWrongId = errors.New("ID is wrong")

// Product is typical struct of product
type Product struct {
	ID          int
	Name        string
	Description string
	Price       uint
	Photo       string
}

// NewProduct creates product checking the name and description
func NewProduct(name, description, photo string, price uint, id int) (*Product, error) {
	if id < 0 {
		return nil, ErrWrongId
	}

	return &Product{Name: name, Description: description, Price: price, Photo: photo, ID: id}, nil
}
