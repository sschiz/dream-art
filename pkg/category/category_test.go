package category

import (
	"testing"

	"github.com/sschiz/dream-art/pkg/product"
)

func TestCategory_AppendProduct(t *testing.T) {
	testProduct, _ := product.NewProduct("кк", "лл", "", 0, 123)
	type fields struct {
		name     string
		products []*product.Product
	}
	type args struct {
		name        string
		description string
		photo       string
		price       uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "EmptyProductSlice",
			fields: fields{
				products: nil,
			},
			args:    args{name: "тряпка", description: "просто тряпка"},
			wantErr: false,
		},
		{
			name: "Non-EmptyProductSlice",
			fields: fields{
				products: []*product.Product{testProduct},
			},
			args:    args{name: "тряпка", description: "просто тряпка"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Category{
				name:     tt.fields.name,
				products: tt.fields.products,
			}

			length := len(c.products)

			if err := c.AppendProduct(tt.args.name, tt.args.description, tt.args.photo, tt.args.price); (err != nil) != tt.wantErr {
				t.Errorf("Category.AppendProduct() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !(length == 0 && c.products[0].Id() == 1 || length != 0 && (c.products[length-1].Id()+1 == c.products[len(c.products)-1].Id())) {
				t.Errorf("Category.AppendProduct(). Appending is wrong")
			}
		})
	}
}
