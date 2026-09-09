package service

import (
	"crypto/subtle"

	"github.com/Cheasezz/fileService/internal/core"
	"github.com/google/uuid"
)

type hashList map[string][]byte

func (m hashList) CompareHash(filename string, userHash []byte) bool {
	serverHash := m[filename]

	return subtle.ConstantTimeCompare(serverHash, userHash) == 1
}

func (s *Service) UserHashList(userID string) (*hashList, error) {
	const op = "service.UserHashList"
	log := s.log.With("op", op)

	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, core.ErrInvalidUUID
	}

	rowHl, err := s.db.GetUserHashList(userID)
	if err != nil {
		log.Error("cant get user hash list: %v", err)
		return nil, err
	}

	hl := hashList(rowHl)

	return &hl, nil
}
