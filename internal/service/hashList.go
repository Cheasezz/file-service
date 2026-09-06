package service

import (
	"github.com/Cheasezz/fileService/internal/core"
	"github.com/google/uuid"
)

type hashList map[string]string

func (m hashList) CompareHash(filename, userHash string) bool {
	serverHash := m[filename]

	return serverHash == userHash
}

func (s *Service) UserHashList(userID string) (*hashList, error) {
	const op = "service.UserManifest"
	log := s.log.With("op", op)

	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, core.ErrInvalidUUID
	}

	rowHl, err := s.db.GetUserHashList(userID)
	if err != nil {
		log.Error("cant get user manifest: %v", err)
		return nil, err
	}

	hl := hashList(rowHl)

	return &hl, nil
}
