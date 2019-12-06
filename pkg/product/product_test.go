package product

import (
	"reflect"
	"testing"
)

func TestNewProduct(t *testing.T) {
	type args struct {
		name        string
		description string
		photo       string
		price       uint
		id          int
	}
	tests := []struct {
		name    string
		args    args
		want    *Product
		wantErr bool
	}{
		{
			name:    "WrongName",
			args:    args{name: "[[[]]]]"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "WrongDescription",
			args:    args{name: "Тряпка", description: "]]]]"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "GoodCase",
			args:    args{name: "тряпка", description: "обычная тряпка", price: 1000},
			want:    &Product{name: "тряпка", description: "обычная тряпка", price: 1000},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProduct(tt.args.name, tt.args.description, tt.args.photo, tt.args.price, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}
