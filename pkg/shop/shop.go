package shop

import (
	"errors"
	"sort"
	"strings"

	"github.com/sschiz/dream-art/pkg/category"
)

type Shop struct {
	categories []category.Category
	admins     []string
	Syncer     *Syncer
}

var (
	ErrAdminAlreadyExists = errors.New("such an admin already exists")
	ErrWrongNickname      = errors.New("wrong nickname")
)

// AppendAdmin appends new admin into slice by admin's nickname.
// Admins's nickname should start with '@'
func (s *Shop) AppendAdmin(name string) error {
	if !strings.HasPrefix(name, "@") {
		return ErrWrongNickname
	}

	i := sort.SearchStrings(s.admins, name)

	if i < len(s.admins) && s.admins[i] == name {
		return ErrAdminAlreadyExists
	} else {
		s.admins = append(s.admins[:i], append([]string{name}, s.admins[i:]...)...)
	}

	return nil
}
