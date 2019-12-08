package actions

import "errors"

// Action is abstract which provides action handling
type Action interface {
	Execute() error                                     // do what you need
	IsDone() bool                                       // check if action is done
	SetDone()                                           // set action that is done, but if all chunks is collected
	AddChunk(chunk interface{}, isLastChunk bool) error // add chunk that is needed
	IsChunksCollected() bool                            // returns all chunks collected
}

var (
	ErrChunksIsNotCollected = errors.New("chunks is not collected")
	ErrActionIsAlreadyDone  = errors.New("action is already done")
)
