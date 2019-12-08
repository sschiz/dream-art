package shop

import "github.com/sschiz/dream-art/pkg/category"

type Shop struct {
	categories []category.Category
	admins     []string
	Syncer     *Syncer
}
