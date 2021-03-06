/*
 * (c) 2019, Matyushkin Alexander <sav3nme@gmail.com>
 * GNU General Public License v3.0+ (see COPYING or https://www.gnu.org/licenses/gpl-3.0.txt)
 */

package category

import (
	"github.com/sschiz/dream-art/pkg/product"
)

// Category is struct which provides typical shop category
type Category struct {
	Name     string             // name of category
	products []*product.Product // list of products that is part of category
}

// NewCategory creates category
func NewCategory(name string) *Category {
	return &Category{Name: name, products: nil}
}

// AppendProduct appends product into internal slice
func (c *Category) AppendProduct(name, description, photo string, price uint) (err error) {
	var item *product.Product
	if len(c.products) == 0 {
		item, err = product.NewProduct(name, description, photo, price, 1)

		if err != nil {
			return err
		}
	} else {
		item, err = product.NewProduct(name, description, photo, price, c.products[len(c.products)-1].ID+1)

		if err != nil {
			return err
		}
	}

	c.products = append(c.products, item)

	return nil
}

// Products returns product list
func (c Category) Products() []*product.Product {
	return c.products
}

// DeleteProduct removes product from list by i
func (c *Category) DeleteProduct(i int) {
	c.products = append(c.products[:i], c.products[i+1:]...)
}
