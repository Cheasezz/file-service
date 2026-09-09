package service

import (
	"io"
	"os"

	"github.com/Cheasezz/fileService/internal/core"
	"github.com/Cheasezz/fileService/pkg/logger"
	"github.com/google/uuid"
)

type DB interface {
	CreateFile(file *core.FileInfo) (*os.File, error)
	OpenFile(file *core.FileInfo) (io.ReadCloser, error)
	GetAllFilesNames(userID string) ([]string, error)
	GetUserHashList(userID string) (map[string][]byte, error)
	UpdateHash(userID, filename string, hash []byte) error
}

type Service struct {
	log logger.Logger
	db  DB
}

func New(db DB, l logger.Logger) *Service {
	return &Service{db: db, log: l}
}

func (s *Service) OpenFile(userID, fileName string) (io.ReadCloser, error) {
	const op = "service.Download"
	log := s.log.With("op", op)

	fileInfo, err := core.NewFileInfo(userID, fileName)
	if err != nil {
		log.Error("cant create file info", err)
		return nil, err
	}

	f, err := s.db.OpenFile(fileInfo)
	if err != nil {
		log.Error("cant open file", err)
		return nil, err
	}

	return f, nil
}

func (s *Service) GetAllFilesNames(userID string) ([]string, error) {
	const op = "service.GetAllFilesNames"
	log := s.log.With("op", op)

	uID, err := uuid.Parse(userID)
	if err != nil {
		log.Error("cant parse uuid", err)
		return nil, core.ErrInvalidUUID
	}

	files, err := s.db.GetAllFilesNames(uID.String())
	if err != nil {
		log.Error("cant get all files names: %v", err)
		return nil, err
	}

	return files, nil
}

func (s *Service) FinishUpload(fw *fileWriter, userID, filename string) error {
	return s.db.UpdateHash(userID, filename, fw.Hash())
}
