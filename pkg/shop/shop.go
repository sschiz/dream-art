package shop

import (
	"errors"
	"strings"
	"sync"

	"github.com/sschiz/dream-art/pkg/category"
)

type Shop struct {
	categories []*category.Category
	Admins     map[string]int64
	syncer     *Syncer
	Mu         *sync.RWMutex
}

var (
	ErrAdminAlreadyExists = errors.New("such an admin already exists")
	ErrWrongNickname      = errors.New("wrong nickname")
	ErrAdminDoesNotExist  = errors.New("admin does not exist")
)

// NewShop creates shop
func NewShop(syncer *Syncer) (*Shop, error) {
	shop := new(Shop)
	shop.Admins = make(map[string]int64)
	shop.Mu = new(sync.RWMutex)
	shop.syncer = syncer
	err := shop.Sync()

	if err != nil {
		return nil, err
	}

	return shop, nil
}

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

	s.Mu.RLock()
	_, ok := s.Admins[name]
	s.Mu.RUnlock()

	if !ok {
		return ErrAdminDoesNotExist
	}

	s.Mu.Lock()
	delete(s.Admins, name)
	s.Mu.Unlock()

	return nil
}

// AppendCategory creates new category and appends into category slice
func (s *Shop) AppendCategory(name string) {
	s.Mu.Lock()
	s.categories = append(s.categories, category.NewCategory(name))
	s.Mu.Unlock()
}

// Categories return all of category list
func (s Shop) Categories() []*category.Category {
	return s.categories
}

// DeleteCategory removes category from list
func (s *Shop) DeleteCategory(i int) {
	s.Mu.Lock()
	s.categories = append(s.categories[:i], s.categories[i+1:]...)
	s.Mu.Unlock()
}

func (s Shop) IsAdmin(name string) bool {
	s.Mu.RLock()
	_, ok := s.Admins[name]
	s.Mu.RUnlock()
	return ok
}

func (s *Shop) Sync() error {
	return s.syncer.Sync(s)
}

// AddChatID adds chatID for admin
func (s *Shop) AddChatID(name string, chatID int64) {
	s.Mu.Lock()
	s.Admins[name] = chatID
	s.Mu.Unlock()
}
