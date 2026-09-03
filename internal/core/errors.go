package core

import "errors"

var (
	ErrInvalidUUID   = errors.New("field user_id must be valid uuid")
	ErrEmptyFileName = errors.New("empty file name")
	ErrFirstMessage  = errors.New("first message must be file info")
	ErrInternal      = errors.New("something went wrong")
)
