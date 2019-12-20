/*
 * (c) 2019, Matyushkin Alexander <sav3nme@gmail.com>
 * GNU General Public License v3.0+ (see COPYING or https://www.gnu.org/licenses/gpl-3.0.txt)
 */

package action

import (
	"errors"

	"github.com/sschiz/dream-art/pkg/product"
	"github.com/sschiz/dream-art/pkg/shop"
)

// Action is abstract which provides action handling
type Action interface {
	Execute(args ...interface{}) error // do what you need
	IsDone() bool                      // check if action is done
	SetDone()                          // set action that is done, but if all chunks is collected
	AddChunk(chunk interface{}) error  // add chunk that is needed
	IsChunksCollected() bool           // returns all chunks collected
	Next() (string, interface{})       // return text and keyboard for next chunk
}

var (
	ErrChunksIsNotCollected   = errors.New("chunks is not collected")
	ErrActionIsAlreadyDone    = errors.New("action is already done")
	ErrActionTypeDoesNotExist = errors.New("such an action type doesn't exist")
	ErrObjectDoesNotExist     = errors.New("such an object doesn't exist")
)

// New creates new action
func New(actionType, object string, shop *shop.Shop) (Action, error) {
	switch actionType {
	case "append":
		switch object {
		case "admin":
			return &AdminAppend{shop: shop}, nil
		case "category":
			return &CategoryAppend{shop: shop}, nil
		case "product":
			return &ProductAppend{shop: shop, product: &product.Product{}}, nil
		default:
			return nil, ErrObjectDoesNotExist
		}
	case "delete":
		switch object {
		case "admin":
			return &AdminDelete{shop: shop}, nil
		case "category":
			return &CategoryDelete{shop: shop}, nil
		case "product":
			return &ProductDelete{shop: shop}, nil
		default:
			return nil, ErrObjectDoesNotExist
		}
	case "buy":
		return &Buy{shop: shop}, nil
	case "change":
		//
	default:
		return nil, ErrActionTypeDoesNotExist
	}

	return nil, nil
}
