package actions

import (
	"errors"

	"github.com/sschiz/dream-art/pkg/product"
	"github.com/sschiz/dream-art/pkg/shop"
)

// Action is abstract which provides action handling
type Action interface {
	Execute() error                   // do what you need
	IsDone() bool                     // check if action is done
	SetDone()                         // set action that is done, but if all chunks is collected
	AddChunk(chunk interface{}) error // add chunk that is needed
	IsChunksCollected() bool          // returns all chunks collected
	Next() (string, interface{})      // return text and keyboard for next chunk
}

var (
	ErrChunksIsNotCollected   = errors.New("chunks is not collected")
	ErrActionIsAlreadyDone    = errors.New("action is already done")
	ErrActionTypeDoesNotExist = errors.New("such an action type doesn't exist")
	ErrObjectDoesNotExist     = errors.New("such an object doesn't exist")
)

// NewAction creates new action
func NewAction(actionType, object string, shop *shop.Shop) (Action, error) {
	switch actionType {
	case "append":
		switch object {
		case "admin":
			return &AdminAppendAction{shop: shop}, nil
		case "category":
			return &CategoryAppendAction{shop: shop}, nil
		case "product":
			return &ProductAppendAction{shop: shop, product: &product.Product{}}, nil
		default:
			return nil, ErrObjectDoesNotExist
		}
	case "delete":
		switch object {
		case "admin":
			return &AdminDeleteAction{shop: shop}, nil
		case "category":
			return &CategoryDeleteAction{shop: shop}, nil
		case "product":
			return &ProductDeleteAction{shop: shop}, nil
		default:
			return nil, ErrObjectDoesNotExist
		}
	case "change":
		//
	default:
		return nil, ErrActionTypeDoesNotExist
	}

	return nil, nil
}
