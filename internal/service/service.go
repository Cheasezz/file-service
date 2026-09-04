package service

import (
	"io"
	"os"

	"github.com/Cheasezz/fileService/internal/core"
	"github.com/Cheasezz/fileService/pkg/logger"
)

type DB interface {
	CreateFile(file *core.FileInfo) (*os.File, error)
	OpenFile(file *core.FileInfo) (io.ReadCloser, error)
}

type Service struct {
	log logger.Logger
	db  DB
}

func New(db DB, l logger.Logger) *Service {
	return &Service{db: db, log: l}
}

func (s *Service) CreateFile(userID, fileName string) (*os.File, error) {
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

	return f, nil
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
