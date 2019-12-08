package actions

import "github.com/sschiz/dream-art/pkg/shop"

type AdminAppendAction struct {
	isDone            bool
	isChunksCollected bool
	shop              *shop.Shop
	adminName         string
}

func (a AdminAppendAction) SetDone() {
	a.isDone = true
}

func (a AdminAppendAction) Execute() error {
	if !a.isChunksCollected {
		return ErrChunksIsNotCollected
	}

	if a.isDone {
		return ErrActionIsAlreadyDone
	}

	err := a.shop.AppendAdmin(a.adminName)
	if err != nil {
		return err
	}

	a.isDone = true

	return nil
}

func (a AdminAppendAction) IsDone() bool {
	return a.isDone
}

func (a AdminAppendAction) IsChunksCollected() bool {
	return a.isChunksCollected
}

func (a AdminAppendAction) AddChunk(chunk interface{}, isLastChunk bool) error {
	a.adminName = chunk.(string)
	a.isChunksCollected = true

	return nil
}
