package shop

import (
	"errors"
	"strings"

	"github.com/sschiz/dream-art/pkg/category"
)

type Shop struct {
	categories []*category.Category
	Admins     map[string]int64
	syncer     *Syncer
}

var (
	ErrAdminAlreadyExists = errors.New("such an admin already exists")
	ErrWrongNickname      = errors.New("wrong nickname")
	ErrAdminDoesNotExist  = errors.New("admin does not exist")
)

// AppendAdmin appends new admin into slice by admin's nickname.
// Admin's nickname should start with '@'
func (s *Shop) AppendAdmin(name string) error {
	if !strings.HasPrefix(name, "@") {
		return ErrWrongNickname
	}

	name = name[1:] // delete '@'

	_, ok := s.Admins[name]

	if ok {
		return ErrAdminAlreadyExists
	} else {
		s.Admins[name] = 0
	}

	return nil
}

// DeleteAdmin removes admin nickname from list of admins
// Admin's nickname should start with '@'
func (s *Shop) DeleteAdmin(name string) error {
	if !strings.HasPrefix(name, "@") {
		return ErrWrongNickname
	}

	name = name[1:]

	_, ok := s.Admins[name]

	if !ok {
		return ErrAdminDoesNotExist
	}

	delete(s.Admins, name)

	return nil
}

// AppendCategory creates new category and appends into category slice
func (s *Shop) AppendCategory(name string) {
	s.categories = append(s.categories, category.NewCategory(name))
}

// Categories return all of category list
func (s Shop) Categories() []*category.Category {
	return s.categories
}

// DeleteCategory removes category from list
func (s *Shop) DeleteCategory(i int) {
	s.categories = append(s.categories[:i], s.categories[i+1:]...)
}

func (s Shop) IsAdmin(name string) bool {
	_, ok := s.Admins[name]
	return ok
}

func (s *Shop) Sync() error {
	return s.syncer.Sync(s)
}
