package service

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"os"

	"github.com/Cheasezz/fileService/internal/core"
)

type fileWriter struct {
	File   *os.File
	hasher hash.Hash
}

func (s *Service) CreateFile(userID, fileName string) (*fileWriter, error) {
	op := "service.CreateFile"
	log := s.log.With("op", op)

	fileInfo, err := core.NewFileInfo(userID, fileName)
	if err != nil {
		log.Error("cant create file info: ", err)
		return nil, core.ErrInternal
	}

	f, err := s.db.CreateFile(fileInfo)
	if err != nil {
		log.Error("cant create file: ", err)
		return nil, core.ErrInternal
	}

	return &fileWriter{File: f, hasher: sha256.New()}, nil
}

func (fw *fileWriter) Write(data []byte) (int, error) {
	fw.hasher.Write(data)
	return fw.File.Write(data)
}

func (fw *fileWriter) Hash() string {
	return hex.EncodeToString([]byte(fw.hasher.Sum(nil)))
}
