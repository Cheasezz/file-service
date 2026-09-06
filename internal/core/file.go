package core

import (
	"github.com/google/uuid"
)

type FileInfo struct {
	UserID uuid.UUID
	Name   string
}

type SyncDecision struct {
	Filename   string
	NeedUpload bool
}

func NewFileInfo(userID, name string) (*FileInfo, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, ErrInvalidUUID
	}

	if name == "" {
		return nil, ErrEmptyFileName
	}

	return &FileInfo{UserID: id, Name: name}, nil
}
