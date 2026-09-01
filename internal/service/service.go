package service

import (
	"os"

	"github.com/Cheasezz/fileService/pkg/logger"
)

type DB interface {
	CreateFile(fileName string) (*os.File, error)
}

type Service struct {
	log logger.Logger
	db  DB
}

func New(db DB, l logger.Logger) *Service {
	return &Service{db: db, log: l}
}

func (s *Service) CreateFile(fileName string) (*os.File, error) {
	op := "service.CreateFile"
	log := s.log.With("op", op)

	f, err := s.db.CreateFile(fileName)
	if err != nil {
		log.Error("cant create file: ", err)
		return nil, err
	}

	return f, nil
}
