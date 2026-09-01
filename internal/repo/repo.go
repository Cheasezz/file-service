package repo

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type FileSystem struct {
	path string
}

func New() (*FileSystem, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	storagePath := filepath.Join(homeDir, ".fileService/uploads")

	err = os.MkdirAll(storagePath, 0o755)
	if err != nil {
		return nil, err
	}

	return &FileSystem{path: storagePath}, nil
}

func (fs *FileSystem) CreateFile(fileName string) (*os.File, error) {
	var f *os.File

	f, err := os.Create(fs.path + "/" + uuid.NewString() + "_" + fileName)
	if err != nil {
		return f, err
	}

	return f, nil
}
