package category

import (
	"github.com/sschiz/dream-art/pkg/product"
)

// Category is struct which provides typical shop category
type Category struct {
	name     string             // name of category
	products []*product.Product // list of products that is part of category
}

// NewCategory creates category
func NewCategory(name string) *Category {
	return &Category{name: name, products: nil}
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
		item, err = product.NewProduct(name, description, photo, price, c.products[len(c.products)-1].Id()+1)

		if err != nil {
			return err
		}
	}

	c.products = append(c.products, item)

	return nil
}
